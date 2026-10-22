//go:build integration && windows

package adopt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// TestDetectOnThisMachine reads the real machine and logs what a take-over
// would replace. It changes nothing and does not start the distro.
func TestDetectOnThisMachine(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	a := &Adopter{
		WSL:    wsl.New(),
		Distro: target.Distro,
		Machine: Windows{
			LegacyDir: filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "docker-wsl"),
			WSLConfig: settings.WSLConfig{Path: filepath.Join(home, ".wslconfig")},
		},
	}
	r, err := a.Detect(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("distro %s, container %q, certificate %v, found %v", r.Distro, r.Container, r.Certificate, r.Found())
	for _, f := range r.Findings {
		t.Logf("  %s: %s (%s)", f.Item, f.Title, f.Detail)
	}
}
