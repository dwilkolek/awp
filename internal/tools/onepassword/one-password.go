package onepassword

import (
	"os/exec"
	"sync"
)

type opClient struct {
	mutex *sync.Mutex
}

var client *opClient

func init() {
	client = &opClient{
		mutex: &sync.Mutex{},
	}
}

func GetOpClient() *opClient {
	return client
}

func (op *opClient) GetDocument(documentId string) []byte {
	cmd := exec.Command("op", "document", "get", documentId)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("Please login to 1password by executing `eval $(op signin)`")
	}
	return out
}

func (op *opClient) GetAwsConfigNote(itemId string) string {
	cmd := exec.Command("op", "item", "get", itemId, "--fields", "notesPlain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("Please login to 1password by executing `eval $(op signin)`")
	}
	return string(out)
}
