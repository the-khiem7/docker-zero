//go:build windows

package setup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/the-khiem7/docker-zero/internal/paths"
)

func TestPathWithout(t *testing.T) {
	cases := []struct {
		name, cur, dir, want string
		found                bool
	}{
		{"middle", `C:\a;C:\b;C:\c`, `C:\b`, `C:\a;C:\c`, true},
		{"only entry", `C:\b`, `C:\b`, ``, true},
		{"case and slash", `C:\a;c:\B\`, `C:\b`, `C:\a`, true},
		{"absent", `C:\a;C:\c`, `C:\b`, `C:\a;C:\c`, false},
		{"empty value", ``, `C:\b`, ``, false},
		{"expandable entry kept", `%USERPROFILE%\bin;C:\b`, `C:\b`, `%USERPROFILE%\bin`, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, found := pathWithout(c.cur, c.dir)
			if got != c.want || found != c.found {
				t.Fatalf("pathWithout(%q, %q) = %q, %v; want %q, %v", c.cur, c.dir, got, found, c.want, c.found)
			}
		})
	}
}

func TestRemoveAfterExitRefusesOtherFolders(t *testing.T) {
	for _, dir := range []string{`C:\Users`, `C:\`, `C:\x\docker-zero"&calc`, `C:\%TEMP%\docker-zero`} {
		if err := removeAfterExit(dir); err == nil {
			t.Errorf("accepted %q", dir)
		}
	}
}

func TestRemoveAfterExitRefusesWhenAnyFolderIsForeign(t *testing.T) {
	if err := removeAfterExit(`C:\x\docker-zero`, `C:\Users`); err == nil {
		t.Error("accepted a foreign second folder")
	}
}

func TestRemoveConfigKeepsTheFolder(t *testing.T) {
	dir := filepath.Join(t.TempDir(), paths.AppDirName)
	if err := os.MkdirAll(filepath.Join(dir, "webview"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, paths.SettingsFile), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := NewHost(HostConfig{ConfigDir: dir})
	for range 2 {
		if err := h.RemoveConfig(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, paths.SettingsFile)); !os.IsNotExist(err) {
		t.Errorf("settings file still there: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "webview")); err != nil {
		t.Errorf("webview folder removed early: %v", err)
	}
}
