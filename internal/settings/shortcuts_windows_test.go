//go:build windows

package settings

import (
	"path/filepath"
	"testing"
)

func TestShortcutRoundTrip(t *testing.T) {
	lnk := filepath.Join(t.TempDir(), "x.lnk")
	exe := `C:\Windows\System32\notepad.exe`
	if err := createShortcut(lnk, exe); err != nil {
		t.Fatal(err)
	}
	target, args, err := ReadShortcutTarget(lnk)
	if err != nil {
		t.Fatal(err)
	}
	if target != exe || args != "" {
		t.Fatalf("got %q %q", target, args)
	}
}
