//go:build plan9

// Plan9 specific instructions that are used to fetch system information in [truefetch](https://github.com/peajack/truefetch)
// thanks diplomat
package main

import (
	"os"
	"os/exec"
	"strings"
)

const (
	RESET = ""
)

func getShell() string {
	return "rc"
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

func getMemory() string {
	return "N/A"
}

func getKernel() string {
	filePath := "/dev/osversion"
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "" // Here, we shouldn't handle the error. But it SHOULD BE handled in Unix. This Function should return an error type
	}
	return string(contentBytes)
}

func getInit() string {
	return ""
}

func getPkgs(_ string) string {
	return "none are needed"
}
