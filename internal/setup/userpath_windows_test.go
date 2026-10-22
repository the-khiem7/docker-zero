//go:build windows

package setup

import "testing"

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
