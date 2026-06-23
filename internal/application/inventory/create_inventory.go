package inventory

import (
	domain "stock-service-version-three/internal/domain/inventory"
)

type CreateInventoryUseCase struct {
	repo domain.InventoryRepository
}

func NewCreateInventoryUseCase(repo domain.InventoryRepository) *CreateInventoryUseCase {
	return &CreateInventoryUseCase{repo: repo}
}

func (uc *CreateInventoryUseCase) Execute(req CreateInventoryRequest) (*CreateInventoryResponse, error) {
	if req.BasePrice <= 0 {
		return nil, domain.ErrInvalidPrice
	}
	if req.Stock < 0 {
		return nil, domain.ErrInvalidStock
	}

	inv, err := domain.NewInventory(req.ProductID, req.WarehouseID, req.BasePrice, req.Stock)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}

	return &CreateInventoryResponse{
		InventoryID:    inv.ID,
		AvailableStock: inv.AvailableStock,
		FinalPrice:     inv.FinalPrice,
	}, nil
}
