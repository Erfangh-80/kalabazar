package inventory

import (
	domain "stock-service-version-three/internal/domain/inventory"
)

type UpdatePriceUseCase struct {
	repo domain.InventoryRepository
}

func NewUpdatePriceUseCase(repo domain.InventoryRepository) *UpdatePriceUseCase {
	return &UpdatePriceUseCase{repo: repo}
}

func (uc *UpdatePriceUseCase) Execute(req UpdatePriceRequest) (*UpdatePriceResponse, error) {
	inv, err := uc.repo.FindByID(req.InventoryID)
	if err != nil {
		return nil, err
	}

	oldPrice := inv.FinalPrice
	inv.ApplyDiscount(req.DiscountPercentage)

	if err := uc.repo.Update(inv); err != nil {
		return nil, err
	}

	return &UpdatePriceResponse{
		InventoryID: inv.ID,
		OldPrice:    oldPrice,
		NewPrice:    inv.FinalPrice,
	}, nil
}
