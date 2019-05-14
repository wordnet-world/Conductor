package models

import (
	"log"

	"github.com/tkanos/gonfig"
)

// Configuration is the overall configuration
type Configuration struct {
	Redis   RedisConfiguration `json:"redis"`
	Wordnet WordnetWorld       `json:"wordnetWorld"`
	Kafka   KafkaConfiguration `json:"kafka"`
	Neo4j   Neo4jConfiguration `json:"neo4j"`
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

// KafkaConfiguration is the configuration of a kafka address
type KafkaConfiguration struct {
	Address string `json:"address"`
}

// Neo4jConfiguration is the configuration for neo4j connection
type Neo4jConfiguration struct {
	URI      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
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
