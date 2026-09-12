package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lucas-de-lima/store-manager-go/model"
)

func TestProductRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Martelo de Thor").
		AddRow(2, "Traje de encolhimento")

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products ORDER BY id").
		WillReturnRows(rows)

	repo := NewProductRepository(db)
	products, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
	if products[0].Name != "Martelo de Thor" {
		t.Errorf("expected 'Martelo de Thor', got '%s'", products[0].Name)
	}
	if products[1].Name != "Traje de encolhimento" {
		t.Errorf("expected 'Traje de encolhimento', got '%s'", products[1].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductRepository_GetAll_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"})

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products ORDER BY id").
		WillReturnRows(rows)

	repo := NewProductRepository(db)
	products, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductRepository_GetByID_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Martelo de Thor")

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	repo := NewProductRepository(db)
	product, err := repo.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}
	if product.ID != 1 || product.Name != "Martelo de Thor" {
		t.Errorf("expected id=1 name='Martelo de Thor', got id=%d name='%s'", product.ID, product.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products WHERE id = ?").
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	repo := NewProductRepository(db)
	product, err := repo.GetByID(999)
	if err != nil {
		t.Fatal(err)
	}

	if product != nil {
		t.Error("expected nil product for non-existent id")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO StoreManager.products \\(name\\) VALUES \\(\\?\\)").
		WithArgs("ProdutoX").
		WillReturnResult(sqlmock.NewResult(4, 1))

	repo := NewProductRepository(db)
	product, err := repo.Create("ProdutoX")
	if err != nil {
		t.Fatal(err)
	}

	if product.ID != 4 {
		t.Errorf("expected id 4, got %d", product.ID)
	}
	if product.Name != "ProdutoX" {
		t.Errorf("expected name 'ProdutoX', got '%s'", product.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Martelo de Thor")

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products WHERE name LIKE").
		WithArgs("%Martelo%").
		WillReturnRows(rows)

	repo := NewProductRepository(db)
	products, err := repo.Search("Martelo")
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 1 {
		t.Errorf("expected 1 product, got %d", len(products))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSaleRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO StoreManager.sales \\(date\\) VALUES \\(NOW\\(\\)\\)").
		WillReturnResult(sqlmock.NewResult(3, 1))
	mock.ExpectPrepare("INSERT INTO StoreManager.sales_products \\(sale_id, product_id, quantity\\) VALUES \\(\\?, \\?, \\?\\)")
	mock.ExpectExec("INSERT INTO StoreManager.sales_products \\(sale_id, product_id, quantity\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(3, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO StoreManager.sales_products \\(sale_id, product_id, quantity\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(3, 2, 5).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewSaleRepository(db)
	saleProducts := []model.SaleProduct{
		{ProductID: 1, Quantity: 1},
		{ProductID: 2, Quantity: 5},
	}
	sale, err := repo.Create(saleProducts)
	if err != nil {
		t.Fatal(err)
	}

	if sale.ID != 3 {
		t.Errorf("expected sale id 3, got %d", sale.ID)
	}
	if len(sale.ItemsSold) != 2 {
		t.Errorf("expected 2 items, got %d", len(sale.ItemsSold))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}