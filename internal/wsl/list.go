package wsl

import (
	"strings"
)

// State is the WSL state of one distro.
type State string

// Distro states as Docker Zero sees them.
const (
	StateRunning      State = "running"
	StateStopped      State = "stopped"
	StateNotInstalled State = "not-installed"
)

// parseNames parses `wsl.exe --list --quiet` output: one distro name per
// line. The quiet form is used because the verbose table's header and state
// column are localised on non-English Windows.
func parseNames(out string) []string {
	var names []string
	for line := range strings.SplitSeq(out, "\n") {
		name := strings.TrimSpace(line)
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
