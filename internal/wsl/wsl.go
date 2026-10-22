// Package wsl runs wsl.exe for Docker Zero: listing distros, starting and
// stopping them, and running bash scripts inside them.
package wsl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// ErrWSLUnavailable is returned when wsl.exe is missing or WSL is not enabled.
var ErrWSLUnavailable = errors.New("WSL is not available")

// Result is the outcome of one process run.
type Result struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

// Runner starts a process and waits for it. Tests replace it.
type Runner func(ctx context.Context, name string, args []string, stdin io.Reader) (Result, error)

// Client runs wsl.exe.
type Client struct {
	exe string
	run Runner
}

// New returns a Client that runs the real wsl.exe.
func New() *Client {
	return &Client{exe: "wsl.exe", run: execRunner}
}

// NewWithRunner returns a Client backed by run, for tests.
func NewWithRunner(run Runner) *Client {
	return &Client{exe: "wsl.exe", run: run}
}

func execRunner(ctx context.Context, name string, args []string, stdin io.Reader) (Result, error) {
	// name is always wsl.exe and args are built by this package, never by a
	// shell, so no argument is interpreted twice.
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // fixed executable, no shell
	// WSL_UTF8=1 asks wsl.exe for UTF-8 instead of UTF-16LE; decodeOutput
	// still accepts both for builds that ignore it.
	cmd.Env = append(os.Environ(), "WSL_UTF8=1")
	cmd.Stdin = stdin
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	hideWindow(cmd)
	err := cmd.Run()
	res := Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		res.ExitCode = exitErr.ExitCode()
		return res, nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return res, ErrWSLUnavailable
	}
	return res, err
}

// CommandError reports a wsl.exe run that exited non-zero.
type CommandError struct {
	Args     []string
	ExitCode int
	Message  string
}

func (e *CommandError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = "no output"
	}
	return fmt.Sprintf("wsl.exe %s: exit code %d: %s", strings.Join(e.Args, " "), e.ExitCode, msg)
}

// wsl runs wsl.exe with args and returns its decoded stdout.
func (c *Client) wsl(ctx context.Context, stdin io.Reader, args ...string) (string, error) {
	res, err := c.run(ctx, c.exe, args, stdin)
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		msg := strings.TrimSpace(decodeOutput(res.Stderr))
		if msg == "" {
			msg = strings.TrimSpace(decodeOutput(res.Stdout))
		}
		return "", &CommandError{Args: args, ExitCode: res.ExitCode, Message: msg}
	}
	return decodeOutput(res.Stdout), nil
}

// Available reports whether WSL is installed and enabled for this user.
func (c *Client) Available(ctx context.Context) error {
	if _, err := c.wsl(ctx, nil, "--status"); err != nil {
		return fmt.Errorf("%w: %w", ErrWSLUnavailable, err)
	}
	return nil
}

// Distros returns the names of all installed distros.
func (c *Client) Distros(ctx context.Context) ([]string, error) {
	out, err := c.wsl(ctx, nil, "--list", "--quiet")
	if err != nil {
		return nil, err
	}
	return parseNames(out), nil
}

// RunningDistros returns the names of the running distros.
func (c *Client) RunningDistros(ctx context.Context) ([]string, error) {
	out, err := c.wsl(ctx, nil, "--list", "--quiet", "--running")
	if err != nil {
		// With no distro running, some WSL builds exit non-zero.
		if _, ok := errors.AsType[*CommandError](err); ok {
			return nil, nil
		}
		return nil, err
	}
	return parseNames(out), nil
}

// DistroState returns the state of distro without starting it.
func (c *Client) DistroState(ctx context.Context, distro string) (State, error) {
	all, err := c.Distros(ctx)
	if err != nil {
		return "", err
	}
	if !containsFold(all, distro) {
		return StateNotInstalled, nil
	}
	running, err := c.RunningDistros(ctx)
	if err != nil {
		return "", err
	}
	if containsFold(running, distro) {
		return StateRunning, nil
	}
	return StateStopped, nil
}

// Start boots distro and returns once it is running. The distro stays up
// only while something keeps it busy or until the WSL idle timeout.
func (c *Client) Start(ctx context.Context, distro string) error {
	_, err := c.wsl(ctx, nil, "--distribution", distro, "--exec", "/bin/true")
	return err
}

// Terminate stops one distro and every process in it.
func (c *Client) Terminate(ctx context.Context, distro string) error {
	_, err := c.wsl(ctx, nil, "--terminate", distro)
	return err
}

// Shutdown stops every distro of the account and the WSL2 VM.
func (c *Client) Shutdown(ctx context.Context) error {
	_, err := c.wsl(ctx, nil, "--shutdown")
	return err
}

// Script runs a bash script inside distro as user. The script goes on stdin
// so no PowerShell or wsl.exe quoting can truncate it. Stdout is returned;
// a non-zero exit becomes a *CommandError carrying stderr.
func (c *Client) Script(ctx context.Context, distro, user, script string) (string, error) {
	// A CRLF script makes bash read "docker\r" as the command name.
	script = strings.ReplaceAll(script, "\r\n", "\n")
	return c.wsl(ctx, strings.NewReader(script),
		"--distribution", distro, "--user", user, "--exec", "/bin/bash", "-s")
}

func containsFold(list []string, s string) bool {
	return slices.ContainsFunc(list, func(v string) bool { return strings.EqualFold(v, s) })
}

// Install downloads the image distribution and registers it as name. The
// image keeps its own name when name is empty or equal to it. location, when
// set, is the folder that receives the virtual disk. The call runs without
// launching the distro, so no first-run user prompt appears.
func (c *Client) Install(ctx context.Context, image, name, location string) error {
	args := []string{"--install", "--distribution", image, "--web-download", "--no-launch"}
	if name != "" && !strings.EqualFold(name, image) {
		args = append(args, "--name", name)
	}
	if location != "" {
		args = append(args, "--location", location)
	}
	_, err := c.wsl(ctx, nil, args...)
	return err
}

// Unregister deletes distro with its disk, every image and every volume.
func (c *Client) Unregister(ctx context.Context, distro string) error {
	_, err := c.wsl(ctx, nil, "--unregister", distro)
	return err
}
