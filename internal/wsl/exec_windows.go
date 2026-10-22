//go:build windows

package wsl

import (
	"os/exec"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
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

// killChildrenOnExit puts this process in a job object that kills every
// member when its last handle closes. The handle is never closed by hand, so
// that happens when this process exits for any reason, including a crash or
// Task Manager. Children inherit the job at creation, so wsl.exe and the
// wsl.exe it starts in turn are covered with no race. Without it, a hold
// left behind by a killed app would keep the distro running forever.
var killChildrenOnExit = sync.OnceFunc(func() {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil { //nolint:gosec // Win32 API takes a pointer to the struct

		_ = windows.CloseHandle(job)
		return
	}
	// Failure leaves holds working; they only survive a crash of this process.
	_ = windows.AssignProcessToJobObject(job, windows.CurrentProcess())
})
