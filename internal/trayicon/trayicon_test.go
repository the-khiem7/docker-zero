package trayicon

import (
	"bytes"
	"image/png"
	"testing"
)

func TestPNGDecodesAndDiffersByLevel(t *testing.T) {
	seen := map[string]Level{}
	for _, l := range []Level{Stopped, Partial, Running, Failed} {
		b := PNG(l)
		img, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			t.Fatalf("level %d: %v", l, err)
		}
		if got := img.Bounds().Dx(); got != size {
			t.Fatalf("width %d", got)
		}
		if prev, dup := seen[string(b)]; dup {
			t.Fatalf("levels %d and %d draw the same icon", prev, l)
		}
		seen[string(b)] = l
	}
}
