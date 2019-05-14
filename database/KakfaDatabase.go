package database

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wordnet-world/Conductor/models"
)

// KafkaBroker is a Broker implementation with Pub/Sub abilities
type KafkaBroker struct {
	producer  *kafka.Producer
	consumer  *kafka.Consumer
	connected bool
	topic     string
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

// Connect establishes a producer
func (broker *KafkaBroker) connect() error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": models.Config.Kafka.Address, "queue.buffering.max.messages": "5", "queue.buffering.max.ms": "300"})
	if err != nil {
		return err
	}

	defer p.Close()
	broker.producer = p
	return nil
}

// Publish uses the kafka producer to publish a message
// cannot be used if connect has not been called
func (broker *KafkaBroker) Publish(message string) error {
	log.Println("Right before is connected")
	if !broker.connected {
		return new(invalidStateError)
	}

	log.Println("Right before Produce")

	broker.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &broker.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	log.Println("Right after produce before flush")

	broker.producer.Flush(15 * 1000)

	log.Println("Right after flush")
	return nil
}

// Subscribe subscribes to a kafka topic using a consumer. Will call the action func
// with whatever message was received everytime consumer consumes
func (broker *KafkaBroker) Subscribe(consumerID string, action func(string)) error {
	log.Println("Right before is connected sub")
	if !broker.connected {
		return new(invalidStateError)
	}

	log.Println("right before consumer creation sub")
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"group.id":          consumerID,
		"bootstrap.servers": models.Config.Kafka.Address,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Right before subscribe to topics")

	c.SubscribeTopics([]string{broker.topic}, nil)
	for {
		log.Println("ReadMessage in sub")
		msg, err := c.ReadMessage(-1)
		log.Printf("Sub message %s\nError: %v\n", msg, err)
		if err == nil {
			action(string(msg.Value))
		}
	}
}

type invalidStateError struct {
	message string
}

func (e *invalidStateError) Error() string {
	return e.message
}
