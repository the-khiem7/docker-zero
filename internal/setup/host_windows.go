//go:build windows

package setup

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/the-khiem7/docker-zero/internal/paths"
	"github.com/the-khiem7/docker-zero/internal/settings"
)

// HostConfig locates what NewHost removes.
type HostConfig struct {
	// InstallDir holds the installed executable.
	InstallDir string
	// ConfigDir holds settings.json.
	ConfigDir string
	// WSLConfig is the user's .wslconfig.
	WSLConfig settings.WSLConfig
	// IdleLineIsOurs reports whether the app set instanceIdleTimeout; a line
	// the user wrote themselves stays.
	IdleLineIsOurs func() bool
	// ReleaseHold lets go of the keepalive session.
	ReleaseHold func()
}

// ErrUnsafeDir is returned for an install folder that is not Docker Zero's own.
var ErrUnsafeDir = errors.New("refusing to delete a folder that is not the Docker Zero install folder")

// NewHost returns the hooks that remove the real per-user items.
func NewHost(c HostConfig) Host {
	return Host{
		ReleaseHold: c.ReleaseHold,
		UntrustCert: func(pem []byte) (bool, error) {
			thumb, _, err := settings.ParseCert(pem)
			if err != nil {
				return false, err
			}
			trusted, err := settings.IsTrusted(thumb)
			if err != nil || !trusted {
				return false, err
			}
			return true, settings.Untrust(thumb)
		},
		RemoveAutoStart: func() error { return settings.WriteAutoStart(settings.AutoStart{}, "") },
		RemoveShortcuts: func() error { return settings.WriteShortcuts(settings.Shortcuts{}, "") },
		RemoveFromPath:  RemoveFromUserPath,
		RemoveIdleLine: func() error {
			if c.IdleLineIsOurs == nil || !c.IdleLineIsOurs() {
				return nil
			}
			_, err := c.WSLConfig.SetIdleDisabled(false)
			return err
		},
		RemoveConfig: func() error { return os.RemoveAll(c.ConfigDir) },
		RemoveInstallDirAfterExit: func() error {
			return removeAfterExit(c.InstallDir)
		},
	}
}

// removeAfterExit starts a detached shell that waits a few seconds, long
// enough for this process to end, then deletes dir. A running executable
// cannot delete its own folder.
func removeAfterExit(dir string) error {
	dir = filepath.Clean(dir)
	// The folder name is the guard against a wrong path; the characters keep
	// the command line below from being reinterpreted by cmd.exe.
	if filepath.Base(dir) != paths.AppDirName || strings.ContainsAny(dir, `"%`) {
		return ErrUnsafeDir
	}
	cmd := exec.Command("cmd.exe") //nolint:gosec // fixed executable; the command line is built below from a checked path
	// CmdLine bypasses Go's argument quoting, which cmd.exe does not parse.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: windows.CREATE_NO_WINDOW | windows.CREATE_NEW_PROCESS_GROUP,
		CmdLine:       `cmd.exe /c "ping -n 5 127.0.0.1 >nul & rmdir /s /q "` + dir + `""`,
	}
	return cmd.Start()
}
