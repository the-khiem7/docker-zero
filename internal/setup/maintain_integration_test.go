//go:build integration

package setup

import (
	"bytes"
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

// TestMaintenance sets up a new distro named in DZ_SETUP_DISTRO, then updates
// Portainer, resets the admin password and uninstalls, checking the effect of
// each on the distro itself. The Windows-side hooks are recorders, so the
// run never touches the real account's auto start, PATH or certificates, and it
// refuses the fixed production name.
func TestMaintenance(t *testing.T) {
	distro := os.Getenv("DZ_SETUP_DISTRO")
	if distro == "" || strings.EqualFold(distro, target.Distro) {
		t.Skip("set DZ_SETUP_DISTRO to a new throwaway distro name")
	}
	client := wsl.New()
	if state, err := client.DistroState(t.Context(), distro); err != nil || state != wsl.StateNotInstalled {
		t.Fatalf("distro %q must not exist before the run: state %q, err %v", distro, state, err)
	}

	loc := t.TempDir()
	t.Cleanup(func() {
		// The uninstall below normally removed it; this covers an early failure.
		if state, _ := client.DistroState(t.Context(), distro); state != wsl.StateNotInstalled {
			_ = client.Unregister(t.Context(), distro)
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
	var hostCalls []string
	record := func(name string) func() error {
		return func() error { hostCalls = append(hostCalls, name); return nil }
	}
	w.Host = Host{
		ReleaseHold:     func() { hostCalls = append(hostCalls, "release") },
		RemoveAutoStart: record("autostart"),
		RemoveShortcuts: record("shortcuts"),
		RemoveIdleLine:  record("idle"),
		RemoveConfig:    record("config"),
		RemoveFromPath:  func(dir string) error { hostCalls = append(hostCalls, "path"); return nil },
		UntrustCert: func(pem []byte) (bool, error) {
			hostCalls = append(hostCalls, "untrust")
			return strings.Contains(string(pem), "BEGIN CERTIFICATE"), nil
		},
	}

	random := func() string {
		b := make([]byte, 12)
		if _, err := rand.Read(b); err != nil {
			t.Fatal(err)
		}
		return hex.EncodeToString(b)
	}
	cr := Credentials{Username: "maint-test", Password: random()}
	if err := w.Run(t.Context(), cr); err != nil {
		t.Fatalf("setup: %v", err)
	}

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
	script := func(s string) string {
		out, err := client.Script(t.Context(), distro, "root", s)
		if err != nil {
			t.Fatalf("script %q: %v", s, err)
		}
		return strings.TrimSpace(out)
	}

	// Update with nothing new leaves the container alone.
	idBefore := script("docker inspect -f '{{.Id}}' portainer\n")
	updated, err := w.UpdatePortainer(t.Context())
	if err != nil || updated {
		t.Fatalf("update on the current image: updated %v, err %v", updated, err)
	}
	if got := script("docker inspect -f '{{.Id}}' portainer\n"); got != idBefore {
		t.Fatal("a current container was recreated")
	}

	// Put Portainer on an image that is not the lts one, then update.
	script(`set -e
printf 'FROM portainer/portainer-ce:lts\nLABEL dz.test=old\n' | docker build -q -t dz-old-portainer - >/dev/null
docker rm -f portainer >/dev/null
docker run -d --name portainer --restart unless-stopped -p 127.0.0.1:19443:9443 -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data -v ` + "/etc/portainer-certs" + `:/certs:ro dz-old-portainer --sslcert /certs/portainer.crt --sslkey /certs/portainer.key >/dev/null
`)
	if err := w.waitReady(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got := script("docker inspect -f '{{.Config.Image}}' portainer\n"); got != "dz-old-portainer" {
		t.Fatalf("test setup: container runs %q", got)
	}
	updated, err = w.UpdatePortainer(t.Context())
	if err != nil || !updated {
		t.Fatalf("update from an older image: updated %v, err %v", updated, err)
	}
	if got := script("docker inspect -f '{{.Config.Image}}' portainer\n"); got != portainerImage {
		t.Fatalf("after the update the container runs %q", got)
	}
	if got := signIn(cr.Username, cr.Password); got != http.StatusOK {
		t.Fatalf("sign-in after the update: HTTP %d; the data volume was not kept", got)
	}

	// Reset the admin password.
	next := Credentials{Username: cr.Username, Password: random()}
	if err := w.ResetAdmin(t.Context(), next); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if got := signIn(next.Username, next.Password); got != http.StatusOK {
		t.Fatalf("sign-in with the new password: HTTP %d", got)
	}
	if got := signIn(cr.Username, cr.Password); got == http.StatusOK {
		t.Fatal("the old password still works")
	}
	if got := script("docker inspect -f '{{.State.Status}}' portainer\n"); got != "running" {
		t.Fatalf("Portainer is %q after the reset", got)
	}

	// Uninstall with a backup.
	backup := filepath.Join(loc, "backup.tar")
	if err := w.Uninstall(t.Context(), UninstallOptions{BackupPath: backup}); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if state, err := client.DistroState(t.Context(), distro); err != nil || state != wsl.StateNotInstalled {
		t.Fatalf("distro after uninstall: state %q, err %v", state, err)
	}
	if fi, err := os.Stat(backup); err != nil || fi.Size() < 1<<20 {
		t.Fatalf("backup: %v, %v", fi, err)
	}
	if _, err := os.Stat(filepath.Join(loc, "bin", "docker.cmd")); !os.IsNotExist(err) {
		t.Fatalf("docker.cmd after the uninstall: %v", err)
	}
	t.Logf("host calls: %v", hostCalls)
	for _, want := range []string{"untrust", "release", "autostart", "shortcuts", "idle", "path", "config"} {
		var found bool
		for _, c := range hostCalls {
			found = found || c == want
		}
		if !found {
			t.Errorf("host hook %q never ran", want)
		}
	}
}
