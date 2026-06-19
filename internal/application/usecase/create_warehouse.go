package usecase

import (
	"time"

	"kalabazar-stock-service/internal/domain/entity"
)

// CreateWarehouseInput contains the data needed to create a new warehouse.
type CreateWarehouseInput struct {
	ID            string
	SellerID      string
	Name          string
	Address       entity.Address
	TotalCapacity int
	AccessType    string
}

// CreateWarehouseOutput contains the result of a warehouse creation.
type CreateWarehouseOutput struct {
	ID            string
	SellerID      string
	Name          string
	Address       entity.Address
	AccessType    string
	TotalCapacity int
	Status        string
	Event         any
	CreatedAt     time.Time
}

// CreateWarehouseUseCase orchestrates the creation of a new warehouse.
type CreateWarehouseUseCase struct {
	repo entity.WarehouseRepository
}

// NewCreateWarehouseUseCase creates a new CreateWarehouseUseCase.
func NewCreateWarehouseUseCase(repo entity.WarehouseRepository) *CreateWarehouseUseCase {
	return &CreateWarehouseUseCase{repo: repo}
}

// Execute creates a new warehouse with the given input.
func (uc *CreateWarehouseUseCase) Execute(input CreateWarehouseInput) (*CreateWarehouseOutput, error) {
	warehouse, err := entity.NewWarehouse(input.ID, input.SellerID, input.Name, input.Address, input.TotalCapacity, input.AccessType)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(warehouse); err != nil {
		return nil, err
	}

	events := warehouse.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &CreateWarehouseOutput{
		ID:            warehouse.ID,
		SellerID:      warehouse.SellerID,
		Name:          warehouse.Name,
		Address:       warehouse.Address,
		AccessType:    warehouse.AccessType,
		TotalCapacity: warehouse.TotalCapacity,
		Status:        string(warehouse.Status),
		Event:         domainEvent,
		CreatedAt:     warehouse.CreatedAt,
	}, nil
}
