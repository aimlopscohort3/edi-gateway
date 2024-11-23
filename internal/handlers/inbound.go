package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"edi-gateway/internal/workflows"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB, kafkaWriter *kafka.Writer, temporalClient client.Client) {
	// Inbound route
	r.HandleFunc("/inbound", func(w http.ResponseWriter, r *http.Request) {
		var transaction workflows.Transaction
		if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		transaction.ID = uuid.New().String()
		transaction.Date = time.Now()
		transaction.Status = "Processed"

		workflowOptions := client.StartWorkflowOptions{
			ID:        "edi-workflow-" + transaction.ID,
			TaskQueue: "EDI_TASK_QUEUE",
		}

		_, err := temporalClient.ExecuteWorkflow(r.Context(), workflowOptions, "EDIWorkflow", transaction)
		if err != nil {
			log.Printf("Failed to start workflow: %v", err)
			http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
			return
		}

		// Publish event to Kafka
		event, _ := json.Marshal(transaction)
		err = kafkaWriter.WriteMessages(r.Context(), kafka.Message{Key: []byte(transaction.ID), Value: event})
		if err != nil {
			log.Printf("Failed to publish Kafka message: %v", err)
			http.Error(w, "Failed to publish Kafka message", http.StatusInternalServerError)
			return
		}

		log.Printf("Workflow started for transaction: %s", transaction.ID)
		fmt.Fprintf(w, "Transaction received: %s\\n", transaction.ID)
	}).Methods("POST")
}