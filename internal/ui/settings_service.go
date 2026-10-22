package ui

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/the-khiem7/docker-zero/internal/control"
	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/target"
)

// SettingsView is what the settings window renders. Auto start, shortcuts,
// the .wslconfig line and the certificate are read from the machine every
// time, so the window never shows a stale stored value.
type SettingsView struct {
	AutoStart      settings.AutoStart   `json:"autoStart"`
	Shortcuts      settings.Shortcuts   `json:"shortcuts"`
	IdleMode       settings.IdleMode    `json:"idleMode"`
	WSLConfigIdle  bool                 `json:"wslConfigIdle"`
	Browser        settings.BrowserKind `json:"browser"`
	BrowserPath    string               `json:"browserPath"`
	ChromePath     string               `json:"chromePath"`
	Cert           settings.Cert        `json:"cert"`
	InstallPath    string               `json:"installPath"`
	Problems       []string             `json:"problems"`
	WSLConfigPath  string               `json:"wslConfigPath"`
	NeedsWSLReboot bool                 `json:"needsWslReboot"`
}

// SettingsService is bound to the settings window.
type SettingsService struct {
	prefs      *Prefs
	ctrl       *control.Controller
	wsl        settings.Scripter
	wslcfg     settings.WSLConfig
	installExe string
	crtPath    string
	app        *application.App
	onChange   func()

	// pendingReboot is set when .wslconfig changed and WSL has not been
	// shut down since, so the change is not active yet.
	pendingReboot bool
}

// NewSettingsService builds the service. installExe is where the app copies
// itself before auto start or shortcuts point at it.
func NewSettingsService(app *application.App, prefs *Prefs, ctrl *control.Controller, w settings.Scripter, wslcfg settings.WSLConfig, installExe string) *SettingsService {
	return &SettingsService{
		app:        app,
		prefs:      prefs,
		ctrl:       ctrl,
		wsl:        w,
		wslcfg:     wslcfg,
		installExe: installExe,
		crtPath:    filepath.Join(filepath.Dir(installExe), "portainer.crt"),
	}
}

// Get reads every setting.
func (s *SettingsService) Get(ctx context.Context) SettingsView {
	cur := s.prefs.Get()
	v := SettingsView{
		IdleMode:       cur.IdleMode,
		Browser:        cur.Browser,
		BrowserPath:    cur.BrowserPath,
		InstallPath:    s.installExe,
		WSLConfigPath:  s.wslcfg.Path,
		NeedsWSLReboot: s.pendingReboot,
	}
	var err error
	problem := func(what string, err error) {
		if err != nil {
			v.Problems = append(v.Problems, what+": "+err.Error())
		}
	}
	v.AutoStart, err = settings.ReadAutoStart()
	problem("auto start", err)
	v.Shortcuts, err = settings.ReadShortcuts()
	problem("shortcuts", err)
	v.WSLConfigIdle, err = s.wslcfg.IdleDisabled()
	problem(".wslconfig", err)
	v.ChromePath, _ = settings.FindChrome(os.Getenv, os.Stat)
	v.Cert, err = s.readCert(ctx)
	if !errors.Is(err, settings.ErrNoCert) {
		problem("certificate", err)
	}
	return v
}

func (s *SettingsService) readCert(ctx context.Context) (settings.Cert, error) {
	pemBytes, err := settings.ReadCertPEM(ctx, s.wsl, target.Distro)
	if err != nil {
		return settings.Cert{}, err
	}
	thumb, notAfter, err := settings.ParseCert(pemBytes)
	if err != nil {
		return settings.Cert{}, err
	}
	trusted, err := settings.IsTrusted(thumb)
	return settings.Cert{Present: true, Thumbprint: thumb, NotAfter: notAfter, Trusted: trusted}, err
}

