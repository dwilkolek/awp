package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
	"github.com/tfmcdigital/aws-web-proxy/internal/proxy"
	"github.com/tfmcdigital/aws-web-proxy/internal/tools/aws"
	"github.com/txn2/txeh"
)

func StartProxy(env string) {
	environment, err := domain.ParseEnvironment(env)
	log.Default().Printf("Starting proxy to %s environemt\n", environment.String())
	if err != nil {
		log.Default().Fatalln(err)
	}
	go gracefulShutdown()
	proxy.StartProxy(environment)

}
func gracefulShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		fmt.Println("Sutting down gracefully.")
		exec.Command(
			"lsof", "-t", "-i", fmt.Sprintf("tcp:%d", domain.SSM_PROXY_PORT), "|", "xargs", "kill",
		).Run()
		exec.Command(
			"lsof", "-t", "-i", fmt.Sprintf("tcp:%d", proxy.WEB_SERVER_PORT), "|", "xargs", "kill",
		).Run()
		os.Exit(0)
	}()
}

func SetupHosts() {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}
	client := aws.GetAwsClient()
	clusters, _ := client.GetEcsClusterMap()
	services := client.GetEcsServices(clusters["dev"])

	knownHosts := []string{}

	hosts.AddHost("127.0.0.1", "awp")

	for service := range services {
		hosts.AddHost("127.0.0.1", service+".service")
		knownHosts = append(knownHosts, service+".service")
	}

	domain.UpdateHosts(knownHosts)
	// hfData := hosts.RenderHostsFile()
	// fmt.Println(hfData)
	err = hosts.Save()
	if err != nil {
		panic(err)
	}
}

func SetupAwsProfile() {
	homedir, _ := os.UserHomeDir()
	credentialsPath := homedir + "/.aws/credentials"
	_, err := os.Stat(credentialsPath)
	if os.IsNotExist(err) {
		os.MkdirAll(homedir+"/.aws", os.ModePerm)
		os.Create(credentialsPath)
	}

	fileBody, _ := ioutil.ReadFile(credentialsPath)
	if strings.Contains(string(fileBody), domain.AWS_PROFILE) {
		log.Default().Println("AWS profile already set")
		return
	}

	f, err := os.OpenFile(credentialsPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()
	if _, err := f.WriteString("\n" + domain.AWS_PROFILE); err != nil {
		log.Println(err)
	}

	log.Default().Println("AWS profile stored")
}

func Version() string {
	return domain.Version
}

func AddDefaultUserHeaders(service string) {
	domain.AddDefaultUserHeaders(service)
}
