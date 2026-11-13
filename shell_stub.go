//go:build !unix

package main

import (
	"os"
	"os/exec"
)

func startShellPTY() (*exec.Cmd, *os.File, error) {
	return nil, nil, os.ErrInvalid
}

func shellErrMessage(err error) string {
	return "failed to start shell\r\n"
}
