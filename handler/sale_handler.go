package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lucas-de-lima/store-manager-go/model"
	"github.com/lucas-de-lima/store-manager-go/service"
)

type SaleHandler struct {
	service *service.SaleService
}

func NewSaleHandler(svc *service.SaleService) *SaleHandler {
	return &SaleHandler{service: svc}
}

func (h *SaleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, `{"message":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

type saleInput struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

func (h *SaleHandler) create(w http.ResponseWriter, r *http.Request) {
	var raw []map[string]any
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil || len(raw) == 0 {
		http.Error(w, `{"message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	var req []saleInput
	for _, item := range raw {
		productID, hasProductID := item["productId"].(float64)
		quantityRaw, hasQuantity := item["quantity"]

		if !hasProductID {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": `"productId" is required`})
			return
		}
		if !hasQuantity {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": `"quantity" is required`})
			return
		}

		q := int(quantityRaw.(float64))
		if q < 1 {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"message": `"quantity" must be greater than or equal to 1`})
			return
		}
		req = append(req, saleInput{ProductID: int(productID), Quantity: q})
	}

	saleProducts := make([]model.SaleProduct, len(req))
	for i, item := range req {
		saleProducts[i] = model.SaleProduct{ProductID: item.ProductID, Quantity: item.Quantity}
	}

	sale, err := h.service.Create(saleProducts)
	if err != nil {
		if err == service.ErrProductNotFound {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "Product not found"})
			return
		}
		http.Error(w, `{"message":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, sale)
}