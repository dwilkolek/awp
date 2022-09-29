package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
	"github.com/tfmcdigital/aws-web-proxy/internal/proxy"
	"github.com/tfmcdigital/aws-web-proxy/internal/tools/aws"
	"github.com/tfmcdigital/aws-web-proxy/internal/tools/onepassword"
	"github.com/txn2/txeh"
)

var envDocumentIdMap = map[domain.Environment]string{
	domain.DEV:  "hgbzp2ptafe75aqxi5bgzkchty",
	domain.DEMO: "bk3jvxag6na2rn6nruorn2m5ri",
	domain.PROD: "czzulyxrkjaivmsq2xo7mgpqqm",
}

const AWS_PROFILE_DOCUMENT = "qt7ixhmfmszawh6c42gdjtx5wq"

func StartProxy(env string) {
	environment, err := domain.ParseEnvironment(env)
	if err != nil {
		log.Default().Fatalln(err)
	}
	proxy.StartProxy(environment, filePathToBastionKey(environment))
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
	client := onepassword.GetOpClient()
	homedir, _ := os.UserHomeDir()
	credentialsPath := homedir + "/.aws/credentials"
	data := strings.ReplaceAll(client.GetAwsConfigNote(AWS_PROFILE_DOCUMENT), "\"", "")
	_, err := os.Stat(credentialsPath)
	if os.IsNotExist(err) {
		os.MkdirAll(homedir+"/.aws", os.ModePerm)
		os.Create(credentialsPath)
	}

	fileBody, _ := ioutil.ReadFile(credentialsPath)
	if strings.Contains(string(fileBody), data) {
		log.Default().Println("AWS profile already set")
		return
	}

	f, err := os.OpenFile(credentialsPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()
	if _, err := f.WriteString("\n" + data); err != nil {
		log.Println(err)
	}
}

func UpdateBastionKeys() {
	client := onepassword.GetOpClient()

	for env, documentId := range envDocumentIdMap {
		keyBytes := client.GetDocument(documentId)
		storeKey(env, keyBytes)
	}
}

func filePathToBastionKey(env domain.Environment) string {
	homedir, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.ssh/%s-bastion-rsa.pem", homedir, env)
}

func storeKey(env domain.Environment, data []byte) {
	err := ioutil.WriteFile(filePathToBastionKey(env), data, 0644)
	if err != nil {
		panic(err)
	}
	log.Default().Printf("Successfuly updated key for %s environment\n", env)
}

func Version() string {
	return domain.Version
}
