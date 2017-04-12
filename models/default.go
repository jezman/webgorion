package models

import (
	"encoding/json"
	"fmt"
	"os"
)

var Conf = ReadConfigFile()

func ReadConfigFile() Config {
	confFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Read configuration file error:", err)

	}
	defer confFile.Close()

	decoder := json.NewDecoder(confFile)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println("JSON decode error:", err)

	}
	return conf

}
