package service

import (
	"github.com/lucas-de-lima/store-manager-go/model"
	"github.com/lucas-de-lima/store-manager-go/repository"
)

type SaleService struct {
	repo          *repository.SaleRepository
	productRepo   *repository.ProductRepository
}

func NewSaleService(repo *repository.SaleRepository, productRepo *repository.ProductRepository) *SaleService {
	return &SaleService{repo: repo, productRepo: productRepo}
}

func (s *SaleService) Create(saleProducts []model.SaleProduct) (*model.Sale, error) {
	productIDs := make(map[int]bool)
	for _, sp := range saleProducts {
		productIDs[sp.ProductID] = true
	}

	for productID := range productIDs {
		product, err := s.productRepo.GetByID(productID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, ErrProductNotFound
		}
	}

	return s.repo.Create(saleProducts)
}