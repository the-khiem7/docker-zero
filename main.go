// Docker Zero is a tray app that bootstraps and operates Docker Engine in a
// WSL2 distro with Portainer CE as the GUI, without admin rights.
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/the-khiem7/docker-zero/internal/control"
	"github.com/the-khiem7/docker-zero/internal/status"
	"github.com/the-khiem7/docker-zero/internal/ui"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var shell *ui.Shell
	app := application.New(application.Options{
		Name:        "Docker Zero",
		Description: "Docker Engine in WSL2 with Portainer, without admin rights",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		// A second launch (shortcut, auto start) must not add a second tray icon.
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "io.github.dockerzero",
			// Launching the app again (shortcut, Start Menu) opens the popup.
			OnSecondInstanceLaunch: func(application.SecondInstanceData) {
				if shell != nil {
					shell.ShowPopup()
				}
			},
		},
		Windows: application.WindowsOptions{
			// The app lives in the tray; hiding every window must not quit it.
			DisableQuitOnLastWindowClosed: true,
		},
	})

	client := wsl.New()
	prober := status.NewProber(client)
	ctrl := control.New(client, prober, app.Browser.OpenURL)
	shell = ui.Install(app, ctrl, prober)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
