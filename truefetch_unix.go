//go:build !plan9

// Unix specific instructions that are used to fetch system information in [truefetch](https://github.com/peajack/truefetch)
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/shirou/gopsutil/v4/process"
	"golang.org/x/sys/unix"
)

// some consts
const (
	RESET = "\033[0;m"
)

func getShellFromEnv() string {
	return path.Base(os.Getenv("SHELL"))
}

func getShell() string {
	pid := os.Getppid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return getShellFromEnv()
	}
	exe, err := proc.Exe()
	if err != nil {
		cmd, err := proc.Cmdline()
		if err != nil {
			return path.Base(cmd)
		}
		return getShellFromEnv()
	}
	return path.Base(exe)
}

func getKernel() string {
	uname := unix.Utsname{}
	err := unix.Uname(&uname)
	if err != nil {
		return "unknown"
	}
	version, _, _ := strings.Cut(string(uname.Release[:]), "-")
	return version
}

// thanks dheison, but it didnt work :(
//func getMemory() string {
//	return getSysinfoField(func(s *unix.Sysinfo_t) string {
//		totalMemory := s.Totalram
//		usedMemory := totalMemory - s.Freeram - s.Bufferram - s.Sharedram
//		return fmt.Sprintf("%dMiB / %dMiB", usedMemory/MIBIBYTE, totalMemory/MIBIBYTE)
//	})
//}

var packageManagers = map[string]string{
	"unknown": "",
	"pacman":  "pacman -Qq",
	"dpkg":    "dpkg -l | tail -n+6",
	// drop this for now, 'cause it's bash-specific
	// "rpm":     "[[ $(which sqlite3 2>/dev/null) && $? -ne 1 ]] && (sqlite3 /var/lib/rpm/rpmdb.sqlite \"select * from Name\") || rpm -qa",
	"rpm":     "rpm -qa",
	"portage": "qlist -IRv",
	"xbps":    "xbps-query -l",
	"apk":     "grep 'P:' /lib/apk/db/installed",
	"flatpak": "flatpak list --app",
	"snap":    "snap list",
	"freebsd": "pkg info",
	"openbsd": "/bin/ls -1 /var/db/pkg/",
	"pkgsrc":  "pkg_info",
	"android": "echo \"$(pm list packages --user 0 2>&1 </dev/null)\" | tr ' ' '\n'",
}

func getPkgs(packageManager string) string {
	neededManagers := map[string]string{}
	if packageManager == "" {
		neededManagers = packageManagers
	} else {
		neededManagers["flatpak"] = packageManagers["flatpak"]
		neededManagers["snap"] = packageManagers["snap"]
		neededManagers[packageManager] = packageManagers[packageManager]
	}

	packageCounts := make(chan string, len(neededManagers))

	var wg sync.WaitGroup
	for manager, command := range neededManagers {
		if command == "" {
			continue
		}
		wait(&wg, func() {
			cmd := exec.Command("/bin/sh", "-c", command)
			stdout, err := cmd.Output()

			if err != nil {
				return
			}

			packageCounts <- fmt.Sprintf(
				"%d (%s)",
				wcL(string(stdout)),
				manager,
			)
		})
	}

	go func() {
		wg.Wait()
		close(packageCounts)
	}()

	var countStrings []string

	for count := range packageCounts {
		countStrings = append(countStrings, count)
	}

	return strings.Join(countStrings, ", ")
}
