package main

import (
	"fmt"
	"os"

	awswebproxy "github.com/tfmcdigital/aws-web-proxy/internal"
	localserver "github.com/tfmcdigital/aws-web-proxy/internal/localserver"
)

func main() {
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "hosts":
		{
			awswebproxy.SetupHosts()
		}
	case "setup":
		{
			awswebproxy.SetupAwsProfile()
			awswebproxy.UpdateBastionKeys()
		}
	case "update-keys":
		{
			awswebproxy.UpdateBastionKeys()
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

			localserver.Start(env)
		}
	case "version":
		{
			fmt.Println("Version " + awswebproxy.Version)
		}
	}
}
