package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/digitalocean/doctl/Godeps/_workspace/src/gopkg.in/yaml.v1"
)

var configPath = ".digitalocean/"
var configFileName = "authorize.yml"
var configApiKey = "v2_oauth_key"

func AbsoluteConfigPath() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	return homeDir + "/" + configPath
}

func AbsoluteConfigFilePath() string {
	return AbsoluteConfigPath() + configFileName
}

func ConfigDirectoryExists() bool {
	_, err := os.Stat(AbsoluteConfigPath())
	return !os.IsNotExist(err)
}

func ConfigFileExists() bool {
	_, err := os.Stat(AbsoluteConfigFilePath())
	return !os.IsNotExist(err)
}

func BuildConfigHash() map[string]string {
	return map[string]string{
		configApiKey: APIKey,
	}
}

func LoadConfigHash() (map[string]string, error) {
	if ConfigFileExists() {
		configFile, err := os.Open(AbsoluteConfigFilePath())
		if err != nil {
			return nil, fmt.Errorf("Error opening config file: %s\n", err)
		}

		configFileData, err := ioutil.ReadAll(configFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading config file: %s\n", err)
		}

		config := map[string]string{}
		yaml.Unmarshal(configFileData, &config)
		return config, nil
	} else {
		return nil, fmt.Errorf("No config file exists")
	}
}
