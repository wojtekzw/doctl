package main

import (
	"log"
	"os"

	"github.com/digitalocean/doctl/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/digitalocean/doctl/Godeps/_workspace/src/gopkg.in/yaml.v1"
)

var AuthorizeCommand = cli.Command{
	Name:   "authorize",
	Usage:  "Save API credentials for re-use later.",
	Action: writeCredentialsToFile,
}

func writeCredentialsToFile(ctx *cli.Context) {
	if APIKey == "" {
		log.Fatalln("Error: Must provide an API key via an environmental variable or flag.")
	}

	if !ConfigDirectoryExists() {
		os.Mkdir(AbsoluteConfigPath(), 0777)
	}

	configFile, err := os.Create(AbsoluteConfigFilePath())
	if err != nil {
		log.Fatalf("Unable to create file: %s\n", err)
	}
	defer configFile.Close()

	configString, err := yaml.Marshal(BuildConfigHash())
	if err != nil {
		log.Fatalf("YAML Encoding Error: %s", err)
	}

	_, err = configFile.Write(configString)
	if err != nil {
		log.Fatalf("Unable to write to file: %s\n", err)
	}

	log.Fatalf("Saved key to %s\n", AbsoluteConfigFilePath())
}
