package awswebproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func SetupAwsProfile() {
	client := OpClient{
		mutex: &sync.Mutex{},
	}
	homedir, _ := os.UserHomeDir()
	credentialsPath := homedir + "/.aws/credentials"
	data := strings.ReplaceAll(client.getAwsConfigNote("qt7ixhmfmszawh6c42gdjtx5wq"), "\"", "")
	_, err := os.Stat(credentialsPath)
	if os.IsNotExist(err) {
		os.MkdirAll(homedir+"/.aws", os.ModePerm)
		os.Create(credentialsPath)
	}

	fileBody, _ := ioutil.ReadFile(credentialsPath)
	if strings.Contains(string(fileBody), data) {
		fmt.Println("AWS profile already set")
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
	client := OpClient{
		mutex: &sync.Mutex{},
	}

	for env, documentId := range envDocumentIdMap {
		keyBytes := client.getDocument(documentId)
		storeKey(env, keyBytes)
	}
}

type OpClient struct {
	mutex *sync.Mutex
}

var envDocumentIdMap = map[string]string{
	"dev":  "hgbzp2ptafe75aqxi5bgzkchty",
	"demo": "bk3jvxag6na2rn6nruorn2m5ri",
	"prod": "czzulyxrkjaivmsq2xo7mgpqqm",
}

func (op *OpClient) getDocument(documentId string) []byte {
	cmd := exec.Command("op", "document", "get", documentId)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("Please login to 1password by executing `eval $(op signin)`")
	}
	return out
}

func (op *OpClient) getAwsConfigNote(itemId string) string {
	cmd := exec.Command("op", "item", "get", itemId, "--fields", "notesPlain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("Please login to 1password by executing `eval $(op signin)`")
	}
	return string(out)
}

func FileName(env string) string {
	homedir, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.ssh/%s-bastion-rsa.pem", homedir, env)
}

func storeKey(env string, data []byte) {
	err := ioutil.WriteFile(FileName(env), data, 0644)
	if err != nil {
		panic(err)
	}
	log.Default().Printf("Successfuly updated key for %s environment\n", env)
}
