package service

import (
	"errors"
	"strings"

	"github.com/lucas-de-lima/store-manager-go/model"
	"github.com/lucas-de-lima/store-manager-go/repository"
)

var (
	ErrNameRequired        = errors.New(`"name" is required`)
	ErrNameTooShort        = errors.New(`"name" length must be at least 5 characters long`)
	ErrProductNotFound     = errors.New("Product not found")
	ErrProductIDRequired   = errors.New(`"productId" is required`)
	ErrQuantityRequired    = errors.New(`"quantity" is required`)
	ErrQuantityInvalid     = errors.New(`"quantity" must be greater than or equal to 1`)
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id int) (*model.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) Search(q string) ([]model.Product, error) {
	return s.repo.Search(q)
}

func (s *ProductService) Create(name string) (*model.Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}
	if len(name) < 5 {
		return nil, ErrNameTooShort
	}
	return s.repo.Create(name)
}