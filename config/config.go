package config

import (
	"encoding/json"
	"os"
)

type Config struct {

	Server struct {
		HOST string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`

	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Dbname   string `json:"dbname"`
		User     string `json:"user"`
		Password string `json:"password"`
		Sslmode  string `json:"sslmode"`
	} `json:"database"`

}

func ConfigInit() error {

	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return err
	}

	// Server configuration
	os.Setenv("SERVER_HOST", config.Server.HOST)
	os.Setenv("SERVER_PORT", config.Server.Port)

	// Database configuration
	os.Setenv("DB_USERNAME", config.Database.User)
	os.Setenv("DB_PASSWORD", config.Database.Password)
	os.Setenv("DB_HOST", config.Database.Host)
	os.Setenv("DB_PORT", config.Database.Port)
	os.Setenv("DB_DATABASE", config.Database.Dbname)
	os.Setenv("DB_SSLMODE", config.Database.Sslmode)

	return nil
}
