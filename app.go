package main

import (
	"fmt"
	"os"

	app "github.com/tfmcdigital/aws-web-proxy/internal/app"
)

func main() {
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "hosts":
		{
			app.SetupHosts()
		}
	case "setup":
		{
			app.SetupAwsProfile()
			app.UpdateBastionKeys()
		}
	case "update-keys":
		{
			app.UpdateBastionKeys()
		}
	case "add-user-headers":
		{
			serviceName := os.Args[2]
			app.AddDefaultUserHeaders(serviceName)
		}
	case "start":
		{
			env := "dev"
			if len(os.Args) > 2 {
				env = os.Args[2]
			}

			app.StartProxy(env)
		}
	case "version":
		{
			fmt.Println("Version " + app.Version())
		}
	}
}
