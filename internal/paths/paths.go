// Package paths resolves the per-user locations Docker Zero installs into.
// Every location is under the user profile so no step needs elevation.
package paths

import (
	"errors"
	"os"
	"path/filepath"
)

// AppDirName is the folder name used under %LOCALAPPDATA%\Programs and %APPDATA%.
const AppDirName = "docker-zero"

// ExeName is the file name the app copies itself to.
const ExeName = "docker-zero.exe"

// SettingsFile is the settings file name inside the config folder.
const SettingsFile = "settings.json"

// ErrNoLocalAppData is returned when LOCALAPPDATA is not set.
var ErrNoLocalAppData = errors.New("LOCALAPPDATA is not set")

// Getenv looks up an environment variable; tests replace it.
type Getenv func(string) string

// InstallDir returns %LOCALAPPDATA%\Programs\docker-zero.
func InstallDir(getenv Getenv) (string, error) {
	base := getenv("LOCALAPPDATA")
	if base == "" {
		return "", ErrNoLocalAppData
	}
	return filepath.Join(base, "Programs", AppDirName), nil
}

// InstalledExe returns the full path of the installed executable.
func InstalledExe(getenv Getenv) (string, error) {
	dir, err := InstallDir(getenv)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ExeName), nil
}

// ConfigDir returns %APPDATA%\docker-zero, falling back to os.UserConfigDir.
func ConfigDir(getenv Getenv) (string, error) {
	if base := getenv("APPDATA"); base != "" {
		return filepath.Join(base, AppDirName), nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, AppDirName), nil
}
