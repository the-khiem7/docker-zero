package ui

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

	"github.com/the-khiem7/docker-zero/internal/control"
	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/setup"
	"github.com/the-khiem7/docker-zero/internal/trayicon"
)

// Window names; the frontend picks its layout from the ?view= query.
const (
	PopupWindow    = "popup"
	CornerWindow   = "corner"
	SettingsWindow = "settings"
	SetupWindow    = "setup"
)

const (
	popupWidth, popupHeight       = 360, 420
	cornerWidth, cornerHeight     = 280, 96
	cornerMargin                  = 12
	settingsWidth, settingsHeight = 520, 680
	setupWidth, setupHeight       = 520, 700
)

// Options are what Install needs.
type Options struct {
	Controller *control.Controller
	Prober     control.Prober
	Prefs      *Prefs
	Scripter   settings.Scripter
	WSLConfig  settings.WSLConfig
	InstallExe string
	// NewWizard makes the Day 1 wizard for one setup run.
	NewWizard func() *setup.Wizard
	// StartDistro starts the runtime quietly after launch, for auto start
	// with the start-distro option.
	StartDistro bool
}

// Shell holds the tray icon and the widget windows.
type Shell struct {
	app        *application.App
	svc        *Service
	prefs      *Prefs
	tray       *application.SystemTray
	popup      *application.WebviewWindow
	corner     *application.WebviewWindow
	settings   *application.WebviewWindow
	setup      *application.WebviewWindow
	cornerItem *application.MenuItem
}

// Install registers the services and creates the tray icon, the popup, the
// corner box and the settings window. Call it before app.Run.
func Install(app *application.App, o Options) *Shell {
	svc := newService(o.Controller, o.Prober)
	svc.startDistro = o.StartDistro
	sh := &Shell{app: app, svc: svc, prefs: o.Prefs}
	svc.shell = sh
	setSvc := NewSettingsService(app, o.Prefs, o.Controller, o.Scripter, o.WSLConfig, o.InstallExe)
	setSvc.onChange = svc.Refresh
	setupSvc := NewSetupService(o.NewWizard)
	setupSvc.onDone = svc.Refresh
	setupSvc.attach(app)
	app.RegisterService(application.NewService(svc))
	app.RegisterService(application.NewService(setSvc))
	app.RegisterService(application.NewService(setupSvc))

	sh.popup = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:            PopupWindow,
		Title:           "Docker Zero",
		Width:           popupWidth,
		Height:          popupHeight,
		Frameless:       true,
		AlwaysOnTop:     true,
		Hidden:          true,
		DisableResize:   true,
		HideOnEscape:    true,
		HideOnFocusLost: true,
		Windows:         application.WindowsWindow{HiddenOnTaskbar: true},
		URL:             "/?view=" + PopupWindow,
	})
	sh.corner = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:          CornerWindow,
		Title:         "Docker Zero status",
		Width:         cornerWidth,
		Height:        cornerHeight,
		Frameless:     true,
		AlwaysOnTop:   true,
		Hidden:        true,
		DisableResize: true,
		Windows:       application.WindowsWindow{HiddenOnTaskbar: true},
		URL:           "/?view=" + CornerWindow,
	})
	sh.settings = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:      SettingsWindow,
		Title:     "Docker Zero settings",
		Width:     settingsWidth,
		Height:    settingsHeight,
		MinWidth:  420,
		MinHeight: 480,
		Hidden:    true,
		URL:       "/?view=" + SettingsWindow,
	})
	sh.setup = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:      SetupWindow,
		Title:     "Set up Docker Zero",
		Width:     setupWidth,
		Height:    setupHeight,
		MinWidth:  420,
		MinHeight: 520,
		Hidden:    true,
		URL:       "/?view=" + SetupWindow,
	})
	for _, w := range []*application.WebviewWindow{sh.popup, sh.corner, sh.settings, sh.setup} {
		name := w.Name()
		// Closing a window only hides it; the app lives in the tray.
		w.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
			w.Hide()
			e.Cancel()
		})
		if name == SettingsWindow || name == SetupWindow {
			continue
		}
		w.OnWindowEvent(events.Common.WindowShow, func(*application.WindowEvent) { svc.setVisible(name, true) })
		w.OnWindowEvent(events.Common.WindowHide, func(*application.WindowEvent) { svc.setVisible(name, false) })
	}

	menu := app.NewMenu()
	menu.Add("Open Portainer").OnClick(func(*application.Context) { _ = svc.OpenPortainer() })
	sh.cornerItem = menu.AddCheckbox("Show corner box", false).OnClick(func(ctx *application.Context) {
		sh.setCorner(ctx.ClickedMenuItem().Checked())
	})
	menu.Add("Set up Docker Zero...").OnClick(func(*application.Context) { sh.showSetup() })
	menu.Add("Settings...").OnClick(func(*application.Context) { sh.showSettings() })
	menu.AddSeparator()
	menu.Add("Quit Docker Zero").OnClick(func(*application.Context) { app.Quit() })

	sh.tray = app.SystemTray.New()
	sh.tray.SetIcon(trayicon.PNG(trayicon.Stopped))
	sh.tray.SetTooltip("Docker Zero")
	sh.tray.SetMenu(menu)
	sh.tray.AttachWindow(sh.popup).WindowOffset(8)
	svc.attach(app, sh.tray)

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		if o.Prefs.Get().CornerBox {
			sh.setCorner(true)
		}
	})
	return sh
}

// ShowPopup opens the tray popup at its tray-anchored position.
func (sh *Shell) ShowPopup() {
	application.InvokeAsync(sh.tray.ShowWindow)
}

func (sh *Shell) showSetup() {
	sh.setup.Show()
	sh.setup.Focus()
}

func (sh *Shell) showSettings() {
	sh.settings.Show()
	sh.settings.Focus()
}

// setCorner shows or hides the corner box in the bottom-right corner of the
// primary screen's work area, and remembers the choice.
func (sh *Shell) setCorner(show bool) {
	sh.cornerItem.SetChecked(show)
	_ = sh.prefs.Update(func(p *settings.Settings) { p.CornerBox = show })
	if !show {
		sh.corner.Hide()
		return
	}
	if scr := sh.app.Screen.GetPrimary(); scr != nil {
		wa := scr.WorkArea
		sh.corner.SetPosition(wa.X+wa.Width-cornerWidth-cornerMargin, wa.Y+wa.Height-cornerHeight-cornerMargin)
	}
	sh.corner.Show()
}
