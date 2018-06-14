package main

// detect windows bits: https://play.golang.org/p/q6YOUS58pf    https://stackoverflow.com/questions/33790814/determining-if-current-process-runs-in-wow64-or-not-in-go
// detect linux bits: 64 should have 'fm' flag in the /proc/cpuinfo

import (
	"runtime"
)

const (
	OS_UNKNOWN = iota
	OS_WINDOWS = iota
	OS_LINUX   = iota
	OS_MACOS   = iota
)

func os_name() int {
	switch runtime.GOOS {
	case "windows":
		return OS_WINDOWS
	case "linux":
		return OS_LINUX
	case "darwin":
		return OS_MACOS
	default:
		return OS_UNKNOWN
	}
}

func os_bits() int {
	switch os_name() {
	case OS_WINDOWS:
		return 1
	case OS_LINUX:
		return 1
	case OS_MACOS:
		return 64
	default:
		return 0
	}
}
