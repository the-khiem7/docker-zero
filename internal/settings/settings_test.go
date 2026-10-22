package settings

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreRoundTripAndDefaults(t *testing.T) {
	st := Store{Path: filepath.Join(t.TempDir(), "dz", "settings.json")}
	got, err := st.Load()
	if err != nil || got != Defaults() {
		t.Fatalf("missing file: %+v, %v", got, err)
	}
	want := Settings{IdleMode: IdleWSLConfig, Browser: BrowserCustom, BrowserPath: `C:\b\b.exe`, CornerBox: true}
	if err := st.Save(want); err != nil {
		t.Fatal(err)
	}
	if got, err := st.Load(); err != nil || got != want {
		t.Fatalf("round trip: %+v, %v", got, err)
	}
}

func TestStoreRejectsAndRepairs(t *testing.T) {
	st := Store{Path: filepath.Join(t.TempDir(), "settings.json")}
	if err := st.Save(Settings{IdleMode: "always", Browser: BrowserDefault}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v", err)
	}
	if err := os.WriteFile(st.Path, []byte(`{"idleMode":"bogus","browser":"chrome"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.IdleMode != IdleKeepalive || got.Browser != BrowserChrome {
		t.Fatalf("got %+v", got)
	}
}

func TestWithIdleOff(t *testing.T) {
	cases := map[string]string{
		"":                                       "[general]\ninstanceIdleTimeout=-1\n",
		"[wsl2]\nmemory=4GB\n":                   "[general]\ninstanceIdleTimeout=-1\n[wsl2]\nmemory=4GB\n",
		"[general]\ninstanceIdleTimeout=60000\n": "[general]\ninstanceIdleTimeout=-1\n",
		"[General]\r\nfoo=1\r\n":                 "[General]\r\ninstanceIdleTimeout=-1\r\nfoo=1\r\n",
		"[wsl2]\ninstanceIdleTimeout=5\n":        "[general]\ninstanceIdleTimeout=-1\n[wsl2]\ninstanceIdleTimeout=5\n",
	}
	for in, want := range cases {
		if got := withIdleOff(in); got != want {
			t.Errorf("withIdleOff(%q) = %q, want %q", in, got, want)
		}
		if !hasIdleOff(withIdleOff(in)) {
			t.Errorf("hasIdleOff after withIdleOff(%q) is false", in)
		}
	}
}

func TestWithoutIdle(t *testing.T) {
	in := "[general]\ninstanceIdleTimeout=-1\n[wsl2]\nmemory=4GB\n"
	if got, want := withoutIdle(in), "[general]\n[wsl2]\nmemory=4GB\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if hasIdleOff("[wsl2]\ninstanceIdleTimeout=-1\n") {
		t.Fatal("key outside [general] counted")
	}
}

func TestSetIdleDisabledFile(t *testing.T) {
	w := WSLConfig{Path: filepath.Join(t.TempDir(), ".wslconfig")}
	changed, err := w.SetIdleDisabled(true)
	if err != nil || !changed {
		t.Fatalf("create: %v %v", changed, err)
	}
	if on, _ := w.IdleDisabled(); !on {
		t.Fatal("not set")
	}
	if changed, _ := w.SetIdleDisabled(true); changed {
		t.Fatal("second set changed the file")
	}
	if changed, err := w.SetIdleDisabled(false); err != nil || !changed {
		t.Fatalf("remove: %v %v", changed, err)
	}
	if _, err := os.Stat(w.Path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("file with only [general] left behind: %v", err)
	}

	if err := os.WriteFile(w.Path, []byte("[wsl2]\nmemory=4GB\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := w.SetIdleDisabled(true); err != nil {
		t.Fatal(err)
	}
	if _, err := w.SetIdleDisabled(false); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(w.Path)
	if string(b) != "[general]\n[wsl2]\nmemory=4GB\n" {
		t.Fatalf("other content changed: %q", b)
	}
}
