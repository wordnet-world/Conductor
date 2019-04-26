package database

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaBroker is a Broker implementation with Pub/Sub abilities
type KafkaBroker struct {
	producer  *kafka.Producer
	consumer  *kafka.Consumer
	connected bool
	topic     string
}

// Connect establishes a producer
func (broker KafkaBroker) connect() error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return err
	}

	defer p.Close()
	broker.producer = p
	return nil
}

// Publish uses the kafka producer to publish a message
// cannot be used if connect has not been called
func (broker KafkaBroker) Publish(message string) error {
	if !broker.connected {
		return new(invalidStateError)
	}

	// Produce messages to topic (asynchronously)
	broker.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &broker.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	broker.producer.Flush(15 * 1000)
	return nil
}

// Subscribe subscribes to a kafka topic using a consumer. Will call the action func
// with whatever message was received everytime consumer consumes
func (broker KafkaBroker) Subscribe(action func(string)) error {
	if !broker.connected {
		return new(invalidStateError)
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		// "group.id":          "myGroup", // need to handle should be different - may work if we don't actually specify it
		"bootstrap.servers": "localhost",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return err
	}

	c.SubscribeTopics([]string{broker.topic}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			action(string(msg.Value))
		}
	}
}

// NewKafkaBroker is a Constructor which attempts to connect to the kafka broker
func NewKafkaBroker(topic string) (*KafkaBroker, error) {
	p := new(KafkaBroker)
	err := p.connect()
	if err != nil {
		return nil, err
	}
	p.topic = topic
	p.connected = true
	return p, nil
}

type invalidStateError struct {
	message string
}

func (e *invalidStateError) Error() string {
	return e.message
}
