//go:build !unix

package main

import (
	"os"
	"os/exec"
)

func newShellCommand() (*exec.Cmd, error) {
	return nil, os.ErrInvalid
}
