package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

var BasePath string = baseAwpPath()

func baseAwpPath() string {
	fmt.Println(Version)
	if Version == "development" {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return path
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exePath := filepath.Dir(ex)
		return exePath
	}
}
