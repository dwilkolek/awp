package awsserviceproxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type aWPConfig struct {
	Version int      `json:"version"`
	Hosts   []string `json:"hosts"`
}

var AWPConfig aWPConfig

func saveAWPConfig() {
	file, _ := json.MarshalIndent(AWPConfig, "", " ")
	_ = ioutil.WriteFile("./config.json", file, 0644)
}

func init() {
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		AWPConfig = aWPConfig{
			Version: 1,
			Hosts:   make([]string, 0),
		}
		return
	}

	// Now let's unmarshall the data into `payload`
	var payload aWPConfig
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	AWPConfig = payload
}
