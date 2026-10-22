package setup

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

// The uninstall steps, in run order.
const (
	StepBackup  StepID = "backup"
	StepTrust   StepID = "trust"
	StepRemove  StepID = "unregister"
	StepWindows StepID = "windows"
	StepData    StepID = "data"
)

// UninstallSteps lists the uninstall steps in run order.
var UninstallSteps = []Step{
	{StepBackup, "Back up the Ubuntu distro"},
	{StepTrust, "Remove the trusted certificate"},
	{StepRemove, "Remove the Ubuntu distro"},
	{StepWindows, "Remove auto start, shortcuts and the docker command"},
	{StepData, "Remove Docker Zero settings"},
}

// ErrBackupExists is returned when the backup file is already there; an
// export must never replace a file the user may still need.
var ErrBackupExists = errors.New("the backup file already exists; choose another name")

// ErrBadBackup is returned for a backup path that is not an absolute .tar file.
var ErrBadBackup = errors.New("the backup must be an absolute path ending in .tar")

// Host performs the Windows-side removals of Uninstall. Each hook must succeed
// when its target is already gone; a nil hook is skipped.
type Host struct {
	// ReleaseHold lets go of the keepalive session before the distro goes.
	ReleaseHold func()
	// UntrustCert removes the certificate in pem from the user's trusted
	// roots and reports whether it was there.
	UntrustCert func(pem []byte) (removed bool, err error)
	// RemoveAutoStart deletes the sign-in start entry.
	RemoveAutoStart func() error
	// RemoveShortcuts deletes the Desktop and Start Menu shortcuts.
	RemoveShortcuts func() error
	// RemoveFromPath takes a folder off the user PATH.
	RemoveFromPath func(dir string) error
	// RemoveIdleLine takes the app's instanceIdleTimeout line out of .wslconfig.
	RemoveIdleLine func() error
	// RemoveConfig deletes the settings folder.
	RemoveConfig func() error
	// RemoveInstallDirAfterExit deletes the folder that holds the running
	// executable once the process has ended.
	RemoveInstallDirAfterExit func() error
}

// UninstallOptions are the choices of one uninstall run.
type UninstallOptions struct {
	// BackupPath, when set, receives a wsl --export of the distro first.
	BackupPath string
}

// Validate checks the options without touching the machine.
func (o UninstallOptions) Validate() error {
	if o.BackupPath == "" {
		return nil
	}
	if !filepath.IsAbs(o.BackupPath) || filepath.Ext(o.BackupPath) != ".tar" {
		return ErrBadBackup
	}
	if _, err := os.Stat(o.BackupPath); err == nil {
		return ErrBackupExists
	}
	return nil
}

// Uninstall removes what the app created: the distro with Portainer and
// every container, the trusted certificate, auto start, shortcuts, the docker
// command and its PATH entry, and the settings. The folder with the running
// executable goes at FinishUninstall. A failed step stops the run before the
// next one, and the run can be repeated: every step skips what is gone.
func (w *Wizard) Uninstall(ctx context.Context, o UninstallOptions) error {
	w.Config = w.Config.withDefaults()
	if err := w.Config.validate(); err != nil {
		return err
	}
	if err := o.Validate(); err != nil {
		return err
	}
	run := map[StepID]func(context.Context) (bool, error){
		StepBackup:  func(ctx context.Context) (bool, error) { return w.stepBackup(ctx, o.BackupPath) },
		StepTrust:   w.stepTrust,
		StepRemove:  w.stepRemoveDistro,
		StepWindows: w.stepWindows,
		StepData:    w.stepData,
	}
	return w.execute(UninstallSteps, func(id StepID) (bool, error) { return run[id](ctx) })
}

// FinishUninstall schedules the removal of the install folder. Call it as the
// app quits after Uninstall succeeded.
func (w *Wizard) FinishUninstall() error {
	if w.Host.RemoveInstallDirAfterExit == nil {
		return nil
	}
	return w.Host.RemoveInstallDirAfterExit()
}

func (w *Wizard) installed(ctx context.Context) (bool, error) {
	state, err := w.WSL.DistroState(ctx, w.Config.Distro)
	return state != wsl.StateNotInstalled, err
}

func (w *Wizard) stepBackup(ctx context.Context, path string) (bool, error) {
	if path == "" {
		return true, nil
	}
	present, err := w.installed(ctx)
	if err != nil || !present {
		return true, err
	}
	w.logf("Exporting %s to %s; Docker and Portainer stop while it runs", w.Config.Distro, path)
	if err := w.WSL.Export(ctx, w.Config.Distro, path); err != nil {
		// A partial tar must not look like a usable backup.
		_ = os.Remove(path)
		return false, err
	}
	return false, nil
}

func (w *Wizard) stepTrust(ctx context.Context) (bool, error) {
	present, err := w.installed(ctx)
	if err != nil || !present || w.Host.UntrustCert == nil {
		return true, err
	}
	pem, err := settings.ReadCertPEM(ctx, w.WSL, w.Config.Distro)
	if errors.Is(err, settings.ErrNoCert) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	removed, err := w.Host.UntrustCert(pem)
	if err != nil {
		// The thumbprint is only known while the distro exists, so stop here
		// instead of removing it.
		return false, fmt.Errorf("the certificate is still trusted and the distro was kept: %w", err)
	}
	if removed {
		w.logf("Removed the Portainer certificate from the trusted roots")
	}
	return !removed, nil
}

func (w *Wizard) stepRemoveDistro(ctx context.Context) (bool, error) {
	present, err := w.installed(ctx)
	if err != nil || !present {
		return true, err
	}
	if w.Host.ReleaseHold != nil {
		w.Host.ReleaseHold()
	}
	w.logf("Removing %s with its disk, images and volumes", w.Config.Distro)
	return false, w.WSL.Unregister(ctx, w.Config.Distro)
}

// stepWindows removes every Windows-side item and reports all failures, so one
// locked item does not leave the others behind.
func (w *Wizard) stepWindows(context.Context) (bool, error) {
	h := w.Host
	var errs []error
	try := func(what string, fn func() error) {
		if fn == nil {
			return
		}
		if err := fn(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", what, err))
		}
	}
	try("auto start", h.RemoveAutoStart)
	try("shortcuts", h.RemoveShortcuts)
	try(".wslconfig", h.RemoveIdleLine)
	if dir := w.Config.BinDir; dir != "" {
		try("docker command", func() error { return removeLauncher(dir) })
		if h.RemoveFromPath != nil {
			try("user PATH", func() error { return h.RemoveFromPath(dir) })
		}
	}
	return false, errors.Join(errs...)
}

// removeLauncher deletes docker.cmd and then the folder if nothing else is in
// it; a folder with other files stays.
func removeLauncher(dir string) error {
	if err := os.Remove(filepath.Join(dir, "docker.cmd")); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	// A folder that still holds other files is not ours to delete.
	_ = os.Remove(dir)
	return nil
}

func (w *Wizard) stepData(context.Context) (bool, error) {
	if w.Host.RemoveConfig == nil {
		return true, nil
	}
	return false, w.Host.RemoveConfig()
}
