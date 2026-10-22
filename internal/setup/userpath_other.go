//go:build !windows

package setup

import "errors"

// AddToUserPath is only implemented on Windows.
func AddToUserPath(string) error {
	return errors.New("changing the user PATH is only supported on Windows")
}
