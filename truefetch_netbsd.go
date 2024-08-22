//go:build !openbsd && !freebsd && netbsd

package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

/*
#include <time.h>
struct timespec get_ts(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return ts;
}
*/
import "C"

func getUptime() string {
	ts := C.get_ts()
	if err != nil {
		return "unknown"
	}
	uptime := time.Duration(ts.Sec * int64(time.Second))
	return fmt.Sprint(uptime)
}

func getMemory() string {
	return "not implemented"
}
