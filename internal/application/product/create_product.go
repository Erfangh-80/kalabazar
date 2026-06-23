package product

import (
	"errors"

	domain "stock-service-version-three/internal/domain/product"
)

type CreateProductUseCase struct {
	repo domain.ProductRepository
}

func NewCreateProductUseCase(repo domain.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{repo: repo}
}

func (uc *CreateProductUseCase) Execute(req CreateProductRequest) (*CreateProductResponse, error) {
	if req.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	p := domain.NewProduct(req.StoreID, req.Title, req.CategoryID, req.Brand)

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return &CreateProductResponse{
		ProductID: p.ID,
		Status:    string(p.Status),
	}, nil
}
