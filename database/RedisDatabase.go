package database

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/wordnet-world/Conductor/models"
)

// RedisDatabase struct, will implement Database
type RedisDatabase struct {
}

// CreateGame allows for game creation when providing a game model, returns the new gameID
func (redisDatabase RedisDatabase) CreateGame(game models.Game) string {
	client := connectToRedis()
	game.ID = generateGameUUID(client)
	return game.ID
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

// SetupDB should be run as the server starts to clear the DB and
// set the counters for uuids
func (redisDatabase RedisDatabase) SetupDB() {
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
		}
	}()

	client := connectToRedis()

	err := client.FlushAll().Err()
	if err != nil {
		log.Panicln(err)
	}
	err = client.Set("game:id", 0, 0).Err()
	if err != nil {
		log.Panicln(err)
	}
	err = client.Set("team:id", 0, 0).Err()
	if err != nil {
		log.Panicln(err)
	}
}

func connectToRedis() *redis.Client {
	connectionString := fmt.Sprintf("%s:%d", models.Config.Redis.Address, models.Config.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: models.Config.Redis.Password,
		DB:       models.Config.Redis.Database,
	})

	testConnection(client)
	return client
}

// checks whether we can connect to the database
func testConnection(client *redis.Client) {
	_, err := client.Ping().Result()
	if err != nil {
		log.Panicln(err)
	}
}

func generateGameUUID(client *redis.Client) string {
	return ""
}

func generateTeamUUID(client *redis.Client) string {
	return ""
}