// SetAutoStart writes or removes the HKCU Run value. The app copies itself
// to its install folder first so the value survives a cleaned Downloads
// folder.
func (s *SettingsService) SetAutoStart(enabled, startDistro bool) error {
	exe := s.installExe
	if enabled {
		var err error
		if exe, err = settings.EnsureInstalled(s.installExe); err != nil {
			return err
		}
	}
	return settings.WriteAutoStart(settings.AutoStart{Enabled: enabled, StartDistro: startDistro}, exe)
}

// SetShortcuts creates or removes the Desktop and Start Menu shortcuts.
func (s *SettingsService) SetShortcuts(desktop, startMenu bool) error {
	exe := s.installExe
	if desktop || startMenu {
		var err error
		if exe, err = settings.EnsureInstalled(s.installExe); err != nil {
			return err
		}
	}
	return settings.WriteShortcuts(settings.Shortcuts{Desktop: desktop, StartMenu: startMenu}, exe)
}

// SetIdleMode stores the idle mode, edits .wslconfig to match and applies
// the hold. It returns true when .wslconfig changed: that change takes
// effect only after WSL shuts down, which the user confirms separately.
func (s *SettingsService) SetIdleMode(ctx context.Context, mode settings.IdleMode) (bool, error) {
	if err := s.prefs.Update(func(p *settings.Settings) { p.IdleMode = mode }); err != nil {
		return false, err
	}
	changed, err := s.wslcfg.SetIdleDisabled(mode == settings.IdleWSLConfig)
	if err != nil {
		return false, err
	}
	if changed {
		s.pendingReboot = !s.pendingReboot
	}
	return s.pendingReboot, errors.Join(s.ctrl.SyncHold(ctx), s.notify())
}

// ApplyWSLConfig shuts WSL down so a .wslconfig change takes effect. The
// window asks for confirmation because every distro of the account stops.
func (s *SettingsService) ApplyWSLConfig(ctx context.Context) error {
	if err := s.ctrl.ShutdownWSL(ctx); err != nil {
		return err
	}
	s.pendingReboot = false
	return s.notify()
}

// SetBrowser stores the browser choice.
func (s *SettingsService) SetBrowser(kind settings.BrowserKind, path string) error {
	return s.prefs.Update(func(p *settings.Settings) {
		p.Browser = kind
		if kind == settings.BrowserCustom {
			p.BrowserPath = path
		}
	})
}

// PickBrowser asks for a browser executable and returns its path, or "".
func (s *SettingsService) PickBrowser() (string, error) {
	return s.app.Dialog.OpenFile().
		SetTitle("Choose a browser").
		AddFilter("Programs", "*.exe").
		CanChooseFiles(true).
		PromptForSingleSelection()
}

// TrustCert adds the distro's certificate to the current user's Root store
// behind Windows' own confirmation dialog.
func (s *SettingsService) TrustCert(ctx context.Context) error {
	pemBytes, err := settings.ReadCertPEM(ctx, s.wsl, target.Distro)
	if err != nil {
		return err
	}
	return settings.Trust(pemBytes, s.crtPath)
}

// UntrustCert removes the distro's certificate from the Root store.
func (s *SettingsService) UntrustCert(ctx context.Context) error {
	c, err := s.readCert(ctx)
	if err != nil {
		return err
	}
	if !c.Trusted {
		return nil
	}
	return settings.Untrust(c.Thumbprint)
}

// RegenerateCert removes trust in the old certificate, creates a new one in
// the distro, restarts Portainer, and trusts the new one if the old one was
// trusted. Each store change shows Windows' confirmation dialog.
func (s *SettingsService) RegenerateCert(ctx context.Context) error {
	old, err := s.readCert(ctx)
	if err != nil && !errors.Is(err, settings.ErrNoCert) {
		return err
	}
	if old.Trusted {
		if err := settings.Untrust(old.Thumbprint); err != nil {
			return err
		}
	}
	if err := settings.Regenerate(ctx, s.wsl, target.Distro); err != nil {
		return err
	}
	if old.Trusted {
		return s.TrustCert(ctx)
	}
	return nil
}

func (s *SettingsService) notify() error {
	if s.onChange != nil {
		s.onChange()
	}
	return nil
}
