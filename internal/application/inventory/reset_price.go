package inventory

import (
	domain "stock-service-version-three/internal/domain/inventory"
)

type ResetPriceUseCase struct {
	repo domain.InventoryRepository
}

func NewResetPriceUseCase(repo domain.InventoryRepository) *ResetPriceUseCase {
	return &ResetPriceUseCase{repo: repo}
}

func (uc *ResetPriceUseCase) Execute(req ResetPriceRequest) (*ResetPriceResponse, error) {
	inv, err := uc.repo.FindByID(req.InventoryID)
	if err != nil {
		return nil, err
	}

	inv.ResetPrice()

	if err := uc.repo.Update(inv); err != nil {
		return nil, err
	}

	return &ResetPriceResponse{
		InventoryID: inv.ID,
		FinalPrice:  inv.FinalPrice,
	}, nil
}
