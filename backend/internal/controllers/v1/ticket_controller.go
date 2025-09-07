package v1_controllers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"backend/utils"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type TicketControllerV1 interface {
	CreateTicket(c *gin.Context)
	UpdateTicket(c *gin.Context)
	DeleteTicket(c *gin.Context)
}

type ticketControllerV1 struct {
	ticketService *services.TicketService
}

func (t *ticketControllerV1) CreateTicket(c *gin.Context) {
	ctx := context.Background()
	ticketDto := dto.CreateTicketDto{}

	if err := c.ShouldBindJSON(&ticketDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	ticket, err := t.ticketService.CreateTicket(ctx, ticketDto)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", gin.H{"ticket": ticket}))
}

func (t *ticketControllerV1) UpdateTicket(c *gin.Context) {
	ctx := context.Background()
	ticketDto := dto.UpdateTicketDto{}

	if err := c.ShouldBindJSON(&ticketDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	ticketIDParam := c.Param("id")
	var ticketID pgtype.UUID
	if err := ticketID.Scan(ticketIDParam); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	ticket, err := t.ticketService.UpdateTicket(ctx, ticketID, ticketDto)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", gin.H{"ticket": ticket}))
}

func (t *ticketControllerV1) DeleteTicket(c *gin.Context) {
	ctx := context.Background()
	ticketIDParam := c.Param("id")
	var ticketID pgtype.UUID
	if err := ticketID.Scan(ticketIDParam); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := t.ticketService.DeleteTicket(ctx, ticketID); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "Ticket Deleted"))
}

func NewTicketControllerV1() TicketControllerV1 {
	return &ticketControllerV1{
		ticketService: services.NewTicketService(),
	}
}
