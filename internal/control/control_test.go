package control

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/status"
	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

type fakeHold struct{ released bool }

func (h *fakeHold) Alive() bool { return !h.released }
func (h *fakeHold) Release()    { h.released = true }

type fakeWSL struct {
	mu    sync.Mutex
	state wsl.State
	log   []string
	block chan struct{}
}

func (f *fakeWSL) record(s string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.log = append(f.log, s)
}

func (f *fakeWSL) Available(context.Context) error { return nil }

func (f *fakeWSL) DistroState(context.Context, string) (wsl.State, error) { return f.state, nil }

func (f *fakeWSL) Script(_ context.Context, distro, user, _ string) (string, error) {
	f.record("script " + distro + " " + user)
	if f.block != nil {
		<-f.block
	}
	return "", nil
}

func (f *fakeWSL) Terminate(_ context.Context, distro string) error {
	f.record("terminate " + distro)
	return nil
}

func (f *fakeWSL) Shutdown(context.Context) error {
	f.record("shutdown")
	return nil
}

type fakeProber struct{ upAfter, calls int }

func (p *fakeProber) Probe(context.Context) status.Snapshot {
	p.calls++
	return status.Snapshot{PortainerUp: p.calls > p.upAfter}
}

type harness struct {
	c      *Controller
	wsl    *fakeWSL
	holds  []*fakeHold
	opened []string
}

func newHarness(state wsl.State, upAfter int) *harness {
	h := &harness{wsl: &fakeWSL{state: state}}
	h.c = &Controller{
		WSL:    h.wsl,
		Prober: &fakeProber{upAfter: upAfter},
		StartHold: func(string) (Holder, error) {
			fh := &fakeHold{}
			h.holds = append(h.holds, fh)
			return fh, nil
		},
		OpenURL:      func(u string) error { h.opened = append(h.opened, u); return nil },
		Wait:         time.Millisecond,
		StartTimeout: time.Second,
	}
	return h
}

func TestStartHoldsEnsuresAndOpens(t *testing.T) {
	h := newHarness(wsl.StateStopped, 2)
	if err := h.c.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	if len(h.holds) != 1 || !h.c.Holding() {
		t.Fatalf("holds %d, holding %v", len(h.holds), h.c.Holding())
	}
	if want := []string{"script " + target.Distro + " root"}; !slices.Equal(h.wsl.log, want) {
		t.Fatalf("log %q", h.wsl.log)
	}
	if !slices.Equal(h.opened, []string{target.PortainerURL}) {
		t.Fatalf("opened %q", h.opened)
	}
	// A second Start reuses the live hold.
	if err := h.c.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	if len(h.holds) != 1 {
		t.Fatalf("holds %d after second start", len(h.holds))
	}
}

func TestStartNotInstalled(t *testing.T) {
	h := newHarness(wsl.StateNotInstalled, 0)
	if err := h.c.Start(t.Context()); !errors.Is(err, ErrNotInstalled) {
		t.Fatalf("got %v", err)
	}
	if len(h.holds) != 0 {
		t.Fatal("hold started for a missing distro")
	}
}

func TestStartTimesOut(t *testing.T) {
	h := newHarness(wsl.StateRunning, 1<<30)
	h.c.StartTimeout = 20 * time.Millisecond
	if err := h.c.Start(t.Context()); !errors.Is(err, ErrTimeout) {
		t.Fatalf("got %v", err)
	}
	if len(h.opened) != 0 {
		t.Fatal("browser opened before Portainer answered")
	}
}

func TestStopAndShutdownReleaseHold(t *testing.T) {
	h := newHarness(wsl.StateRunning, 0)
	if err := h.c.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := h.c.StopDistro(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !h.holds[0].released || h.c.Holding() {
		t.Fatal("stop kept the hold")
	}
	if err := h.c.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := h.c.ShutdownWSL(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !h.holds[1].released {
		t.Fatal("shutdown kept the hold")
	}
	want := []string{
		"script " + target.Distro + " root",
		"terminate " + target.Distro,
		"script " + target.Distro + " root",
		"shutdown",
	}
	if !slices.Equal(h.wsl.log, want) {
		t.Fatalf("log %q", h.wsl.log)
	}
}

func TestOneActionAtATime(t *testing.T) {
	h := newHarness(wsl.StateRunning, 0)
	h.wsl.block = make(chan struct{})
	done := make(chan error)
	go func() { done <- h.c.Start(context.Background()) }()
	for h.c.Busy() == "" {
		time.Sleep(time.Millisecond)
	}
	if err := h.c.StopDistro(t.Context()); !errors.Is(err, ErrBusy) {
		t.Fatalf("got %v", err)
	}
	close(h.wsl.block)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
