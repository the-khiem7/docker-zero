package ui

import (
	"sync"

	"github.com/the-khiem7/docker-zero/internal/settings"
)

// Prefs is the in-memory copy of the stored settings, shared by the
// services and the browser opener.
type Prefs struct {
	store settings.Store

	mu  sync.Mutex
	cur settings.Settings
}

// LoadPrefs reads the stored settings; on error it keeps the defaults and
// returns the error for logging.
func LoadPrefs(store settings.Store) (*Prefs, error) {
	cur, err := store.Load()
	return &Prefs{store: store, cur: cur}, err
}

// Get returns the current settings.
func (p *Prefs) Get() settings.Settings {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cur
}

// Update changes the settings and saves them. A failed save leaves the
// in-memory copy unchanged.
func (p *Prefs) Update(fn func(*settings.Settings)) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	next := p.cur
	fn(&next)
	if err := p.store.Save(next); err != nil {
		return err
	}
	p.cur = next
	return nil
}
