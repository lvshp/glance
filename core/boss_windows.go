//go:build windows

package core

import (
	"os/exec"
	"syscall"
)

func buildBossCommand(command string) *exec.Cmd {
	cmd := exec.Command("cmd.exe", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	return cmd
}
