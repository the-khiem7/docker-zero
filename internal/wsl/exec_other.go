//go:build !windows

package wsl

import "os/exec"

func hideWindow(*exec.Cmd) {}

func startInHoldJob(cmd *exec.Cmd) error { return cmd.Start() }
