package v1_routes

import (
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func invoiceRoutes(r *gin.RouterGroup) {

	invoice := r.Group("/invoice")
	invoice.Use(middleware.ValidationErrorHandler())
	invoice.Use(middleware.DBErrorHandler())
	// invoiceController := v1_controllers.NewInvoiceControllerV1()
	// invoice.POST("", invoiceController.CreateInvoice)
	// invoice.PUT("", invoiceController.UpdateInvoice)
	// invoice.DELETE("/:id", invoiceController.DeleteInvoice)
	// invoice.GET("/:id", invoiceController.GetInvoiceByID)
	// invoice.GET("", invoiceController.ListAllInvoices)
	// invoice.POST("/payment", invoiceController.CreatePayment)
	// invoice.GET("/payment/:id", invoiceController.GetPaymentByID)
	// invoice.GET("/payments", invoiceController.ListAllPayments)

}
