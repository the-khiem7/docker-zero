package settings

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

// The key that stops WSL from stopping idle distros. It lives under
// [general] and applies to every WSL2 distro of the account once WSL has
// been shut down and started again.
const (
	idleKey     = "instanceIdleTimeout"
	idleValue   = "-1"
	generalHead = "[general]"
)

// WSLConfig edits %USERPROFILE%\.wslconfig.
type WSLConfig struct {
	Path string
}

// IdleDisabled reports whether the file sets instanceIdleTimeout=-1.
func (w WSLConfig) IdleDisabled() (bool, error) {
	b, err := os.ReadFile(w.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return hasIdleOff(string(b)), nil
}

// SetIdleDisabled adds or removes the line and reports whether the file
// changed. A file left with nothing but an empty [general] is deleted.
func (w WSLConfig) SetIdleDisabled(on bool) (bool, error) {
	b, err := os.ReadFile(w.Path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return false, err
	}
	old := string(b)
	var updated string
	if on {
		updated = withIdleOff(old)
	} else {
		updated = withoutIdle(old)
	}
	if updated == old {
		return false, nil
	}
	if strings.TrimSpace(updated) == "" || strings.TrimSpace(updated) == generalHead {
		if err := os.Remove(w.Path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}
		return true, nil
	}
	return true, os.WriteFile(w.Path, []byte(updated), 0o600) //nolint:gosec // Path is the user's own .wslconfig, set by the app
}

func isSection(line string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]")
}

func isGeneral(line string) bool {
	return strings.EqualFold(strings.TrimSpace(line), generalHead)
}

func isIdleLine(line string) bool {
	k, _, ok := strings.Cut(line, "=")
	return ok && strings.EqualFold(strings.TrimSpace(k), idleKey)
}

// lines splits text keeping the file's line ending style.
func lines(text string) ([]string, string) {
	eol := "\n"
	if strings.Contains(text, "\r\n") {
		eol = "\r\n"
	}
	text = strings.TrimSuffix(text, eol)
	if text == "" {
		return nil, eol
	}
	return strings.Split(text, eol), eol
}

func hasIdleOff(text string) bool {
	ls, _ := lines(text)
	inGeneral := false
	for _, l := range ls {
		if isSection(l) {
			inGeneral = isGeneral(l)
			continue
		}
		if inGeneral && isIdleLine(l) {
			_, v, _ := strings.Cut(l, "=")
			return strings.TrimSpace(v) == idleValue
		}
	}
	return false
}

// withIdleOff sets instanceIdleTimeout=-1 in [general], replacing another
// value and creating the section when missing.
func withIdleOff(text string) string {
	ls, eol := lines(text)
	line := idleKey + "=" + idleValue
	inGeneral, generalAt := false, -1
	for i, l := range ls {
		if isSection(l) {
			inGeneral = isGeneral(l)
			if inGeneral {
				generalAt = i
			}
			continue
		}
		if inGeneral && isIdleLine(l) {
			ls[i] = line
			return strings.Join(ls, eol) + eol
		}
	}
	if generalAt >= 0 {
		ls = append(ls[:generalAt+1], append([]string{line}, ls[generalAt+1:]...)...)
	} else {
		ls = append([]string{generalHead, line}, ls...)
	}
	return strings.Join(ls, eol) + eol
}

// withoutIdle drops the instanceIdleTimeout line from [general].
func withoutIdle(text string) string {
	ls, eol := lines(text)
	out := ls[:0:0]
	inGeneral := false
	for _, l := range ls {
		if isSection(l) {
			inGeneral = isGeneral(l)
		} else if inGeneral && isIdleLine(l) {
			continue
		}
		out = append(out, l)
	}
	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, eol) + eol
}
