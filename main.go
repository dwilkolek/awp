package main

import (
	"os"

	awsserviceproxy "github.com/tfmcdigital/aws-service-proxy/internal"
)

func main() {
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "hosts":
		{
			awsserviceproxy.Setup()
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
			awsserviceproxy.Start(env)
		}
	}
}
