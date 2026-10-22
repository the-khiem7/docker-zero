//go:build windows

package wsl

import (
	"os/exec"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// createNoWindow keeps wsl.exe from flashing a console window when it
	// is started by a GUI-subsystem process.
	createNoWindow = 0x08000000
	// createSuspended starts the process with its first thread paused.
	createSuspended = 0x00000004
)

func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow,
	}
}

// holdJob kills every member when its last handle closes. The handle is
// never closed by hand, so that happens when this process exits for any
// reason, including a crash or Task Manager. Without it, a hold left behind
// by a killed app keeps the distro running forever. Only holds join it: a
// browser the app opens must outlive the app.
var holdJob = sync.OnceValues(func() (windows.Handle, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	//nolint:gosec // the Win32 API takes a pointer to the struct
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil {
		_ = windows.CloseHandle(job)
		return 0, err
	}
	return job, nil
})

var ntResumeProcess = windows.NewLazySystemDLL("ntdll.dll").NewProc("NtResumeProcess")

// startInHoldJob starts cmd suspended, puts it in holdJob, then lets it run.
// wsl.exe starts a second wsl.exe; joining before the first instruction
// runs means that child inherits the job too, with no race. If the job
// cannot be used the process still runs; it only survives a crash.
func startInHoldJob(cmd *exec.Cmd) error {
	hideWindow(cmd)
	cmd.SysProcAttr.CreationFlags |= createSuspended
	if err := cmd.Start(); err != nil {
		return err
	}
	h, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE|windows.PROCESS_SUSPEND_RESUME,
		false, uint32(cmd.Process.Pid)) //nolint:gosec // Windows PIDs fit in uint32
	if err != nil {
		_ = cmd.Process.Kill()
		return err
	}
	defer func() { _ = windows.CloseHandle(h) }()
	if job, err := holdJob(); err == nil {
		_ = windows.AssignProcessToJobObject(job, h)
	}
	if r, _, _ := ntResumeProcess.Call(uintptr(h)); r != 0 {
		_ = cmd.Process.Kill()
		return windows.NTStatus(r) //nolint:gosec // NTSTATUS is 32 bits wide
	}
	return nil
}
