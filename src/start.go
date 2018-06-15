package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const ZIP_STORE string = "/.cache/image.zip"
const WORK_DIR string = "/work/"

var args []string

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			time.Sleep(30 * time.Second)
		}
	}()
	args = os.Args[1:]

	remoteUrl := getRemoteUrl()
	if REMOTE != "" {
		if remoteUrl != "" {
			panic("--image-url parameter can't be used because predefined url is " + REMOTE)
		}
		remoteUrl = REMOTE
	} else {
		if remoteUrl == "" {
			panic("--image-url parameter should be defined because there is no predefined url")
		}
	}

	localDir := getMiniDir() + getApplicationId(remoteUrl)
	if needUpdateZip(remoteUrl, localDir) {
		updateZip(remoteUrl, localDir)
		extractZip(localDir)
	}

	fmt.Print("Starting... ")
	var cmd *exec.Cmd
	switch os_name() {
	case OS_WINDOWS:
		params := append([]string{"/c", "start.cmd"}, args...)
		cmd = exec.Command("cmd.exe", params...)
		break
	default:
		cmd = exec.Command("./start.sh", args...)
		break
	}
	cmd.Dir = localDir + WORK_DIR
	err := cmd.Start()
	if err != nil {
		panic("Error execution startup script: " + err.Error())
	}
	fmt.Println("Done")
	time.Sleep(3 * time.Second)
}

func getRemoteUrl() string {
	for i, a := range args {
		if strings.HasPrefix(a, "--image-url") {
			args = append(args[0:i], args[i+1:]...)
			return a[12:]
		}
	}
	return ""
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

func getApplicationId(url string) string {
	var r string
	if APPLICATION_ID != "" {
		if match, _ := regexp.MatchString("[A-Za-z0-9\\.\\-_]", APPLICATION_ID); !match {
			panic("APPLICATION_ID must contain only A-Z, a-z, 0-9 and '.', '-', '_' chars")
		}
		r = APPLICATION_ID
	} else {
		re1 := regexp.MustCompile("[^A-Za-z0-9\\.\\-_]")
		re2 := regexp.MustCompile("_{2,}")
		r = re1.ReplaceAllString(url, "_")
		r = re2.ReplaceAllString(r, "_")
	}
	return r
}
