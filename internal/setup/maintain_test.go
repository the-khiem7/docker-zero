package setup

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// runningPortainer answers the container state and the image comparison.
func runningPortainer(image string) func(string) (string, bool) {
	return func(script string) (string, bool) {
		switch script {
		case portainerStateScript:
			return "running\n", true
		case currentScript:
			return image + "\n", true
		}
		return "", false
	}
}

func recreations(h *harness) []string {
	var runs []string
	for _, s := range h.fake.scripts {
		if strings.Contains(s, "docker run -d --name portainer") {
			runs = append(runs, s)
		}
	}
	return runs
}

func TestUpdatePortainerKeepsACurrentContainer(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	h.fake.answer = runningPortainer("current")
	updated, err := h.w.UpdatePortainer(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if updated || len(recreations(h)) != 0 {
		t.Fatalf("updated %v, recreations %d; a current image must not restart Portainer", updated, len(recreations(h)))
	}
	if !contains(h.fake.scripts, "docker pull "+portainerImage) {
		t.Fatal("the image was not pulled")
	}
}

func TestUpdatePortainerRecreatesOnANewImage(t *testing.T) {
	h := newHarness(t, Config{Port: 19443})
	h.fake.state = wsl.StateRunning
	h.fake.answer = runningPortainer("newer")
	updated, err := h.w.UpdatePortainer(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	runs := recreations(h)
	if !updated || len(runs) != 1 {
		t.Fatalf("updated %v, recreations %d", updated, len(runs))
	}
	run := runs[0]
	if !strings.Contains(run, "docker rm -f portainer") || !strings.Contains(run, "-v portainer_data:/data") ||
		strings.Contains(run, "env-file") || !strings.Contains(run, "-p 127.0.0.1:19443:9443") {
		t.Errorf("recreation: %s", run)
	}
	if h.srv.tokenSeen != "" {
		t.Errorf("an update must not open the first-run window, token %q", h.srv.tokenSeen)
	}
}

func TestUpdatePortainerNeedsAContainer(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	if _, err := h.w.UpdatePortainer(t.Context()); !errors.Is(err, ErrNoPortainer) {
		t.Fatalf("got %v", err)
	}
	if contains(h.fake.scripts, "docker pull "+portainerImage) {
		t.Fatal("pulled an image for a container that does not exist")
	}
}

func TestResetAdminSendsOnlyAHash(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	h.srv.adminExists = true
	h.srv.adminBody = map[string]string{"Username": admin.Username, "Password": admin.Password}
	h.fake.answer = runningPortainer("current")
	if err := h.w.ResetAdmin(t.Context(), admin); err != nil {
		t.Fatal(err)
	}
	if !h.srv.loginOK {
		t.Fatal("no sign-in with the new password")
	}
	var reset string
	for _, s := range h.fake.scripts {
		if strings.Contains(s, resetImage) {
			reset = s
		}
		if strings.Contains(s, admin.Password) {
			t.Fatalf("the password reached the distro: %s", s)
		}
	}
	if reset == "" {
		t.Fatal("the reset helper did not run")
	}
	_, rest, _ := strings.Cut(reset, "--password-hash '")
	hash, _, _ := strings.Cut(rest, "'")
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(admin.Password)); err != nil {
		t.Fatalf("the hash does not match the password: %v", err)
	}
	// The container is stopped before the helper and started again even if
	// the helper fails.
	stop := strings.Index(reset, "docker stop portainer")
	run := strings.Index(reset, "docker run --rm")
	if stop < 0 || run < stop || !strings.Contains(reset, "trap 'docker start portainer") {
		t.Errorf("reset order: %s", reset)
	}
	for _, e := range h.events {
		if strings.Contains(e.Message, hash) || strings.Contains(e.Message, admin.Password) {
			t.Fatalf("secret in event %+v", e)
		}
	}
}

func TestResetAdminReportsAWrongUsername(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	h.srv.adminExists = true
	h.srv.adminBody = map[string]string{"Username": "someone-else", "Password": admin.Password}
	h.fake.answer = runningPortainer("current")
	err := h.w.ResetAdmin(t.Context(), admin)
	if err == nil || !strings.Contains(err.Error(), "check the username") {
		t.Fatalf("got %v", err)
	}
}

func TestResetAdminRejectsBadInput(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	h.fake.answer = runningPortainer("current")
	if err := h.w.ResetAdmin(t.Context(), Credentials{Username: "a", Password: strings.Repeat("s", 5)}); !errors.Is(err, ErrCredentials) {
		t.Fatalf("short password: %v", err)
	}
	if err := h.w.ResetAdmin(t.Context(), Credentials{Username: "a", Password: strings.Repeat("q", 73)}); !errors.Is(err, ErrPasswordTooLong) {
		t.Fatalf("long password: %v", err)
	}
	if len(h.fake.scripts) != 0 {
		t.Fatalf("a rejected password reached the distro: %d scripts", len(h.fake.scripts))
	}
}

func TestResetAdminNeedsAContainer(t *testing.T) {
	h := newHarness(t, Config{})
	h.fake.state = wsl.StateRunning
	if err := h.w.ResetAdmin(t.Context(), admin); !errors.Is(err, ErrNoPortainer) {
		t.Fatalf("got %v", err)
	}
}

func TestRedactHidesTheSecret(t *testing.T) {
	err := redact(errors.New("helper said $2a$10$abc failed"), "$2a$10$abc")
	if strings.Contains(err.Error(), "abc") {
		t.Fatalf("got %v", err)
	}
	if redact(nil, "x") != nil {
		t.Fatal("nil error changed")
	}
}

func contains(scripts []string, want string) bool {
	return slices.ContainsFunc(scripts, func(s string) bool { return strings.Contains(s, want) })
}
