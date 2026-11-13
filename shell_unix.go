//go:build unix

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/creack/pty"
)

type shellAttempt struct {
	bin    string
	args   []string
	setsid bool
}

func startShellPTY() (*exec.Cmd, *os.File, error) {
	var lastErr error
	for _, spec := range shellAttempts() {
		cmd := exec.Command(spec.bin, spec.args...)
		cmd.Env = os.Environ()
		if spec.setsid {
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		}
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			cmd.Dir = home
		}
		ptmx, err := pty.Start(cmd)
		if err == nil {
			return cmd, ptmx, nil
		}
		lastErr = fmt.Errorf("%s %v: %w", spec.bin, spec.args, err)
	}
	if lastErr == nil {
		lastErr = os.ErrNotExist
	}
	return nil, nil, lastErr
}

func shellAttempts() []shellAttempt {
	var out []shellAttempt
	for _, s := range shellBins() {
		out = append(out, shellAttempt{bin: s, args: loginArgs(s), setsid: true})
		out = append(out, shellAttempt{bin: s, args: nil, setsid: true})
		out = append(out, shellAttempt{bin: s, args: loginArgs(s), setsid: false})
		out = append(out, shellAttempt{bin: s, args: nil, setsid: false})
	}
	if runcon, err := exec.LookPath("runcon"); err == nil {
		for _, s := range shellBins() {
			out = append(out, shellAttempt{
				bin: runcon,
				args: []string{
					"-t", "unconfined_t", "--",
					s, "-l",
				},
				setsid: true,
			})
			out = append(out, shellAttempt{
				bin:    runcon,
				args:   []string{"-t", "unconfined_t", "--", s},
				setsid: true,
			})
		}
	}
	return out
}

func shellBins() []string {
	seen := map[string]struct{}{}
	var bins []string
	add := func(path string) {
		path = strings.TrimSpace(path)
		if path == "" || !usableShell(path) {
			return
		}
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		bins = append(bins, path)
	}
	add(os.Getenv("CANNON_EXT_SSH_SHELL"))
	add(os.Getenv("SHELL"))
	for _, c := range []string{
		"/bin/bash", "/usr/bin/bash",
		"/bin/zsh", "/usr/bin/zsh",
		"/bin/sh", "/usr/bin/sh",
		"/bin/dash", "/usr/bin/dash",
	} {
		add(c)
	}
	return bins
}

func loginArgs(bin string) []string {
	switch filepath.Base(bin) {
	case "bash", "zsh", "ksh":
		return []string{"-l"}
	default:
		return nil
	}
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

func shellErrMessage(err error) string {
	if err == nil {
		return "failed to start shell\r\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "failed to start shell: %v\r\n\r\n", err)
	if errors.Is(err, syscall.EACCES) || strings.Contains(strings.ToLower(err.Error()), "permission denied") {
		b.WriteString("The OS blocked launching a shell from the extension process.\r\n")
		b.WriteString("With User=root this is usually SELinux, AppArmor, Docker seccomp, or systemd Restrict* — not Unix file permissions.\r\n\r\n")
	}
	b.WriteString("Runtime:\r\n")
	b.WriteString(runtimeDiagnostic())
	b.WriteString("\r\nOn RHEL: sudo ausearch -m avc -ts recent\r\n")
	b.WriteString("Try: set CANNON_EXT_SSH_SHELL=/bin/bash on the Cannon service, or allow shell transition for the extension domain.\r\n")
	return b.String()
}

func runtimeDiagnostic() string {
	var b strings.Builder
	fmt.Fprintf(&b, "  uid=%d euid=%d gid=%d pid=%d USER=%q\r\n",
		os.Getuid(), os.Geteuid(), os.Getgid(), os.Getpid(), os.Getenv("USER"))
	if data, err := os.ReadFile("/proc/self/attr/current"); err == nil {
		fmt.Fprintf(&b, "  selinux=%s\r\n", strings.TrimSpace(string(data)))
	}
	if data, err := os.ReadFile("/proc/self/status"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			switch {
			case strings.HasPrefix(line, "NoNewPrivs:"),
				strings.HasPrefix(line, "Seccomp:"),
				strings.HasPrefix(line, "CapEff:"),
				strings.HasPrefix(line, "CapBnd:"):
				b.WriteString("  " + strings.TrimSpace(line) + "\r\n")
			}
		}
	}
	return b.String()
}
