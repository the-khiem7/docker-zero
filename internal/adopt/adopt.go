// Package adopt finds a Docker Engine and Portainer setup made by hand and
// takes it over: the distro, container and data stay exactly as they are, and
// the legacy launchers are replaced by Docker Zero's own.
package adopt

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// Item is one legacy part that can be replaced.
type Item string

// The replaceable items.
const (
	// ItemShortcuts is the Portainer.lnk launcher on the Desktop and in the Start Menu.
	ItemShortcuts Item = "shortcuts"
	// ItemStartup is the Startup folder shortcut that starts the distro at sign-in.
	ItemStartup Item = "startup"
	// ItemCommand is the manual docker command folder and its PATH entry.
	ItemCommand Item = "command"
	// ItemIdle is the .wslconfig idle line, which Docker Zero then manages.
	ItemIdle Item = "idle"
)

// Finding is a legacy part found on the machine.
type Finding struct {
	Item   Item   `json:"item"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Report says what exists. Container and Certificate are read only while the
// distro runs, so looking never starts it.
type Report struct {
	// Distro is the state of the Docker Zero distro: not-installed, stopped or running.
	Distro string `json:"distro"`
	// Container is the Portainer container state, "missing", or "" when unknown.
	Container   string    `json:"container"`
	Certificate bool      `json:"certificate"`
	Findings    []Finding `json:"findings"`
}

// Found reports whether there is a manual setup to take over.
func (r Report) Found() bool {
	if r.Distro == string(wsl.StateNotInstalled) {
		return false
	}
	return len(r.Findings) > 0 || (r.Container != "" && r.Container != "missing")
}

// Machine is the Windows side of the legacy setup.
type Machine interface {
	// Find reads the legacy parts without changing anything.
	Find() ([]Finding, error)
	// Replace puts Docker Zero's equivalent in place, then removes the legacy
	// part. The replacement is made first, so a failure leaves the old one.
	Replace(item Item) error
}

// WSL is the part of *wsl.Client the adopter uses.
type WSL interface {
	DistroState(ctx context.Context, distro string) (wsl.State, error)
	Script(ctx context.Context, distro, user, script string) (string, error)
}

// Adopter detects and takes over a manual setup.
type Adopter struct {
	WSL     WSL
	Machine Machine
	Distro  string
}

// inspectScript prints the container state and whether the certificate is there.
const inspectScript = `docker inspect -f '{{.State.Status}}' ` + target.PortainerContainer + ` 2>/dev/null || echo missing
if [ -s ` + settings.CertDir + `/portainer.crt ]; then echo certificate; else echo no-certificate; fi
`

// Detect looks for a manual setup. It changes nothing and does not start the distro.
func (a *Adopter) Detect(ctx context.Context) (Report, error) {
	state, err := a.WSL.DistroState(ctx, a.Distro)
	if err != nil {
		return Report{}, err
	}
	r := Report{Distro: string(state), Findings: []Finding{}}
	if state == wsl.StateNotInstalled {
		return r, nil
	}
	if r.Findings, err = a.Machine.Find(); err != nil {
		return r, err
	}
	if r.Findings == nil {
		r.Findings = []Finding{}
	}
	if state != wsl.StateRunning {
		return r, nil
	}
	out, err := a.WSL.Script(ctx, a.Distro, "root", inspectScript)
	if err != nil {
		return r, err
	}
	lines := strings.Fields(out)
	if len(lines) >= 2 {
		r.Container = lines[0]
		r.Certificate = lines[1] == "certificate"
	}
	return r, nil
}

// Adopt replaces the chosen legacy parts that are still there and returns the
// ones it replaced. A failed item does not stop the others; the errors are
// joined. The distro and the container are never touched.
func (a *Adopter) Adopt(ctx context.Context, items []Item) ([]Item, error) {
	state, err := a.WSL.DistroState(ctx, a.Distro)
	if err != nil {
		return nil, err
	}
	if state == wsl.StateNotInstalled {
		return nil, errors.New("there is no " + a.Distro + " distro to take over")
	}
	found, err := a.Machine.Find()
	if err != nil {
		return nil, err
	}
	var done []Item
	var errs []error
	for _, f := range found {
		if !slices.Contains(items, f.Item) {
			continue
		}
		if err := a.Machine.Replace(f.Item); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", f.Title, err))
			continue
		}
		done = append(done, f.Item)
	}
	return done, errors.Join(errs...)
}

// legacyPortainer reports whether a Portainer.lnk is the one the manual
// guide's launcher script created.
func legacyPortainer(args string) bool {
	return strings.Contains(strings.ToLower(args), "open-portainer.ps1")
}

// legacyKeepalive reports whether a Startup shortcut is the manual guide's
// hidden "wsl.exe --exec true" that starts the distro at sign-in.
func legacyKeepalive(args string) bool {
	a := strings.ToLower(strings.TrimSpace(args))
	return strings.Contains(a, "wsl.exe") && strings.HasSuffix(a, "--exec true")
}

// removeKnown deletes the named files below dir, then each folder on the way
// up to dir if it is empty. A file that is not named stays, so a folder the
// user also uses is never emptied.
func removeKnown(dir string, files ...string) error {
	var errs []error
	folders := map[string]bool{dir: true}
	for _, f := range files {
		p := filepath.Join(dir, f)
		if err := os.Remove(p); err != nil && !errors.Is(err, fs.ErrNotExist) {
			errs = append(errs, err)
		}
		for d := filepath.Dir(p); d != dir && len(d) > len(dir); d = filepath.Dir(d) {
			folders[d] = true
		}
	}
	// Deepest first, so a child goes before its parent.
	for len(folders) > 0 {
		deepest := ""
		for d := range folders {
			if len(d) > len(deepest) {
				deepest = d
			}
		}
		delete(folders, deepest)
		_ = os.Remove(deepest)
	}
	return errors.Join(errs...)
}
