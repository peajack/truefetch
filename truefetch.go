// truefetch - simple fetch-alike program
package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

// ansi colors
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
	col1, col2, col3, col4, col5, col6, col7, col8 string
	color                                          string
	packageManager                                 string
}

type Result struct {
	name   string
	result string
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

func getMemory() string {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "0 MB/0 MB (∞ %)"
	}
	return fmt.Sprintf("%v MB/%v MB (%v%%)", v.Used/1024/1024, v.Total/1024/1024, math.Round(v.UsedPercent))
}

func getUptime() string {
	u, err := host.Uptime()
	if err != nil {
		return "∞ "
	}
	return fmt.Sprint(time.Duration(u * uint64(time.Second)))
}

func getInit() string {
	proc, err := process.NewProcess(1)
	if err != nil {
		return "unknown"
	}
	cmdline, err := proc.Cmdline()
	if err != nil {
		return "unknown"
	}
	exe := path.Base(cmdline)

	if exe == "runit" {
		return "runit"
	} else if exe == "launchd" {
		return "launchd"
	} else if _, err := os.Stat("/run/systemd/system"); err == nil {
		return "systemd"
	} else if _, err := os.Stat("/run/s6/current"); err == nil {
		return "s6"
	} else if exe == "init" {
		if _, err := os.Stat("/etc/init.d"); err == nil {
			if doesExist("openrc") {
				return "openrc"
			}
			return "SysV-style"
		} else if _, err := os.Stat("/etc/rc.d"); err == nil {
			return "BSD-style rc.d"
		}
	}

	return "unknown"
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
%[10]s %[6]s      INIT%[9]s %[16]s
%[10]s %[7]s      PKGS%[9]s %[17]s
%[10]s %[8]s    MEMORY%[9]s %[18]s

`
	info := make(chan Result, 7)

	var wg sync.WaitGroup
	wait(&wg, func() { info <- Result{"user", getUser()} })
	wait(&wg, func() { info <- Result{"krnl", getKernel()} })
	wait(&wg, func() { info <- Result{"uptime", getUptime()} })
	wait(&wg, func() { info <- Result{"sh", getShell()} })
	wait(&wg, func() { info <- Result{"init", getInit()} })
	wait(&wg, func() { info <- Result{"mem", getMemory()} })
	wait(&wg, func() { info <- Result{"pkgs", getPkgs(logo.packageManager)} })

	go func() {
		wg.Wait()
		close(info)
	}()

	results := map[string]string{}
	for result := range info {
		results[result.name] = result.result
	}

	fmt.Printf(
		format,
		logo.col1, logo.col2, logo.col3, logo.col4,
		logo.col5, logo.col6, logo.col7, logo.col8,
		RESET, logo.color,
		results["user"], osName.name,
		results["krnl"], results["uptime"],
		results["sh"], results["init"],
		results["pkgs"], results["mem"],
	)
}
