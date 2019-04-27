package database

import (
	"github.com/wordnet-world/Conductor/models"
)

// Database is the interface for normal database
type Database interface {
	TestConnection() bool
	CreateGame() error
	GetGames() []models.Game
	GetTeams(game models.Game) []models.Team
}
