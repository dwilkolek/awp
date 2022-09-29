package domain

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type configInfo struct {
	Version int      `json:"version"`
	Hosts   []string `json:"hosts"`
}

var config configInfo

func UpdateHosts(hosts []string) {
	config.Hosts = hosts
	saveAWPConfig()
}

func GetConfig() configInfo {
	return config
}

func saveAWPConfig() {
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(configPath(), file, 0644)
}

func configPath() string {
	return BasePath + "/config.json"
}

func init() {
	content, err := ioutil.ReadFile(configPath())
	if err != nil {
		config = configInfo{
			Version: 1,
			Hosts:   make([]string, 0),
		}
		return
	}

	// Now let's unmarshall the data into `payload`
	var payload configInfo
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	config = payload
}
