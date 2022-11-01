//go:build windows

package main

import (
	"os"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ProcessGroup on windows is just a HANDLE
type ProcessGroup windows.Handle

// NewProcessGroup creates a process group
func NewProcessGroup() (ProcessGroup, error) {
	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return 0, err
	}

	return ProcessGroup(handle), nil
}

// Kill closes all processes attached to this group and their sub-processes
func (g ProcessGroup) Kill() error {
	if g != 0 {
		return windows.CloseHandle(windows.Handle(g))
	}
	return nil
}

// SetPgidToCmd is a noop on Windows
func (g ProcessGroup) SetPgidToCmd(cmd *exec.Cmd) {}

// AddProcess add processes to this group which will allow you to close their sub-processes aswell
func (g ProcessGroup) AddProcess(p *os.Process) error {
	type process struct {
		Pid    int
		Handle uintptr
	}
	if g != 0 {
		return windows.AssignProcessToJobObject(
			windows.Handle(g),
			windows.Handle((*process)(unsafe.Pointer(p)).Handle))

	}
	return nil
}
