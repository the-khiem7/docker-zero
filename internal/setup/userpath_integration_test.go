//go:build integration && windows

package setup

import (
	"testing"

	"golang.org/x/sys/windows/registry"
)

// TestUserPathRoundTrip adds a folder to the real user PATH, checks it is
// listed, removes it, and checks the value and its type are back exactly as
// they were. The original value is restored even when the test fails midway.
func TestUserPathRoundTrip(t *testing.T) {
	k, err := registry.OpenKey(registry.CURRENT_USER, envKey, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = k.Close() }()
	before, typ, err := k.GetStringValue(pathValue)
	if err != nil {
		t.Skipf("no user PATH value to round-trip: %v", err)
	}
	t.Cleanup(func() {
		if typ == registry.SZ {
			_ = k.SetStringValue(pathValue, before)
		} else {
			_ = k.SetExpandStringValue(pathValue, before)
		}
	})

	dir := t.TempDir()
	if err := AddToUserPath(dir); err != nil {
		t.Fatal(err)
	}
	during, _, err := k.GetStringValue(pathValue)
	if err != nil {
		t.Fatal(err)
	}
	if _, found := pathWithout(during, dir); !found {
		t.Fatalf("%s not on the PATH after AddToUserPath", dir)
	}
	if err := RemoveFromUserPath(dir); err != nil {
		t.Fatal(err)
	}
	after, afterTyp, err := k.GetStringValue(pathValue)
	if err != nil {
		t.Fatal(err)
	}
	if after != before || afterTyp != typ {
		t.Fatalf("PATH changed by the round trip:\n before %q (%d)\n after  %q (%d)", before, typ, after, afterTyp)
	}
	if err := RemoveFromUserPath(dir); err != nil {
		t.Fatalf("a second removal must succeed: %v", err)
	}
}
