package linox

import (
	"os/exec"
	"syscall"
)

// SupportsUnprivilegedCloneNewUser returns true if
// the Linux kernel allows unprivileged users to call the clone()
// syscall with `CLONE_NEWUSER`.
// It is useful, for example to establish whether the Electron 5.0+ suid sandbox
// can be used, or if it needs to be disabled.
// cf. https://github.com/electron/electron/issues/17972
func SupportsUnprivilegedCloneNewUser() bool {
	cmd := exec.Command("/bin/true")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Cloneflags = syscall.CLONE_NEWUSER
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
