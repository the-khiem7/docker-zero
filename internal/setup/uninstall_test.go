package setup

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// uninstallHarness is a harness whose distro exists and whose Host hooks
// record into the fake's call list.
func uninstallHarness(t *testing.T) *harness {
	t.Helper()
	h := newHarness(t, Config{Distro: "dz-throwaway"})
	h.fake.state = wsl.StateRunning
	if err := os.MkdirAll(h.binDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(h.binDir, "docker.cmd"), []byte("@echo off\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	record := func(name string) func() error {
		return func() error { h.fake.calls = append(h.fake.calls, name); return nil }
	}
	h.w.Host = Host{
		ReleaseHold:     func() { h.fake.calls = append(h.fake.calls, "release") },
		RemoveAutoStart: record("autostart"),
		RemoveShortcuts: record("shortcuts"),
		RemoveIdleLine:  record("idle"),
		RemoveConfig:    record("config"),
		RemoveFromPath: func(dir string) error {
			h.fake.calls = append(h.fake.calls, "path:"+dir)
			return nil
		},
		UntrustCert: func(pem []byte) (bool, error) {
			h.fake.calls = append(h.fake.calls, "untrust")
			return strings.Contains(string(pem), "BEGIN CERTIFICATE"), nil
		},
	}
	return h
}

func TestUninstallRemovesEverythingInOrder(t *testing.T) {
	h := uninstallHarness(t)
	backup := filepath.Join(t.TempDir(), "backup.tar")
	if err := h.w.Uninstall(t.Context(), UninstallOptions{BackupPath: backup}); err != nil {
		t.Fatal(err)
	}
	// The backup comes first and the certificate goes before the distro, whose
	// certificate it needs to find; the hold is released before the removal.
	want := []string{"export", "untrust", "release", "unregister", "autostart", "shortcuts", "idle", "path:" + h.binDir, "config"}
	if !slices.Equal(h.fake.calls, want) {
		t.Fatalf("calls %q, want %q", h.fake.calls, want)
	}
	if !slices.Equal(h.fake.unregistered, []string{"dz-throwaway"}) {
		t.Fatalf("unregistered %q", h.fake.unregistered)
	}
	if _, err := os.Stat(backup); err != nil {
		t.Fatalf("backup missing: %v", err)
	}
	if _, err := os.Stat(h.binDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the docker command folder is still there: %v", err)
	}
	for _, s := range UninstallSteps {
		if got := h.status(s.ID); got != StatusDone {
			t.Errorf("step %s: %q, want done", s.ID, got)
		}
	}
}

func TestUninstallRunsAgainAfterTheDistroIsGone(t *testing.T) {
	h := uninstallHarness(t)
	if err := h.w.Uninstall(t.Context(), UninstallOptions{}); err != nil {
		t.Fatal(err)
	}
	h.events, h.fake.calls = nil, nil
	if err := h.w.Uninstall(t.Context(), UninstallOptions{}); err != nil {
		t.Fatalf("second run: %v", err)
	}
	for _, id := range []StepID{StepBackup, StepTrust, StepRemove} {
		if got := h.status(id); got != StatusSkipped {
			t.Errorf("step %s: %q, want skipped", id, got)
		}
	}
	if len(h.fake.unregistered) != 1 {
		t.Fatalf("unregistered %d times", len(h.fake.unregistered))
	}
}

func TestUninstallKeepsTheDistroWhenTheCertificateStays(t *testing.T) {
	h := uninstallHarness(t)
	h.w.Host.UntrustCert = func([]byte) (bool, error) { return false, errors.New("declined in the Windows dialog") }
	err := h.w.Uninstall(t.Context(), UninstallOptions{})
	if err == nil || !strings.Contains(err.Error(), "declined") {
		t.Fatalf("got %v", err)
	}
	if len(h.fake.unregistered) != 0 {
		t.Fatal("the distro was removed while its certificate is still trusted")
	}
	if got := h.status(StepWindows); got != "" {
		t.Fatalf("a later step ran: %q", got)
	}
}

func TestUninstallStopsWhenTheBackupFails(t *testing.T) {
	h := uninstallHarness(t)
	h.fake.exportErr = errors.New("disk full")
	backup := filepath.Join(t.TempDir(), "backup.tar")
	if err := h.w.Uninstall(t.Context(), UninstallOptions{BackupPath: backup}); err == nil {
		t.Fatal("no error")
	}
	if len(h.fake.unregistered) != 0 {
		t.Fatal("the distro was removed without its backup")
	}
	if _, err := os.Stat(backup); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a partial backup was left behind: %v", err)
	}
}

func TestUninstallOptionsValidate(t *testing.T) {
	dir := t.TempDir()
	taken := filepath.Join(dir, "taken.tar")
	if err := os.WriteFile(taken, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	cases := map[string]struct {
		path string
		want error
	}{
		"none":     {"", nil},
		"fine":     {filepath.Join(dir, "new.tar"), nil},
		"relative": {"backup.tar", ErrBadBackup},
		"not tar":  {filepath.Join(dir, "backup.zip"), ErrBadBackup},
		"exists":   {taken, ErrBackupExists},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := (UninstallOptions{BackupPath: c.path}).Validate(); !errors.Is(err, c.want) {
				t.Fatalf("got %v, want %v", err, c.want)
			}
		})
	}
}

func TestUninstallWindowsStepReportsEveryFailure(t *testing.T) {
	h := uninstallHarness(t)
	h.w.Host.RemoveAutoStart = func() error { return errors.New("registry locked") }
	err := h.w.Uninstall(t.Context(), UninstallOptions{})
	if err == nil || !strings.Contains(err.Error(), "auto start: registry locked") {
		t.Fatalf("got %v", err)
	}
	// The other items were still removed.
	for _, want := range []string{"shortcuts", "idle", "path:" + h.binDir} {
		if !slices.Contains(h.fake.calls, want) {
			t.Errorf("%s skipped after the auto start failure", want)
		}
	}
	if slices.Contains(h.fake.calls, "config") {
		t.Error("settings removed although a step failed")
	}
}

func TestRemoveLauncherLeavesOtherFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(dir, "mine.txt")
	for _, f := range []string{filepath.Join(dir, "docker.cmd"), filepath.Join(dir, "docker"), other} {
		if err := os.WriteFile(f, nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := removeLauncher(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(other); err != nil {
		t.Fatalf("a file that is not ours was removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "docker.cmd")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("docker.cmd stays: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "docker")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("docker stays: %v", err)
	}
	if err := removeLauncher(filepath.Join(dir, "gone")); err != nil {
		t.Fatalf("a missing folder is not an error: %v", err)
	}
}

func TestFinishUninstallSchedulesTheInstallFolder(t *testing.T) {
	h := uninstallHarness(t)
	called := false
	h.w.Host.RemoveInstallDirAfterExit = func() error { called = true; return nil }
	if err := h.w.FinishUninstall(); err != nil || !called {
		t.Fatalf("called %v, err %v", called, err)
	}
	if err := (&Wizard{}).FinishUninstall(); err != nil {
		t.Fatalf("no hook: %v", err)
	}
}
