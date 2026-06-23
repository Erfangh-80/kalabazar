package inventory_test

import (
	"testing"

	appInventory "stock-service-version-three/internal/application/inventory"
	domain "stock-service-version-three/internal/domain/inventory"
)

func TestResetInventoryPrice_Success(t *testing.T) {
	repo := newMockInventoryRepository()
	updatePriceUC := appInventory.NewUpdatePriceUseCase(repo)
	uc := appInventory.NewResetPriceUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	updatePriceUC.Execute(appInventory.UpdatePriceRequest{
		InventoryID: createResp.InventoryID, DiscountPercentage: 20,
	})

	req := appInventory.ResetPriceRequest{
		InventoryID: createResp.InventoryID,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.InventoryID != createResp.InventoryID {
		t.Errorf("InventoryID = %d, want %d", resp.InventoryID, createResp.InventoryID)
	}
	if resp.FinalPrice != 100000 {
		t.Errorf("FinalPrice = %d, want 100000", resp.FinalPrice)
	}

	saved, _ := repo.FindByID(createResp.InventoryID)
	if saved.FinalPrice != 100000 {
		t.Errorf("saved FinalPrice = %d, want 100000", saved.FinalPrice)
	}
}

func TestResetInventoryPrice_NoPriceChange(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewResetPriceUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.ResetPriceRequest{
		InventoryID: createResp.InventoryID,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.FinalPrice != 100000 {
		t.Errorf("FinalPrice = %d, want 100000", resp.FinalPrice)
	}
}

func TestResetInventoryPrice_InventoryNotFound(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewResetPriceUseCase(repo)

	req := appInventory.ResetPriceRequest{
		InventoryID: 999,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}
