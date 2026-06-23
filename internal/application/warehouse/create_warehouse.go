package warehouse

import (
	domain "stock-service-version-three/internal/domain/warehouse"
)

type CreateWarehouseUseCase struct {
	repo domain.WarehouseRepository
}

func NewCreateWarehouseUseCase(repo domain.WarehouseRepository) *CreateWarehouseUseCase {
	return &CreateWarehouseUseCase{repo: repo}
}

func (uc *CreateWarehouseUseCase) Execute(req CreateWarehouseRequest) (*CreateWarehouseResponse, error) {
	w, err := domain.NewWarehouse(req.Name, req.Capacity)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(w); err != nil {
		return nil, err
	}

	return &CreateWarehouseResponse{
		WarehouseID: w.ID,
		Name:        w.Name,
	}, nil
}
