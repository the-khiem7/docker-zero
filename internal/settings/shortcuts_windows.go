//go:build windows

package settings

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows"
)

const shortcutName = "Docker Zero.lnk"

// Shortcuts says which Docker Zero shortcuts exist.
type Shortcuts struct {
	Desktop   bool `json:"desktop"`
	StartMenu bool `json:"startMenu"`
}

func shortcutPaths() (desktop, startMenu string, err error) {
	d, err := windows.KnownFolderPath(windows.FOLDERID_Desktop, 0)
	if err != nil {
		return "", "", err
	}
	p, err := windows.KnownFolderPath(windows.FOLDERID_Programs, 0)
	if err != nil {
		return "", "", err
	}
	return filepath.Join(d, shortcutName), filepath.Join(p, shortcutName), nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadShortcuts checks both shortcut files.
func ReadShortcuts() (Shortcuts, error) {
	d, p, err := shortcutPaths()
	if err != nil {
		return Shortcuts{}, err
	}
	return Shortcuts{Desktop: exists(d), StartMenu: exists(p)}, nil
}

// WriteShortcuts creates or removes each shortcut so the files match s.
func WriteShortcuts(s Shortcuts, exe string) error {
	d, p, err := shortcutPaths()
	if err != nil {
		return err
	}
	return errors.Join(setShortcut(d, s.Desktop, exe), setShortcut(p, s.StartMenu, exe))
}

func setShortcut(path string, want bool, exe string) error {
	if !want {
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return nil
	}
	return createShortcut(path, exe)
}

// withShell runs fn with the WScript.Shell COM object.
func withShell(fn func(shell *ole.IDispatch) error) error {
	// COM is initialised per OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		var oleErr *ole.OleError
		// S_FALSE: already initialised on this thread, which is fine.
		if !errors.As(err, &oleErr) || oleErr.Code() != 1 {
			return err
		}
	}
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer unknown.Release()
	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer shell.Release()
	return fn(shell)
}

// ReadShortcutTarget returns the target and the arguments of the .lnk at path.
func ReadShortcutTarget(path string) (target, args string, err error) {
	err = withShell(func(shell *ole.IDispatch) error {
		res, err := oleutil.CallMethod(shell, "CreateShortcut", path)
		if err != nil {
			return err
		}
		link := res.ToIDispatch()
		defer link.Release()
		t, err := oleutil.GetProperty(link, "TargetPath")
		if err != nil {
			return err
		}
		a, err := oleutil.GetProperty(link, "Arguments")
		if err != nil {
			return err
		}
		target, args = t.ToString(), a.ToString()
		return nil
	})
	return target, args, err
}

// createShortcut writes a .lnk through the WScript.Shell COM object.
func createShortcut(path, exe string) error {
	return withShell(func(shell *ole.IDispatch) error {
		res, err := oleutil.CallMethod(shell, "CreateShortcut", path)
		if err != nil {
			return err
		}
		link := res.ToIDispatch()
		defer link.Release()
		for name, value := range map[string]string{
			"TargetPath":       exe,
			"WorkingDirectory": filepath.Dir(exe),
			"IconLocation":     exe + ",0",
			"Description":      "Docker Engine in WSL2 with Portainer",
		} {
			if _, err := oleutil.PutProperty(link, name, value); err != nil {
				return err
			}
		}
		_, err = oleutil.CallMethod(link, "Save")
		return err
	})
}
