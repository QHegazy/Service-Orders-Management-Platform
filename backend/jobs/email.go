package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jung-kurt/gofpdf"
)

const TypePDFInvoice = "pdf:invoice"

type InvoicePayload struct {
	UserID string
	Amount float64
}

var jobLogger *log.Logger

func init() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
		return
	}

	// Open or create the job log file
	logFile, err := os.OpenFile("logs/jobs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open job log file: %v", err)
		return
	}

	// Create logger with timestamp
	jobLogger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}

func EnqueueInvoice(client *asynq.Client, userID string, amount float64) error {
	payload := InvoicePayload{
		UserID: userID,
		Amount: amount,
	}

	if jobLogger != nil {
		jobLogger.Printf("ENQUEUE: Starting to enqueue PDF invoice job for user %s, amount $%.2f", userID, amount)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		if jobLogger != nil {
			jobLogger.Printf("ENQUEUE_ERROR: Failed to marshal payload for user %s: %v", userID, err)
		}
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypePDFInvoice, payloadBytes)
	_, err = client.Enqueue(task, asynq.MaxRetry(3))

	if err != nil {
		if jobLogger != nil {
			jobLogger.Printf("ENQUEUE_ERROR: Failed to enqueue task for user %s: %v", userID, err)
		}
		return err
	}

	if jobLogger != nil {
		jobLogger.Printf("ENQUEUE_SUCCESS: Successfully enqueued PDF invoice job for user %s", userID)
	}

	return nil
}

func HandlePDFTask(ctx context.Context, t *asynq.Task) error {
	startTime := time.Now()

	if jobLogger != nil {
		jobLogger.Printf("TASK_START: PDF invoice task started at %s", startTime.Format(time.RFC3339))
	}

	var payload InvoicePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		if jobLogger != nil {
			jobLogger.Printf("TASK_ERROR: Failed to unmarshal payload: %v", err)
		}
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if jobLogger != nil {
		jobLogger.Printf("TASK_PROCESSING: Processing PDF invoice for user %s, amount $%.2f", payload.UserID, payload.Amount)
	}

	// Ensure invoices directory exists
	if err := os.MkdirAll("invoices", 0755); err != nil {
		if jobLogger != nil {
			jobLogger.Printf("TASK_ERROR: Failed to create invoices directory for user %s: %v", payload.UserID, err)
		}
		return fmt.Errorf("failed to create invoices directory: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice for User %s", payload.UserID))
	pdf.Ln(12)
	pdf.Cell(40, 10, fmt.Sprintf("Amount: $%.2f", payload.Amount))

	filePath := fmt.Sprintf("invoices/%s.pdf", payload.UserID)
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		duration := time.Since(startTime)
		if jobLogger != nil {
			jobLogger.Printf("TASK_ERROR: Failed to generate PDF for user %s after %v: %v", payload.UserID, duration, err)
		}
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	duration := time.Since(startTime)
	if jobLogger != nil {
		jobLogger.Printf("TASK_SUCCESS: Successfully generated invoice %s for user %s in %v", filePath, payload.UserID, duration)
	}

	fmt.Printf("📄 Generated invoice: %s\n", filePath)
	return nil
}
