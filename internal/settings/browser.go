package settings

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// ErrChromeNotFound is returned when Google Chrome is not installed.
var ErrChromeNotFound = errors.New("google Chrome not found")

// ChromeCandidates lists where Google Chrome installs. Chrome is started by
// full path because the App Paths entry for chrome.exe can point at another
// Chromium-based browser.
func ChromeCandidates(getenv func(string) string) []string {
	var out []string
	for _, base := range []string{getenv("ProgramFiles"), getenv("ProgramFiles(x86)"), getenv("LOCALAPPDATA")} {
		if base != "" {
			out = append(out, filepath.Join(base, "Google", "Chrome", "Application", "chrome.exe"))
		}
	}
	return out
}

// FindChrome returns the first existing Chrome executable.
func FindChrome(getenv func(string) string, stat func(string) (os.FileInfo, error)) (string, error) {
	for _, p := range ChromeCandidates(getenv) {
		if fi, err := stat(p); err == nil && !fi.IsDir() {
			return p, nil
		}
	}
	return "", ErrChromeNotFound
}

// Opener opens URLs according to the browser setting.
type Opener struct {
	// Default opens a URL in the system default browser.
	Default func(url string) error
	// Current returns the browser setting at the time of the call.
	Current func() Settings
}

// Open opens url in the chosen browser.
func (o Opener) Open(url string) error {
	s := o.Current()
	switch s.Browser {
	case BrowserChrome:
		exe, err := FindChrome(os.Getenv, os.Stat)
		if err != nil {
			return err
		}
		return startDetached(exe, url)
	case BrowserCustom:
		return startDetached(s.BrowserPath, url)
	default:
		return o.Default(url)
	}
}

func startDetached(exe, url string) error {
	cmd := exec.Command(exe, url) //nolint:gosec // the user picked this browser; the URL is the fixed Portainer page
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
