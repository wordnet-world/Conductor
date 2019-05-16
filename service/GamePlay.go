package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/wordnet-world/Conductor/database"

	"github.com/gorilla/websocket"
	"github.com/wordnet-world/Conductor/models"
)

// PlayGame initiates the basic gameplay loop
func PlayGame(ws *websocket.Conn, teamID string) {

	cache := database.GetCacheDatabase()
	graph := database.GetGraphDatabase()
	err := graph.Connect(models.Config.Neo4j.URI, models.Config.Neo4j.Username, models.Config.Neo4j.Password)
	if err != nil {
		log.Panicln(err)
	}
	defer graph.Close()
	consumerID := cache.GetConsumerID()

	broker, err := database.GetBroker(teamID)
	if err != nil {
		log.Panicln(err)
	}

	go broker.Subscribe(consumerID, createConsumerFunction(ws))

	for {
		var msg models.WordGuess
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Panicln(err)
		}
		broker.Publish(jsonMsg)
	}
}

func createConsumerFunction(ws *websocket.Conn) func(string) {
	return func(message string) {
		err := ws.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func processGuess(msg models.WordGuess, teamID string, cache database.CacheDatabase, graph database.GraphDatabase) {
	var graphUpdate models.GraphUpdate
	if cache.IsFound(msg.Guess, teamID) {
		nodeID := cache.IsPeriphery(msg.Guess, teamID)
		if nodeID != -1 {

			neighbors, err := graph.GetNeighborsNodeID(strconv.FormatInt(nodeID, 10))
			if err != nil {
				log.Panicln(err)
			}

		}
	} else {
		graphUpdate = models.GraphUpdate{
			Guess:              msg.Guess,
			Correct:            false,
			NewNodeID:          -1,
			ConnectingNodeID:   -1,
			NewNodeText:        "",
			ConnectingNodeText: "",
			UndiscoveredNodes:  0,
		}
	}
}
