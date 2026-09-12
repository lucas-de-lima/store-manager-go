package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/lucas-de-lima/store-manager-go/model"
	"github.com/lucas-de-lima/store-manager-go/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := strings.TrimPrefix(r.URL.Path, "/products")
		idStr = strings.TrimPrefix(idStr, "/")
		if idStr == "" || idStr == "search" {
			if idStr == "search" {
				h.search(w, r)
				return
			}
			h.getAll(w, r)
			return
		}
		h.getByID(w, r, idStr)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, `{"message":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, `{"message":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	if products == nil {
		products = []model.Product{}
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"message":"Invalid id"}`, http.StatusBadRequest)
		return
	}
	product, err := h.service.GetByID(id)
	if err != nil {
		if err == service.ErrProductNotFound {
			http.Error(w, `{"message":"Product not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"message":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	product, err := h.service.Create(req.Name)
	if err != nil {
		switch err {
		case service.ErrNameRequired:
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		case service.ErrNameTooShort:
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"message": err.Error()})
		default:
			http.Error(w, `{"message":"Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	products, err := h.service.Search(q)
	if err != nil {
		http.Error(w, `{"message":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, products)
}