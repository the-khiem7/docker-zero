//go:build integration && windows

package setup

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/paths"
)

// TestRemoveAfterExitDeletesFolders runs the real detached cmd.exe on two
// throwaway folders and checks both are gone once its delay has passed.
func TestRemoveAfterExitDeletesFolders(t *testing.T) {
	base := t.TempDir()
	dirs := []string{filepath.Join(base, "a", paths.AppDirName), filepath.Join(base, "b", paths.AppDirName)}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(d, "sub"), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "sub", "f.txt"), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := removeAfterExit(dirs...); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(20 * time.Second)
	for _, d := range dirs {
		for {
			if _, err := os.Stat(d); os.IsNotExist(err) {
				break
			}
			if time.Now().After(deadline) {
				t.Fatalf("%s still exists", d)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
