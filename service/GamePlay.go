package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

		graphUpdate := processGuess(msg, teamID, cache, graph)

		jsonMsg, err := json.Marshal(graphUpdate)
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

func processGuess(msg models.WordGuess, teamID string, cache database.CacheDatabase, graph database.GraphDatabase) models.UpdateMessage {
	updateMessage := models.UpdateMessage{
		Type:             "update",
		Guess:            msg.Guess,
		Correct:          false,
		NewNodeID:        -1,
		NewNodeText:      "",
		NewNodeNeighbors: []int64{},
	}
	if !cache.IsFound(strings.ToLower(msg.Guess), teamID) {
		log.Printf("Guess not in Found, Guess:%s\n", msg.Guess)
		nodeID := cache.IsPeriphery(strings.ToLower(msg.Guess), teamID)
		if nodeID != -1 {
			node := models.Node{
				ID:   nodeID,
				Text: graph.GetNodeText(nodeID),
			}
			log.Printf("Node found in periphery: %v\n", node)
			neighbors, err := graph.GetNeighborsNodeID(nodeID)
			if err != nil {
				log.Panicln(err)
			}

			neighborIDs := make([]int64, len(neighbors))
			for i, n := range neighbors {
				neighborIDs[i] = n.ID
			}

			log.Printf("Retrieved neighbors from graph: %v\n", neighbors)

			resultNodes, foundNodes := cache.UpdateCache(node, neighbors, teamID)
			log.Printf("ResultNodes: %v\n", resultNodes)
			log.Printf("FoundNodes: %v\n", foundNodes)

			/*if len(foundNodes) > 1 {
				log.Panicln("WE HAVE A CYCLE IN THE GRAPH!!!")
			}*/

			updateMessage = models.UpdateMessage{
				Type:             "update",
				Guess:            msg.Guess,
				Correct:          true,
				NewNodeID:        node.ID,
				NewNodeText:      node.Text,
				NewNodeNeighbors: neighborIDs,
			}
		}
	}
	return updateMessage
}
