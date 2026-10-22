package wsl

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

// decodeOutput turns wsl.exe's own output into a string. wsl.exe writes
// UTF-16LE without a BOM unless WSL_UTF8=1 is honoured, and older builds
// ignore that variable, so both forms must be accepted.
func decodeOutput(b []byte) string {
	b = bytes.TrimPrefix(b, []byte{0xff, 0xfe})
	if !looksUTF16LE(b) {
		return string(bytes.TrimPrefix(b, []byte{0xef, 0xbb, 0xbf}))
	}
	if len(b)%2 == 1 {
		b = b[:len(b)-1]
	}
	u := make([]uint16, len(b)/2)
	for i := range u {
		u[i] = binary.LittleEndian.Uint16(b[2*i:])
	}
	return string(utf16.Decode(u))
}

// looksUTF16LE reports whether most odd bytes are zero, which is true for
// ASCII-range text encoded as UTF-16LE and false for UTF-8.
func looksUTF16LE(b []byte) bool {
	if len(b) < 2 {
		return false
	}
	zeros, pairs := 0, 0
	for i := 1; i < len(b); i += 2 {
		pairs++
		if b[i] == 0 {
			zeros++
		}
	}
	return zeros*2 > pairs
}
