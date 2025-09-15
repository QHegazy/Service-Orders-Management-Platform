package v1_controllers

import (
	"net/http"

	"backend/internal/dto"
	"backend/internal/services"
	"backend/jobs"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type InvoiceControllerV1 interface {
	CreateInvoice(c *gin.Context)
}

type invoiceControllerV1 struct {
	invoiceService *services.InvoiceService
	asynqClient    *asynq.Client
}

func (i *invoiceControllerV1) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid_request", err.Error()))
		return
	}

	client := i.asynqClient
	var createdClient bool
	if client == nil {
		client = asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
		createdClient = true
	}

	// Ensure we close the client we created here (don't close injected client).
	if createdClient {
		defer client.Close()
	}
	if err := jobs.EnqueueInvoice(client, req.TicketID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("enqueue_failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("success", "Invoice created and PDF generation enqueued"))
}

func NewInvoiceControllerV1() InvoiceControllerV1 {
	return &invoiceControllerV1{
		invoiceService: services.NewInvoiceService(),
	}
}
