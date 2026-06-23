package inventory

import (
	domain "stock-service-version-three/internal/domain/inventory"
)

type HandleOrderPaidUseCase struct {
	repo domain.InventoryRepository
}

func NewHandleOrderPaidUseCase(repo domain.InventoryRepository) *HandleOrderPaidUseCase {
	return &HandleOrderPaidUseCase{repo: repo}
}

func (uc *HandleOrderPaidUseCase) Execute(req HandleOrderPaidRequest) (*HandleOrderPaidResponse, error) {
	inv, err := uc.repo.FindByID(req.InventoryID)
	if err != nil {
		return nil, err
	}

	if err := inv.ReserveStock(req.Quantity); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(inv); err != nil {
		return nil, err
	}

	return &HandleOrderPaidResponse{
		AvailableStock: inv.AvailableStock,
		ReservedStock:  inv.ReservedStock,
	}, nil
}
