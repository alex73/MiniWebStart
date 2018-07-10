package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func listDir(dirPath string) map[string]fileinfo {
	result := make(map[string]fileinfo)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			result = nil
			return nil
		}
		if err != nil {
			panic(fmt.Sprintf("Error read directory: %v", err.Error()))
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(dirPath, path)
			if err != nil {
				panic(fmt.Sprintf("Error read directory: %v", err.Error()))
			}
			rel = filepath.ToSlash(rel)
			result[rel] = fileinfo{size: info.Size(), lastModified: info.ModTime().Unix(), mode: uint32(info.Mode())}
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("Error read directory: %v", err.Error()))
	}

	return result
}

func listFile(filePath string) fileinfo {
	info, err := os.Stat(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error read file info: %v", err.Error()))
	}

	return fileinfo{size: info.Size(), lastModified: info.ModTime().Unix(), mode: uint32(info.Mode())}
}

func localPermStr(pathTo string, existMode uint32, newMode string) {
	if newMode == "" {
		newMode = "0644"
	}
	modeOct, err := strconv.ParseUint(newMode, 8, 32)
	if err != nil {
		panic(fmt.Sprintf("Wrong mode for file: %v, error: %v", newMode, err.Error()))
	}
	localPermMode(pathTo, existMode, os.FileMode(modeOct))
}

func localPermMode(pathTo string, existMode uint32, modeOct os.FileMode) {
	if uint32(modeOct) != existMode {
		err := os.Chmod(pathTo, os.FileMode(modeOct))
		if err != nil {
			panic(fmt.Sprintf("Error set file mode %v: %v", pathTo, err.Error()))
		}
	}
}

func localCopy(pathFrom string, pathTo string, mode string) {
	from, err := os.Open(pathFrom)
	if err != nil {
		panic(fmt.Sprintf("Error copy file %v: %v", pathFrom, err.Error()))
	}
	defer from.Close()

	if err = os.MkdirAll(filepath.Dir(pathTo), os.ModePerm); err != nil {
		panic(fmt.Sprintf("Error copy file %v: %v", pathFrom, err.Error()))
	}

	if mode == "" {
		mode = "0644"
	}
	modeOct, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		panic(fmt.Sprintf("Wrong mode for file: %v, error: %v", mode, err.Error()))
	}
	to, err := os.OpenFile(pathTo, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(modeOct))
	if err != nil {
		panic(fmt.Sprintf("Error copy file %v: %v", pathFrom, err.Error()))
	}

	_, err = io.Copy(to, from)
	if err != nil {
		panic(fmt.Sprintf("Error copy file %v: %v", pathFrom, err.Error()))
	}
	to.Close()

	info, err := os.Stat(pathFrom)
	if err != nil {
		panic(fmt.Sprintf("Error read file info: %v", err.Error()))
	}
	err = os.Chtimes(pathTo, time.Now(), info.ModTime())
	if err != nil {
		panic(fmt.Sprintf("Error set file permissions: %v", err.Error()))
	}
}

func removeEmptyDirs(path string) {
	path = filepath.Clean(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("Error list directory %v: %v", path, err.Error()))
	}

	for _, file := range files {
		if file.IsDir() {
			removeEmptyDirs(path + "/" + file.Name())
		}
	}

	files, err = ioutil.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("Error list directory %v: %v", path, err.Error()))
	}
	if len(files) == 0 {
		os.Remove(path)
	}
}

func mustBeRelative(path string) {
	if strings.Contains(path, "..") {
		panic(fmt.Sprintf("Wrong relative path: %v", path))
	}
	if strings.HasPrefix(path, " ") || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
		panic(fmt.Sprintf("Wrong relative path: %v", path))
	}
}
