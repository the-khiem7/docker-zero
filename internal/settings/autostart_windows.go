//go:build windows

package settings

import (
	"errors"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	runKey   = `Software\Microsoft\Windows\CurrentVersion\Run`
	runValue = "Docker Zero"
)

// StartDistroFlag makes an auto-started app also start the distro.
const StartDistroFlag = "--start-distro"

// AutoStart is the HKCU Run value that starts the app at sign-in.
type AutoStart struct {
	Enabled     bool `json:"enabled"`
	StartDistro bool `json:"startDistro"`
}

// runCommand builds the Run value: the quoted exe plus the optional flag.
func runCommand(exe string, startDistro bool) string {
	cmd := `"` + exe + `"`
	if startDistro {
		cmd += " " + StartDistroFlag
	}
	return cmd
}

// ReadAutoStart reads the Run value.
func ReadAutoStart() (AutoStart, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return AutoStart{}, err
	}
	defer func() { _ = k.Close() }()
	v, _, err := k.GetStringValue(runValue)
	if errors.Is(err, registry.ErrNotExist) {
		return AutoStart{}, nil
	}
	if err != nil {
		return AutoStart{}, err
	}
	return AutoStart{Enabled: true, StartDistro: strings.Contains(v, StartDistroFlag)}, nil
}

// WriteAutoStart sets or deletes the Run value for exe.
func WriteAutoStart(a AutoStart, exe string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func() { _ = k.Close() }()
	if !a.Enabled {
		if err := k.DeleteValue(runValue); err != nil && !errors.Is(err, registry.ErrNotExist) {
			return err
		}
		return nil
	}
	return k.SetStringValue(runValue, runCommand(exe, a.StartDistro))
}
