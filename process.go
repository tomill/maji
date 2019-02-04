package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	if len(p.command) == 0 {
		return nil
	}

	if len(p.command) > 1 {
		p.Cmd = exec.Command(p.command[0], p.command[1:]...)
	} else if runtime.GOOS == "windows" {
		p.Cmd = exec.Command("C:\\Windows\\System32\\cmd.exe", "/c", p.command[0]) // possibly works
	} else {
		p.Cmd = exec.Command("sh", "-c", p.command[0])
	}

	p.Cmd.Stdout = os.Stdout
	p.Cmd.Stderr = os.Stderr
	p.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return p.Cmd.Start()
}

func (p *process) Stop() {
	if p.Cmd != nil && p.Cmd.Process != nil && runtime.GOOS != "windows" {
		_ = syscall.Kill(-p.Process.Pid, syscall.SIGKILL)
	}
}

func (p *process) String() string {
	if len(p.command) == 0 {
		return "(noop)"
	}

	s := fmt.Sprintf("`%s`", strings.Join(p.command, " "))
	if p.Cmd != nil && p.Cmd.Process != nil {
		s += fmt.Sprintf(" (pid: %d)", p.Process.Pid)
	}

	return s
}
