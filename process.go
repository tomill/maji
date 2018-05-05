package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type process struct {
	*exec.Cmd
	command []string
}

func NewProcess(command []string) *process {
	return &process{command: command}
}

func (p *process) Start() error {
	p.Cmd = exec.Command(p.command[0], p.command[1:]...)
	p.Cmd.Stdout = os.Stdout
	p.Cmd.Stderr = os.Stderr
	p.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return p.Cmd.Start()
}

func (p *process) Stop() {
	if p.Cmd != nil && p.Cmd.Process != nil {
		// TODO: this works only in *nix world
		syscall.Kill(-p.Process.Pid, syscall.SIGKILL)
	}
}

func (p *process) String() string {
	s := fmt.Sprintf("`%s`", strings.Join(p.command, " "))
	if p.Cmd != nil && p.Cmd.Process != nil {
		s += fmt.Sprintf(" (pid: %d)", p.Process.Pid)
	}

	return s
}
