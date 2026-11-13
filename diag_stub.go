//go:build !unix

package main

func maybeRunDiag(args []string) bool {
	return false
}
