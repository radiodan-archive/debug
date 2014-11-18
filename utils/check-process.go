package utils

import (
	"os"
	"syscall"
)

func CheckProcess(pid int64) (running bool) {
	running = false

	// this isn't the pid we're looking for
	if pid == 0 {
		return
	}

	// see if that process exists
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return
	}

	// see if the process responds to sig 0
	err = process.Signal(syscall.Signal(0))
	running = (err == nil)

	return
}
