package awswebproxy

import (
	"os"
	"path/filepath"
)

var Version string

func baseAwpPath() string {
	if Version == "development" {
		return "."
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exePath := filepath.Dir(ex)
		return exePath
	}
}
