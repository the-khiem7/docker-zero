// Docker Zero is a tray app that bootstraps and operates Docker Engine in a
// WSL2 distro with Portainer CE as the GUI, without admin rights.
package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/the-khiem7/docker-zero/internal/control"
	"github.com/the-khiem7/docker-zero/internal/paths"
	"github.com/the-khiem7/docker-zero/internal/settings"
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

	configDir, err := paths.ConfigDir(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	installExe, err := paths.InstalledExe(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	prefs, err := ui.LoadPrefs(settings.Store{Path: filepath.Join(configDir, "settings.json")})
	if err != nil {
		app.Logger.Warn("settings unreadable, using defaults", "error", err)
	}

	client := wsl.New()
	prober := status.NewProber(client)
	opener := settings.Opener{Default: app.Browser.OpenURL, Current: prefs.Get}
	ctrl := control.New(client, prober, opener.Open)
	ctrl.KeepAlive = func() bool { return prefs.Get().IdleMode == settings.IdleKeepalive }

	shell = ui.Install(app, ui.Options{
		Controller:  ctrl,
		Prober:      prober,
		Prefs:       prefs,
		Scripter:    client,
		WSLConfig:   settings.WSLConfig{Path: filepath.Join(home, ".wslconfig")},
		InstallExe:  installExe,
		StartDistro: slices.Contains(os.Args[1:], settings.StartDistroFlag),
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
