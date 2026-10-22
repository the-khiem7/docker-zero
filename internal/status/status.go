// Package status reports the live state of WSL, the distro, Docker Engine
// and Portainer without starting anything.
package status

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// Snapshot is one status reading.
type Snapshot struct {
	WSLAvailable     bool      `json:"wslAvailable"`
	Distro           wsl.State `json:"distro"`
	Docker           string    `json:"docker"`
	Portainer        string    `json:"portainer"`
	PortainerUp      bool      `json:"portainerUp"`
	PortainerVersion string    `json:"portainerVersion,omitempty"`
	Error            string    `json:"error,omitempty"`
	CheckedAt        time.Time `json:"checkedAt"`
}

// Values of Snapshot.Docker and Snapshot.Portainer besides the ones systemd
// and Docker print themselves ("active", "running", "exited", ...).
const (
	Unknown = "unknown"
	Missing = "missing"
)

// WSL is the part of *wsl.Client the prober uses.
type WSL interface {
	Available(ctx context.Context) error
	DistroState(ctx context.Context, distro string) (wsl.State, error)
	Script(ctx context.Context, distro, user, script string) (string, error)
}

// Prober takes status snapshots.
type Prober struct {
	WSL       WSL
	HTTP      *http.Client
	StatusURL string
	Now       func() time.Time
}

// NewProber returns a Prober for the real runtime.
func NewProber(c WSL) *Prober {
	return &Prober{
		WSL:       c,
		HTTP:      loopbackClient(),
		StatusURL: target.PortainerStatusURL,
		Now:       time.Now,
	}
}

// loopbackClient talks only to Portainer on 127.0.0.1. Certificate checks
// are skipped because the probe may run before the locally trusted
// certificate exists; it sends no credentials and reads only the public
// version endpoint.
func loopbackClient() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // loopback-only health probe, see above
			Proxy:           nil,
		},
	}
}

// probeScript prints the Docker service state and the Portainer container
// state on two lines. It runs as root so it does not depend on group
// membership of the default user.
const probeScript = `systemctl is-active docker 2>/dev/null || true
docker inspect -f '{{.State.Status}}' ` + target.PortainerContainer + ` 2>/dev/null || echo ` + Missing + `
`

// Probe reads the current status. It never starts the distro: the Docker
// and Portainer checks run only when the distro is already running.
func (p *Prober) Probe(ctx context.Context) Snapshot {
	s := Snapshot{
		Distro:    wsl.StateNotInstalled,
		Docker:    Unknown,
		Portainer: Unknown,
		CheckedAt: p.Now(),
	}
	if err := p.WSL.Available(ctx); err != nil {
		s.Error = err.Error()
		return s
	}
	s.WSLAvailable = true

	state, err := p.WSL.DistroState(ctx, target.Distro)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	s.Distro = state
	if state != wsl.StateRunning {
		return s
	}

	out, err := p.WSL.Script(ctx, target.Distro, "root", probeScript)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	s.Docker, s.Portainer = parseProbe(out)

	if s.Portainer == "running" {
		version, err := p.portainerVersion(ctx)
		if err != nil {
			s.Error = err.Error()
		} else {
			s.PortainerUp = true
			s.PortainerVersion = version
		}
	}
	return s
}

func parseProbe(out string) (docker, portainer string) {
	docker, portainer = Unknown, Unknown
	lines := strings.Split(strings.TrimSpace(strings.ReplaceAll(out, "\r", "")), "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
		docker = strings.TrimSpace(lines[0])
	}
	if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
		portainer = strings.TrimSpace(lines[1])
	}
	return docker, portainer
}

var errPortainerStatus = errors.New("portainer status endpoint")

func (p *Prober) portainerVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.StatusURL, http.NoBody)
	if err != nil {
		return "", err
	}
	resp, err := p.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errPortainerStatus, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: HTTP %d", errPortainerStatus, resp.StatusCode)
	}
	var body struct {
		Version string `json:"Version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("%w: %w", errPortainerStatus, err)
	}
	return body.Version, nil
}
