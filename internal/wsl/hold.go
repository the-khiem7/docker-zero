package wsl

import (
	"errors"
	"os"
	"os/exec"
	"sync"
)

// Hold is a hidden `sleep infinity` session inside a distro. While it runs,
// WSL does not stop the distro for being idle; releasing it lets the distro
// follow the normal idle timeout again.
type Hold struct {
	mu   sync.Mutex
	cmd  *exec.Cmd
	done chan struct{}
}

// StartHold starts the session in distro. The distro boots if it was stopped.
func (c *Client) StartHold(distro string) (*Hold, error) {
	killChildrenOnExit()
	cmd := exec.Command(c.exe, "--distribution", distro, "--exec", "/bin/sleep", "infinity") //nolint:gosec // fixed executable, no shell
	cmd.Env = append(os.Environ(), "WSL_UTF8=1")
	hideWindow(cmd)
	if err := cmd.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, ErrWSLUnavailable
		}
		return nil, err
	}
	h := &Hold{cmd: cmd, done: make(chan struct{})}
	go func() {
		_ = cmd.Wait()
		close(h.done)
	}()
	return h, nil
}

// Alive reports whether the session is still running. It ends on its own
// when the distro is terminated or WSL shuts down.
func (h *Hold) Alive() bool {
	if h == nil {
		return false
	}
	select {
	case <-h.done:
		return false
	default:
		return true
	}
}

// Release ends the session and waits for wsl.exe to exit.
func (h *Hold) Release() {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.Alive() {
		_ = h.cmd.Process.Kill()
	}
	<-h.done
}
