package models

// Game is a model for a single game instance
type Game struct {
	Name      string   `json:"name"`
	ID        string   `json:"gameId"`
	StartNode string   `json:"startNode"` // probably unique id, will need a function that returns me this when querying neo4j
	TimeLimit int      `json:"timeLimit"` // probably in seconds or minutes
	Teams     []string `json:"teams"`
}

// Team is a model of a single team in a Game
type Team struct {
	Name  string `json:"name"`
	ID    string `json:"teamid"`
	Score string `json:"score"`
}
