package ui

import (
	"context"
	"time"

	"github.com/the-khiem7/docker-zero/internal/adopt"
)

// adoptTimeout bounds a take-over; it only edits shortcuts, folders and the registry.
const adoptTimeout = time.Minute

// AdoptService is bound to the setup window.
type AdoptService struct {
	adopter *adopt.Adopter
	guard   *opGuard
	// onDone runs after a take-over, to refresh the status.
	onDone func()
}

// NewAdoptService builds the service.
func NewAdoptService(a *adopt.Adopter, guard *opGuard) *AdoptService {
	return &AdoptService{adopter: a, guard: guard}
}

// Detect looks for a setup made by hand. It changes nothing and does not
// start the distro.
func (s *AdoptService) Detect(ctx context.Context) (adopt.Report, error) {
	ctx, cancel := context.WithTimeout(ctx, adoptTimeout)
	defer cancel()
	return s.adopter.Detect(ctx)
}

// Adopt replaces the chosen legacy parts and returns the ones it replaced.
// The distro, the container and Portainer's data are not touched.
func (s *AdoptService) Adopt(ctx context.Context, items []string) ([]string, error) {
	if err := s.guard.acquire("take over the existing setup"); err != nil {
		return nil, err
	}
	defer s.guard.release()
	ctx, cancel := context.WithTimeout(ctx, adoptTimeout)
	defer cancel()
	chosen := make([]adopt.Item, len(items))
	for i, it := range items {
		chosen[i] = adopt.Item(it)
	}
	done, err := s.adopter.Adopt(ctx, chosen)
	out := make([]string, len(done))
	for i, d := range done {
		out[i] = string(d)
	}
	if len(done) > 0 && s.onDone != nil {
		s.onDone()
	}
	return out, err
}
