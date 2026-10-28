package setup

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// fakeWSL answers the wizard's calls from fields the test sets.
type fakeWSL struct {
	state     wsl.State
	installed []string
	restarts  int
	scripts   []string
	certPEM   []byte
	// answer overrides the default reply for a root script; ok false falls through.
	answer       func(script string) (out string, ok bool)
	installErr   error
	exportErr    error
	exported     []string
	unregistered []string
	// calls records the order of backup, removal and host actions.
	calls []string
}

func (f *fakeWSL) Export(_ context.Context, _, file string) error {
	f.calls = append(f.calls, "export")
	f.exported = append(f.exported, file)
	if f.exportErr != nil {
		// A real export leaves a partial file behind when it fails.
		_ = os.WriteFile(file, []byte("partial"), 0o600)
		return f.exportErr
	}
	return os.WriteFile(file, []byte("tar"), 0o600)
}

func (f *fakeWSL) Unregister(_ context.Context, distro string) error {
	f.calls = append(f.calls, "unregister")
	f.unregistered = append(f.unregistered, distro)
	f.state = wsl.StateNotInstalled
	return nil
}

func (f *fakeWSL) Available(context.Context) error { return nil }

func (f *fakeWSL) DistroState(context.Context, string) (wsl.State, error) {
	return f.state, nil
}

func (f *fakeWSL) Install(_ context.Context, _, name, _ string) error {
	if f.installErr != nil {
		return f.installErr
	}
	f.installed = append(f.installed, name)
	f.state = wsl.StateRunning
	return nil
}

func (f *fakeWSL) Terminate(context.Context, string) error { f.restarts++; return nil }
func (f *fakeWSL) Start(context.Context, string) error     { return nil }

func (f *fakeWSL) Script(_ context.Context, _, _, script string) (string, error) {
	f.scripts = append(f.scripts, script)
	if f.answer != nil {
		if out, ok := f.answer(script); ok {
			return out, nil
		}
	}
	switch {
	case strings.Contains(script, "ps -p 1"):
		return "systemd\n", nil
	case strings.HasPrefix(script, "cat ") && strings.Contains(script, "portainer.crt"):
		return string(f.certPEM), nil
	case script == portainerStateScript:
		return "missing\n", nil
	}
	return "", nil
}

// portainerFake is a TLS handler that behaves like the Portainer endpoints
// the wizard uses.
type portainerFake struct {
	mu          sync.Mutex
	adminExists bool
	tokenSeen   string
	adminBody   map[string]string
	loginOK     bool
}

