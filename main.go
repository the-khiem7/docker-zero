// Docker Zero is a tray app that bootstraps and operates Docker Engine in a
// WSL2 distro with Portainer CE as the GUI, without admin rights.
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Docker Zero",
		Description: "Docker Engine in WSL2 with Portainer, without admin rights",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		// A second launch (shortcut, auto start) must not add a second tray icon.
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "io.github.dockerzero",
		},
		Windows: application.WindowsOptions{
			// The app lives in the tray; hiding the popup must not quit it.
			DisableQuitOnLastWindowClosed: true,
		},
	})

	popup := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:            "popup",
		Title:           "Docker Zero",
		Width:           360,
		Height:          420,
		Frameless:       true,
		AlwaysOnTop:     true,
		Hidden:          true,
		DisableResize:   true,
		HideOnEscape:    true,
		HideOnFocusLost: true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		URL: "/",
	})
	popup.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		popup.Hide()
		e.Cancel()
	})

	menu := app.NewMenu()
	menu.Add("Quit Docker Zero").OnClick(func(*application.Context) { app.Quit() })

	tray := app.SystemTray.New()
	tray.SetTooltip("Docker Zero")
	tray.SetMenu(menu)
	tray.AttachWindow(popup).WindowOffset(8)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
