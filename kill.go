//go:build !windows

// Assumes *Nix

package main

import (
	"os"
	"os/exec"
	"syscall"
)

// ProcessGroup on *nix is just a slice of processes
type ProcessGroup struct {
	ps []*os.Process
}

// NewProcessGroup creates a process group
func NewProcessGroup() (ProcessGroup, error) {
	return ProcessGroup{make([]*os.Process, 0)}, nil
}

// Kill closes all processes attached to this group and their sub-processes
func (g *ProcessGroup) Kill() error {
	for _, p := range g.ps {
		syscall.Kill(-p.Pid, syscall.SIGKILL)
	}
	g.ps = nil
	return nil
}

// SetPgidToCmd must be called on *nix so that child processes get the same PGID as their parent
func (g *ProcessGroup) SetPgidToCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// AddProcess add processes to this group which will allow you to close their sub-processes aswell
func (g *ProcessGroup) AddProcess(p *os.Process) error {
	g.ps = append(g.ps, p)
	return nil
}
