package ui

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

	"github.com/the-khiem7/docker-zero/internal/control"
	"github.com/the-khiem7/docker-zero/internal/trayicon"
)

// Window names; the frontend picks its layout from the ?view= query.
const (
	PopupWindow  = "popup"
	CornerWindow = "corner"
)

const (
	popupWidth, popupHeight   = 360, 400
	cornerWidth, cornerHeight = 280, 96
	cornerMargin              = 12
)

// Shell holds the tray icon and the two widget windows.
type Shell struct {
	app        *application.App
	svc        *Service
	tray       *application.SystemTray
	popup      *application.WebviewWindow
	corner     *application.WebviewWindow
	cornerItem *application.MenuItem
}

// Install registers the service and creates the tray icon, the popup and
// the corner box. Call it before app.Run.
func Install(app *application.App, ctrl *control.Controller, prober control.Prober) *Shell {
	svc := newService(ctrl, prober)
	sh := &Shell{app: app, svc: svc}
	svc.shell = sh
	app.RegisterService(application.NewService(svc))

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
	for _, w := range []*application.WebviewWindow{sh.popup, sh.corner} {
		name := w.Name()
		// Closing a widget only hides it; the app lives in the tray.
		w.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
			w.Hide()
			e.Cancel()
		})
		w.OnWindowEvent(events.Common.WindowShow, func(*application.WindowEvent) { svc.setVisible(name, true) })
		w.OnWindowEvent(events.Common.WindowHide, func(*application.WindowEvent) { svc.setVisible(name, false) })
	}

	menu := app.NewMenu()
	menu.Add("Open Portainer").OnClick(func(*application.Context) { _ = svc.OpenPortainer() })
	sh.cornerItem = menu.AddCheckbox("Show corner box", false).OnClick(func(ctx *application.Context) {
		sh.setCorner(ctx.ClickedMenuItem().Checked())
	})
	menu.AddSeparator()
	menu.Add("Quit Docker Zero").OnClick(func(*application.Context) { app.Quit() })

	sh.tray = app.SystemTray.New()
	sh.tray.SetIcon(trayicon.PNG(trayicon.Stopped))
	sh.tray.SetTooltip("Docker Zero")
	sh.tray.SetMenu(menu)
	sh.tray.AttachWindow(sh.popup).WindowOffset(8)
	svc.attach(app, sh.tray)
	return sh
}

// ShowPopup opens the tray popup at its tray-anchored position.
func (sh *Shell) ShowPopup() {
	application.InvokeAsync(sh.tray.ShowWindow)
}

// setCorner shows or hides the corner box in the bottom-right corner of the
// primary screen's work area.
func (sh *Shell) setCorner(show bool) {
	sh.cornerItem.SetChecked(show)
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
