// Package control implements Docker Zero's system-level Day 2 actions:
// start the runtime, stop the distro, shut down WSL. Container work stays
// in Portainer.
package control

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/the-khiem7/docker-zero/internal/status"
	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// Errors returned by actions.
var (
	ErrBusy         = errors.New("another action is still running")
	ErrNotInstalled = errors.New(target.Distro + " is not installed; run setup first")
	ErrTimeout      = errors.New("timed out waiting for Portainer")
)

// WSL is the part of *wsl.Client the controller uses.
type WSL interface {
	status.WSL
	Terminate(ctx context.Context, distro string) error
	Shutdown(ctx context.Context) error
}

// Holder keeps the distro from idling out while alive.
type Holder interface {
	Alive() bool
	Release()
}

// Prober takes status snapshots.
type Prober interface {
	Probe(ctx context.Context) status.Snapshot
}

// Controller runs one action at a time.
type Controller struct {
	WSL    WSL
	Prober Prober
	// StartHold starts a keepalive session in a distro.
	StartHold func(distro string) (Holder, error)
	// OpenURL opens Portainer in the browser.
	OpenURL func(url string) error
	// Wait is the pause between readiness checks; tests shorten it.
	Wait time.Duration
	// StartTimeout bounds how long Start waits for Portainer.
	StartTimeout time.Duration
	// KeepAlive reports whether the app should hold the distro (the
	// keepalive idle mode). Nil means always.
	KeepAlive func() bool

	mu   sync.Mutex
	busy string
	hold Holder
}

// New returns a Controller with production timings.
func New(c *wsl.Client, p Prober, openURL func(string) error) *Controller {
	return &Controller{
		WSL:          c,
		StartHold:    func(d string) (Holder, error) { return c.StartHold(d) },
		Prober:       p,
		OpenURL:      openURL,
		Wait:         time.Second,
		StartTimeout: 90 * time.Second,
	}
}

// Busy returns the name of the running action, or "".
func (c *Controller) Busy() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.busy
}

func (c *Controller) begin(action string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.busy != "" {
		return fmt.Errorf("%w: %s", ErrBusy, c.busy)
	}
	c.busy = action
	return nil
}

func (c *Controller) end() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.busy = ""
}

// ensureScript starts the Docker service and the Portainer container. The
// container has restart policy unless-stopped, so this matters only after
// someone stopped it by hand.
const ensureScript = `set -e
systemctl start docker
docker start ` + target.PortainerContainer + ` >/dev/null
`

// Start boots the distro, keeps it running with a hold in keepalive mode,
// makes sure Docker and Portainer run, waits for Portainer to answer, then
// opens it unless quiet is set.
func (c *Controller) Start(ctx context.Context, quiet bool) error {
	if err := c.begin("start"); err != nil {
		return err
	}
	defer c.end()

	state, err := c.WSL.DistroState(ctx, target.Distro)
	if err != nil {
		return err
	}
	if state == wsl.StateNotInstalled {
		return ErrNotInstalled
	}
	if c.keepAlive() {
		if err := c.ensureHold(); err != nil {
			return err
		}
	}
	if _, err := c.WSL.Script(ctx, target.Distro, "root", ensureScript); err != nil {
		return err
	}
	if err := c.waitPortainer(ctx); err != nil {
		return err
	}
	if quiet {
		return nil
	}
	return c.OpenURL(target.PortainerURL)
}

// SyncHold applies the idle mode to a distro that is already running:
// in keepalive mode it takes a hold, otherwise it releases one. It never
// starts a stopped distro.
func (c *Controller) SyncHold(ctx context.Context) error {
	if !c.keepAlive() {
		c.releaseHold()
		return nil
	}
	state, err := c.WSL.DistroState(ctx, target.Distro)
	if err != nil || state != wsl.StateRunning {
		return err
	}
	return c.ensureHold()
}

func (c *Controller) keepAlive() bool {
	return c.KeepAlive == nil || c.KeepAlive()
}

// OpenPortainer opens the Portainer page without starting anything.
func (c *Controller) OpenPortainer() error {
	return c.OpenURL(target.PortainerURL)
}

// StopDistro releases the hold and terminates the distro.
func (c *Controller) StopDistro(ctx context.Context) error {
	if err := c.begin("stop"); err != nil {
		return err
	}
	defer c.end()
	c.releaseHold()
	return c.WSL.Terminate(ctx, target.Distro)
}

// ShutdownWSL releases the hold and shuts down every distro and the VM.
// The caller must have confirmed it with the user.
func (c *Controller) ShutdownWSL(ctx context.Context) error {
	if err := c.begin("shutdown"); err != nil {
		return err
	}
	defer c.end()
	c.releaseHold()
	return c.WSL.Shutdown(ctx)
}

// Close releases the hold when the app exits, so the distro can idle out.
func (c *Controller) Close() {
	c.releaseHold()
}

// Holding reports whether the app currently keeps the distro running.
func (c *Controller) Holding() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hold != nil && c.hold.Alive()
}

func (c *Controller) ensureHold() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.hold != nil && c.hold.Alive() {
		return nil
	}
	h, err := c.StartHold(target.Distro)
	if err != nil {
		return err
	}
	c.hold = h
	return nil
}

func (c *Controller) releaseHold() {
	c.mu.Lock()
	h := c.hold
	c.hold = nil
	c.mu.Unlock()
	if h != nil {
		h.Release()
	}
}

func (c *Controller) waitPortainer(ctx context.Context) error {
	ctx, cancel := context.WithTimeoutCause(ctx, c.StartTimeout, ErrTimeout)
	defer cancel()
	for {
		if s := c.Prober.Probe(ctx); s.PortainerUp {
			return nil
		}
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case <-time.After(c.Wait):
		}
	}
}
