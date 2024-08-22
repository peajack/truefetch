//go:build freebsd

// FreeBSD-specific calls
package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

func getUptime() string {
	ts := unix.Timespec{}
	err := unix.ClockGettime(unix.CLOCK_MONOTONIC, &ts)
	if err != nil {
		return "unknown"
	}
	uptime := time.Duration(ts.Sec * int64(time.Second)) // yes, int64, not int32 as LSP would say
	return fmt.Sprint(uptime)
}

func getMemory() string {
	return "not implemented"
}
