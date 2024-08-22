//go:build netbsd

package main

import (
	"fmt"
	"time"
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
	uptime := time.Duration(int64(ts.tv_sec) * int64(time.Second))
	return fmt.Sprint(uptime)
}

func getMemory() string {
	return "not implemented"
}
