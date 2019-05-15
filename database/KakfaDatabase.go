package database

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/wordnet-world/Conductor/models"
)

// KafkaBroker is a Broker implementation with Pub/Sub abilities
type KafkaBroker struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
	topic    string
}

// NewKafkaBroker is a Constructor which attempts to connect to the kafka broker
func NewKafkaBroker(topic string) (*KafkaBroker, error) {
	p := new(KafkaBroker)
	p.topic = topic
	return p, nil
}

// Publish uses the kafka producer to publish a message
// cannot be used if connect has not been called
func (broker *KafkaBroker) Publish(message string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":            models.Config.Kafka.Address,
		"queue.buffering.max.messages": "5",
		"queue.buffering.max.ms":       "300",
	})

	deliveryChan := make(chan kafka.Event)

	fmt.Println(message)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &broker.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		return err
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	p.Flush(15 * 1000)
	p.Close()
	return nil
}

// Subscribe subscribes to a kafka topic using a consumer. Will call the action func
// with whatever message was received everytime consumer consumes
func (broker *KafkaBroker) Subscribe(consumerID string, action func(string)) error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"group.id":          consumerID,
		"bootstrap.servers": models.Config.Kafka.Address,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Println(err)
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
