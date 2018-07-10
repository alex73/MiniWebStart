// detect linux bits: 64 should have 'lm' flag in the /proc/cpuinfo

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var linux_bits_value int = -1

func os_bits() int {
	if linux_bits_value >= 0 {
		return linux_bits_value
	}

	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		panic(fmt.Sprintf("Can't get linux bits: /proc/cpuinfo read error: %v", err.Error()))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "flags") {
			p := strings.Index(line, ":")
			flags := strings.Split(line[p+1:], " ")
			for _, f := range flags {
				if f == "lm" {
					linux_bits_value = 64
					return linux_bits_value
				}
			}
			linux_bits_value = 32
			return linux_bits_value
		}
	}
	linux_bits_value = 0
	return linux_bits_value
}
