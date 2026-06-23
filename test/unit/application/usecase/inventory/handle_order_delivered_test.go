package inventory_test

import (
	"testing"

	appInventory "stock-service-version-three/internal/application/inventory"
	domain "stock-service-version-three/internal/domain/inventory"
)

func TestHandleOrderDelivered_Success(t *testing.T) {
	repo := newMockInventoryRepository()
	handlePaidUC := appInventory.NewHandleOrderPaidUseCase(repo)
	uc := appInventory.NewHandleOrderDeliveredUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	handlePaidUC.Execute(appInventory.HandleOrderPaidRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    10,
	})

	req := appInventory.HandleOrderDeliveredRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    4,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StockOut != 4 {
		t.Errorf("StockOut = %d, want 4", resp.StockOut)
	}

	saved, _ := repo.FindByID(createResp.InventoryID)
	if saved.ReservedStock != 6 {
		t.Errorf("ReservedStock = %d, want 6", saved.ReservedStock)
	}
	if saved.StockOut != 4 {
		t.Errorf("StockOut = %d, want 4", saved.StockOut)
	}
}

func TestHandleOrderDelivered_InventoryNotFound(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderDeliveredUseCase(repo)

	req := appInventory.HandleOrderDeliveredRequest{
		InventoryID: 999,
		Quantity:    5,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}

func TestHandleOrderDelivered_InsufficientReservedStock(t *testing.T) {
	repo := newMockInventoryRepository()
	handlePaidUC := appInventory.NewHandleOrderPaidUseCase(repo)
	uc := appInventory.NewHandleOrderDeliveredUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	handlePaidUC.Execute(appInventory.HandleOrderPaidRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    2,
	})

	req := appInventory.HandleOrderDeliveredRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    5,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInsufficientReservedStock {
		t.Errorf("expected ErrInsufficientReservedStock, got %v", err)
	}
}

func TestHandleOrderDelivered_ZeroQuantity(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderDeliveredUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.HandleOrderDeliveredRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    0,
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StockOut != 0 {
		t.Errorf("StockOut = %d, want 0", resp.StockOut)
	}
}

func TestHandleOrderDelivered_NegativeQuantity(t *testing.T) {
	repo := newMockInventoryRepository()
	uc := appInventory.NewHandleOrderDeliveredUseCase(repo)

	createUC := appInventory.NewCreateInventoryUseCase(repo)
	createResp, _ := createUC.Execute(appInventory.CreateInventoryRequest{
		ProductID: 10, WarehouseID: 1, BasePrice: 100000, Stock: 50,
	})

	req := appInventory.HandleOrderDeliveredRequest{
		InventoryID: createResp.InventoryID,
		Quantity:    -1,
	}

	_, err := uc.Execute(req)
	if err != domain.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}
