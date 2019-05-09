package database

import (
	"github.com/wordnet-world/Conductor/models"
)

// Broker is an interface which allows you to Publish a message
// and subscribe to a particular topic with an action to take
type Broker interface {
	Connect()
	Publish(message string)
	Subscribe(topic string, action func(string))
}

// Database is the interface for normal database
type Database interface {
	CreateGame(game models.Game) string // the game id for the game just created
	GetGames() []models.Game
	GetTeams(gameID string) []models.Team
	DeleteGame(gameID string)
	SetupDB()
}

// GetDatabase returns the default database type
func GetDatabase() Database {
	return RedisDatabase{}
}
