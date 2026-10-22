//go:build integration

package status

import (
	"testing"

	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// TestProbeReal reads the real machine's status. It changes nothing and
// never starts the distro. Run with: go test -tags integration ./internal/status
func TestProbeReal(t *testing.T) {
	s := NewProber(wsl.New()).Probe(t.Context())
	t.Logf("%+v", s)
	if !s.WSLAvailable {
		t.Fatalf("WSL not available: %s", s.Error)
	}
}
