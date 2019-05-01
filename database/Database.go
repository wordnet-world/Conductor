package database

import (
	"github.com/wordnet-world/Conductor/models"
)

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
