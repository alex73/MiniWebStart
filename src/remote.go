package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func needUpdateZip(remoteUrl string, localDir string) bool {
	fmt.Print("Checking for updates... ")

	resp1, err1 := http.Head(remoteUrl)
	if err1 != nil {
		panic(fmt.Sprintf("Error server response: %v", err1.Error()))
	}
	if resp1.StatusCode != 200 {
		panic(fmt.Sprintf("Error server response: %v", resp1.StatusCode))
	}
	httpLen, err2 := strconv.ParseInt(resp1.Header.Get("Content-Length"), 10, 64)
	if err2 != nil {
		panic(fmt.Sprintf("Wrong server length response: %v", err2.Error()))
	}
	fi, err3 := os.Stat(localDir + ZIP_STORE)
	var localLen int64 = -1
	if !os.IsNotExist(err3) {
		if err3 != nil {
			panic(fmt.Sprintf("Wrong local file stat: %v", err3.Error()))
		}
		localLen = fi.Size()
	}

	if localLen != httpLen {
		fmt.Println("need to update")
		return true
	} else {
		fmt.Println("no need")
		return false
	}
}

func updateZip(remoteUrl string, localDir string) {
	fmt.Print("Updating... ")
	zipPath := localDir + ZIP_STORE
	if errMkDir := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); errMkDir != nil {
		panic(fmt.Sprintf("Error create directory: %v", errMkDir.Error()))
	}
	out, err1 := os.Create(zipPath)
	if err1 != nil {
		panic(fmt.Sprintf("Error create local file: %v", err1.Error()))
	}
	defer out.Close()
	resp, err2 := http.Get(remoteUrl)
	if err2 != nil {
		panic(fmt.Sprintf("Error server response: %v", err2.Error()))
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Error server response: %v", resp.StatusCode))
	}
	defer resp.Body.Close()

	var downloaded int64 = 0
	for {
		count, err3 := io.CopyN(out, resp.Body, 4*1024*1024)
		if err3 != nil {
			if err3 == io.EOF {
				break
			}
			panic(fmt.Sprintf("Error download: %v", err3.Error()))
		}
		downloaded += count
		fmt.Printf("%v MiB ", downloaded/1024/1024)
	}
	fmt.Println("Done")
}
