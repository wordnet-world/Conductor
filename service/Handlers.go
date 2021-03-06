package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wordnet-world/Conductor/database"
	"github.com/wordnet-world/Conductor/models"
	//"github.com/google/go-cmp/cmp"
)

// HeartBeat end point to verify connection
func HeartBeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I'm Alive!")
}

// AdminPasswordCheck determines if the client can access the admin pages
func AdminPasswordCheck(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	verifyPassword(r)

	fmt.Fprintln(w, models.CreateHTTPResponse(nil, "Correct AdminPassword", true).ToJSON())
}

// StartGame asks neo4j for the root node, publishes a graphUpdate to each team's topic
// start some sort of timer which will send an endgame message
// Input will be the gameID to start
func StartGame(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	gameIDarray, ok := r.URL.Query()["gameID"]
	if !ok || len(gameIDarray) < 1 {
		log.Panicln("No query parameter 'gameID' specified")
	}

	cache := database.GetCacheDatabase()
	graph := database.GetGraphDatabase()
	err := graph.Connect(models.Config.Neo4j.URI, models.Config.Neo4j.Username, models.Config.Neo4j.Password)
	if err != nil {
		log.Panicln(err)
	}

	defer graph.Close()

	root, err := graph.GetRoot()
	if err != nil {
		log.Panicln(err)
	}

	neighbors, err := graph.GetNeighbors(root)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Neighbors: %v\n", neighbors)

	games := cache.GetGame([]string{"teams", "timeLimit"}, gameIDarray[0])
	timeLimit := games["timeLimit"].(int)
	teams := games["teams"].([]models.Team) //Forgive me
	teamIDs := make([]string, len(teams))
	for i, team := range teams {
		teamIDs[i] = team.ID
	}
	log.Printf("TeamIDs to populate graph cache for: %v\n", teamIDs)
	cache.SetupTeamCaches(teamIDs, root, neighbors)

	neighborIDs := make([]int64, len(neighbors))
	for i, n := range neighbors {
		neighborIDs[i] = n.ID
	}
	startMessage := models.StartMessage{
		Type:          "start",
		RootID:        root.ID,
		RootText:      root.Text,
		RootNeighbors: neighborIDs,
	}

	for _, teamID := range teamIDs {
		broker, err := database.GetBroker(teamID)
		if err != nil {
			log.Panicln(err)
		}
		msg, err := json.Marshal(startMessage)
		if err != nil {
			log.Panicln(err)
		}
		err = broker.Publish(msg)
		if err != nil {
			log.Panicln(err)
		}
	}

	cache.UpdateGame(gameIDarray[0], map[string]interface{}{
		"status": "in-progress",
	})

	go endGameHandler(cache, gameIDarray[0], teamIDs, timeLimit)

	successMessage := fmt.Sprintf("Started Game %s", gameIDarray[0])
	fmt.Fprintln(w, models.CreateHTTPResponse(nil, successMessage, true).ToJSON())
}

func endGameHandler(cache database.CacheDatabase, gameID string, teamIDs []string, timeLimit int) {
	time.Sleep(time.Second * time.Duration(timeLimit))
	endMessage := models.EndMessage{
		Type: "end",
	}
	for _, teamID := range teamIDs {
		broker, err := database.GetBroker(teamID)
		if err != nil {
			log.Panicln(err)
		}
		msg, err := json.Marshal(endMessage)
		if err != nil {
			log.Panicln(err)
		}
		err = broker.Publish(msg)
		if err != nil {
			log.Panicln(err)
		}
	}
	cache.UpdateGame(gameID, map[string]interface{}{
		"status": "complete",
	})
}

// JoinGame attempts to upgrade the connection into a websocket and initiates GamePlay logic
func JoinGame(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln(err)
	}
	defer ws.Close()

	teamIDarray, ok := r.URL.Query()["teamID"]
	if !ok || len(teamIDarray) < 1 {
		log.Panicln("No query parameter 'teamID' specified")
	}

	PlayGame(ws, teamIDarray[0])
}

