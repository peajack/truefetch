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
	"time"

	"golang.org/x/sys/unix"
)

const (
	RESET = "\033[0;m"
)

func getShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return path.Base(shell)
	}
	return "" // TODO: Try to get the parent process of truefetch, which should be the shell. ($0). Only applies to this file and truefetch_windows.go if implemented in the future
}

func getUname() string {
	return getUtsnameField(func(u *unix.Utsname) []byte {
		return u.Sysname[:]
	})
}

func getKernel() string {
	return getUtsnameField(func(u *unix.Utsname) []byte {
		return u.Release[:]
	})
}

func getUptime() string {
	sysinfo := unix.Sysinfo_t{}
	unix.Sysinfo(&sysinfo)
	uptime := time.Duration(sysinfo.Uptime * int64(time.Second))
	return fmt.Sprint(uptime)
}

func getUtsnameField(fieldFunc func(*unix.Utsname) []byte) string {
	u := unix.Utsname{}
	if err := unix.Uname(&u); err != nil {
		return ""
	}
	result := fieldFunc(&u)
	return strings.TrimRight(string(result), "\x00")
}

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
	"freebsd": "pkg info | wc -l | tr -d ' '",
	"openbsd": "/bin/ls -1 /var/db/pkg/ | wc -l | tr -d ' '",
	"plan9":   "#nope#",
}

func getPkgs(packageManager string) string {
	neededManagers := map[string]string{}
	if packageManager == "" {
		neededManagers = packageManagers
	} else if packageManagers[packageManager] == "#nope#" {
		return ""
	} else {
		neededManagers["flatpak"] = packageManagers["flatpak"]
		neededManagers["snap"] = packageManagers["snap"]
		neededManagers[packageManager] = packageManagers[packageManager]
	}

	packageCounts := make(chan string, len(neededManagers))

	var wg sync.WaitGroup
	for manager, command := range neededManagers {
		wg.Add(1)
		go func(manager, command string) {
			defer wg.Done()
			cmd := exec.Command("sh", "-c", command)
			stdout, err := cmd.Output()

			if err != nil {
				return
			}

			count := wcL(string(stdout))

			packageCounts <- fmt.Sprintf(
				"%d (%s)",
				count,
				manager,
			)
		}(manager, command)
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
