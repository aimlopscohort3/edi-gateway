package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// EDIWorkflow orchestrates the processing of an EDI transaction.
func EDIWorkflow(ctx workflow.Context, transaction Transaction) error {
	// Initialize Temporal workflow logger
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting EDI workflow", "TransactionID", transaction.ID)

	// Define activity options with a retry policy
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5, // Timeout for activity execution
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 2, // Initial retry interval
			BackoffCoefficient: 2.0,             // Exponential backoff factor
			MaximumInterval:    time.Minute,     // Maximum retry interval
			MaximumAttempts:    3,               // Maximum retry attempts
		},
	}

	// Attach activity options to the workflow context
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Step 1: Save transaction to the database
	logger.Info("Executing SaveToDatabaseActivity", "TransactionID", transaction.ID)
	err := workflow.ExecuteActivity(ctx, SaveToDatabaseActivity, transaction).Get(ctx, nil)
	if err != nil {
		logger.Error("SaveToDatabaseActivity failed", "Error", err, "TransactionID", transaction.ID)
		return err
	}
	logger.Info("SaveToDatabaseActivity completed successfully", "TransactionID", transaction.ID)

	// Step 2: Publish transaction to Kafka
	logger.Info("Executing PublishToKafkaActivity", "TransactionID", transaction.ID)
	err = workflow.ExecuteActivity(ctx, PublishToKafkaActivity, transaction).Get(ctx, nil)
	if err != nil {
		logger.Error("PublishToKafkaActivity failed", "Error", err, "TransactionID", transaction.ID)
		return err
	}
	logger.Info("PublishToKafkaActivity completed successfully", "TransactionID", transaction.ID)

	// Workflow completed successfully
	logger.Info("EDI workflow completed successfully", "TransactionID", transaction.ID)
	return nil
}