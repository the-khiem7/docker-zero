// Package ui connects the Day 2 controller and the status probe to the Wails
// windows, the tray icon and the frontend.
package ui

import (
	"context"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/the-khiem7/docker-zero/internal/control"
	"github.com/the-khiem7/docker-zero/internal/status"
	"github.com/the-khiem7/docker-zero/internal/trayicon"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// StatusEvent carries a View to every window.
const StatusEvent = "status"

func init() {
	application.RegisterEvent[View](StatusEvent)
}

// Poll timings: a visible widget refreshes every pollVisible; for burstFor
// after an action it refreshes every second so the change shows at once.
const (
	pollVisible = 5 * time.Second
	burstFor    = 15 * time.Second
	tick        = time.Second
)

// View is what the frontend renders.
type View struct {
	status.Snapshot
	Busy    string `json:"busy"`
	Holding bool   `json:"holding"`
	Level   string `json:"level"`
	Corner  bool   `json:"corner"`
}

// Service is bound to the frontend. Its exported methods become bindings.
type Service struct {
	ctrl   *control.Controller
	prober control.Prober
	// emit and setIcon are set by attach; nil-safe until then.
	emit    func(View)
	setIcon func([]byte)
	shell   *Shell

	mu       sync.Mutex
	last     View
	lastAt   time.Time
	visible  map[string]bool
	burstEnd time.Time
	kick     chan struct{}
	cancel   context.CancelFunc
	level    trayicon.Level
	iconSet  bool
}

func newService(ctrl *control.Controller, prober control.Prober) *Service {
	return &Service{
		ctrl:    ctrl,
		prober:  prober,
		visible: map[string]bool{},
		kick:    make(chan struct{}, 1),
	}
}

// attach wires the service to the app's events and the tray icon.
func (s *Service) attach(app *application.App, tray *application.SystemTray) {
	s.emit = func(v View) { app.Event.Emit(StatusEvent, v) }
	s.setIcon = func(b []byte) { tray.SetIcon(b) }
}

// ServiceStartup starts the poll loop.
func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	go s.loop(ctx)
	return nil
}

// ServiceShutdown stops polling and releases the keepalive session so the
// distro can idle out after the app exits.
func (s *Service) ServiceShutdown() error {
	if s.cancel != nil {
		s.cancel()
	}
	s.ctrl.Close()
	return nil
}

// setVisible records whether a widget window is on screen. A newly shown
// window gets fresh status at once.
func (s *Service) setVisible(window string, visible bool) {
	s.mu.Lock()
	s.visible[window] = visible
	s.mu.Unlock()
	if visible {
		s.Refresh()
		return
	}
	// Re-send the last reading so the other window sees the new
	// visibility, without probing for a window nobody looks at.
	s.publish(s.Status().Snapshot)
}

// Refresh asks the loop for a status reading now.
func (s *Service) Refresh() {
	select {
	case s.kick <- struct{}{}:
	default:
	}
}

// Status returns the last reading without probing.
func (s *Service) Status() View {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last
}

// Start boots the distro, ensures Portainer and opens it.
func (s *Service) Start(ctx context.Context) error {
	return s.act(func() error { return s.ctrl.Start(ctx) })
}

// StopDistro terminates the distro.
func (s *Service) StopDistro(ctx context.Context) error {
	return s.act(func() error { return s.ctrl.StopDistro(ctx) })
}

// ShutdownWSL stops every distro and the WSL VM. The frontend asks the user
// to confirm first.
func (s *Service) ShutdownWSL(ctx context.Context) error {
	return s.act(func() error { return s.ctrl.ShutdownWSL(ctx) })
}

// OpenPortainer opens the Portainer page.
func (s *Service) OpenPortainer() error {
	return s.ctrl.OpenPortainer()
}

// SetCornerBox shows or hides the always-on-top corner box.
func (s *Service) SetCornerBox(show bool) {
	if s.shell != nil {
		application.InvokeAsync(func() { s.shell.setCorner(show) })
	}
}

// HidePopup hides the tray popup, for example after opening Portainer.
func (s *Service) HidePopup() {
	if s.shell != nil {
		application.InvokeAsync(func() { s.shell.popup.Hide() })
	}
}

func (s *Service) act(fn func() error) error {
	s.burst()
	s.publish(s.Status().Snapshot) // show the busy state right away
	err := fn()
	s.burst()
	return err
}

func (s *Service) burst() {
	s.mu.Lock()
	s.burstEnd = time.Now().Add(burstFor)
	s.mu.Unlock()
	s.Refresh()
}

func (s *Service) loop(ctx context.Context) {
	s.probe(ctx)
	t := time.NewTicker(tick)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.kick:
			s.probe(ctx)
		case <-t.C:
			if s.due(time.Now()) {
				s.probe(ctx)
			}
		}
	}
}

// due reports whether a tick should probe: always during a burst, every
// pollVisible while a widget is visible, never otherwise.
func (s *Service) due(now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if now.Before(s.burstEnd) {
		return true
	}
	for _, v := range s.visible {
		if v {
			return now.Sub(s.lastAt) >= pollVisible
		}
	}
	return false
}

func (s *Service) probe(ctx context.Context) {
	s.publish(s.prober.Probe(ctx))
}

func (s *Service) publish(snap status.Snapshot) {
	level := LevelOf(snap)
	v := View{Snapshot: snap, Busy: s.ctrl.Busy(), Holding: s.ctrl.Holding(), Level: levelName(level)}

	s.mu.Lock()
	v.Corner = s.visible[CornerWindow]
	s.last, s.lastAt = v, time.Now()
	iconChanged := !s.iconSet || s.level != level
	s.level, s.iconSet = level, true
	emit, setIcon := s.emit, s.setIcon
	s.mu.Unlock()

	if emit != nil {
		emit(v)
	}
	if iconChanged && setIcon != nil {
		setIcon(trayicon.PNG(level))
	}
}

// LevelOf maps a snapshot to the tray dot.
func LevelOf(s status.Snapshot) trayicon.Level {
	switch {
	case !s.WSLAvailable || s.Distro == wsl.StateNotInstalled:
		return trayicon.Failed
	case s.PortainerUp:
		return trayicon.Running
	case s.Distro == wsl.StateRunning:
		return trayicon.Partial
	default:
		return trayicon.Stopped
	}
}

func levelName(l trayicon.Level) string {
	switch l {
	case trayicon.Running:
		return "running"
	case trayicon.Partial:
		return "partial"
	case trayicon.Failed:
		return "failed"
	default:
		return "stopped"
	}
}