func (p *portainerFake) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/system/status", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"Version":"2.45.2"}`))
	})
	mux.HandleFunc("GET /api/users/admin/check", func(w http.ResponseWriter, _ *http.Request) {
		p.mu.Lock()
		defer p.mu.Unlock()
		if p.adminExists {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
	mux.HandleFunc("POST /api/users/admin/init", func(w http.ResponseWriter, r *http.Request) {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.tokenSeen = r.Header.Get("X-Setup-Token")
		_ = json.NewDecoder(r.Body).Decode(&p.adminBody)
		p.adminExists = true
		_, _ = w.Write([]byte(`{}`))
	})
	mux.HandleFunc("POST /api/auth", func(w http.ResponseWriter, r *http.Request) {
		var b map[string]string
		_ = json.NewDecoder(r.Body).Decode(&b)
		p.mu.Lock()
		defer p.mu.Unlock()
		if p.adminExists && b["username"] == p.adminBody["Username"] && b["password"] == p.adminBody["Password"] {
			p.loginOK = true
			_, _ = w.Write([]byte(`{"jwt":"x"}`))
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
	})
	return mux
}

type harness struct {
	w      *Wizard
	fake   *fakeWSL
	srv    *portainerFake
	events []Event
	binDir string
	path   []string
}

func newHarness(t *testing.T, cfg Config) *harness {
	t.Helper()
	pf := &portainerFake{}
	ts := httptest.NewTLSServer(pf.handler())
	t.Cleanup(ts.Close)
	h := &harness{
		srv:    pf,
		fake:   &fakeWSL{state: wsl.StateNotInstalled, certPEM: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ts.Certificate().Raw})},
		binDir: filepath.Join(t.TempDir(), "bin"),
	}
	cfg.BinDir = h.binDir
	cfg.AddToPath = true
	tokens := []string{"aa11", "bb22", "cc33"}
	h.w = &Wizard{
		WSL:          h.fake,
		Config:       cfg.withDefaults(),
		Emit:         func(e Event) { h.events = append(h.events, e) },
		AddPath:      func(dir string) error { h.path = append(h.path, dir); return nil },
		APIBase:      ts.URL + "/api",
		Wait:         time.Millisecond,
		ReadyTimeout: 2 * time.Second,
		NewToken: func() (string, error) {
			tok := tokens[0]
			tokens = tokens[1:]
			return tok, nil
		},
	}
	return h
}

func (h *harness) status(step StepID) Status {
	var got Status
	for _, e := range h.events {
		if e.Step == step && e.Status != StatusLog && e.Status != StatusRunning {
			got = e.Status
		}
	}
	return got
}

// admin is a valid account; its secret is built at run time so no literal
// credential sits in the source.
var admin = Credentials{Username: "tester", Password: strings.Repeat("pw", 8)}

func TestRunFreshInstall(t *testing.T) {
	h := newHarness(t, Config{Distro: "dz-throwaway", Port: 19443})
	if err := h.w.Run(t.Context(), admin); err != nil {
		t.Fatal(err)
	}
	for _, s := range Steps {
		if got := h.status(s.ID); got != StatusDone {
			t.Errorf("step %s: %q, want done", s.ID, got)
		}
	}
	if !slices.Equal(h.fake.installed, []string{"dz-throwaway"}) {
		t.Fatalf("installed %q", h.fake.installed)
	}
	if h.srv.tokenSeen != "aa11" {
		t.Fatalf("init sent token %q, want the one the container started with", h.srv.tokenSeen)
	}
	if h.srv.adminBody["Username"] != "tester" || !h.srv.loginOK {
		t.Fatalf("admin %v, login ok %v", h.srv.adminBody, h.srv.loginOK)
	}

	var runs []string
	for _, s := range h.fake.scripts {
		if strings.Contains(s, "docker run -d --name portainer") {
			runs = append(runs, s)
		}
	}
	if len(runs) != 2 {
		t.Fatalf("%d container creations, want 2 (with the token, then without)", len(runs))
	}
	if !strings.Contains(runs[0], "'aa11'") || !strings.Contains(runs[0], "--env-file") || strings.Contains(runs[0], "docker rm -f") {
		t.Errorf("first creation: %s", runs[0])
	}
	if strings.Contains(runs[1], "aa11") || strings.Contains(runs[1], "env-file") || !strings.Contains(runs[1], "docker rm -f portainer") {
		t.Errorf("recreation keeps the token or skips the removal: %s", runs[1])
	}
	if !strings.Contains(runs[0], "-p 127.0.0.1:19443:9443") {
		t.Errorf("port not applied: %s", runs[0])
	}

	// Neither the token nor the password may reach the window or the distro.
	for _, e := range h.events {
		if strings.Contains(e.Message, "aa11") || strings.Contains(e.Message, admin.Password) {
			t.Fatalf("secret in event %+v", e)
		}
	}
	for _, s := range h.fake.scripts {
		if strings.Contains(s, admin.Password) {
			t.Fatalf("password sent into the distro: %s", s)
		}
	}

	cmd, err := os.ReadFile(filepath.Join(h.binDir, "docker.cmd"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cmd), `"DOCKER_WSL_DISTRO=dz-throwaway"`) || !strings.Contains(string(cmd), "\r\n") {
		t.Errorf("wrapper: %q", cmd)
	}
	sh, err := os.ReadFile(filepath.Join(h.binDir, "docker"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(sh), "DOCKER_WSL_DISTRO:=dz-throwaway") || strings.Contains(string(sh), "\r") {
		t.Errorf("bash wrapper: %q", sh)
	}
	if !slices.Equal(h.path, []string{h.binDir}) {
		t.Errorf("PATH additions %q", h.path)
	}
}

func TestRunSkipsExistingWork(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	h.srv.adminExists = true
	h.fake.answer = func(script string) (string, bool) {
		switch {
		case strings.Contains(script, "docker-ce"), strings.HasPrefix(script, "set -e\nif id -u"):
			return markerPresent + "\n", true
		case script == portainerStateScript:
			return "running\n", true
		}
		return "", false
	}
	if err := h.w.Run(t.Context(), admin); err != nil {
		t.Fatal(err)
	}
	for _, id := range []StepID{StepDistro, StepDocker, StepUser, StepPortainer, StepAdmin} {
		if got := h.status(id); got != StatusSkipped {
			t.Errorf("step %s: %q, want skipped", id, got)
		}
	}
	if len(h.fake.installed) != 0 || h.fake.restarts != 0 {
		t.Errorf("reinstalled: %q, restarts %d", h.fake.installed, h.fake.restarts)
	}
	for _, s := range h.fake.scripts {
		if strings.Contains(s, "docker run") {
			t.Fatalf("recreated the container: %s", s)
		}
	}
}

func TestAdminReopensFirstRunForAnOlderContainer(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	h.fake.answer = func(script string) (string, bool) {
		if script == portainerStateScript {
			return "running\n", true
		}
		return "", false
	}
	if err := h.w.Run(t.Context(), admin); err != nil {
		t.Fatal(err)
	}
	if h.srv.tokenSeen != "aa11" {
		t.Fatalf("init sent token %q, want the token made for the restarted container", h.srv.tokenSeen)
	}
	if !h.srv.loginOK {
		t.Fatal("no sign-in with the new admin")
	}
}

func TestRunStopsAtTheFirstFailure(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.installErr = errors.New("download failed")
	err := h.w.Run(t.Context(), admin)
	if err == nil || !strings.Contains(err.Error(), "download failed") {
		t.Fatalf("got %v", err)
	}
	if got := h.status(StepDistro); got != StatusFailed {
		t.Fatalf("distro step %q", got)
	}
	if got := h.status(StepDocker); got != "" {
		t.Fatalf("a later step ran: %q", got)
	}
	if len(h.fake.scripts) != 0 {
		t.Fatalf("scripts ran after the failure: %d", len(h.fake.scripts))
	}
}

func TestRunRejectsBadInput(t *testing.T) {
	long := strings.Repeat("q", 14)
	cases := map[string]struct {
		cfg  Config
		cr   Credentials
		want error
	}{
		"short password":  {Config{}, Credentials{Username: "a", Password: strings.Repeat("s", 5)}, ErrCredentials},
		"empty username":  {Config{}, Credentials{Password: long}, ErrCredentials},
		"shell in user":   {Config{User: "dev; rm -rf /"}, admin, ErrBadName},
		"shell in distro": {Config{Distro: "x$(id)"}, admin, ErrBadName},
		"bad port":        {Config{Port: 70000}, admin, ErrBadName},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			h := newHarness(t, Config{})
			h.w.Config = c.cfg
			if err := h.w.Run(t.Context(), c.cr); !errors.Is(err, c.want) {
				t.Fatalf("got %v, want %v", err, c.want)
			}
			if len(h.fake.scripts) != 0 || len(h.events) != 0 {
				t.Fatal("work started before the input was checked")
			}
		})
	}
}

func TestPortainerRunScriptKeepsTheTokenOffTheCommandLine(t *testing.T) {
	script := portainerRunScript(9443, "deadbeef", false)
	_, runLine, _ := strings.Cut(script, "docker run")
	if strings.Contains(runLine, "deadbeef") {
		t.Fatalf("token on the docker run line: %s", runLine)
	}
	if !strings.Contains(script, "trap 'rm -f /run/dz-portainer.env' EXIT") {
		t.Fatalf("env file not removed: %s", script)
	}
}
