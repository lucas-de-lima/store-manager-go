package repository

import (
	"database/sql"
	"fmt"

	"github.com/lucas-de-lima/store-manager-go/model"
)

type SaleRepository struct {
	db *sql.DB
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) Create(saleProducts []model.SaleProduct) (*model.Sale, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO StoreManager.sales (date) VALUES (NOW())")
	if err != nil {
		return nil, err
	}
	saleID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO StoreManager.sales_products (sale_id, product_id, quantity) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, sp := range saleProducts {
		if _, err := stmt.Exec(saleID, sp.ProductID, sp.Quantity); err != nil {
			return nil, fmt.Errorf("failed to insert sale product: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Sale{ID: int(saleID), ItemsSold: saleProducts}, nil
}