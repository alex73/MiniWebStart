package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func listRemote(remoteUrl string) fileinfo {
	fmt.Print("Checking '" + remoteUrl + "' for updates... ")

	resp, err := http.Head(remoteUrl)
	if err != nil {
		panic(fmt.Sprintf("Error server response: %v", err.Error()))
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Error server response: %v", resp.StatusCode))
	}
	httpLenStr := resp.Header.Get("Content-Length")
	var httpLen int64
	if httpLenStr == "" {
		httpLen = -1
	} else {
		httpLen, err = strconv.ParseInt(httpLenStr, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Wrong server length response: %v, error: %v", httpLenStr, err.Error()))
		}
	}
	tm := parseTime(resp.Header.Get("Last-Modified"))
	fmt.Println("Done")
	return fileinfo{size: httpLen, lastModified: tm}
}

func updateFromRemote(remoteUrl string, localPath string) {
	fmt.Print("Updating from ", remoteUrl, "... ")

	lastModified := downloadFromRemote(remoteUrl, localPath)

	if lastModified >= 0 {
		err4 := os.Chtimes(localPath, time.Now(), time.Unix(lastModified, 0))
		if err4 != nil {
			panic(fmt.Sprintf("Error save local file: %v", err4.Error()))
		}
	}
	fmt.Println("Done")
}

func downloadFromRemote(remoteUrl string, localPath string) int64 {
	if errMkDir := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); errMkDir != nil {
		panic(fmt.Sprintf("Error create directory: %v", errMkDir.Error()))
	}
	out, err1 := os.Create(localPath)
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

	var result = parseTime(resp.Header.Get("Last-Modified"))

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

  return result
}

func parseTime(t string) int64 {
	if t == "" {
		return -1
	}
	lastModified, err := time.Parse(time.RFC1123, t)
	if err != nil {
		panic(fmt.Sprintf("Wrong date from server: %v, error: %v", t, err.Error()))
	}
	return lastModified.Unix()
}
