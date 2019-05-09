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

// CacheDatabase is the interface for normal database
type CacheDatabase interface {
	CreateGame(game models.CreateGame) string // the game id for the game just created
	GetGames() []models.CacheGame
	GetTeams(gameID string) []models.Team
	DeleteGame(gameID string)
	SetupDB()
}

// Graph is an interface abstraction for a graph store
type Graph interface {
}

// GetCacheDatabase returns the default cache database type
func GetCacheDatabase() CacheDatabase {
	return RedisDatabase{}
}
