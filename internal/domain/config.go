package domain

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/exp/slices"
)

var lock = &sync.Mutex{}

type configInfo struct {
	Version          int                          `json:"version"`
	Hosts            []string                     `json:"hosts"`
	HeaderOverwrites map[string]map[string]string `json:"headers"`
}

var config configInfo

func UpdateHosts(hosts []string) {
	lock.Lock()
	config.Hosts = hosts
	saveAWPConfig()
	lock.Unlock()
}

func GetConfig() configInfo {
	return config
}

func AddDefaultUserHeaders(service string) {
	lock.Lock()
	if !slices.Contains(config.Hosts, service) {
		log.Panicf("Service %s is not valid.\n", service)
	}
	config.HeaderOverwrites[service] = map[string]string{}
	config.HeaderOverwrites[service]["USER-UUID"] = "toBase64:9e69a631-8aa7-4c67-81ed-aadf8b9e4efa"
	config.HeaderOverwrites[service]["USER-EMAIL"] = "toBase64:bruce@technipfmc.com"
	config.HeaderOverwrites[service]["USER-NAME"] = "toBase64:Bruce, the bad guy"
	config.HeaderOverwrites[service]["USER-ROLES"] = "toBase64:USER,ADMIN"

	saveAWPConfig()

	lock.Unlock()
}

func saveAWPConfig() {
	file, _ := json.MarshalIndent(config, "", " ")
	log.Printf("Saving config(version = %d) ... %s\n", config.Version, configPath())
	err := ioutil.WriteFile(configPath(), file, 0644)
	if err != nil {
		panic(err)
	}
}

func configPath() string {
	return BasePath + "/config.json"
}

func init() {
	lock.Lock()
	emptyConfig := configInfo{
		Version:          2,
		Hosts:            make([]string, 0),
		HeaderOverwrites: make(map[string]map[string]string),
	}
	content, err := ioutil.ReadFile(configPath())
	if err != nil {
		config = emptyConfig
	}

	// Now let's unmarshall the data into `payload`
	var payload configInfo
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Printf("Cannot read config file from %s. Creating new one. \n", configPath())
		config = emptyConfig
		saveAWPConfig()
	} else {
		config = payload
	}

	if config.Version == 1 {
		config.Version = 2
		config.HeaderOverwrites = make(map[string]map[string]string)
		saveAWPConfig()
	}
	lock.Unlock()
}
