package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/wordnet-world/Conductor/models"
)

// RedisDatabase struct, will implement Database
type RedisDatabase struct {
}

// CreateGame allows for game creation when providing a game model, returns the new gameID
func (redisDatabase RedisDatabase) CreateGame(game models.CreateGame) string {
	// TODO add a recover to delete hash on err, and then pass along the panic
	// if the delete fails notify the user of the id and that it should be deleted
	// or that the database is in an uncertain state
	client := connectToRedis()
	defer client.Close()
	game.ID = generateUUID(client, "game:id")
	gameKey := fmt.Sprintf("game:%s", game.ID)
	teamIDs := generateTeams(client, game.Teams)
	gameFieldsMap := map[string]interface{}{
		"gameID":    game.ID,
		"name":      game.Name,
		"startNode": game.StartNode,
		"timeLimit": game.TimeLimit,
		"teamIDs":   teamIDs,
		"status":    "waiting",
		"startTime": 0,
	}
	log.Printf("CreateGame fields: %v\n", gameFieldsMap)
	err := client.HMSet(gameKey, gameFieldsMap).Err()
	checkErr(err)
	err = client.SAdd("games", gameKey).Err()
	checkErr(err)
	return game.ID
}

// GetGames returns a slice of Game objects
func (redisDatabase RedisDatabase) GetGames(fields []string) []map[string]interface{} {
	client := connectToRedis()
	defer client.Close()
	keys, err := client.SMembers("games").Result()
	if err != nil {
		log.Panicln(err)
	}
	if len(keys) < 1 {
		log.Panicln("No games found")
	}
	log.Printf("keys: %v\n", keys)

	for i, field := range fields {
		if field == "teams" {
			fields[i] = "teamIDs"
			break
		}
	}
	games := make([]map[string]interface{}, len(keys))

	for i, key := range keys {
		game, err := client.HMGet(key, fields...).Result()
		if err != nil {
			log.Panicln(err)
		}
		games[i] = make(map[string]interface{})
		for j, field := range fields {
			if field == "teamIDs" {
				var teamIDs []string
				err = json.Unmarshal([]byte(game[j].(string)), &teamIDs)
				if err != nil {
					log.Panicln(err)
				}
				games[i]["teams"] = getTeamData(client, teamIDsToKeys(teamIDs))
			} else if field == "timeLimit" {
				games[i][field], err = strconv.Atoi(game[j].(string))
				if err != nil {
					log.Panicln(err)
				}
			} else if field == "startTime" {
				games[i][field], err = strconv.Atoi(game[j].(string))
				if err != nil {
					log.Panicln(err)
				}
			} else {
				games[i][field] = game[j]
			}
		}
	}

	return games
}

// GetGame is like ListGames but only for the provided gameID
func (redisDatabase RedisDatabase) GetGame(fields []string, gameID string) map[string]interface{} {
	client := connectToRedis()
	defer client.Close()

	key := fmt.Sprintf("game:%s", gameID)

	game := make(map[string]interface{})
	for i, field := range fields {
		if field == "teams" {
			fields[i] = "teamIDs"
			break
		}
	}

	for i, field := range fields {
		if field == "teams" {
			fields[i] = "teamIDs"
			break
		}
	}

	gameObj, err := client.HMGet(key, fields...).Result()
	for j, field := range fields {
		if field == "teamIDs" {
			var teamIDs []string
			err = json.Unmarshal([]byte(gameObj[j].(string)), &teamIDs)
			if err != nil {
				log.Panicln(err)
			}
			game["teams"] = getTeamData(client, teamIDsToKeys(teamIDs))
		} else if field == "timeLimit" {
			game[field], err = strconv.Atoi(gameObj[j].(string))
			if err != nil {
				log.Panicln(err)
			}
		} else if field == "startTime" {
			game[field], err = strconv.Atoi(gameObj[j].(string))
			if err != nil {
				log.Panicln(err)
			}
		} else {
			game[field] = gameObj[j]
		}
	}

	return game
}

// GetTeams returns a slice of Team objects for a given Game
func (redisDatabase RedisDatabase) GetTeams() []models.Team {
	client := connectToRedis()
	defer client.Close()
	teamKeys, err := client.SMembers("teams").Result()
	checkErr(err)

	teams := getTeamData(client, teamKeys)
	return teams
}

// GetTeam is like GetTeams but for a single teamID
func (redisDatabase RedisDatabase) GetTeam(teamID string) models.Team {
	client := connectToRedis()
	defer client.Close()
	teamIDs := []string{teamID}
	teams := getTeamData(client, teamIDsToKeys(teamIDs))
	return teams[0]
}

