package database

import (
	"github.com/wordnet-world/Conductor/models"
)

// Broker is an interface which allows you to Publish a message
// and subscribe to a particular topic with an action to take
type Broker interface {
	Publish(message []byte) error
	Subscribe(consumerID string, action func(string)) error
}

// CacheDatabase is the interface for normal database
type CacheDatabase interface {
	CreateGame(game models.CreateGame) string // the game id for the game just created
	UpdateGame(gameID string, updates map[string]interface{})
	GetGames(fields []string) []map[string]interface{}
	GetGame(fields []string, gameID string) map[string]interface{}
	GetTeams() []models.Team
	GetTeam(teamID string) models.Team
	DeleteGame(gameID string) bool
	GetConsumerID() string
	SetupTeamCaches(teamIDs []string, root models.Node, neighbors []models.Node)
	IsFound(guess string, teamID string) bool
	IsPeriphery(guess string, teamID string) int64
	UpdateCache(newNode models.Node, neighbors []models.Node, teamID string) ([]models.Node, []models.Node)
	SetupDB()
}

// GraphDatabase is an interface abstraction for a graph store
type GraphDatabase interface {
	Connect(uri, username, password string) error
	PopulateDummy(uri, username, password string) error
	Close()
	GetNeighbors(models.Node) ([]models.Node, error)
	GetNeighborsNodeID(nodeID int64) ([]models.Node, error)
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

// GetBroker returns an implementation of the Broker Interface
func GetBroker(topic string) (Broker, error) {
	return NewKafkaBroker(topic)
}
