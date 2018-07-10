// detect windows
// bits: https://play.golang.org/p/q6YOUS58pf
// https://stackoverflow.com/questions/33790814/determining-if-current-process-runs-in-wow64-or-not-in-go
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var windows_bits_value int = -1

// only for 32-bit subsystem in windows
func os_bits() int {
	if windows_bits_value >= 0 {
		return windows_bits_value
	}

	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		panic(fmt.Sprintf("Can't get windows bits: kernel32 load error: %v", err.Error()))
	}
	defer dll.Release()
	proc, err := dll.FindProc("IsWow64Process")
	if err != nil || proc == nil {
		windows_bits_value = 32
		return windows_bits_value
	}

	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		panic(fmt.Sprintf("Can't get windows bits: GetCurrentProcess call error: %v", err.Error()))
	}

	var result bool
	proc.Call(uintptr(handle), uintptr(unsafe.Pointer(&result)))

	if result {
		windows_bits_value = 64
	} else {
		windows_bits_value = 32
	}
	return windows_bits_value
}
