package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func extractZip(localDir string) {
	fmt.Print("Unpacking updates... ")
	err := os.RemoveAll(localDir + WORK_DIR)
	if err != nil {
		panic(fmt.Sprintf("Error remove dir: %v", err.Error()))
	}

	r, err := zip.OpenReader(localDir + ZIP_STORE)
	if err != nil {
		panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
		}
		defer rc.Close()

		outFilePath := filepath.Join(localDir+WORK_DIR, f.Name)

		if !f.FileInfo().IsDir() {
			if err = os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm); err != nil {
				panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
			}

			outFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			if err != nil {
				panic(fmt.Sprintf("Error extract zip: %v", err.Error()))
			}
		}
	}
	fmt.Println("Done")
}
