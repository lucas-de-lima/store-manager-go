package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lucas-de-lima/store-manager-go/repository"
	"github.com/lucas-de-lima/store-manager-go/service"
)

func newTestProductHandler(t *testing.T) (*ProductHandler, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	h := NewProductHandler(productSvc)
	return h, mock, func() { db.Close() }
}

func newTestSaleHandler(t *testing.T) (*SaleHandler, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	productRepo := repository.NewProductRepository(db)
	saleRepo := repository.NewSaleRepository(db)
	saleSvc := service.NewSaleService(saleRepo, productRepo)
	h := NewSaleHandler(saleSvc)
	return h, mock, func() { db.Close() }
}

func TestProductHandler_GetAll(t *testing.T) {
	h, mock, cleanup := newTestProductHandler(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Martelo de Thor").
		AddRow(2, "Traje de encolhimento")

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products ORDER BY id").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var products []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil {
		t.Fatal(err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductHandler_GetAll_Empty(t *testing.T) {
	h, mock, cleanup := newTestProductHandler(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT id, name FROM StoreManager.products ORDER BY id").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductHandler_GetByID_Found(t *testing.T) {
	h, mock, cleanup := newTestProductHandler(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Martelo de Thor")

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var product map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &product); err != nil {
		t.Fatal(err)
	}
	if product["name"] != "Martelo de Thor" {
		t.Errorf("expected 'Martelo de Thor', got '%s'", product["name"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	h, mock, cleanup := newTestProductHandler(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products WHERE id = ?").
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "Product not found" {
		t.Errorf("expected 'Product not found', got '%s'", body["message"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductHandler_Create_Success(t *testing.T) {
	h, mock, cleanup := newTestProductHandler(t)
	defer cleanup()

	mock.ExpectExec("INSERT INTO StoreManager.products \\(name\\) VALUES \\(\\?\\)").
		WithArgs("ProdutoX").
		WillReturnResult(sqlmock.NewResult(4, 1))

	body := `{"name":"ProdutoX"}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var product map[string]any
	json.Unmarshal(w.Body.Bytes(), &product)
	if product["id"] != float64(4) {
		t.Errorf("expected id 4, got %v", product["id"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestProductHandler_Create_NoName(t *testing.T) {
	h, _, cleanup := newTestProductHandler(t)
	defer cleanup()

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_Create_NameTooShort(t *testing.T) {
	h, _, cleanup := newTestProductHandler(t)
	defer cleanup()

	body := `{"name":"abc"}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", w.Code)
	}
}

func TestSaleHandler_Create_Success(t *testing.T) {
	h, mock, cleanup := newTestSaleHandler(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Martelo de Thor")
	mock.ExpectQuery("SELECT id, name FROM StoreManager.products WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO StoreManager.sales \\(date\\) VALUES \\(NOW\\(\\)\\)").
		WillReturnResult(sqlmock.NewResult(3, 1))
	mock.ExpectPrepare("INSERT INTO StoreManager.sales_products \\(sale_id, product_id, quantity\\) VALUES \\(\\?, \\?, \\?\\)")
	mock.ExpectExec("INSERT INTO StoreManager.sales_products \\(sale_id, product_id, quantity\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(3, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := `[{"productId":1,"quantity":1}]`
	req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSaleHandler_Create_NoProductID(t *testing.T) {
	h, _, cleanup := newTestSaleHandler(t)
	defer cleanup()

	body := `[{"quantity":1}]`
	req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != `"productId" is required` {
		t.Errorf("expected productId required message, got '%s'", resp["message"])
	}
}

func TestSaleHandler_Create_NoQuantity(t *testing.T) {
	h, _, cleanup := newTestSaleHandler(t)
	defer cleanup()

	body := `[{"productId":1}]`
	req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != `"quantity" is required` {
		t.Errorf("expected quantity required message, got '%s'", resp["message"])
	}
}

func TestSaleHandler_Create_QuantityZero(t *testing.T) {
	h, _, cleanup := newTestSaleHandler(t)
	defer cleanup()

	body := `[{"productId":1,"quantity":0}]`
	req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != `"quantity" must be greater than or equal to 1` {
		t.Errorf("expected quantity >= 1 message, got '%s'", resp["message"])
	}
}

func TestSaleHandler_Create_ProductNotFound(t *testing.T) {
	h, mock, cleanup := newTestSaleHandler(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, name FROM StoreManager.products WHERE id = ?").
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	body := `[{"productId":999,"quantity":1}]`
	req := httptest.NewRequest(http.MethodPost, "/sales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "Product not found" {
		t.Errorf("expected 'Product not found', got '%s'", resp["message"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}