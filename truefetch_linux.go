//go:build linux

// Linux-specific calls
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

func getUptime() string {
	sysinfo := unix.Sysinfo_t{}
	err := unix.Sysinfo(&sysinfo)
	if err != nil {
		return "unknown"
	}
	uptime := time.Duration(sysinfo.Uptime * int64(time.Second))
	return fmt.Sprint(uptime)
}

func processMemory(str string) int {
	str = strings.TrimSuffix(str, "kB")
	str = strings.TrimSpace(str)
	kb, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("%v", err)
		return 0
	}
	return kb / 1024
}

func getMemory() string {
	var memTotal, memAvail int
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return "unknown"
	}
	data := make([]byte, 150)
	_, err = file.Read(data)
	if err != nil {
		return "unknown"
	}
	meminfo := string(data[:])
	for _, line := range strings.Split(meminfo, "\n") {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		} else if key == "MemTotal" {
			memTotal = processMemory(value)
		} else if key == "MemAvailable" {
			memAvail = processMemory(value)
		}
	}
	return fmt.Sprintf("%d MiB / %d MiB", memTotal-memAvail, memTotal)
}
