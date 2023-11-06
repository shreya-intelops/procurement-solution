package services

import (
	"github.com/shreya-intelops/procurement-solution/invoice/pkg/rest/server/daos"
	"github.com/shreya-intelops/procurement-solution/invoice/pkg/rest/server/models"
)

type InvoiceService struct {
	invoiceDao *daos.InvoiceDao
}

func NewInvoiceService() (*InvoiceService, error) {
	invoiceDao, err := daos.NewInvoiceDao()
	if err != nil {
		return nil, err
	}
	return &InvoiceService{
		invoiceDao: invoiceDao,
	}, nil
}

func (invoiceService *InvoiceService) CreateInvoice(invoice *models.Invoice) (*models.Invoice, error) {
	return invoiceService.invoiceDao.CreateInvoice(invoice)
}

func (invoiceService *InvoiceService) ListInvoices() ([]*models.Invoice, error) {
	return invoiceService.invoiceDao.ListInvoices()
}

func (invoiceService *InvoiceService) GetInvoice(id int64) (*models.Invoice, error) {
	return invoiceService.invoiceDao.GetInvoice(id)
}

func (invoiceService *InvoiceService) UpdateInvoice(id int64, invoice *models.Invoice) (*models.Invoice, error) {
	return invoiceService.invoiceDao.UpdateInvoice(id, invoice)
}

func (invoiceService *InvoiceService) DeleteInvoice(id int64) error {
	return invoiceService.invoiceDao.DeleteInvoice(id)
}
