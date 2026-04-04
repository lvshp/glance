//go:build !windows

package core

import (
	"os"
	"os/exec"
	"strings"
)

func buildBossCommand(command string) *exec.Cmd {
	shell := strings.TrimSpace(os.Getenv("SHELL"))
	if shell == "" {
		shell = "/bin/sh"
	}
	return exec.Command(shell, "-lc", command)
}
