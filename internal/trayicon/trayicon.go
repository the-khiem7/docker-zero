// Package trayicon draws the tray icon: a blue tile with a white "Z" and a
// status dot, rendered at start-up so no image files need shipping.
package trayicon

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
)

// Level is the overall health the dot shows.
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
	tile  = color.NRGBA{0x00, 0x67, 0xc0, 0xff}
	glyph = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	ring  = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	dots  = map[Level]color.NRGBA{
		Stopped: {0x8a, 0x8a, 0x8a, 0xff},
		Partial: {0xf7, 0xb5, 0x00, 0xff},
		Running: {0x16, 0xc6, 0x0c, 0xff},
		Failed:  {0xe8, 0x11, 0x23, 0xff},
	}
)

// PNG returns the icon for level as PNG bytes.
func PNG(level Level) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			if inRoundRect(fx, fy, 1, 1, 30, 30, 6) {
				img.SetNRGBA(x, y, tile)
			}
			if onZ(fx, fy) {
				img.SetNRGBA(x, y, glyph)
			}
			d := math.Hypot(fx-25, fy-25)
			switch {
			case d <= 5:
				img.SetNRGBA(x, y, dots[level])
			case d <= 7:
				img.SetNRGBA(x, y, ring)
			}
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img) // encoding to a bytes.Buffer cannot fail
	return buf.Bytes()
}

func inRoundRect(x, y, x0, y0, x1, y1, r float64) bool {
	if x < x0 || x > x1 || y < y0 || y > y1 {
		return false
	}
	cx := math.Max(x0+r, math.Min(x, x1-r))
	cy := math.Max(y0+r, math.Min(y, y1-r))
	return math.Hypot(x-cx, y-cy) <= r
}

// onZ draws a Z from (8,8) to (22,22) with a stroke of about 3.5 px.
func onZ(x, y float64) bool {
	const left, right, top, bottom, stroke = 8.0, 22.0, 8.0, 22.0, 3.5
	if x < left || x > right {
		return false
	}
	if y >= top && y <= top+stroke || y >= bottom-stroke && y <= bottom {
		return true
	}
	// Diagonal from top-right to bottom-left: distance to the line.
	ax, ay, bx, by := right, top+stroke, left, bottom-stroke
	dist := math.Abs((by-ay)*x-(bx-ax)*y+bx*ay-by*ax) / math.Hypot(by-ay, bx-ax)
	return y > top && y < bottom && dist <= stroke/2+0.5
}
