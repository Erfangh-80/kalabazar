package warehouse

import (
	domain "stock-service-version-three/internal/domain/warehouse"
)

type LinkRepository interface {
	SaveLink(link *domain.StoreWarehouseLink) error
}

type LinkWarehouseToStoreUseCase struct {
	linkRepo LinkRepository
}

func NewLinkWarehouseToStoreUseCase(linkRepo LinkRepository) *LinkWarehouseToStoreUseCase {
	return &LinkWarehouseToStoreUseCase{linkRepo: linkRepo}
}

func (uc *LinkWarehouseToStoreUseCase) Execute(req LinkWarehouseRequest) (*LinkWarehouseResponse, error) {
	link, err := domain.NewStoreWarehouseLink(req.StoreID, req.WarehouseID, req.Type)
	if err != nil {
		return nil, err
	}

	if err := uc.linkRepo.SaveLink(link); err != nil {
		return nil, err
	}

	return &LinkWarehouseResponse{
		StoreID:     link.StoreID,
		WarehouseID: link.WarehouseID,
	}, nil
}
