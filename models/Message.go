package models

// Message represents a message to go through the Broker and the WebSocket
type Message struct {
	Text string `json:"text"`
}
