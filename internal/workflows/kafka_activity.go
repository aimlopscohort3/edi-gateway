// Placeholder for kafka_activity.go
package workflows

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// PublishToKafkaActivity publishes a transaction to a Kafka topic.
func PublishToKafkaActivity(ctx context.Context, kafkaWriter *kafka.Writer, transaction Transaction) error {
	log.Printf("Publishing transaction %s to Kafka topic: %s", transaction.ID, kafkaWriter.Topic)

	// Marshal the transaction into JSON
	event, err := json.Marshal(transaction)
	if err != nil {
		log.Printf("Error marshalling transaction %s: %v", transaction.ID, err)
		return err
	}

	// Publish the message to Kafka
	err = kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(transaction.ID),
		Value: event,
	})
	if err != nil {
		log.Printf("Error publishing transaction %s to Kafka: %v", transaction.ID, err)
		return err
	}

	log.Printf("Transaction %s published successfully to Kafka", transaction.ID)
	return nil
}