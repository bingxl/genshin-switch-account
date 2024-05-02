package backend

import (
	"os/exec"
	"syscall"
)

func RunInBack(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
