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

// GraphDatabase is an interface abstraction for a graph store
type GraphDatabase interface {
	Connect(uri, username, password string) error
	GetNeighbors(models.Node) ([]models.Node, error)
	GetRoot() (models.Node, error)
}

// GetCacheDatabase returns the default cache database type
func GetCacheDatabase() CacheDatabase {
	return RedisDatabase{}
}

// GetGraphDatabase returns the default graph database
func GetGraphDatabase() GraphDatabase {
	return NewNeo4jDatabase()
}
