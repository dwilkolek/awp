package domain

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/exp/slices"
)

const SSM_PROXY_PORT = 15077
const REGION = "eu-west-1"
const PROFILE = "awp"
const PARAMETERS = "{\"host\":[\"access-management.service\"], \"portNumber\":[\"80\"], \"localPortNumber\":[\"15077\"]}"
const AWS_PROFILE = "[awp]\nsso_start_url = https://technipfmc.awsapps.com/start\nsso_region = eu-west-1\nsso_account_id = 835811189142\nsso_role_name = AWSSSO-DeveloperAccess\nregion = eu-west-1"

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
	serviceKey := service
	if !slices.Contains(config.Hosts, serviceKey) {
		if !slices.Contains(config.Hosts, service+".service") {
			log.Panicf("Service %s is not valid.\n", service)
		} else {
			serviceKey = service + ".service"
		}
	}
	config.HeaderOverwrites[serviceKey] = map[string]string{}
	config.HeaderOverwrites[serviceKey]["USER-UUID"] = "toBase64:9e69a631-8aa7-4c67-81ed-aadf8b9e4efa"
	config.HeaderOverwrites[serviceKey]["USER-EMAIL"] = "toBase64:bruce@technipfmc.com"
	config.HeaderOverwrites[serviceKey]["USER-NAME"] = "toBase64:Bruce, the bad guy"
	config.HeaderOverwrites[serviceKey]["USER-ROLES"] = "toBase64:USER,ADMIN"

	saveAWPConfig()

	lock.Unlock()
}

func saveAWPConfig() {
	file, _ := json.MarshalIndent(config, "", " ")
	log.Printf("Saving config version %d... %s\n%+v\n", config.Version, configPath(), config)
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
	defer lock.Unlock()

	log.Default().Printf("Using config file: %s\n", configPath())
	content, err := ioutil.ReadFile(configPath())
	if err != nil {
		config = configInfo{
			Version:          2,
			Hosts:            make([]string, 0),
			HeaderOverwrites: make(map[string]map[string]string),
		}
		saveAWPConfig()
		return
	}

	// Now let's unmarshall the data into `payload`
	var payload configInfo
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	config = payload

	if config.Version == 1 {
		config.Version = 2
		config.HeaderOverwrites = make(map[string]map[string]string)
		saveAWPConfig()
	}
}
