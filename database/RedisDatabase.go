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
	// TODO add a recover to delete hash on err, and then pass along the panic
	// if the delete fails notify the user of the id and that it should be deleted
	// or that the database is in an uncertain state
	client := connectToRedis()
	game.ID = generateUUID(client, "game:id")
	gameKey := fmt.Sprintf("game:%s", game.ID)
	teamIDs := generateTeams(client, game.Teams)
	gameFieldsMap := map[string]interface{}{
		"game_id":    game.ID,
		"name":       game.Name,
		"start_node": game.StartNode,
		"time_limit": game.TimeLimit,
		"teams":      teamIDs,
	}
	err := client.HMSet(gameKey, gameFieldsMap).Err()
	checkErr(err)

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
	checkErr(err)
	err = client.Set("game:id", 0, 0).Err()
	checkErr(err)
	err = client.Set("team:id", 0, 0).Err()
	checkErr(err)
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
	checkErr(err)
}

func generateUUID(client *redis.Client, key string) string {
	result, err := client.Incr(key).Result()
	checkErr(err)
	return string(result)
}

func generateTeams(client *redis.Client, names []string) []string {
	// TODO add a recover to delete the hash on err
	teamIDs := make([]string, len(names))
	for i, name := range names {
		teamIDs[i] = generateUUID(client, "team:id")
		teamKey := fmt.Sprintf("team:%s", teamIDs[i])
		teamFieldsMap := map[string]interface{}{"team_id": teamIDs[i], "name": name, "score": 0}
		err := client.HMSet(teamKey, teamFieldsMap)
		checkErr(err)
	}
	return teamIDs
}

func checkErr(err interface{}) {
	if err != nil {
		log.Panicln(err)
	}
}
