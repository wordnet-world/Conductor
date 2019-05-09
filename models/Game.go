package models

// CreateGame is a model for a single game instance for the CreateGame endpoint
type CreateGame struct {
	Name      string   `json:"name"`
	ID        string   `json:"gameID"`
	StartNode string   `json:"startNode"` // probably unique id, will need a function that returns me this when querying neo4j
	TimeLimit int      `json:"timeLimit"` // probably in seconds or minutes
	Teams     []string `json:"teams"`
}

// CacheGame is the model for retreiving games from cache db
type CacheGame struct {
	Name      string   `json:"name"`
	ID        string   `json:"gameID"`
	StartNode string   `json:"startNode"`
	TimeLimit int      `json:"timeLimit"`
	TeamIDs   []string `json:"teamIDs"`
	Status    string   `json:"status"`
	StartTime int      `json:"startTime"` // This should be unix time probably
}

// Team is a model of a single team in a Game
type Team struct {
	Name  string `json:"name"`
	ID    string `json:"teamID"`
	Score string `json:"score"`
}

// ListGameFields is a model for the request to ListGames
type ListGameFields struct {
	Fields []string `json:"fields"`
}
