// Package setup runs the Day 1 wizard: it installs the WSL distro, Docker
// Engine and Portainer, and creates the Portainer admin account, one
// resumable step at a time. It also holds the maintenance actions that act on
// the same distro: update Portainer, reset the admin password and uninstall.
package setup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/the-khiem7/docker-zero/internal/settings"
	"github.com/the-khiem7/docker-zero/internal/target"
	"github.com/the-khiem7/docker-zero/internal/wsl"
)

const (
	// imageName is the distro image the wizard installs.
	imageName = "Ubuntu-24.04"
	// portainerImage is pinned to the long-term-support tag.
	portainerImage = "portainer/portainer-ce:lts"
	// MinPasswordLen is Portainer's minimum admin password length.
	MinPasswordLen = 12
	// binDirName is the folder under the install dir that holds docker.cmd.
	binDirName = "bin"
)

// Status is the state of a step in an Event.
type Status string

// Event statuses. StatusLog carries a line of progress text for the step.
const (
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusSkipped Status = "skipped"
	StatusFailed  Status = "failed"
	StatusLog     Status = "log"
)

// StepID names a wizard step.
type StepID string

// The wizard steps, in run order.
const (
	StepWSL       StepID = "wsl"
	StepDistro    StepID = "distro"
	StepDocker    StepID = "docker"
	StepUser      StepID = "user"
	StepCert      StepID = "certificate"
	StepPortainer StepID = "portainer"
	StepAdmin     StepID = "admin"
	StepWrapper   StepID = "wrapper"
)

// Step is a wizard step as the window lists it.
type Step struct {
	ID    StepID `json:"id"`
	Title string `json:"title"`
}

// Steps lists the steps in run order.
var Steps = []Step{
	{StepWSL, "Check WSL"},
	{StepDistro, "Install the Ubuntu distro"},
	{StepDocker, "Install Docker Engine"},
	{StepUser, "Create the Linux user"},
	{StepCert, "Create the local certificate"},
	{StepPortainer, "Start Portainer"},
	{StepAdmin, "Create the Portainer admin"},
	{StepWrapper, "Add the docker command"},
}

// Event reports progress to the window. Messages never carry the setup
// token or the password.
type Event struct {
	Step    StepID `json:"step"`
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
}

// Errors returned by Run.
var (
	ErrWSLMissing  = errors.New("WSL is not available: ask IT to enable WSL, then run setup again")
	ErrBadName     = errors.New("invalid distro, user or port setting")
	ErrCredentials = errors.New("invalid admin credentials")
	ErrNoSystemd   = errors.New("the distro did not start with systemd, which Docker Engine needs")
	ErrBusy        = errors.New("setup is already running")
)

var namePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,31}$`)

// Config selects what the wizard installs. The zero value is the fixed
// production target; tests point it at a throwaway distro.
type Config struct {
	// Distro is the name the distro is registered under.
	Distro string
	// User is the default Linux user.
	User string
	// Location is the folder for the distro disk; empty keeps the WSL default.
	Location string
	// Port is the host port Portainer is published on, on 127.0.0.1.
	Port int
	// BinDir receives docker.cmd.
	BinDir string
	// AddToPath puts BinDir on the user PATH.
	AddToPath bool
}

func (c Config) withDefaults() Config {
	if c.Distro == "" {
		c.Distro = target.Distro
	}
	if c.User == "" {
		c.User = target.User
	}
	if c.Port == 0 {
		c.Port = target.PortainerPort
	}
	return c
}

func (c Config) validate() error {
	// Distro and user names are placed into shell scripts, so they must stay
	// within a small safe alphabet.
	if !namePattern.MatchString(strings.ToLower(c.Distro)) || !namePattern.MatchString(c.User) ||
		c.Port < 1 || c.Port > 65535 {
		return ErrBadName
	}
	return nil
}

// Credentials are the Portainer admin the user typed.
type Credentials struct {
	Username string
	Password string
}

// Validate checks the credentials against Portainer's rules without
// repeating them in the error.
func (c Credentials) Validate() error {
	switch {
	case strings.TrimSpace(c.Username) == "":
		return fmt.Errorf("%w: the username is empty", ErrCredentials)
	case len(c.Password) < MinPasswordLen:
		return fmt.Errorf("%w: the password needs at least %d characters", ErrCredentials, MinPasswordLen)
	}
	return nil
}

// WSL is the part of *wsl.Client the wizard uses.
type WSL interface {
	Available(ctx context.Context) error
	DistroState(ctx context.Context, distro string) (wsl.State, error)
	Install(ctx context.Context, image, name, location string) error
	Export(ctx context.Context, distro, file string) error
	Unregister(ctx context.Context, distro string) error
	Terminate(ctx context.Context, distro string) error
	Start(ctx context.Context, distro string) error
	Script(ctx context.Context, distro, user, script string) (string, error)
}

// Wizard runs the setup steps.
type Wizard struct {
	WSL    WSL
	Config Config
	// Emit receives progress; nil discards it.
	Emit func(Event)
	// AddPath puts a folder on the user PATH.
	AddPath func(dir string) error
	// APIBase overrides the Portainer API address, for tests.
	APIBase string
	// Wait is the pause between readiness checks; ReadyTimeout bounds them.
	Wait         time.Duration
	ReadyTimeout time.Duration
	// NewToken makes the one-time setup token; tests replace it.
	NewToken func() (string, error)
	// Host performs the Windows-side removals of Uninstall.
	Host Host

	cur   StepID
	token string
	api   *api
}

// New returns a Wizard for the real runtime.
func New(w WSL, cfg Config) *Wizard {
	return &Wizard{
		WSL:          w,
		Config:       cfg.withDefaults(),
		AddPath:      AddToUserPath,
		Wait:         time.Second,
		ReadyTimeout: 120 * time.Second,
		NewToken:     newToken,
	}
}

// DefaultBinDir returns the folder under the install dir that holds
// docker.cmd, from the installed executable path.
func DefaultBinDir(installExe string) string {
	return filepath.Join(filepath.Dir(installExe), binDirName)
}

type stepFunc func(ctx context.Context, cr Credentials) (skipped bool, err error)

// Run executes every step in order and stops at the first failure. Each step
// checks what already exists, so a run after a failure continues instead of
// starting over.
func (w *Wizard) Run(ctx context.Context, cr Credentials) error {
	w.Config = w.Config.withDefaults()
	if err := w.Config.validate(); err != nil {
		return err
	}
	if err := cr.Validate(); err != nil {
		return err
	}
	w.token = ""
	w.api = nil
	run := map[StepID]stepFunc{
		StepWSL:       w.stepWSL,
		StepDistro:    w.stepDistro,
		StepDocker:    w.stepDocker,
		StepUser:      w.stepUser,
		StepCert:      w.stepCert,
		StepPortainer: w.stepPortainer,
		StepAdmin:     w.stepAdmin,
		StepWrapper:   w.stepWrapper,
	}
	return w.execute(Steps, func(id StepID) (bool, error) { return run[id](ctx, cr) })
}

// execute runs steps in order, reporting each one, and stops at the first
// failure.
func (w *Wizard) execute(steps []Step, run func(StepID) (skipped bool, err error)) error {
	for _, s := range steps {
		w.cur = s.ID
		w.emit(StatusRunning, "")
		skipped, err := run(s.ID)
		if err != nil {
			w.emit(StatusFailed, err.Error())
			return fmt.Errorf("%s: %w", s.Title, err)
		}
		if skipped {
			w.emit(StatusSkipped, "")
		} else {
			w.emit(StatusDone, "")
		}
	}
	return nil
}

func (w *Wizard) emit(st Status, msg string) {
	if w.Emit != nil {
		w.Emit(Event{Step: w.cur, Status: st, Message: msg})
	}
}

func (w *Wizard) logf(format string, args ...any) {
	w.emit(StatusLog, fmt.Sprintf(format, args...))
}

// root runs a script in the distro as root and logs its output.
func (w *Wizard) root(ctx context.Context, script string) (string, error) {
	out, err := w.WSL.Script(ctx, w.Config.Distro, "root", script)
	if text := strings.TrimSpace(out); text != "" && text != markerPresent {
		for line := range strings.SplitSeq(text, "\n") {
			w.logf("%s", strings.TrimRight(line, "\r"))
		}
	}
	return out, err
}

func (w *Wizard) stepWSL(ctx context.Context, _ Credentials) (bool, error) {
	if err := w.WSL.Available(ctx); err != nil {
		return false, fmt.Errorf("%w (%w)", ErrWSLMissing, err)
	}
	return false, nil
}

func (w *Wizard) stepDistro(ctx context.Context, _ Credentials) (bool, error) {
	d := w.Config.Distro
	state, err := w.WSL.DistroState(ctx, d)
	if err != nil {
		return false, err
	}
	existed := state != wsl.StateNotInstalled
	if !existed {
		w.logf("Downloading and registering %s as %s", imageName, d)
		if err := w.WSL.Install(ctx, imageName, d, w.Config.Location); err != nil {
			return false, err
		}
	}
	out, err := w.root(ctx, "ps -p 1 -o comm=\n")
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(out) == "systemd" {
		return existed, nil
	}
	// Docker Engine runs as a systemd service, so enable systemd and restart.
	w.logf("Enabling systemd")
	if _, err := w.root(ctx, enableSystemdScript); err != nil {
		return false, err
	}
	if err := w.restart(ctx); err != nil {
		return false, err
	}
	out, err = w.root(ctx, "ps -p 1 -o comm=\n")
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(out) != "systemd" {
		return false, ErrNoSystemd
	}
	return false, nil
}

// restart stops the distro and boots it again so wsl.conf takes effect.
func (w *Wizard) restart(ctx context.Context) error {
	if err := w.WSL.Terminate(ctx, w.Config.Distro); err != nil {
		return err
	}
	return w.WSL.Start(ctx, w.Config.Distro)
}

func (w *Wizard) stepDocker(ctx context.Context, _ Credentials) (bool, error) {
	out, err := w.root(ctx, installDockerScript)
	if err != nil {
		return false, err
	}
	present := strings.TrimSpace(out) == markerPresent
	if _, err := w.root(ctx, dnsScript); err != nil {
		return false, err
	}
	return present, nil
}

func (w *Wizard) stepUser(ctx context.Context, _ Credentials) (bool, error) {
	out, err := w.root(ctx, userScript(w.Config.User))
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(out) == markerPresent {
		return true, nil
	}
	// The default user in wsl.conf applies from the next boot.
	if err := w.restart(ctx); err != nil {
		return false, err
	}
	check := "n=0; until docker info >/dev/null 2>&1 || [ $n -ge 30 ]; do sleep 1; n=$((n+1)); done\ndocker info >/dev/null\n"
	if _, err := w.WSL.Script(ctx, w.Config.Distro, w.Config.User, check); err != nil {
		return false, fmt.Errorf("user %s cannot use Docker: %w", w.Config.User, err)
	}
	w.logf("User %s can use Docker", w.Config.User)
	return false, nil
}

func (w *Wizard) stepCert(ctx context.Context, _ Credentials) (bool, error) {
	return false, settings.Generate(ctx, w.WSL, w.Config.Distro)
}

func (w *Wizard) stepWrapper(_ context.Context, _ Credentials) (bool, error) {
	dir := w.Config.BinDir
	if dir == "" {
		return true, nil
	}
	if err := WriteLauncher(dir, w.Config.Distro); err != nil {
		return false, err
	}
	w.logf("Wrote %s", filepath.Join(dir, "docker.cmd"))
	if w.Config.AddToPath && w.AddPath != nil {
		if err := w.AddPath(dir); err != nil {
			return false, fmt.Errorf("add %s to the user PATH: %w", dir, err)
		}
		w.logf("Added it to the user PATH; open a new terminal to use docker")
	}
	return false, nil
}

// WriteLauncher writes docker.cmd (cmd, PowerShell) and docker (Git Bash) into
// dir, creating the folder. The wrappers forward docker to the engine in distro.
func WriteLauncher(dir, distro string) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "docker"), []byte(dockerSh(distro)), 0o644); err != nil { //nolint:gosec // a launcher script, readable by the user's shell
		return err
	}
	return os.WriteFile(filepath.Join(dir, "docker.cmd"), []byte(dockerCmd(distro)), 0o644) //nolint:gosec // a launcher script, readable by the user's shell
}
