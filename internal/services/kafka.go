package services

import (
	"log"

	"edi-gateway/internal/config"
	"github.com/segmentio/kafka-go"
)

// InitKafka initializes the Kafka writer using the provided configuration.
func InitKafka(cfg config.Config) *kafka.Writer {
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{cfg.KafkaBroker},
		Topic:   "edi_topic",
	})
	log.Println("Kafka writer initialized successfully.")
	return kafkaWriter
}