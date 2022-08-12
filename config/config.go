package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

type UnderscoreAIConfig struct {
	HFAPIKey          string `json:"hf_api_key"`
	PromptContextPath string `json:"prompt_context_path"`
}

var Config UnderscoreAIConfig

func LoadConfig() {
	// which os is running
	currentOS := runtime.GOOS

	// get home directory of current user
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// if os is windows
	if currentOS == "windows" {
		// load config from windows
		LoadConfigWindows(homeDir)
	} else if currentOS == "linux" {
		// load config from linux
		LoadConfigLinux(homeDir)
	} else if currentOS == "darwin" {
		// load config from mac
		LoadConfigMac(homeDir)
	} else {
		// load config from other os
		log.Fatal("Unsupported OS")
	}
}

func LoadConfigWindows(homeDir string) {
	// get config file path
	configFilePath := homeDir + "\\.underscoreai.json"

	// load config from file
	LoadConfigFromFile(configFilePath)
}

func LoadConfigLinux(homeDir string) {
	// get config file path
	configFilePath := homeDir + "/.underscoreai.json"

	// load config from file
	LoadConfigFromFile(configFilePath)
}

func LoadConfigMac(homeDir string) {
	// get config file path
	configFilePath := homeDir + "/.underscoreai.json"

	// load config from file
	LoadConfigFromFile(configFilePath)
}

func LoadConfigFromFile(configFilePath string) {
	// check is file exists, if not create one
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// create empty file
		dummyConfig := UnderscoreAIConfig{
			HFAPIKey:          "#HuggingFace API Key#",
			PromptContextPath: "~/.underscoreai/prompt_context_linux",
		}

		data, err := json.MarshalIndent(dummyConfig, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(configFilePath, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// open config file and read content
	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	// read config file content
	configFileContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal config file content
	err = json.Unmarshal(configFileContent, &Config)
	if err != nil {
		log.Fatal(err)
	}

}
