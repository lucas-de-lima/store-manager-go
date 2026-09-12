package repository

import (
	"database/sql"

	"github.com/lucas-de-lima/store-manager-go/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name FROM StoreManager.products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	row := r.db.QueryRow("SELECT id, name FROM StoreManager.products WHERE id = ?", id)
	var p model.Product
	if err := row.Scan(&p.ID, &p.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Search(q string) ([]model.Product, error) {
	query := "SELECT id, name FROM StoreManager.products WHERE name LIKE ? ORDER BY id"
	rows, err := r.db.Query(query, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Create(name string) (*model.Product, error) {
	result, err := r.db.Exec("INSERT INTO StoreManager.products (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &model.Product{ID: int(id), Name: name}, nil
}