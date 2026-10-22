//go:build windows

package wsl

import (
	"os/exec"
	"syscall"
)

// createNoWindow keeps wsl.exe from flashing a console window when it is
// started by a GUI-subsystem process.
const createNoWindow = 0x08000000

func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow,
	}
}
