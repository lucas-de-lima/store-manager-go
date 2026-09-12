package handler

import (
	"net/http"

	"github.com/lucas-de-lima/store-manager-go/repository"
	"github.com/lucas-de-lima/store-manager-go/service"
)

type Router struct {
	ProductHandler http.Handler
	SaleHandler    http.Handler
}

func NewRouter(productRepo *repository.ProductRepository, saleRepo *repository.SaleRepository) *Router {
	productSvc := service.NewProductService(productRepo)
	saleSvc := service.NewSaleService(saleRepo, productRepo)
	return &Router{
		ProductHandler: NewProductHandler(productSvc),
		SaleHandler:    NewSaleHandler(saleSvc),
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path

	if path == "/" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if len(path) >= 9 && path[:9] == "/products" {
		rt.ProductHandler.ServeHTTP(w, r)
		return
	}

	if len(path) >= 6 && path[:6] == "/sales" {
		rt.SaleHandler.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}