package inventory_test

import (
	"testing"

	appInventory "stock-service-version-three/internal/application/inventory"
	domain "stock-service-version-three/internal/domain/inventory"
)

func TestCreateInventory_Success(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewCreateInventoryUseCase(repo)

	req := appInventory.CreateInventoryRequest{
		ProductID:   10,
		WarehouseID: 1,
		BasePrice:   100000,
		Stock:       50,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.InventoryID == 0 {
		t.Error("InventoryID should not be zero")
	}
	if resp.AvailableStock != 50 {
		t.Errorf("AvailableStock = %d, want 50", resp.AvailableStock)
	}
	if resp.FinalPrice != 100000 {
		t.Errorf("FinalPrice = %d, want 100000", resp.FinalPrice)
	}

	saved, err := repo.FindByID(resp.InventoryID)
	if err != nil {
		t.Fatalf("failed to find saved inventory: %v", err)
	}
	if saved.ProductID != 10 {
		t.Errorf("ProductID = %d, want 10", saved.ProductID)
	}
}

func TestCreateInventory_ZeroBasePrice(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewCreateInventoryUseCase(repo)

	req := appInventory.CreateInventoryRequest{
		ProductID:   10,
		WarehouseID: 1,
		BasePrice:   0,
		Stock:       50,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestCreateInventory_NegativeBasePrice(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewCreateInventoryUseCase(repo)

	req := appInventory.CreateInventoryRequest{
		ProductID:   10,
		WarehouseID: 1,
		BasePrice:   -1000,
		Stock:       50,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestCreateInventory_NegativeStock(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewCreateInventoryUseCase(repo)

	req := appInventory.CreateInventoryRequest{
		ProductID:   10,
		WarehouseID: 1,
		BasePrice:   100000,
		Stock:       -5,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInvalidStock {
		t.Errorf("expected ErrInvalidStock, got %v", err)
	}
}

func TestCreateInventory_ZeroStock(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewCreateInventoryUseCase(repo)

	req := appInventory.CreateInventoryRequest{
		ProductID:   10,
		WarehouseID: 1,
		BasePrice:   100000,
		Stock:       0,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.AvailableStock != 0 {
		t.Errorf("AvailableStock = %d, want 0", resp.AvailableStock)
	}
}
