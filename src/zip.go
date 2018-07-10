package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func updateFromZip(zipPath string, outputBaseDir string, outputDirPrefix string, localFiles map[string]fileinfo, useModes bool) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
	}
	defer r.Close()

	for _, f := range r.File {
		mustBeRelative(f.Name)
		if !f.FileInfo().IsDir() {
			var needUpdate bool
			outFilePath := outputBaseDir + outputDirPrefix + f.Name
			if localInfo, ok := localFiles[outputDirPrefix+f.Name]; ok {
				delete(localFiles, outputDirPrefix+f.Name)
				needUpdate = int64(f.UncompressedSize64) != localInfo.size || f.Modified.Unix() != localInfo.lastModified
				if !needUpdate && useModes && localInfo.mode != uint32(f.Mode()) {
					localPermMode(outFilePath, localInfo.mode, f.Mode())
				}
			} else {
				needUpdate = true
			}
			if needUpdate {
				rc, err := f.Open()
				if err != nil {
					panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
				}
				defer rc.Close()
				if strings.Contains(f.Name, "..") {
					panic(fmt.Sprintf("Wrong file name in zip: %v", f.Name))
				}

				mustBeRelative(f.Name)
				if err = os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm); err != nil {
					panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
				}

				var outFile *os.File
				if useModes {
					outFile, err = os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				} else {
					outFile, err = os.Create(outFilePath)
				}
				if err != nil {
					panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
				}

				_, err = io.Copy(outFile, rc)
				if err != nil {
					panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
				}
				outFile.Close()
				err = os.Chtimes(outFilePath, time.Now(), time.Unix(f.Modified.Unix(), 0))
				if err != nil {
					panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
				}
			}
		}
	}
}
