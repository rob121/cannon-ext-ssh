package main

import (
	"fmt"
	"os"
)

func maybeRunDiag(args []string) bool {
	for _, a := range args {
		if a != "--test-shell" && a != "--test-pty" {
			continue
		}
		fmt.Print(runtimeDiagnostic())
		_, ptmx, err := startShellPTY()
		if err != nil {
			fmt.Println(shellErrMessage(err))
			os.Exit(1)
		}
		fmt.Println("shell started ok (PTY open)")
		_ = ptmx.Close()
		os.Exit(0)
	}
	return false
}
