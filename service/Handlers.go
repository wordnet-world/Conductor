package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

// JoinGame attempts to upgrade the connection into a websocket and initiates GamePlay logic
func JoinGame(w http.ResponseWriter, r *http.Request) {
	// TODO need to pick a team for this connection, probably through url parameters
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// Incredibly stupid connection that just echos the response with a slight addition
	for {
		var msg models.WordGuess
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		fmt.Println(msg)
		msg.Guess = msg.Guess + " but a Response~"
		ws.WriteJSON(msg)
	}
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
	// TODO Will need to have special handling if the string Teams is specified in fields

	// Check admin password
	verifyPassword(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln("Could not read the body of the message")
	}
	log.Printf("Received body: %s\n", string(body))

	fields := models.ListGameFields{}
	err = json.Unmarshal(body, &fields)
	if err != nil {
		log.Panicln("Could not unmarshall body into fields object")
	}
	log.Printf("This is the json %v\n", fields)

	db := database.GetCacheDatabase()

	if len(fields.Fields) == 0 {
		games := db.GetGames([]string{"gameID", "name", "startNode", "timeLimit", "status", "startTime", "teams"})
		log.Printf("Here are the games%v\n", games)
		fmt.Fprintln(w, models.CreateHTTPResponse(nil, games, true).ToJSON())
	} else {
		games := db.GetGames(fields.Fields)
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
	// TODO Will need to have special handling if the string Teams is specified in fields

	// Check admin password
	verifyPassword(r)
}

func verifyPassword(r *http.Request) {
	adminPassword := r.Header["Adminpassword"] // TODO: Need to capitalize P and hopefully it'll still work, Postman sends this
	if len(adminPassword) != 1 {
		log.Panicln("Malformed header 'AdminPassword'")
	} else if adminPassword[0] != models.Config.Wordnet.AdminPassword {
		log.Panicln("Incorrect Admin Password")
	}
}
