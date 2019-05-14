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

	consumer, err := database.GetBroker(teamID) // _ is the broker
	if err != nil {
		log.Panicln(err)
	}
	log.Println(consumer)

	go consumer.Subscribe(consumerID, createConsumerFunction(ws))

	producer, err := database.GetBroker(teamID)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(producer)

	for {
		var msg models.WordGuess
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		fmt.Println(msg)
		msg.Guess = msg.Guess + " but a Response~"
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Panicln(err)
		}
		producer.Publish(string(jsonMsg))
		//ws.WriteJSON(msg)
	}
}

func createConsumerFunction(ws *websocket.Conn) func(string) {
	return func(message string) {
		log.Printf("handler message:%s\n", message)
		ws.WriteJSON(message)
	}
}
