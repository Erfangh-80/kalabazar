package inventory_test

import (
	"testing"

	appInventory "stock-service-version-three/internal/application/inventory"
	domain "stock-service-version-three/internal/domain/inventory"
)

func TestUpdatePriceFromCampaign_Success(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewUpdatePriceUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.UpdatePriceRequest{
		InventoryID:        createResp.InventoryID,
		DiscountPercentage: 20,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.InventoryID != createResp.InventoryID {
		t.Errorf("InventoryID = %d, want %d", resp.InventoryID, createResp.InventoryID)
	}
	if resp.OldPrice != 100000 {
		t.Errorf("OldPrice = %d, want 100000", resp.OldPrice)
	}
	if resp.NewPrice != 80000 {
		t.Errorf("NewPrice = %d, want 80000", resp.NewPrice)
	}
}

func TestUpdatePriceFromCampaign_InventoryNotFound(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewUpdatePriceUseCase(repo)

	req := appInventory.UpdatePriceRequest{
		InventoryID:        999,
		DiscountPercentage: 20,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}

func TestUpdatePriceFromCampaign_ZeroDiscount(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewUpdatePriceUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.UpdatePriceRequest{
		InventoryID:        createResp.InventoryID,
		DiscountPercentage: 0,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.NewPrice != 100000 {
		t.Errorf("NewPrice = %d, want 100000", resp.NewPrice)
	}
}

func TestUpdatePriceFromCampaign_FullDiscount(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewUpdatePriceUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.UpdatePriceRequest{
		InventoryID:        createResp.InventoryID,
		DiscountPercentage: 100,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.NewPrice != 0 {
		t.Errorf("NewPrice = %d, want 0", resp.NewPrice)
	}
}
