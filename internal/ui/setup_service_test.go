package ui

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/setup"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// noWSL reports WSL as unavailable, so the wizard fails at its first step.
type noWSL struct{}

func (noWSL) Available(context.Context) error { return errors.New("no wsl") }
func (noWSL) DistroState(context.Context, string) (wsl.State, error) {
	return wsl.StateNotInstalled, nil
}
func (noWSL) Install(context.Context, string, string, string) error { return nil }
func (noWSL) Export(context.Context, string, string) error          { return nil }
func (noWSL) Unregister(context.Context, string) error              { return nil }
func (noWSL) Terminate(context.Context, string) error               { return nil }
func (noWSL) Start(context.Context, string) error                   { return nil }
func (noWSL) Script(context.Context, string, string, string) (string, error) {
	return "", nil
}

func newTestSetupService() (*SetupService, chan SetupView) {
	views := make(chan SetupView, 64)
	s := NewSetupService(func() *setup.Wizard {
		return &setup.Wizard{WSL: noWSL{}}
	}, &opGuard{})
	s.emit = func(v SetupView) { views <- v }
	return s, views
}

func TestSetupStartRejectsBadCredentials(t *testing.T) {
	s, _ := newTestSetupService()
	if err := s.Start("", strings.Repeat("x", 14)); !errors.Is(err, setup.ErrCredentials) {
		t.Fatalf("empty username: %v", err)
	}
	if err := s.Start("admin", "short"); !errors.Is(err, setup.ErrCredentials) {
		t.Fatalf("short password: %v", err)
	}
	if s.State().Running {
		t.Fatal("setup started with invalid credentials")
	}
}

func TestSetupFailureReachesTheView(t *testing.T) {
	s, views := newTestSetupService()
	if err := s.Start("admin", strings.Repeat("x", 14)); err != nil {
		t.Fatal(err)
	}
	var last SetupView
	timeout := time.After(5 * time.Second)
	for !last.Done && last.Error == "" {
		select {
		case last = <-views:
		case <-timeout:
			t.Fatal("no final view")
		}
	}
	if last.Running || last.Done {
		t.Fatalf("running %v, done %v", last.Running, last.Done)
	}
	if !strings.Contains(last.Error, "WSL is not available") {
		t.Fatalf("error %q", last.Error)
	}
	if got := last.Steps[0]; got.ID != setup.StepWSL || got.Status != setup.StatusFailed {
		t.Fatalf("first step %+v", got)
	}
	if got := last.Steps[1].Status; got != "" {
		t.Fatalf("a later step shows %q", got)
	}
}
