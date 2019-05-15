package models

// GraphUpdate represents a request to upgrade the graph in the front end
// NewNode represents the new node to be added to the graph
// ConnectingNode represents the node that the NewNode is connected to
// If ConnectingNode is nil, this is the root node
// Guess represents what was guessed for updating the list of guesses
// Correct represents whether or not the guess was correct - if correct false, only Guess will be populated
type GraphUpdate struct {
	Guess              string `json:"guess"`
	Correct            bool   `json:"correct"`
	NewNodeID          int64  `json:"newNodeId"`
	ConnectingNodeID   int64  `json:"connectingNodeId"`
	NewNodeText        string `json:"newNodeText"`
	ConnectingNodeText string `json:"connectingNodeText"`
	UndiscoveredNodes  int    `json:"undiscoveredNodes"`
}

// WordGuess represents a guess of a word
type WordGuess struct {
	Guess string `json:"guess"`
}
