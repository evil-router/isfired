package config

import (
	"os"
	"encoding/json"
	"fmt"
)

type Configuration struct {
	DB_Host  string
	DB_User  string
	DB_Pass  string
	DB_Port  string
	DB_Name  string
}

var Config  Configuration

func GetConfig (path string) (error){
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	fmt.Println(configuration.DB_Host)
	Config = configuration
	return nil
}
