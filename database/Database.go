package database

// Broker is an interface which allows you to Publish a message
// and subscribe to a particular topic with an action to take
type Broker interface {
	Connect()
	Publish(message string)
	Subscribe(topic string, action func(string))
}
