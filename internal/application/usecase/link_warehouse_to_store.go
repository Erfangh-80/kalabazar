package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrInvalidRelationType           = errors.New("relation type must be 'primary' or 'secondary'")
	ErrWarehouseAlreadyLinkedToStore = errors.New("warehouse is already linked to this store")
)

// LinkWarehouseToStoreInput contains the data needed to link a warehouse to a store.
type LinkWarehouseToStoreInput struct {
	WarehouseID  string
	StoreID      string
	RelationType string
}

// LinkWarehouseToStoreOutput contains the result of linking a warehouse to a store.
type LinkWarehouseToStoreOutput struct {
	WarehouseID  string
	StoreID      string
	RelationType string
	Event        any
}

// LinkWarehouseToStoreUseCase orchestrates linking a warehouse to a store.
type LinkWarehouseToStoreUseCase struct {
	repo entity.WarehouseRepository
}

// NewLinkWarehouseToStoreUseCase creates a new LinkWarehouseToStoreUseCase.
func NewLinkWarehouseToStoreUseCase(repo entity.WarehouseRepository) *LinkWarehouseToStoreUseCase {
	return &LinkWarehouseToStoreUseCase{repo: repo}
}

// Execute links a warehouse to a store with the given relation type.
func (uc *LinkWarehouseToStoreUseCase) Execute(input LinkWarehouseToStoreInput) (*LinkWarehouseToStoreOutput, error) {
	if input.RelationType != "primary" && input.RelationType != "secondary" {
		return nil, ErrInvalidRelationType
	}

	warehouse, err := uc.repo.FindByID(input.WarehouseID)
	if err != nil {
		return nil, err
	}

	if warehouse.StoreID == input.StoreID {
		return nil, ErrWarehouseAlreadyLinkedToStore
	}

	if err := warehouse.LinkToStore(input.StoreID, input.RelationType); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(warehouse); err != nil {
		return nil, err
	}

	events := warehouse.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &LinkWarehouseToStoreOutput{
		WarehouseID:  warehouse.ID,
		StoreID:      warehouse.StoreID,
		RelationType: warehouse.RelationType,
		Event:        domainEvent,
	}, nil
}
