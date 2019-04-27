package database

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/tkanos/gonfig"
	"github.com/wordnet-world/Conductor/models"
)

// RedisDatabase struct, will implement Database
type RedisDatabase struct {
}

// CreateGame allows for game creation when providing a game model
func (redisDatabase RedisDatabase) CreateGame(game models.Game) string {
	client := connectToRedis()
	log.Println(testConnection(client))
	return ""
}

// GetGames returns a slice of Game objects
func (redisDatabase RedisDatabase) GetGames() []models.Game {
	return nil
}

// GetTeams returns a slice of Team objects for a given Game
func (redisDatabase RedisDatabase) GetTeams(gameID string) []models.Team {
	return nil
}

// DeleteGame deletes the Game matching the provided gameID
func (redisDatabase RedisDatabase) DeleteGame(gameID string) {

}

func connectToRedis() *redis.Client {
	config := models.Configuration{}
	err := gonfig.GetConf("./config/conductor-conf.json", &config)
	if err != nil {
		log.Panicln(err)
	}
	connectionString := fmt.Sprintf("%s:%d", config.Redis.Address, config.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: config.Redis.Password,
		DB:       config.Redis.Database,
	})

	return client
}

// TestConnection checks whether we can connect to the database
func testConnection(client *redis.Client) bool {
	_, err := client.Ping().Result()
	if err != nil {
		log.Panicln(err)
		return false
	}
	return true
}
