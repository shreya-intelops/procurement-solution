package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shreya-intelops/procurement-solution/invoice/pkg/rest/server/daos/clients/sqls"
	"github.com/shreya-intelops/procurement-solution/invoice/pkg/rest/server/models"
	"github.com/shreya-intelops/procurement-solution/invoice/pkg/rest/server/services"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"strconv"
)

type InvoiceController struct {
	invoiceService *services.InvoiceService
}

func NewInvoiceController() (*InvoiceController, error) {
	invoiceService, err := services.NewInvoiceService()
	if err != nil {
		return nil, err
	}
	return &InvoiceController{
		invoiceService: invoiceService,
	}, nil
}

func (invoiceController *InvoiceController) CreateInvoice(context *gin.Context) {
	// validate input
	var input models.Invoice
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// trigger invoice creation
	invoiceCreated, err := invoiceController.invoiceService.CreateInvoice(&input)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, invoiceCreated)
}

func (invoiceController *InvoiceController) ListInvoices(context *gin.Context) {
	// trigger all invoices fetching
	invoices, err := invoiceController.invoiceService.ListInvoices()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, invoices)
}

func (invoiceController *InvoiceController) FetchInvoice(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger invoice fetching
	invoice, err := invoiceController.invoiceService.GetInvoice(id)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request.Context())
		currentSpan.SetAttributes(attribute.String("invoice.id", strconv.FormatInt(invoice.Id, 10)))
	}

	context.JSON(http.StatusOK, invoice)
}

func (invoiceController *InvoiceController) UpdateInvoice(context *gin.Context) {
	// validate input
	var input models.Invoice
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger invoice update
	if _, err := invoiceController.invoiceService.UpdateInvoice(id, &input); err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (invoiceController *InvoiceController) DeleteInvoice(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger invoice deletion
	if err := invoiceController.invoiceService.DeleteInvoice(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}
