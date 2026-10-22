// Package settings stores Docker Zero's preferences and applies the ones
// that change the machine: auto start, shortcuts, the WSL idle mode, the
// browser, and the locally trusted Portainer certificate. Every change is
// per-user; nothing needs elevation.
package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// IdleMode says how the distro is kept from idling out.
type IdleMode string

// Idle modes. Keepalive is the default because it changes nothing outside
// the app; WSLConfig affects every WSL2 distro of the account.
const (
	IdleKeepalive IdleMode = "keepalive"
	IdleWSLConfig IdleMode = "wslconfig"
	IdleDefault   IdleMode = "default"
)

// BrowserKind selects the browser Open Portainer uses.
type BrowserKind string

// Browser kinds.
const (
	BrowserDefault BrowserKind = "default"
	BrowserChrome  BrowserKind = "chrome"
	BrowserCustom  BrowserKind = "custom"
)

// Settings are the stored preferences. Auto start, shortcuts and the
// certificate are read from the machine instead, so they cannot drift.
type Settings struct {
	IdleMode    IdleMode    `json:"idleMode"`
	Browser     BrowserKind `json:"browser"`
	BrowserPath string      `json:"browserPath,omitempty"`
	CornerBox   bool        `json:"cornerBox,omitzero"`
}

// Defaults returns the settings of a fresh install.
func Defaults() Settings {
	return Settings{IdleMode: IdleKeepalive, Browser: BrowserDefault}
}

// ErrInvalid reports a value outside the allowed set.
var ErrInvalid = errors.New("invalid setting")

// Validate checks enumerated values.
func (s Settings) Validate() error {
	switch s.IdleMode {
	case IdleKeepalive, IdleWSLConfig, IdleDefault:
	default:
		return fmt.Errorf("%w: idle mode %q", ErrInvalid, s.IdleMode)
	}
	switch s.Browser {
	case BrowserDefault, BrowserChrome:
	case BrowserCustom:
		if s.BrowserPath == "" {
			return fmt.Errorf("%w: custom browser needs a path", ErrInvalid)
		}
	default:
		return fmt.Errorf("%w: browser %q", ErrInvalid, s.Browser)
	}
	return nil
}

// Store reads and writes settings.json.
type Store struct {
	Path string
}

// Load returns the stored settings, or the defaults when none are stored.
// Unknown or invalid values fall back to the defaults field by field.
func (st Store) Load() (Settings, error) {
	s := Defaults()
	b, err := os.ReadFile(st.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return Defaults(), fmt.Errorf("read %s: %w", st.Path, err)
	}
	if s.Validate() != nil {
		d := Defaults()
		if (Settings{IdleMode: s.IdleMode, Browser: d.Browser}).Validate() != nil {
			s.IdleMode = d.IdleMode
		}
		if (Settings{IdleMode: d.IdleMode, Browser: s.Browser, BrowserPath: s.BrowserPath}).Validate() != nil {
			s.Browser, s.BrowserPath = d.Browser, ""
		}
	}
	return s, nil
}

// Save writes s atomically.
func (st Store) Save(s Settings) error {
	if err := s.Validate(); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(st.Path), 0o700); err != nil {
		return err
	}
	tmp := st.Path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, st.Path)
}
