package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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

// UpdateGame updates the given fields with new values
func (redisDatabase RedisDatabase) UpdateGame(gameID string, updates map[string]interface{}) {
	client := connectToRedis()
	defer client.Close()

	key := fmt.Sprintf("game:%s", gameID)

	err := client.HMSet(key, updates).Err()
	if err != nil {
		log.Panicln(err)
	}
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
	deleteTeamCaches(client, teamIDs)
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

	for _, node := range neighbors {
		addNodeToCache(client, node)
	}

	addNodeToCache(client, root)

	for _, teamID := range teamIDs {
		addNodesToPeriphery(client, teamID, neighbors)
		addNodeToFound(client, teamID, root)
	}
}

func deleteTeamCaches(client *redis.Client, teamIDs []string) {
	for _, t := range teamIDs {
		periphiKey := fmt.Sprintf("periph:%s", t)
		knownKey := fmt.Sprintf("known:%s", t)
		err := client.Del(periphiKey, knownKey).Err()
		if err != nil {
			log.Panicln(err)
		}
	}
}

func addNodeToFound(client *redis.Client, teamID string, node models.Node) {
	knownKey := fmt.Sprintf("known:%s", teamID)
	err := client.SAdd(knownKey, node.ID).Err()
	if err != nil {
		log.Panicln(err)
	}
}

func addNodesToPeriphery(client *redis.Client, teamID string, nodes []models.Node) {
	periphiKey := fmt.Sprintf("periph:%s", teamID)
	nodesAsInterfaces := make([]interface{}, len(nodes))
	for i, v := range nodes {
		nodesAsInterfaces[i] = v.ID
	}
	err := client.SAdd(periphiKey, nodesAsInterfaces...).Err()
	if err != nil {
		log.Panicln(err)
	}
}

func addNodeToCache(client *redis.Client, node models.Node) {
	idKey := fmt.Sprintf("nodeID:%d", node.ID)
	textKey := fmt.Sprintf("nodeText:%s", strings.ToLower(node.Text))
	err := client.SetNX(idKey, strings.ToLower(node.Text), 0).Err()
	if err != nil {
		log.Panicln(err)
	}
	err = client.SetNX(textKey, node.ID, 0).Err()
	if err != nil {
		log.Panicln(err)
	}
}

func removeNodesFromPeriphery(client *redis.Client, teamID string, nodes []models.Node) {
	periphiKey := fmt.Sprintf("periph:%s", teamID)
	nodeIDs := make([]interface{}, len(nodes))
	for i, node := range nodes {
		nodeIDs[i] = node.ID
	}
	err := client.SRem(periphiKey, nodeIDs...).Err()
	if err != nil {
		log.Panicln(err)
	}
}

// IsFound returns true if the guess is already in the graph
func (redisDatabase RedisDatabase) IsFound(guess string, teamID string) bool {
	client := connectToRedis()
	defer client.Close()

	guessKey := fmt.Sprintf("nodeText:%s", guess)
	nodeID, err := client.Get(guessKey).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Panicln(err)
	}

	knownKey := fmt.Sprintf("known:%s", teamID)
	found, err := client.SIsMember(knownKey, nodeID).Result()
	if err != nil {
		log.Panicln(err)
	}
	return found
}

// IsPeriphery returns the nodeID if the guess is in the periphery
func (redisDatabase RedisDatabase) IsPeriphery(guess string, teamID string) int64 {
	client := connectToRedis()
	defer client.Close()

	guessKey := fmt.Sprintf("nodeText:%s", guess)
	nodeID, err := client.Get(guessKey).Result()
	if err == redis.Nil {
		return -1
	} else if err != nil {
		log.Panicln(err)
	}

	periphiKey := fmt.Sprintf("periph:%s", teamID)
	periph, err := client.SIsMember(periphiKey, nodeID).Result()
	if err != nil {
		log.Panicln(err)
	}
	if periph {
		nodeIDInt64, err := strconv.ParseInt(nodeID, 10, 64)
		if err != nil {
			log.Panicln(err)
		}
		return nodeIDInt64
	}
	return -1
}

// UpdateCache gets the diff of new neighbors and found nodes, updates the periphery, and updates found, and returning the diff
// resultNodes first, foundNodes second in return
func (redisDatabase RedisDatabase) UpdateCache(newNode models.Node, neighbors []models.Node, teamID string) ([]models.Node, []models.Node) {
	client := connectToRedis()
	defer client.Close()

	knownKey := fmt.Sprintf("known:%s", teamID)

	resultNodes := make([]models.Node, 0, len(neighbors))
	foundNodes := make([]models.Node, 0, len(neighbors))
	for _, node := range neighbors {
		known, err := client.SIsMember(knownKey, node.ID).Result()
		if err != nil {
			log.Panicln(err)
		}
		if !known {
			resultNodes = append(resultNodes, node)
		} else {
			foundNodes = append(foundNodes, node)
		}
	}

	if len(resultNodes) != 0 {
		addNodesToPeriphery(client, teamID, resultNodes)
	}
	addNodeToFound(client, teamID, newNode)
	removeNodesFromPeriphery(client, teamID, []models.Node{newNode})
	for _, node := range neighbors {
		addNodeToCache(client, node)
	}

	return resultNodes, foundNodes
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
