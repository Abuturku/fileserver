package fileserver

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Port      string
	ServerCrt string
	ServerKey string
}

func GetConfig() Configuration {
	file, _ := os.Open("src/de/vorlesung/projekt/landrive/fileserver/conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
	//fmt.Println(configuration.Users) // output: [UserA, UserB]
}
