package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wordnet-world/Conductor/database"

	"github.com/gorilla/websocket"
	"github.com/wordnet-world/Conductor/models"
)

// PlayGame initiates the basic gameplay loop
func PlayGame(ws *websocket.Conn, teamID string) {

	db := database.GetCacheDatabase()
	consumerID := db.GetConsumerID()

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
