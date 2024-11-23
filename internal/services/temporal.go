package services

import (
	"log"

	"edi-gateway/internal/config"
	"edi-gateway/internal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func InitTemporal(cfg config.Config) (client.Client, worker.Worker) {
	temporalClient, err := client.Dial(client.Options{HostPort: cfg.TemporalHostPort})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	log.Println("Temporal client initialized successfully.")

	temporalWorker := worker.New(temporalClient, "EDI_TASK_QUEUE", worker.Options{})

	// Register workflows and activities
	temporalWorker.RegisterWorkflow(workflows.EDIWorkflow)
	temporalWorker.RegisterActivity(workflows.SaveToDatabaseActivity)
	temporalWorker.RegisterActivity(workflows.PublishToKafkaActivity)

	log.Println("Temporal worker initialized successfully.")
	return temporalClient, temporalWorker
}