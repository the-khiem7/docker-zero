//go:build windows

package adopt

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"

	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/setup"
)

const (
	portainerLink = "Portainer.lnk"
	startupLink   = "WSL Docker keepalive.lnk"
)

// Windows is the machine the manual guide set up, read and changed per user.
type Windows struct {
	// InstallExe is where Docker Zero copies itself.
	InstallExe string
	// BinDir receives Docker Zero's docker.cmd.
	BinDir string
	// Distro is the distro the docker command forwards to.
	Distro string
	// LegacyDir is the manual guide's folder, %LOCALAPPDATA%\Programs\docker-wsl.
	LegacyDir string
	// WSLConfig is the user's .wslconfig.
	WSLConfig settings.WSLConfig
	// IdleAdopted reports whether Docker Zero already manages the idle line.
	IdleAdopted func() bool
	// AdoptIdle makes Docker Zero manage the idle line.
	AdoptIdle func() error
}

type legacyPaths struct {
	desktop, startMenu, startup string
}

func knownPaths() (legacyPaths, error) {
	d, err := windows.KnownFolderPath(windows.FOLDERID_Desktop, 0)
	if err != nil {
		return legacyPaths{}, err
	}
	m, err := windows.KnownFolderPath(windows.FOLDERID_Programs, 0)
	if err != nil {
		return legacyPaths{}, err
	}
	s, err := windows.KnownFolderPath(windows.FOLDERID_Startup, 0)
	if err != nil {
		return legacyPaths{}, err
	}
	return legacyPaths{desktop: d, startMenu: m, startup: s}, nil
}

// legacyLink returns path when it is a shortcut that match accepts.
func legacyLink(path string, match func(args string) bool) (string, bool) {
	if _, err := os.Stat(path); err != nil {
		return "", false
	}
	_, args, err := settings.ReadShortcutTarget(path)
	if err != nil || !match(args) {
		return "", false
	}
	return path, true
}

func (w Windows) legacyDesktopMenu(p legacyPaths) (desktop, menu string) {
	desktop, _ = legacyLink(filepath.Join(p.desktop, portainerLink), legacyPortainer)
	menu, _ = legacyLink(filepath.Join(p.startMenu, portainerLink), legacyPortainer)
	return desktop, menu
}

// Find reads the legacy parts.
func (w Windows) Find() ([]Finding, error) {
	p, err := knownPaths()
	if err != nil {
		return nil, err
	}
	var out []Finding
	if d, m := w.legacyDesktopMenu(p); d != "" || m != "" {
		detail := "Portainer.lnk on the " + where(d != "", m != "")
		out = append(out, Finding{ItemShortcuts, "Replace the Portainer shortcuts with Docker Zero shortcuts", detail})
	}
	if s, ok := legacyLink(filepath.Join(p.startup, startupLink), legacyKeepalive); ok {
		out = append(out, Finding{ItemStartup, "Replace the sign-in start shortcut with Docker Zero's auto start", s})
	}
	if _, err := os.Stat(w.LegacyDir); err == nil {
		out = append(out, Finding{ItemCommand, "Replace the manual docker command with Docker Zero's", w.LegacyDir})
	}
	if idle, err := w.WSLConfig.IdleDisabled(); err == nil && idle && (w.IdleAdopted == nil || !w.IdleAdopted()) {
		out = append(out, Finding{ItemIdle, "Let Docker Zero manage the .wslconfig idle line", w.WSLConfig.Path})
	}
	return out, nil
}

func where(desktop, menu bool) string {
	switch {
	case desktop && menu:
		return "Desktop and in the Start Menu"
	case desktop:
		return "Desktop"
	}
	return "Start Menu"
}

// Replace puts Docker Zero's equivalent in place of one legacy part.
func (w Windows) Replace(item Item) error {
	p, err := knownPaths()
	if err != nil {
		return err
	}
	switch item {
	case ItemShortcuts:
		return w.replaceShortcuts(p)
	case ItemStartup:
		return w.replaceStartup(p)
	case ItemCommand:
		return w.replaceCommand()
	case ItemIdle:
		if w.AdoptIdle == nil {
			return nil
		}
		return w.AdoptIdle()
	}
	return errors.New("unknown item " + string(item))
}

func (w Windows) replaceShortcuts(p legacyPaths) error {
	d, m := w.legacyDesktopMenu(p)
	exe, err := settings.EnsureInstalled(w.InstallExe)
	if err != nil {
		return err
	}
	// Keep the Docker Zero shortcuts that already exist.
	cur, err := settings.ReadShortcuts()
	if err != nil {
		return err
	}
	want := settings.Shortcuts{Desktop: cur.Desktop || d != "", StartMenu: cur.StartMenu || m != ""}
	if err := settings.WriteShortcuts(want, exe); err != nil {
		return err
	}
	return removeFiles(d, m)
}

func (w Windows) replaceStartup(p legacyPaths) error {
	link, _ := legacyLink(filepath.Join(p.startup, startupLink), legacyKeepalive)
	exe, err := settings.EnsureInstalled(w.InstallExe)
	if err != nil {
		return err
	}
	if err := settings.WriteAutoStart(settings.AutoStart{Enabled: true, StartDistro: true}, exe); err != nil {
		return err
	}
	return removeFiles(link)
}

func (w Windows) replaceCommand() error {
	if err := setup.WriteLauncher(w.BinDir, w.Distro); err != nil {
		return err
	}
	if err := setup.AddToUserPath(w.BinDir); err != nil {
		return err
	}
	if err := setup.RemoveFromUserPath(filepath.Join(w.LegacyDir, "bin")); err != nil {
		return err
	}
	return removeKnown(w.LegacyDir, "open-portainer.ps1", "portainer.ico", "portainer.crt", filepath.Join("bin", "docker.cmd"), filepath.Join("bin", "docker.ps1"))
}

func removeFiles(paths ...string) error {
	var errs []error
	for _, p := range paths {
		if p == "" {
			continue
		}
		if err := os.Remove(p); err != nil && !errors.Is(err, fs.ErrNotExist) {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
