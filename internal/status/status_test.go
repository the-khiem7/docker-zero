package status

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/the-khiem7/docker-zero/internal/wsl"
)

type fakeWSL struct {
	availErr error
	state    wsl.State
	script   string
	scripted bool
}

func (f *fakeWSL) Available(context.Context) error { return f.availErr }

func (f *fakeWSL) DistroState(context.Context, string) (wsl.State, error) { return f.state, nil }

func (f *fakeWSL) Script(context.Context, string, string, string) (string, error) {
	f.scripted = true
	return f.script, nil
}

func newProber(t *testing.T, f *fakeWSL, handler http.HandlerFunc) *Prober {
	t.Helper()
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)
	p := NewProber(f)
	p.StatusURL = srv.URL + "/api/system/status"
	p.Now = func() time.Time { return time.Unix(0, 0) }
	return p
}

func ok(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`{"Version":"2.45.2","InstanceID":"x"}`))
}

func TestProbeAllUp(t *testing.T) {
	f := &fakeWSL{state: wsl.StateRunning, script: "active\nrunning\n"}
	s := newProber(t, f, ok).Probe(t.Context())
	want := Snapshot{
		WSLAvailable: true, Distro: wsl.StateRunning, Docker: "active", Portainer: "running",
		PortainerUp: true, PortainerVersion: "2.45.2", CheckedAt: time.Unix(0, 0),
	}
	if s != want {
		t.Fatalf("got %+v\nwant %+v", s, want)
	}
}

func TestProbeStoppedDistroIsNotStarted(t *testing.T) {
	f := &fakeWSL{state: wsl.StateStopped}
	s := newProber(t, f, ok).Probe(t.Context())
	if f.scripted {
		t.Fatal("probe ran a script in a stopped distro, which would start it")
	}
	if s.Distro != wsl.StateStopped || s.Docker != Unknown || s.PortainerUp {
		t.Fatalf("got %+v", s)
	}
}

func TestProbeWSLUnavailable(t *testing.T) {
	f := &fakeWSL{availErr: wsl.ErrWSLUnavailable}
	s := newProber(t, f, ok).Probe(t.Context())
	if s.WSLAvailable || s.Error == "" {
		t.Fatalf("got %+v", s)
	}
}

func TestProbeContainerMissing(t *testing.T) {
	f := &fakeWSL{state: wsl.StateRunning, script: "active\nmissing\n"}
	s := newProber(t, f, ok).Probe(t.Context())
	if s.Portainer != Missing || s.PortainerUp {
		t.Fatalf("got %+v", s)
	}
}

func TestProbePortainerHTTPError(t *testing.T) {
	f := &fakeWSL{state: wsl.StateRunning, script: "active\nrunning\n"}
	s := newProber(t, f, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}).Probe(t.Context())
	if s.PortainerUp || s.Error == "" {
		t.Fatalf("got %+v", s)
	}
}

func TestParseProbe(t *testing.T) {
	d, p := parseProbe("inactive\r\n")
	if d != "inactive" || p != Unknown {
		t.Fatalf("got %q %q", d, p)
	}
}