// CreateGame will create a game with the specified configuration
func CreateGame(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	// Check admin password
	verifyPassword(r)

	db := database.GetCacheDatabase()

	game := models.CreateGame{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln("Could not read the body of the message")
	}
	log.Printf("Received body: %s\n", string(body))
	err = json.Unmarshal(body, &game)
	if err != nil {
		log.Panicln("Could not unmarshall body into game object")
	}
	log.Printf("This is the json %v\n", game)

	// TODO use graph database to get a random root node
	// from the graph for the start node
	game.StartNode = "startNode123"
	gameID := db.CreateGame(game)
	fmt.Fprintf(w, models.CreateHTTPResponse(nil, map[string]interface{}{"gameID": gameID}, true).ToJSON())
}

// DeleteGame will delete the game with the matching id
func DeleteGame(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	// Check admin password
	verifyPassword(r)

	db := database.GetCacheDatabase()

	gameIDs, ok := r.URL.Query()["gameID"]
	if !ok || len(gameIDs) < 1 {
		log.Panicln("No query parameter 'gameID' specified")
	}

	result := db.DeleteGame(gameIDs[0])
	fmt.Fprintln(w, models.CreateHTTPResponse(nil, gameIDs[0], result).ToJSON())
}

// ListGames will return an array of games
// For convenience I will return everything if the no fields are specified
func ListGames(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()
	// TODO Will need to have special handling if the string Teams is specified in field

	fieldParam, ok := r.URL.Query()["fields"]
	if !ok || len(fieldParam) < 1 {
		log.Panicln("No query parameter 'fields' specified")
	}

	fields := strings.Split(fieldParam[0], ",")

	db := database.GetCacheDatabase()

	if len(fields) == 0 {
		games := db.GetGames([]string{"gameID", "name", "startNode", "timeLimit", "status", "startTime", "teams"})
		log.Printf("Here are the games%v\n", games)
		fmt.Fprintln(w, models.CreateHTTPResponse(nil, games, true).ToJSON())
	} else {
		games := db.GetGames(fields)
		log.Printf("Here are the games: %v\n", games)
		fmt.Fprintln(w, models.CreateHTTPResponse(nil, games, true).ToJSON())
	}
}

// GameInfo is like ListGames but for a single game id in the query params
func GameInfo(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	fieldParam, ok := r.URL.Query()["fields"]
	if !ok || len(fieldParam) < 1 {
		log.Panicln("No query parameter 'fields' specified")
	}

	fields := strings.Split(fieldParam[0], ",")

	gameIDs, ok := r.URL.Query()["gameID"]
	if !ok || len(gameIDs) < 1 {
		log.Panicln("No query parameter 'gameID' specified")
	}

	db := database.GetCacheDatabase()
	game := db.GetGame(fields, gameIDs[0])
	log.Printf("Requested game info: %v\n", game)
	fmt.Fprintln(w, models.CreateHTTPResponse(nil, game, true).ToJSON())
}

// ListTeams is like ListGames except it takes no fields
func ListTeams(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	db := database.GetCacheDatabase()
	teams := db.GetTeams()
	log.Printf("Teams List: %v\n", teams)
	fmt.Fprintln(w, models.CreateHTTPResponse(nil, teams, true).ToJSON())
}

// TeamInfo is like ListTeams but for a provided teamID
func TeamInfo(w http.ResponseWriter, r *http.Request) {
	// TODO consider refactoring to use recover package
	defer func() {
		if recovery := recover(); recovery != nil {
			log.Println(recovery)
			fmt.Fprintln(w, models.CreateHTTPResponse(recovery, nil, false).ToJSON())
		}
	}()

	teamIDs, ok := r.URL.Query()["teamID"]
	if !ok || len(teamIDs) < 1 {
		log.Panicln("No query parameter 'gameID' specified")
	}

	db := database.GetCacheDatabase()
	team := db.GetTeam(teamIDs[0])
	log.Printf("Team info: %v\n", team)
	fmt.Fprintln(w, models.CreateHTTPResponse(nil, team, true).ToJSON())
}

func verifyPassword(r *http.Request) {
	adminPassword := r.Header["Adminpassword"] // TODO: Need to capitalize P and hopefully it'll still work, Postman sends this
	if len(adminPassword) != 1 {
		log.Panicln("Malformed header 'AdminPassword'")
	} else if adminPassword[0] != models.Config.Wordnet.AdminPassword {
		log.Panicln("Incorrect Admin Password")
	}
}
