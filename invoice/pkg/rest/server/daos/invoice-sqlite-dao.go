package daos

import (
	"database/sql"
	"errors"
	"github.com/shreya-intelops/procurement-solution/invoice/pkg/rest/server/daos/clients/sqls"
	"github.com/shreya-intelops/procurement-solution/invoice/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type InvoiceDao struct {
	sqlClient *sqls.SQLiteClient
}

func migrateInvoices(r *sqls.SQLiteClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS invoices(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        
		Amount REAL NOT NULL,
        CONSTRAINT id_unique_key UNIQUE (Id)
	)
	`
	_, err1 := r.DB.Exec(query)
	return err1
}

func NewInvoiceDao() (*InvoiceDao, error) {
	sqlClient, err := sqls.InitSqliteDB()
	if err != nil {
		return nil, err
	}
	err = migrateInvoices(sqlClient)
	if err != nil {
		return nil, err
	}
	return &InvoiceDao{
		sqlClient,
	}, nil
}

func (invoiceDao *InvoiceDao) CreateInvoice(m *models.Invoice) (*models.Invoice, error) {
	insertQuery := "INSERT INTO invoices(Amount)values(?)"
	res, err := invoiceDao.sqlClient.DB.Exec(insertQuery, m.Amount)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id

	log.Debugf("invoice created")
	return m, nil
}

func (invoiceDao *InvoiceDao) ListInvoices() ([]*models.Invoice, error) {
	selectQuery := "SELECT * FROM invoices"
	rows, err := invoiceDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var invoices []*models.Invoice
	for rows.Next() {
		m := models.Invoice{}
		if err = rows.Scan(&m.Id, &m.Amount); err != nil {
			return nil, err
		}
		invoices = append(invoices, &m)
	}
	if invoices == nil {
		invoices = []*models.Invoice{}
	}

	log.Debugf("invoice listed")
	return invoices, nil
}

func (invoiceDao *InvoiceDao) GetInvoice(id int64) (*models.Invoice, error) {
	selectQuery := "SELECT * FROM invoices WHERE Id = ?"
	row := invoiceDao.sqlClient.DB.QueryRow(selectQuery, id)
	m := models.Invoice{}
	if err := row.Scan(&m.Id, &m.Amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	log.Debugf("invoice retrieved")
	return &m, nil
}

func (invoiceDao *InvoiceDao) UpdateInvoice(id int64, m *models.Invoice) (*models.Invoice, error) {
	if id == 0 {
		return nil, errors.New("invalid invoice ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	invoice, err := invoiceDao.GetInvoice(id)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE invoices SET Amount = ? WHERE Id = ?"
	res, err := invoiceDao.sqlClient.DB.Exec(updateQuery, m.Amount, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sqls.ErrUpdateFailed
	}

	log.Debugf("invoice updated")
	return m, nil
}

func (invoiceDao *InvoiceDao) DeleteInvoice(id int64) error {
	deleteQuery := "DELETE FROM invoices WHERE Id = ?"
	res, err := invoiceDao.sqlClient.DB.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sqls.ErrDeleteFailed
	}

	log.Debugf("invoice deleted")
	return nil
}
