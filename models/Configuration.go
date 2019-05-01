package models

import (
	"log"

	"github.com/tkanos/gonfig"
)

// Configuration is the overall configuration
type Configuration struct {
	Redis   RedisConfiguration `json:"redis"`
	Wordnet WordnetWorld       `json:"wordnetWorld"`
}

// RedisConfiguration for connection
type RedisConfiguration struct {
	Port     int    `json:"port"`
	Address  string `json:"address"`
	Database int    `json:"db"`
	Password string `json:"password"`
}

// WordnetWorld is the configuration pertaining to the Conductor application
type WordnetWorld struct {
	AdminPassword string `json:"adminPassword"`
}

// Config is contains the loaded configuration for the Conductor Application
var Config Configuration

func init() {
	config := Configuration{}
	err := gonfig.GetConf("./config/conductor-conf.json", &config)
	if err != nil {
		log.Panicln(err)
	}
	Config = config
}
