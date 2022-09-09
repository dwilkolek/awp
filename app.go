package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/creativeprojects/go-selfupdate"
	awsserviceproxy "github.com/tfmcdigital/aws-web-proxy/internal"
)

var version string

func main() {
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "hosts":
		{
			awsserviceproxy.SetupHosts()
		}
	case "setup":
		{
			awsserviceproxy.SetupAwsProfile()
			awsserviceproxy.UpdateBastionKeys()
		}
	case "update-keys":
		{
			awsserviceproxy.UpdateBastionKeys()
		}
	case "start":
		{
			env := "dev"
			if len(os.Args) > 2 {
				env = os.Args[2]
			}
			if env != "dev" && env != "demo" && env != "prod" {
				panic("Do not recognize that environment: " + env)
			}
			awsserviceproxy.StartWebServer()
			awsserviceproxy.Start(env)
		}
	case "version":
		{
			fmt.Println("Version " + version)
		}
	case "update":
		{
			latest, found, err := selfupdate.DetectLatest("creativeprojects/resticprofile")
			if err != nil {
				fmt.Printf("error occurred while detecting version: %v", err)
				return
			}
			if !found {
				fmt.Printf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
			}

			if latest.LessOrEqual(version) {
				log.Printf("Current version (%s) is the latest", version)
				return
			}

			exe, err := os.Executable()
			if err != nil {
				fmt.Printf("could not locate executable path")
				return
			}
			if err := selfupdate.UpdateTo(latest.AssetURL, latest.AssetName, exe); err != nil {
				fmt.Printf("error occurred while updating binary: %v", err)
				return
			}
			log.Printf("Successfully updated to version %s", latest.Version())
		}
	}
}
