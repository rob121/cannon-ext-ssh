//go:build unix

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func newShellCommand() (*exec.Cmd, error) {
	bin, args := resolveShell()
	if bin == "" {
		return nil, os.ErrNotExist
	}
	cmd := exec.Command(bin, args...)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		cmd.Dir = home
	}
	return cmd, nil
}

func resolveShell() (string, []string) {
	var candidates []string
	if s := strings.TrimSpace(os.Getenv("SHELL")); usableShell(s) {
		candidates = append(candidates, s)
	}
	candidates = append(candidates, "/bin/bash", "/bin/zsh", "/bin/sh")
	for _, c := range candidates {
		if !usableShell(c) {
			continue
		}
		switch filepath.Base(c) {
		case "bash", "zsh":
			return c, []string{"-l"}
		default:
			return c, nil
		}
	}
	return "", nil
}

func usableShell(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	lower := strings.ToLower(path)
	if strings.Contains(lower, "nologin") || strings.HasSuffix(lower, "/false") {
		return false
	}
	st, err := os.Stat(path)
	return err == nil && st.Mode()&0111 != 0
}
