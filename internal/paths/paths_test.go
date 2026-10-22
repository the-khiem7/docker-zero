package paths

import (
	"errors"
	"path/filepath"
	"testing"
)

func env(m map[string]string) Getenv {
	return func(k string) string { return m[k] }
}

func TestInstalledExe(t *testing.T) {
	got, err := InstalledExe(env(map[string]string{"LOCALAPPDATA": `C:\Users\u\AppData\Local`}))
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(`C:\Users\u\AppData\Local`, "Programs", "docker-zero", "docker-zero.exe")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestInstallDirMissingEnv(t *testing.T) {
	if _, err := InstallDir(env(nil)); !errors.Is(err, ErrNoLocalAppData) {
		t.Fatalf("got %v, want ErrNoLocalAppData", err)
	}
}

func TestConfigDir(t *testing.T) {
	got, err := ConfigDir(env(map[string]string{"APPDATA": `C:\Users\u\AppData\Roaming`}))
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(`C:\Users\u\AppData\Roaming`, "docker-zero"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
