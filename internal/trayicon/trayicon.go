// Package trayicon holds the tray icon: the app logo at 32 px, its ring
// recoloured per status level, rendered from the brand source and embedded.
package trayicon

import _ "embed"

// Level is the overall health the icon shows.
type Level int

// Levels, from nothing running to everything up.
const (
	Stopped Level = iota
	Partial
	Running
	Failed
)

const size = 32

var (
	//go:embed stopped.png
	stoppedPNG []byte
	//go:embed partial.png
	partialPNG []byte
	//go:embed running.png
	runningPNG []byte
	//go:embed failed.png
	failedPNG []byte
)

// PNG returns the icon for level as PNG bytes; an unknown level shows Stopped.
func PNG(level Level) []byte {
	switch level {
	case Partial:
		return partialPNG
	case Running:
		return runningPNG
	case Failed:
		return failedPNG
	default:
		return stoppedPNG
	}
}
