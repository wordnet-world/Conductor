package models

// Game is a model for a single game instance
type Game struct {
	Name      string
	ID        string
	GraphID   string
	StartNode string // probably unique id
	TimeLimit int    // probably in seconds or minutes
	Teams     []Team
}

// Team is a model of a single team in a Game
type Team struct {
	Name    string
	ID      string
	Score   string
	Players []Player
}

// Player is a model of a single player on a Team in a Game
type Player struct {
	Name  string
	ID    string
	Score int
}
