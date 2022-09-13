package awswebproxy

import (
	"os"
	"path/filepath"
)

var Version string = "development"

func BaseAwpPath() string {
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