// DeleteGame deletes the Game matching the provided gameID
// Currently not handling partial deletes
func (redisDatabase RedisDatabase) DeleteGame(gameID string) bool {
	client := connectToRedis()
	defer client.Close()
	gameKey := fmt.Sprintf("game:%s", gameID)
	redisTeamIDs, err := client.HGet(gameKey, "teamIDs").Result()
	checkErr(err)
	var teamIDs []string
	err = json.Unmarshal([]byte(redisTeamIDs), &teamIDs)
	checkErr(err)
	deleteTeams(client, teamIDs)
	err = client.HDel(gameKey, "gameID", "name", "startNode", "timeLimit", "teamIDs", "status", "startTime").Err()
	checkErr(err)
	err = client.SRem("games", gameKey).Err()

	return true
}

// GetConsumerID returns a unique id string for creating consumers
func (redisDatabase RedisDatabase) GetConsumerID() string {
	client := connectToRedis()
	defer client.Close()
	consumerID := generateUUID(client, "consumer:id")
	return consumerID
}

// SetupTeamCaches is used in the StartGame endpoint and initializes teams caches
func (redisDatabase RedisDatabase) SetupTeamCaches(teamIDs []string, root models.Node, neighbors []models.Node) {
	client := connectToRedis()
	defer client.Close()

	nodesAsInterfaces := make([]interface{}, len(neighbors))
	for i, v := range neighbors {
		nodesAsInterfaces[i] = v.ID
		key := fmt.Sprintf("node:%d", v.ID)
		err := client.Set(key, v.Text, 0).Err()
		if err != nil {
			log.Panicln(err)
		}
	}

	rootKey := fmt.Sprintf("node:%d", root.ID)
	err := client.Set(rootKey, root.Text, 0).Err()
	if err != nil {
		log.Panicln(err)
	}

	for _, teamID := range teamIDs {
		periphiKey := fmt.Sprintf("periph:%s", teamID)
		knownKey := fmt.Sprintf("known:%s", teamID)
		err := client.SAdd(periphiKey, nodesAsInterfaces...).Err()
		if err != nil {
			log.Panicln(err)
		}
		err = client.SAdd(knownKey, root.ID).Err()
		if err != nil {
			log.Panicln(err)
		}
	}
}

// SetupDB should be run as the server starts to clear the DB and
// set the counters for uuids
func (redisDatabase RedisDatabase) SetupDB() {
	/*defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
		}
	}()*/

	client := connectToRedis()
	defer client.Close()

	err := client.FlushAll().Err()
	checkErr(err)
	err = client.Set("game:id", 0, 0).Err()
	checkErr(err)
	err = client.Set("team:id", 0, 0).Err()
	checkErr(err)
	err = client.Set("consumer:id", 0, 0).Err()
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
	log.Printf("Generated ID:%d for Key:%s\n", result, key)
	return strconv.FormatInt(result, 10)
}

func generateTeams(client *redis.Client, names []string) string {
	// TODO add a recover to delete the hash on err
	teamIDs := make([]string, len(names))
	for i, name := range names {
		teamIDs[i] = generateUUID(client, "team:id")
		teamKey := fmt.Sprintf("team:%s", teamIDs[i])
		teamFieldsMap := map[string]interface{}{"teamID": teamIDs[i], "name": name, "score": 0}
		err := client.HMSet(teamKey, teamFieldsMap).Err()
		checkErr(err)
		err = client.SAdd("teams", teamKey).Err()
		checkErr(err)
	}
	result, _ := json.Marshal(teamIDs)
	return string(result)
}

// teamKeys should be the full "team:teamID" strings, so I can reuse the function
func getTeamData(client *redis.Client, teamKeys []string) []models.Team {
	teams := make([]models.Team, len(teamKeys))
	for i, teamKey := range teamKeys {
		team, err := client.HGetAll(teamKey).Result()
		checkErr(err)
		score, err := strconv.Atoi(team["score"])
		checkErr(err)
		teams[i].ID = team["teamID"]
		teams[i].Name = team["name"]
		teams[i].Score = score
	}
	return teams
}

func teamIDsToKeys(teamIDs []string) []string {
	keys := make([]string, len(teamIDs))
	for i, teamID := range teamIDs {
		keys[i] = fmt.Sprintf("team:%s", teamID)
	}
	return keys
}

func deleteTeams(client *redis.Client, teamIDs []string) {
	for _, ID := range teamIDs {
		teamKey := fmt.Sprintf("team:%s", ID)
		err := client.HDel(teamKey, "teamID", "name", "score").Err()
		checkErr(err)
	}
}

func checkErr(err interface{}) {
	if err != nil {
		log.Panicln(err)
	}
}
