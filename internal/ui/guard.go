package ui

import (
	"fmt"
	"sync"

	"github.com/the-khiem7/docker-zero/internal/setup"
)

// opGuard lets one long action run at a time: setup, an update, an admin
// reset and the uninstall all act on the same distro.
type opGuard struct {
	mu   sync.Mutex
	name string
}

// acquire marks name as running, or fails naming the action in progress.
func (g *opGuard) acquire(name string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.name != "" {
		return fmt.Errorf("%w: %s", setup.ErrBusy, g.name)
	}
	g.name = name
	return nil
}

func (g *opGuard) release() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.name = ""
}
