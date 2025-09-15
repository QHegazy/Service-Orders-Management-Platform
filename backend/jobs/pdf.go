package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jung-kurt/gofpdf"
)

const TypePDFInvoice = "pdf:invoice"

type InvoicePayload struct {
	TicketID string  `json:"ticket_id"`
	Amount   float64 `json:"amount"`
}

var (
	jobLogger *log.Logger
	logFile   *os.File
	once      sync.Once
)

func InitLogger(path string) error {
	var err error
	once.Do(func() {
		if err = os.MkdirAll("logs", 0755); err != nil {
			return
		}
		logFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}
		jobLogger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	})
	return err
}

func CloseLogger() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

func EnqueueInvoice(client *asynq.Client, ticketID string, amount float64) error {
	payload := InvoicePayload{
		TicketID: ticketID,
		Amount:   amount,
	}

	if jobLogger != nil {
		jobLogger.Printf("ENQUEUE: Enqueue PDF invoice job for ticket %s amount=%.2f", ticketID, amount)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		if jobLogger != nil {
			jobLogger.Printf("ENQUEUE_ERROR: marshal payload for ticket %s: %v", ticketID, err)
		}
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypePDFInvoice, payloadBytes)
	_, err = client.Enqueue(task, asynq.MaxRetry(3))
	if err != nil {
		if jobLogger != nil {
			jobLogger.Printf("ENQUEUE_ERROR: enqueue task for ticket %s: %v", ticketID, err)
		}
		return err
	}

	if jobLogger != nil {
		jobLogger.Printf("ENQUEUE_SUCCESS: enqueued PDF invoice job for ticket %s", ticketID)
	}
	return nil
}

func HandlePDFTask(ctx context.Context, t *asynq.Task) error {
	startTime := time.Now()
	if jobLogger != nil {
		jobLogger.Printf("TASK_START: %s at %s", TypePDFInvoice, startTime.Format(time.RFC3339))
	}

	var payload InvoicePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		if jobLogger != nil {
			jobLogger.Printf("TASK_ERROR: unmarshal payload: %v", err)
		}
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if jobLogger != nil {
		jobLogger.Printf("TASK_PROCESSING: ticket=%s amount=%.2f", payload.TicketID, payload.Amount)
	}

	// Ensure invoices directory exists
	if err := os.MkdirAll("./invoices", 0755); err != nil {
		if jobLogger != nil {
			jobLogger.Printf("TASK_ERROR: create invoices dir: %v", err)
		}
		return fmt.Errorf("failed to create invoices directory: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, fmt.Sprintf("Invoice for Ticket %s", payload.TicketID), "", 1, "L", false, 0, "")
	pdf.Ln(4)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 8, fmt.Sprintf("Amount: $%.2f", payload.Amount), "", 1, "L", false, 0, "")

	filePath := fmt.Sprintf("invoices/%s.pdf", payload.TicketID)
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		duration := time.Since(startTime)
		if jobLogger != nil {
			jobLogger.Printf("TASK_ERROR: failed to generate PDF for ticket %s after %v: %v", payload.TicketID, duration, err)
		}
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	duration := time.Since(startTime)
	if jobLogger != nil {
		jobLogger.Printf("TASK_SUCCESS: generated invoice %s for ticket %s in %v", filePath, payload.TicketID, duration)
	}
	return nil
}
