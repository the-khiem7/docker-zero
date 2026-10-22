//go:build windows

package setup

import (
	"errors"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	envKey    = `Environment`
	pathValue = "Path"
)

// AddToUserPath appends dir to the user PATH in HKCU, keeping the value
// type (expandable or plain) the account already uses. It does nothing when
// dir is already listed. Only processes started afterwards see the change.
func AddToUserPath(dir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, envKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func() { _ = k.Close() }()
	cur, typ, err := k.GetStringValue(pathValue)
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	for entry := range strings.SplitSeq(cur, ";") {
		if strings.EqualFold(filepath.Clean(entry), filepath.Clean(dir)) {
			return nil
		}
	}
	next := strings.TrimRight(cur, ";")
	if next != "" {
		next += ";"
	}
	next += dir
	if typ == registry.SZ {
		err = k.SetStringValue(pathValue, next)
	} else {
		err = k.SetExpandStringValue(pathValue, next)
	}
	if err != nil {
		return err
	}
	broadcastEnvChange()
	return nil
}

// RemoveFromUserPath takes dir off the user PATH in HKCU, keeping the value
// type and every other entry. It does nothing when dir is not listed.
func RemoveFromUserPath(dir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, envKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func() { _ = k.Close() }()
	cur, typ, err := k.GetStringValue(pathValue)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	next, found := pathWithout(cur, dir)
	if !found {
		return nil
	}
	if typ == registry.SZ {
		err = k.SetStringValue(pathValue, next)
	} else {
		err = k.SetExpandStringValue(pathValue, next)
	}
	if err != nil {
		return err
	}
	broadcastEnvChange()
	return nil
}

// pathWithout returns the PATH value without the entries equal to dir, and
// whether any was there.
func pathWithout(cur, dir string) (string, bool) {
	var keep []string
	found := false
	for entry := range strings.SplitSeq(cur, ";") {
		if entry != "" && strings.EqualFold(filepath.Clean(entry), filepath.Clean(dir)) {
			found = true
			continue
		}
		keep = append(keep, entry)
	}
	return strings.Join(keep, ";"), found
}

// broadcastEnvChange tells running programs, Explorer first, that the
// environment changed so new terminals pick the PATH up.
func broadcastEnvChange() {
	const (
		hwndBroadcast   = 0xffff
		wmSettingChange = 0x001A
		smtoAbortIfHung = 0x0002
	)
	name, err := windows.UTF16PtrFromString("Environment")
	if err != nil {
		return
	}
	proc := windows.NewLazySystemDLL("user32.dll").NewProc("SendMessageTimeoutW")
	var result uintptr
	// The result is advisory: a hung window must not fail the setup step.
	_, _, _ = proc.Call(hwndBroadcast, wmSettingChange, 0, uintptr(unsafe.Pointer(name)), smtoAbortIfHung, 2000, uintptr(unsafe.Pointer(&result))) //nolint:gosec // syscall arguments are pointers to locals that outlive the call
}
