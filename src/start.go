package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const CACHE_DIR string = "/cache/"
const WORK_DIR string = "/soft/"

var args []string
var desc MWSXML
var remoteUrl string
var remoteBase string
var localDir string

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			time.Sleep(30 * time.Second)
		}
	}()

	args = os.Args[1:]

	remoteUrl = getRemoteUrl()
	if REMOTE != "" {
		if remoteUrl != "" {
			panic("--remote parameter can't be used because predefined url is " + REMOTE)
		}
		remoteUrl = REMOTE
	} else {
		if remoteUrl == "" {
			panic("--remote parameter should be defined because there is no predefined url")
		}
	}
	remoteBase = remoteUrl[0 : strings.LastIndex(remoteUrl, "/")+1]

	localDir = getMiniDir() + getApplicationId(remoteUrl)

	updateFromRemote(remoteUrl, localDir+"/description.xml")
	desc = parseXml(localDir + "/description.xml")
	refreshCacheFiles()
	refreshWorkingFiles()

	fmt.Print("Starting... ", localDir+WORK_DIR)
	startup := getStartupFile(desc)

	var cmd *exec.Cmd
	cmd = exec.Command(startup, args...)
	cmd.Dir = localDir + WORK_DIR
	fmt.Print("Starting ",startup,"...")
	err := cmd.Start()
	if err != nil {
		panic("Error execution startup script : " + err.Error())
	}
	time.Sleep(3 * time.Second)
}

/**
 * It downloads files from remote server to local cache.
 * Exist files with the same size and modified date will be skipped.
 * Files in cache, that not exist on remote server will be removed.
 */
func refreshCacheFiles() {
	localFiles := listDir(localDir + CACHE_DIR)
	for _, res := range getResources(desc) {
		for _, f := range res.File {
			mustBeRelative(f.Href)
			remoteInfo := listRemote(remoteBase + f.Href)
			var needUpdate bool
			if localInfo, ok := localFiles[f.Href]; ok {
				delete(localFiles, f.Href)
				needUpdate = remoteInfo.size != localInfo.size || remoteInfo.lastModified != localInfo.lastModified
			} else {
				needUpdate = true
			}
			if needUpdate {
				updateFromRemote(remoteBase+f.Href, localDir+CACHE_DIR+f.Href)
			}
		}
		for _, u := range res.Unpack {
			mustBeRelative(u.Href)
			remoteInfo := listRemote(remoteBase + u.Href)
			var needUpdate bool
			if localInfo, ok := localFiles[u.Href]; ok {
				delete(localFiles, u.Href)
				needUpdate = remoteInfo.size != localInfo.size || remoteInfo.lastModified != localInfo.lastModified
			} else {
				needUpdate = true
			}
			if needUpdate {
				updateFromRemote(remoteBase+u.Href, localDir+CACHE_DIR+u.Href)
			}
		}
	}
	for f, _ := range localFiles {
		os.Remove(localDir + CACHE_DIR + f)
	}
	removeEmptyDirs(localDir + CACHE_DIR)
}

/**
 * It unpacks files from local cache into working area.
 * Exist files with the same size and modified date will be skipped.
 * Files in working area, that not exist in cache will be removed.
 */
func refreshWorkingFiles() {
	localFiles := listDir(localDir + WORK_DIR)
	for _, res := range getResources(desc) {
		for _, f := range res.File {
			mustBeRelative(f.Href)
			mustBeRelative(f.ToFile)
			cacheInfo := listFile(localDir + CACHE_DIR + f.Href)
			var needUpdate bool
			var localPath string
			if f.ToFile == "" {
				localPath = f.Href
			} else {
				localPath = f.ToFile
			}
			var oldMode uint32 = 0
			if localInfo, ok := localFiles[localPath]; ok {
				delete(localFiles, localPath)
				needUpdate = cacheInfo.size != localInfo.size || cacheInfo.lastModified != localInfo.lastModified
				oldMode = localInfo.mode
			} else {
				needUpdate = true
			}
			if needUpdate {
				localCopy(localDir+CACHE_DIR+f.Href, localDir+WORK_DIR+localPath, f.Mode)
			} else {
				localPermStr(localDir+WORK_DIR+localPath, oldMode, f.Mode)
			}
		}
		for _, u := range res.Unpack {
			mustBeRelative(u.Href)
			mustBeRelative(u.ToDir)
			toDir := strings.Replace(strings.TrimSpace(u.ToDir), "\\", "/", -1)
			if toDir != "" && !strings.HasSuffix(toDir, "/") {
				toDir += "/"
			}
			updateFromZip(localDir+CACHE_DIR+u.Href, localDir+WORK_DIR, toDir, localFiles, u.UseModes)
		}
	}
	for f, _ := range localFiles {
		os.Remove(localDir + WORK_DIR + f)
	}
	removeEmptyDirs(localDir + WORK_DIR)
}

/**
 * Retrieve remote url from command line.
 */
func getRemoteUrl() string {
	for i, a := range args {
		if strings.HasPrefix(a, "--remote=") {
			args = append(args[0:i], args[i+1:]...)
			return a[9:]
		}
	}
	return ""
}

/**
 * Get application working directory, depends on OS.
 */
func getMiniDir() string {
	var localDir = ""
	switch os_name() {
	case "windows":
		localDir = os.Getenv("APPDATA")
		if localDir == "" {
			panic("Environment variable APPDATA is not defined")
		}
		if _, err := os.Stat(localDir); os.IsNotExist(err) {
			panic("Directory " + localDir + " is not exist")
		}
		break
	case "linux":
		localDir = os.Getenv("HOME")
		if localDir == "" {
			panic("Environment variable HOME is not defined")
		}
		localDir += "/.local/share"
		if _, err := os.Stat(localDir); os.IsNotExist(err) {
			panic("Directory " + localDir + " is not exist")
		}
		break
	case "darwin":
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

/**
 * Get application id from predefined constants or constructs from remote url.
 */
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
