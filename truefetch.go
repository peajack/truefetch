// truefetch - simple fetch-alike program
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"sync"
)

const (
	BBLACK   = "\033[1;30m"
	BGRAY    = "\033[1;90m"
	BRED     = "\033[1;31m"
	BGREEN   = "\033[1;32m"
	BYELLOW  = "\033[1;33m"
	BBLUE    = "\033[1;34m"
	BMAGENTA = "\033[1;35m"
	BCYAN    = "\033[1;36m"
	BWHITE   = "\033[1;37m"
	BLACK    = "\033[0;30m"
	RED      = "\033[0;31m"
	GREEN    = "\033[0;32m"
	YELLOW   = "\033[0;33m"
	BLUE     = "\033[0;34m"
	MAGENTA  = "\033[0;35m"
	CYAN     = "\033[0;36m"
	WHITE    = "\033[0;37m"
	BITAL    = "\033[1;3m"
)

var prettyNames = map[string]string{
	"freebsd":   "FreeBSD",
	"openbsd":   "OpenBSD",
	"netbsd":    "NetBSD",
	"dragonfly": "DragonflyBSD",
	"darwin":    "macOS",
	"ios":       "iOS",
	"plan9":     "Plan9",
	"android":   "Android",
	"windows":   "Windows",
}

// OSName - container for os name
type OSName struct {
	name string
	id   string
}

// Logo - container for logo
type Logo struct {
	Col1, Col2, Col3, Col4, Col5, Col6, Col7, Col8 string
	Color                                          string
	PackageManager                                 string
}

// Result - way to grab results from goroutines
type Result struct {
	Name   string
	Result string
}

func wait(wg *sync.WaitGroup, routine func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		routine()
	}()
}

func wcL(s string) int {
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

func doesExist(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func getUser() string {
	if currentUser, err := user.Current(); err == nil {
		return currentUser.Username
	}
	return "who are you?"
}

func getOS() (osNames OSName) {
	osNames = OSName{"Unknown", "_UNKNOWN_"}
	switch runtime.GOOS {
	case "linux":
		osReleaseFile := "/etc/os-release"
		if customFile := os.Getenv("TRUEFETCH_OSRELEASE"); customFile != "" {
			osReleaseFile = customFile
		}
		file, err := os.ReadFile(osReleaseFile)
		if err != nil {
			return
		}
		text := string(file[:])
		for _, line := range strings.Split(text, "\n") {
			key, value, found := strings.Cut(line, "=")
			value = strings.Trim(value, "\"")
			if !found {
				continue
			} else if key == "NAME" {
				osNames.name = value
			} else if key == "ID" {
				if _, present := getLogo(value); present == true {
					osNames.id = value
				}
			}
		}
	default:
		osNames.id = runtime.GOOS
		osNames.name = prettyNames[runtime.GOOS]
	}
	return
}

func main() {
	osName := getOS()
	logo, _ := getLogo(osName.id)

	format := `
%[10]s %[1]s      USER%[9]s %[11]s
%[10]s %[2]s        OS%[9]s %[12]s
%[10]s %[3]s    KERNEL%[9]s %[13]s
%[10]s %[4]s    UPTIME%[9]s %[14]s
%[10]s %[5]s     SHELL%[9]s %[15]s
%[10]s %[6]s    MEMORY%[9]s %[16]s
%[10]s %[7]s      %[19]s %[9]s%[17]s
%[10]s %[8]s      %[20]s %[9]s%[18]s

`
	info := make(chan Result, 7)

	var wg sync.WaitGroup
	wait(&wg, func() { info <- Result{"user", getUser()} })
	wait(&wg, func() { info <- Result{"krnl", getKernel()} })
	wait(&wg, func() { info <- Result{"uptime", getUptime()} })
	wait(&wg, func() { info <- Result{"sh", getShell()} })
	wait(&wg, func() { info <- Result{"init", getInit()} })
	wait(&wg, func() { info <- Result{"mem", getMemory()} })
	wait(&wg, func() { info <- Result{"pkgs", getPkgs(logo.PackageManager)} })

	go func() {
		wg.Wait()
		close(info)
	}()

	results := map[string]string{}
	for result := range info {
		results[result.Name] = result.Result
	}

	var havePkgs, haveInit string
	if results["pkgs"] != "" {
		havePkgs = "PKGS"
	}
	if results["init"] != "" {
		haveInit = "INIT"
	}

	reset := RESET
	color := logo.Color
	if colors := os.Getenv("TRUEFETCH_NOCOLORS"); colors != "" {
		reset = ""
		color = ""
	}

	fmt.Printf(
		format,
		logo.Col1, logo.Col2, logo.Col3, logo.Col4,
		logo.Col5, logo.Col6, logo.Col7, logo.Col8,
		reset, color,
		results["user"], osName.name,
		results["krnl"], results["uptime"],
		results["sh"], results["mem"],
		results["pkgs"], results["init"],
		havePkgs, haveInit,
	)
}
