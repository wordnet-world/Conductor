package database

import (
	"github.com/wordnet-world/Conductor/models"
)

// Database is the interface for normal database
type CacheDatabase interface {
	CreateGame(game models.CreateGame) string // the game id for the game just created
	GetGames() []models.CacheGame
	GetTeams(gameID string) []models.Team
	DeleteGame(gameID string)
	SetupDB()
}

// GetDatabase returns the default database type
func GetCacheDatabase() CacheDatabase {
	return RedisDatabase{}
}
