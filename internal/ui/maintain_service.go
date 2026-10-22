package ui

import (
	"context"
	"time"

	"github.com/the-khiem7/docker-zero/internal/setup"
)

const (
	// updateTimeout bounds an update; the image pull dominates it.
	updateTimeout = 20 * time.Minute
	// resetTimeout bounds an admin reset.
	resetTimeout = 5 * time.Minute
)

// MaintainService is bound to the settings window. It runs the actions that
// change Portainer inside the distro.
type MaintainService struct {
	newWizard func() *setup.Wizard
	guard     *opGuard
	// onDone runs after an action that changed the runtime.
	onDone func()
}

// NewMaintainService builds the service. newWizard makes a fresh wizard for each action.
func NewMaintainService(newWizard func() *setup.Wizard, guard *opGuard) *MaintainService {
	return &MaintainService{newWizard: newWizard, guard: guard}
}

// UpdatePortainer pulls the pinned Portainer image and recreates the container
// when it changed. It returns true when the container was replaced.
func (s *MaintainService) UpdatePortainer(ctx context.Context) (bool, error) {
	if err := s.guard.acquire("update Portainer"); err != nil {
		return false, err
	}
	defer s.guard.release()
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()
	updated, err := s.newWizard().UpdatePortainer(ctx)
	if err == nil && updated && s.onDone != nil {
		s.onDone()
	}
	return updated, err
}

// ResetAdmin sets a new password for the Portainer admin account and checks
// it by signing in. The password is not stored.
func (s *MaintainService) ResetAdmin(ctx context.Context, username, secret string) error {
	cr := setup.Credentials{Username: username, Password: secret}
	if err := cr.Validate(); err != nil {
		return err
	}
	if err := s.guard.acquire("reset the admin password"); err != nil {
		return err
	}
	defer s.guard.release()
	ctx, cancel := context.WithTimeout(ctx, resetTimeout)
	defer cancel()
	err := s.newWizard().ResetAdmin(ctx, cr)
	if s.onDone != nil {
		s.onDone()
	}
	return err
}
