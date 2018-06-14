package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			time.Sleep(30 * time.Second)
		}
	}()
	localDir := getMiniDir() + APPLICATION_ID
	if needUpdateZip(REMOTE, localDir) {
		updateZip(REMOTE, localDir)
		extractZip(localDir)
	}

	fmt.Print("Starting... ")
	var cmd *exec.Cmd
	switch os_name() {
	case OS_WINDOWS:
		params := []string{"/c", "start.cmd"}
		params = append(params, (os.Args[1:])...)
		cmd = exec.Command("cmd.exe", params...)
		break
	default:
		cmd = exec.Command("start.sh", os.Args[1:]...)
		break
	}
	cmd.Dir = localDir + WORK_DIR
	err := cmd.Start()
	if err != nil {
		panic("Error execution startup script")
	}
	fmt.Println("Done")
	time.Sleep(3 * time.Second)
}

/*
 * Get application working directory, depends on OS.
 */
func getMiniDir() string {
	var localDir = ""
	switch os_name() {
	case OS_WINDOWS:
		localDir = os.Getenv("APPDATA")
		if localDir == "" {
			panic("Environment variable APPDATA is not defined")
		}
		if _, err := os.Stat(localDir); os.IsNotExist(err) {
			panic("Directory " + localDir + " is not exist")
		}
		break
	case OS_LINUX:
		localDir = os.Getenv("HOME")
		if localDir == "" {
			panic("Environment variable HOME is not defined")
		}
		localDir += "/.local/share"
		if _, err := os.Stat(localDir); os.IsNotExist(err) {
			panic("Directory " + localDir + " is not exist")
		}
		break
	case OS_MACOS:
		localDir = os.Getenv("HOME")
		if localDir == "" {
			panic("Environment variable HOME is not defined")
		}
		localDir += "/Library/Application Support"
		if _, err := os.Stat(localDir); os.IsNotExist(err) {
			panic("Directory " + localDir + " is not exist")
		}
	default:
		panic("Unknown OS")
	}

	return localDir + "/MiniWebStart/"
}
