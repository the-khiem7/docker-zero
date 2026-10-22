package adopt

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/the-khiem7/docker-zero/internal/wsl"
)

type fakeWSL struct {
	state   wsl.State
	out     string
	scripts int
}

func (f *fakeWSL) DistroState(context.Context, string) (wsl.State, error) { return f.state, nil }
func (f *fakeWSL) Script(context.Context, string, string, string) (string, error) {
	f.scripts++
	return f.out, nil
}

type fakeMachine struct {
	found    []Finding
	replaced []Item
	fail     Item
}

func (m *fakeMachine) Find() ([]Finding, error) { return m.found, nil }
func (m *fakeMachine) Replace(it Item) error {
	if it == m.fail {
		return errors.New("locked")
	}
	m.replaced = append(m.replaced, it)
	return nil
}

func TestDetectNeverStartsAStoppedDistro(t *testing.T) {
	w := &fakeWSL{state: wsl.StateStopped}
	m := &fakeMachine{found: []Finding{{Item: ItemStartup}}}
	r, err := (&Adopter{WSL: w, Machine: m, Distro: "d"}).Detect(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if w.scripts != 0 {
		t.Fatal("a script ran in a stopped distro, which would have started it")
	}
	if !r.Found() || r.Container != "" {
		t.Fatalf("report %+v", r)
	}
}

func TestDetectReadsTheContainerOfARunningDistro(t *testing.T) {
	w := &fakeWSL{state: wsl.StateRunning, out: "running\ncertificate\n"}
	r, err := (&Adopter{WSL: w, Machine: &fakeMachine{}, Distro: "d"}).Detect(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if r.Container != "running" || !r.Certificate || !r.Found() {
		t.Fatalf("report %+v", r)
	}
	if r.Findings == nil {
		t.Fatal("findings must be an empty list, not null, for the window")
	}
}

func TestNothingToAdoptWithoutADistroOrAContainer(t *testing.T) {
	absent, _ := (&Adopter{WSL: &fakeWSL{state: wsl.StateNotInstalled}, Machine: &fakeMachine{found: []Finding{{Item: ItemIdle}}}, Distro: "d"}).Detect(t.Context())
	if absent.Found() {
		t.Fatal("found a setup without a distro")
	}
	empty, _ := (&Adopter{WSL: &fakeWSL{state: wsl.StateRunning, out: "missing\nno-certificate\n"}, Machine: &fakeMachine{}, Distro: "d"}).Detect(t.Context())
	if empty.Found() {
		t.Fatal("found a setup with no container and no legacy part")
	}
}

func TestAdoptReplacesOnlyChosenFoundItems(t *testing.T) {
	m := &fakeMachine{found: []Finding{{Item: ItemShortcuts, Title: "s"}, {Item: ItemStartup, Title: "u"}, {Item: ItemIdle, Title: "i"}}}
	a := &Adopter{WSL: &fakeWSL{state: wsl.StateRunning}, Machine: m, Distro: "d"}
	done, err := a.Adopt(t.Context(), []Item{ItemShortcuts, ItemIdle, ItemCommand})
	if err != nil {
		t.Fatal(err)
	}
	want := []Item{ItemShortcuts, ItemIdle}
	if !slices.Equal(done, want) || !slices.Equal(m.replaced, want) {
		t.Fatalf("done %v, replaced %v, want %v (startup not chosen, command not found)", done, m.replaced, want)
	}
}

func TestAdoptContinuesPastAFailure(t *testing.T) {
	m := &fakeMachine{fail: ItemShortcuts, found: []Finding{{Item: ItemShortcuts, Title: "shortcuts"}, {Item: ItemStartup, Title: "startup"}}}
	a := &Adopter{WSL: &fakeWSL{state: wsl.StateRunning}, Machine: m, Distro: "d"}
	done, err := a.Adopt(t.Context(), []Item{ItemShortcuts, ItemStartup})
	if err == nil || !slices.Equal(done, []Item{ItemStartup}) {
		t.Fatalf("done %v, err %v", done, err)
	}
}

func TestAdoptNeedsTheDistro(t *testing.T) {
	m := &fakeMachine{found: []Finding{{Item: ItemIdle}}}
	a := &Adopter{WSL: &fakeWSL{state: wsl.StateNotInstalled}, Machine: m, Distro: "d"}
	if _, err := a.Adopt(t.Context(), []Item{ItemIdle}); err == nil || len(m.replaced) != 0 {
		t.Fatalf("err %v, replaced %v", err, m.replaced)
	}
}

func TestLegacyMatchers(t *testing.T) {
	if !legacyPortainer(`-NoProfile -File "C:\x\docker-wsl\open-portainer.ps1"`) || legacyPortainer(`-File other.ps1`) {
		t.Error("portainer launcher match")
	}
	if !legacyKeepalive(`--headless wsl.exe -d Ubuntu-24.04 --exec true`) || legacyKeepalive(`--headless wsl.exe -d Ubuntu-24.04`) || legacyKeepalive(``) {
		t.Error("keepalive match")
	}
}

func TestRemoveKnownKeepsOtherFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "docker-wsl")
	write := func(rel string) {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("open-portainer.ps1")
	write(filepath.Join("bin", "docker.cmd"))
	if err := removeKnown(dir, "open-portainer.ps1", filepath.Join("bin", "docker.cmd"), "portainer.ico"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("an emptied legacy folder stays: %v", err)
	}

	write("open-portainer.ps1")
	write("mine.txt")
	write(filepath.Join("bin", "docker.cmd"))
	if err := removeKnown(dir, "open-portainer.ps1", filepath.Join("bin", "docker.cmd")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "mine.txt")); err != nil {
		t.Fatalf("a file that is not ours was removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "bin")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the emptied bin folder stays: %v", err)
	}
}
