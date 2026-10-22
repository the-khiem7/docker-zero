//go:build integration

package setup

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// TestFullSetup runs the whole wizard against a new distro whose name is in
// DZ_SETUP_DISTRO, checks that Portainer accepts the typed credentials, runs
// the wizard a second time to prove it resumes, and unregisters the distro.
// It refuses the fixed production name so it can never touch a real setup.
func TestFullSetup(t *testing.T) {
	distro := os.Getenv("DZ_SETUP_DISTRO")
	if distro == "" || strings.EqualFold(distro, target.Distro) {
		t.Skip("set DZ_SETUP_DISTRO to a new throwaway distro name")
	}
	client := wsl.New()
	if state, err := client.DistroState(t.Context(), distro); err != nil || state != wsl.StateNotInstalled {
		t.Fatalf("distro %q must not exist before the run: state %q, err %v", distro, state, err)
	}

	// The temp folder is created first so the cleanup below runs, and the
	// distro is unregistered, before the folder is removed.
	loc := t.TempDir()
	t.Cleanup(func() {
		// The test context is already canceled when cleanups run.
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := client.Unregister(ctx, distro); err != nil {
			t.Logf("unregister %s: %v", distro, err)
		}
	})

	const port = 19443
	w := New(client, Config{
		Distro:   distro,
		Location: filepath.Join(loc, "disk"),
		Port:     port,
		BinDir:   filepath.Join(loc, "bin"),
	})
	w.Emit = func(e Event) {
		if e.Status == StatusLog {
			t.Logf("  %s", e.Message)
			return
		}
		t.Logf("[%s] %s %s", e.Step, e.Status, e.Message)
	}

	secret := make([]byte, 12)
	if _, err := rand.Read(secret); err != nil {
		t.Fatal(err)
	}
	cr := Credentials{Username: "wizard-test", Password: hex.EncodeToString(secret)}

	start := time.Now()
	if err := w.Run(t.Context(), cr); err != nil {
		t.Fatalf("first run: %v", err)
	}
	t.Logf("first run took %s", time.Since(start).Round(time.Second))

	// An independent client, not the wizard's, signs in with the typed
	// credentials. The certificate is the wizard's own leaf, so it is only
	// skipped here, not trusted.
	hc := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, //nolint:gosec // loopback test against the throwaway distro
	}
	signIn := func(user, pass string) int {
		body, _ := json.Marshal(map[string]string{"username": user, "password": pass})
		resp, err := hc.Post("https://127.0.0.1:19443/api/auth", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		return resp.StatusCode
	}
	if got := signIn(cr.Username, cr.Password); got != http.StatusOK {
		t.Fatalf("sign-in with the typed credentials: HTTP %d", got)
	}
	if got := signIn(cr.Username, cr.Password+"x"); got == http.StatusOK {
		t.Fatal("sign-in with a wrong password succeeded")
	}

	// The token must not remain in the container configuration.
	out, err := client.Script(t.Context(), distro, "root",
		"docker inspect -f '{{json .Config.Env}} {{json .Args}}' "+target.PortainerContainer+"\n")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "SETUP_TOKEN") || strings.Contains(out, "setup-token") {
		t.Fatalf("setup token left in the container config: %s", out)
	}

	// A second run resumes: nothing is installed again, the admin stays.
	var again []Event
	w.Emit = func(e Event) { again = append(again, e) }
	if err := w.Run(t.Context(), cr); err != nil {
		t.Fatalf("second run: %v", err)
	}
	for _, id := range []StepID{StepDistro, StepDocker, StepUser, StepPortainer, StepAdmin} {
		var got Status
		for _, e := range again {
			if e.Step == id && e.Status != StatusLog && e.Status != StatusRunning {
				got = e.Status
			}
		}
		if got != StatusSkipped {
			t.Errorf("second run, step %s: %q, want skipped", id, got)
		}
	}
	if got := signIn(cr.Username, cr.Password); got != http.StatusOK {
		t.Fatalf("sign-in after the second run: HTTP %d", got)
	}
	if _, err := os.Stat(filepath.Join(loc, "bin", "docker.cmd")); err != nil {
		t.Fatal(err)
	}
}
