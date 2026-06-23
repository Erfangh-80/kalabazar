package inventory

import (
	domain "stock-service-version-three/internal/domain/inventory"
)

type HandleOrderDeliveredUseCase struct {
	repo domain.InventoryRepository
}

func NewHandleOrderDeliveredUseCase(repo domain.InventoryRepository) *HandleOrderDeliveredUseCase {
	return &HandleOrderDeliveredUseCase{repo: repo}
}

func (uc *HandleOrderDeliveredUseCase) Execute(req HandleOrderDeliveredRequest) (*HandleOrderDeliveredResponse, error) {
	inv, err := uc.repo.FindByID(req.InventoryID)
	if err != nil {
		return nil, err
	}

	if err := inv.FinalizeSale(req.Quantity); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(inv); err != nil {
		return nil, err
	}

	return &HandleOrderDeliveredResponse{
		StockOut: inv.StockOut,
	}, nil
}
