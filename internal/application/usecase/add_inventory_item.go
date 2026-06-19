package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrStoreNotActive              = errors.New("store is not active")
	ErrWarehouseNotLinkedToStore   = errors.New("warehouse is not linked to the specified store")
	ErrCategoryAccessNotApproved   = errors.New("category access has not been approved for this store")
)

// AddInventoryItemInput contains the data needed to add a new inventory item.
type AddInventoryItemInput struct {
	ID          string
	StoreID     string
	WarehouseID string
	ProductID   string
	CategoryID  string
	BasePrice   float64
	InstantQty  int
	SaleModel   string
	Condition   string
	MinOrderQty int
	MaxOrderQty *int
	Attributes  map[string]string
}

// AddInventoryItemOutput contains the result of adding an inventory item.
type AddInventoryItemOutput struct {
	ID               string
	StoreID          string
	WarehouseID      string
	ProductID        string
	VendorSaleStatus string
	SystemSaleStatus string
	InstantQty       int
	BasePrice        float64
	FinalPrice       float64
	Event            any
}

// AddInventoryItemUseCase orchestrates adding a new inventory item.
type AddInventoryItemUseCase struct {
	storeRepo         entity.StoreRepository
	warehouseRepo     entity.WarehouseRepository
	storeCategoryRepo entity.StoreCategoryRepository
	inventoryRepo     entity.InventoryRepository
}

// NewAddInventoryItemUseCase creates a new AddInventoryItemUseCase.
func NewAddInventoryItemUseCase(
	storeRepo entity.StoreRepository,
	warehouseRepo entity.WarehouseRepository,
	storeCategoryRepo entity.StoreCategoryRepository,
	inventoryRepo entity.InventoryRepository,
) *AddInventoryItemUseCase {
	return &AddInventoryItemUseCase{
		storeRepo:         storeRepo,
		warehouseRepo:     warehouseRepo,
		storeCategoryRepo: storeCategoryRepo,
		inventoryRepo:     inventoryRepo,
	}
}

// Execute adds a new inventory item after validating preconditions.
func (uc *AddInventoryItemUseCase) Execute(input AddInventoryItemInput) (*AddInventoryItemOutput, error) {
	store, err := uc.storeRepo.FindByID(input.StoreID)
	if err != nil {
		return nil, err
	}
	if store.Status != entity.StoreStatusActive {
		return nil, ErrStoreNotActive
	}

	warehouse, err := uc.warehouseRepo.FindByID(input.WarehouseID)
	if err != nil {
		return nil, err
	}
	if warehouse.StoreID != input.StoreID {
		return nil, ErrWarehouseNotLinkedToStore
	}

	sc, err := uc.storeCategoryRepo.FindByStoreIDAndCategoryID(input.StoreID, input.CategoryID)
	if err != nil {
		return nil, err
	}
	if sc == nil || sc.Status != entity.StoreCategoryStatusApproved {
		return nil, ErrCategoryAccessNotApproved
	}

	inv, err := entity.NewInventory(input.ID, input.StoreID, input.WarehouseID, input.ProductID,
		input.BasePrice, input.InstantQty, input.SaleModel, input.Condition,
		input.MinOrderQty, input.MaxOrderQty, input.Attributes)
	if err != nil {
		return nil, err
	}

	if err := uc.inventoryRepo.Save(inv); err != nil {
		return nil, err
	}

	events := inv.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &AddInventoryItemOutput{
		ID:               inv.ID,
		StoreID:          inv.StoreID,
		WarehouseID:      inv.WarehouseID,
		ProductID:        inv.ProductID,
		VendorSaleStatus: string(inv.VendorSaleStatus),
		SystemSaleStatus: string(inv.SystemSaleStatus),
		InstantQty:       inv.InstantQty,
		BasePrice:        inv.BasePrice,
		FinalPrice:       inv.FinalPrice,
		Event:            domainEvent,
	}, nil
}
