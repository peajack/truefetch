//go:build plan9

// Unix specific instructions that are used to fetch system information in [truefetch](https://github.com/peajack/truefetch)
package main

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	RESET = ""
)

func getUname() string {
	return os.Getenv("sysname")
}

func getUptime() string {
	cmd := exec.Command("uptime")
	stdout, err := cmd.Output()
	if err != nil {
		return ""
	}
	uptime := strings.ReplaceAll(string(stdout), "\n", "")
	fields := strings.Fields(uptime)
	return strings.Join(fields[2:], " ")
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
