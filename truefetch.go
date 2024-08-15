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
	RESET    = "\033[0;m"
	BITAL    = "\033[1;3m"
)

var packageManagers = map[string]string{
	"unknown": "",
	"pacman":  "pacman -Qq",
	"dpkg":    "dpkg -l | tail -n+6",
	"rpm":     "[[ $(which sqlite3 2>/dev/null) && $? -ne 1 ]] && (sqlite3 /var/lib/rpm/rpmdb.sqlite \"select * from Name\") || rpm -qa",
	"portage": "qlist -IRv",
	"xbps":    "xbps-query -l",
	"apk":     "grep 'P:' /lib/apk/db/installed",
	"flatpak": "flatpak list --app",
	"snap":    "snap list",
}

type OSName struct {
	name string
	id   string
}

type Logo struct {
	col1, col2, col3, col4, col5, col6, col7, col8 string
	color                                          string
	packageManager                                 string
}

func getUser() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	} else if username := os.Getenv("USERNAME"); username != "" {
		return username
	} else {
		return "who are you?"
	}
}

func getShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return path.Base(shell)
	} else {
		return "unknown"
	}
}

func getKernel() string {
	u := unix.Utsname{}
	unix.Uname(&u)
	return string(u.Release[:])
}

func getOS() (osNames OSName) {
	osNames = OSName{"Unknown", "_UNKNOWN_"}
	file, err := os.ReadFile("/etc/os-release")
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
	return
}

func getUptime() string {
	sysinfo := unix.Sysinfo_t{}
	unix.Sysinfo(&sysinfo)
	uptime := time.Duration(sysinfo.Uptime * int64(time.Second))
	return fmt.Sprint(uptime)
}

