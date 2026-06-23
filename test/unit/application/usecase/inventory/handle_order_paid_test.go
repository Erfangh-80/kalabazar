package inventory_test

import (
	"testing"

	appInventory "stock-service-version-three/internal/application/inventory"
	domain "stock-service-version-three/internal/domain/inventory"
)

func TestHandleOrderPaid_Success(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderPaidUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.HandleOrderPaidRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    5,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.AvailableStock != 45 {
		t.Errorf("AvailableStock = %d, want 45", resp.AvailableStock)
	}
	if resp.ReservedStock != 5 {
		t.Errorf("ReservedStock = %d, want 5", resp.ReservedStock)
	}
}

func TestHandleOrderPaid_InventoryNotFound(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderPaidUseCase(repo)

	req := appInventory.HandleOrderPaidRequest{
		InventoryID: 999,
		Quantity:    5,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}

func TestHandleOrderPaid_InsufficientStock(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderPaidUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 5,
	})

	req := appInventory.HandleOrderPaidRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    10,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInsufficientStock {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}
}

func TestHandleOrderPaid_ZeroQuantity(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderPaidUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.HandleOrderPaidRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    0,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.AvailableStock != 50 {
		t.Errorf("AvailableStock = %d, want 50", resp.AvailableStock)
	}
	if resp.ReservedStock != 0 {
		t.Errorf("ReservedStock = %d, want 0", resp.ReservedStock)
	}
}

func TestHandleOrderPaid_NegativeQuantity(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderPaidUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.HandleOrderPaidRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    -1,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}
