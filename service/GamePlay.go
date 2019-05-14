package service

import (
	"fmt"
	"log"

	"github.com/wordnet-world/Conductor/database"

	"github.com/gorilla/websocket"
	"github.com/wordnet-world/Conductor/models"
)

// PlayGame initiates the basic gameplay loop
func PlayGame(ws *websocket.Conn, teamID string) {
	_, err := database.GetBroker(teamID) // _ is the broker
	if err != nil {
		log.Panicln(err)
	}

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