func getLogo(id string) (Logo, bool) {
	logos := map[string]Logo{
		"arch": {
			`      /\      `,
			`     /  \     `,
			`    /\   \    `,
			`   /      \   `,
			`  /   ,,   \  `,
			` /   |  |  -\ `,
			`/_-''    ''-_\`,
			`              `,
			BCYAN,
			"pacman",
		},
		"archbang": {
			`          ____`,
			`      /\ /   /`,
			`     /  /   / `,
			`    /   / /   `,
			`   /   /_/\   `,
			`  /   __   \  `,
			` /   /_/\   \ `,
			`/_-''    ''-_\`,
			BCYAN,
			"pacman",
		},
		"arcolinux": {
			`              `,
			`      /\      `,
			`     /  \     `,
			`    / /\ \    `,
			`   / /  \ \   `,
			`  / /    \ \  `,
			` / / _____\ \ `,
			"/_/  `----.\\_\\",
			BBLUE,
			"pacman",
		},
		"opensuse-leap": {
			`  _______  `,
			`__|   __ \ `,
			`     / .\ \`,
			`     \__/ |`,
			`   _______|`,
			`   \_______`,
			`__________/`,
			`           `,
			BGREEN,
			"rpm",
		},
		"debian": {
			`  _____  `,
			` /  __ \ `,
			`|  /    |`,
			`|  \\___-`,
			`-_       `,
			`  --_    `,
			`         `,
			`         `,
			BRED,
			"dpkg",
		},
		"fedora": {
			`      _____   `,
			`     /   __)\ `,
			`     |  /  \ \`,
			`  ___|  |__/ /`,
			` / (_    _)_/ `,
			`/ /  |  |     `,
			`\ \__/  |     `,
			` \ (_____/    `,
			BBLUE,
			"rpm",
		},
		"gentoo": {
			`   _-----_   `,
			`  (       \  `,
			`  \    0   \ `,
			`   \        )`,
			`   /      _/ `,
			`  (     _-   `,
			`  \____-     `,
			`             `,
			BMAGENTA,
			"portage",
		},
		"ubuntu": {
			`           `,
			`         _ `,
			`     ---(_)`,
			` _/  ---  \`,
			`(_) |   |  `,
			`  \  --- _/`,
			`     ---(_)`,
			`           `,
			BRED,
			"dpkg",
		},
		"linuxmint": {
			` _____________ `,
			`|_            \`,
			` |  | _____  | `,
			` |  | | | |  | `,
			` |  | | | |  | `,
			` |  \_____/  | `,
			` \___________/ `,
			`               `,
			BGREEN,
			"dpkg",
		},
		"manjaro": {
			` ________  __ `,
			`|       | |  |`,
			`|   ____| |  |`,
			`|  |  __  |  |`,
			`|  | |  | |  |`,
			`|  | |  | |  |`,
			`|  | |  | |  |`,
			`|__| |__| |__|`,
			BGREEN,
			"pacman",
		},
		"artix": {
			`      /\      `,
			`     /  \     `,
			`    /''.,\    `,
			`   /     ',   `,
			`  /      ',\  `,
			` /   ,.''.  \ `,
			`/.,''     ''.\`,
			`              `,
			BCYAN,
			"pacman",
		},
		"alpine": {
			`    /\ /\    `,
			`   /  \  \   `,
			`  /    \  \  `,
			` /      \  \ `,
			`/        \  \`,
			`          \  `,
			`             `,
			`             `,
			BBLUE,
			"apk",
		},
		"void": {
			`      _____    `,
			`   _  \____ -  `,
			`  / / ____ \ \ `,
			` / / /    \ \ \`,
			` | |  VOID  | |`,
			` \ \ \____/ / /`,
			`  \ \____  /_/ `,
			`   -_____\     `,
			BGREEN,
			"xbps",
		},
		"_UNKNOWN_": {
			`     ___     `,
			` ___/   \___ `,
			`/   '---'   \`,
			`--_______--' `,
			`     / \     `,
			`    /   \    `,
			`   /     \   `,
			`             `,
			BWHITE,
			"unknown",
		},
		"linux": {
			`    .--.   `,
			`   |o_o |  `,
			`   |:_/ |  `,
			`  /   \ \  `,
			` (|     | )`,
			`/'\_   _/'\`,
			`\___)=(___/`,
			`           `,
			BWHITE,
			"unknown",
		},
	}
	logos["opensuse-tumbleweed"] = logos["opensuse-leap"]
	logo, present := logos[id]
	return logo, present
}

func wcL(s string) int {
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
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
		wg.Add(1)
		go func() {
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
		}()
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

func doesExist(command string) bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("which %s", command))
	err := cmd.Run()
	return err == nil
}

func getInit() string {
	cmd := exec.Command("sh", "-c", "ps -eo comm= | head -n 1")
	stdout, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	exe := strings.TrimSuffix(string(stdout), "\n")

	if exe == "runit" {
		return "runit"
	}
	if exe == "init" {
		if doesExist("openrc") {
			return "openrc"
		}
	}
	if exe == "launchd" {
		return "launchd"
	}
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return "systemd"
	}
	if _, err := os.Stat("/run/s6/current"); err == nil {
		return "s6"
	}

	return "unknown"
}

func main() {
	osName := getOS()
	logo, _ := getLogo(osName.id)
	pkgs := getPkgs(logo.packageManager)
	format := `
%[10]s %[1]s      USER%[9]s %[11]s
%[10]s %[2]s        OS%[9]s %[12]s
%[10]s %[3]s    KERNEL%[9]s %[13]s
%[10]s %[4]s    UPTIME%[9]s %[14]s
%[10]s %[5]s     SHELL%[9]s %[15]s
%[10]s %[6]s      PKGS%[9]s %[16]s
%[10]s %[7]s      INIT%[9]s %[17]s
%[10]s %[8]s %[9]s
    `
	fmt.Printf(format, logo.col1, logo.col2, logo.col3, logo.col4, logo.col5, logo.col6, logo.col7, logo.col8, RESET, logo.color, getUser(), osName.name, getKernel(), getUptime(), getShell(), pkgs, getInit())
	fmt.Print("\n")
}
