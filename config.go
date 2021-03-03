package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	Backends []string `json:"backends"`
}

func readConfig() (config, error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return config{}, err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return config{}, err
	}

	var conf config

	json.Unmarshal(bytes, &conf)
	return conf, nil
}
