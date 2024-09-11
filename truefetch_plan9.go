//go:build plan9

// Plan9 specific instructions that are used to fetch system information in [truefetch](https://github.com/peajack/truefetch)
// thanks diplomat
package main

import (
	"os"
	"path"
)

const (
	RESET = ""
)

func getUname() string {
	return os.Getenv("sysname")
}

func getShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return path.Base(shell)
	}
	return "rc"
}

func getKernel() string {
	filePath := "/dev/osversion"
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "" // Here, we shouldn't handle the error. But it SHOULD BE handled in Unix. This Function should return an error type
	}
	return string(contentBytes)
}

func getPkgs(_ string) string {
	return "None are needed"
}
