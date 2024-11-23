package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"edi-gateway/internal/config"
	"edi-gateway/internal/handlers"
	"edi-gateway/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Config{
		DatabaseDSN:      "host=postgres user=postgres password=postgres dbname=edi_gateway port=5432 sslmode=disable",
		KafkaBroker:      "broker:9092",
		TemporalHostPort: "temporal-nginx:7233",
	}

	// Initialize services
	db := services.InitDatabase(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	kafkaWriter := services.InitKafka(cfg)
	defer kafkaWriter.Close()

	temporalClient, temporalWorker := services.InitTemporal(cfg)
	defer temporalClient.Close()

	// Run Temporal worker
	go func() {
		stopCh := make(chan interface{})
		go func() {
			<-services.InterruptCh() // Wait for an OS signal
			close(stopCh)            // Signal the worker to stop
		}()
		if err := temporalWorker.Run(stopCh); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	// Set up HTTP routes
	r := mux.NewRouter()
	handlers.RegisterRoutes(r, db, kafkaWriter, temporalClient)
	r.Handle("/metrics", promhttp.Handler())

	// Handle OS signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server := &http.Server{
		Addr:    ":8086",
		Handler: r,
	}

	go func() {
		log.Println("Server running on port 8086")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down gracefully...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped.")
}