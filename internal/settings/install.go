package settings

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// EnsureInstalled copies the running exe to dest unless it already runs
// from there, so auto start and shortcuts keep working when the download
// folder is cleaned up. It returns the path to point them at.
func EnsureInstalled(dest string) (string, error) {
	self, err := os.Executable()
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(self); err == nil {
		self = resolved
	}
	if strings.EqualFold(filepath.Clean(self), filepath.Clean(dest)) {
		return dest, nil
	}
	return dest, copyFile(self, dest)
}

func copyFile(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return err
	}
	in, err := os.Open(src) //nolint:gosec // src is this process's own executable
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	tmp := dest + ".new"
	out, err := os.Create(tmp) //nolint:gosec // dest is under the user's LOCALAPPDATA
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	// Rename fails while the installed copy runs; the caller reports it.
	return os.Rename(tmp, dest)
}
