package inventory_test

import (
	"testing"

	"stock-service-version-three/internal/domain/inventory"
)

func TestNewInventory_SetsInitialValues(t *testing.T) {
	inv, err := inventory.NewInventory(10, 1, 1200000, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.ProductID != 10 {
		t.Errorf("ProductID = %d, want 10", inv.ProductID)
	}
	if inv.WarehouseID != 1 {
		t.Errorf("WarehouseID = %d, want 1", inv.WarehouseID)
	}
	if inv.BasePrice != 1200000 {
		t.Errorf("BasePrice = %d, want 1200000", inv.BasePrice)
	}
	if inv.FinalPrice != 1200000 {
		t.Errorf("FinalPrice = %d, want 1200000", inv.FinalPrice)
	}
	if inv.AvailableStock != 50 {
		t.Errorf("AvailableStock = %d, want 50", inv.AvailableStock)
	}
	if inv.ReservedStock != 0 {
		t.Errorf("ReservedStock = %d, want 0", inv.ReservedStock)
	}
	if inv.StockOut != 0 {
		t.Errorf("StockOut = %d, want 0", inv.StockOut)
	}
	if inv.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if inv.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestNewInventory_ZeroStock(t *testing.T) {
	inv, err := inventory.NewInventory(10, 1, 1200000, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.AvailableStock != 0 {
		t.Errorf("AvailableStock = %d, want 0", inv.AvailableStock)
	}
}

func TestNewInventory_NegativeStock(t *testing.T) {
	_, err := inventory.NewInventory(10, 1, 1200000, -5)
	if err != inventory.ErrInvalidStock {
		t.Errorf("expected ErrInvalidStock, got %v", err)
	}
}

func TestNewInventory_ZeroBasePrice(t *testing.T) {
	_, err := inventory.NewInventory(10, 1, 0, 50)
	if err != inventory.ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestNewInventory_NegativeBasePrice(t *testing.T) {
	_, err := inventory.NewInventory(10, 1, -1000, 50)
	if err != inventory.ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestNewInventory_EmitsEvents(t *testing.T) {
	inv, err := inventory.NewInventory(10, 1, 1200000, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := inv.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	created, ok := events[0].(inventory.InventoryCreatedEvent)
	if !ok {
		t.Fatalf("expected InventoryCreatedEvent, got %T", events[0])
	}
	if created.InventoryID != 0 {
		t.Errorf("InventoryID = %d, want 0", created.InventoryID)
	}
	if created.ProductID != 10 {
		t.Errorf("ProductID = %d, want 10", created.ProductID)
	}
	if created.AvailableStock != 50 {
		t.Errorf("AvailableStock = %d, want 50", created.AvailableStock)
	}
	if created.FinalPrice != 1200000 {
		t.Errorf("FinalPrice = %d, want 1200000", created.FinalPrice)
	}

	stockIn, ok := events[1].(inventory.StockInEvent)
	if !ok {
		t.Fatalf("expected StockInEvent, got %T", events[1])
	}
	if stockIn.InventoryID != 0 {
		t.Errorf("InventoryID = %d, want 0", stockIn.InventoryID)
	}
	if stockIn.Quantity != 50 {
		t.Errorf("Quantity = %d, want 50", stockIn.Quantity)
	}
}

func TestApplyDiscount(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	inv.ApplyDiscount(15)

	if inv.FinalPrice != 1020000 {
		t.Errorf("FinalPrice = %d, want 1020000", inv.FinalPrice)
	}
}

func TestApplyDiscount_FullDiscount(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	inv.ApplyDiscount(100)

	if inv.FinalPrice != 0 {
		t.Errorf("FinalPrice = %d, want 0", inv.FinalPrice)
	}
}

func TestApplyDiscount_ZeroDiscount(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	inv.ApplyDiscount(0)

	if inv.FinalPrice != 1200000 {
		t.Errorf("FinalPrice = %d, want 1200000", inv.FinalPrice)
	}
}

func TestApplyDiscount_EmitsPriceUpdated(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	inv.ApplyDiscount(15)

	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	evt, ok := events[0].(inventory.PriceUpdatedEvent)
	if !ok {
		t.Fatalf("expected PriceUpdatedEvent, got %T", events[0])
	}
	if evt.InventoryID != 0 {
		t.Errorf("InventoryID = %d, want 0", evt.InventoryID)
	}
	if evt.OldPrice != 1200000 {
		t.Errorf("OldPrice = %d, want 1200000", evt.OldPrice)
	}
	if evt.NewPrice != 1020000 {
		t.Errorf("NewPrice = %d, want 1020000", evt.NewPrice)
	}
}

func TestReserveStock_Success(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	err := inv.ReserveStock(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.AvailableStock != 48 {
		t.Errorf("AvailableStock = %d, want 48", inv.AvailableStock)
	}
	if inv.ReservedStock != 2 {
		t.Errorf("ReservedStock = %d, want 2", inv.ReservedStock)
	}
}

func TestReserveStock_InsufficientStock(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 5)
	_ = inv.Events()

	err := inv.ReserveStock(10)
	if err != inventory.ErrInsufficientStock {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}

	if inv.AvailableStock != 5 {
		t.Errorf("AvailableStock should remain 5, got %d", inv.AvailableStock)
	}
	if inv.ReservedStock != 0 {
		t.Errorf("ReservedStock should remain 0, got %d", inv.ReservedStock)
	}
}

func TestReserveStock_ZeroQuantity(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	err := inv.ReserveStock(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.AvailableStock != 50 {
		t.Errorf("AvailableStock = %d, want 50", inv.AvailableStock)
	}
	if inv.ReservedStock != 0 {
		t.Errorf("ReservedStock = %d, want 0", inv.ReservedStock)
	}
}

func TestReserveStock_NegativeQuantity(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	err := inv.ReserveStock(-1)
	if err != inventory.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestReserveStock_EmitsReservedEvent(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	_ = inv.ReserveStock(3)

	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	evt, ok := events[0].(inventory.ReservedEvent)
	if !ok {
		t.Fatalf("expected ReservedEvent, got %T", events[0])
	}
	if evt.InventoryID != 0 {
		t.Errorf("InventoryID = %d, want 0", evt.InventoryID)
	}
	if evt.Quantity != 3 {
		t.Errorf("Quantity = %d, want 3", evt.Quantity)
	}
}

func TestFinalizeSale_Success(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(2)
	_ = inv.Events()

	err := inv.FinalizeSale(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.ReservedStock != 0 {
		t.Errorf("ReservedStock = %d, want 0", inv.ReservedStock)
	}
	if inv.StockOut != 2 {
		t.Errorf("StockOut = %d, want 2", inv.StockOut)
	}
	if inv.AvailableStock != 48 {
		t.Errorf("AvailableStock = %d, want 48", inv.AvailableStock)
	}
}

func TestFinalizeSale_InsufficientReserved(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(2)
	_ = inv.Events()

	err := inv.FinalizeSale(5)
	if err != inventory.ErrInsufficientReservedStock {
		t.Errorf("expected ErrInsufficientReservedStock, got %v", err)
	}

	if inv.ReservedStock != 2 {
		t.Errorf("ReservedStock should remain 2, got %d", inv.ReservedStock)
	}
	if inv.StockOut != 0 {
		t.Errorf("StockOut should remain 0, got %d", inv.StockOut)
	}
}

func TestFinalizeSale_PartialDelivery(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(10)
	_ = inv.Events()

	err := inv.FinalizeSale(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.ReservedStock != 6 {
		t.Errorf("ReservedStock = %d, want 6", inv.ReservedStock)
	}
	if inv.StockOut != 4 {
		t.Errorf("StockOut = %d, want 4", inv.StockOut)
	}
}

func TestFinalizeSale_ZeroQuantity(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(2)
	_ = inv.Events()

	err := inv.FinalizeSale(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.ReservedStock != 2 {
		t.Errorf("ReservedStock = %d, want 2", inv.ReservedStock)
	}
	if inv.StockOut != 0 {
		t.Errorf("StockOut = %d, want 0", inv.StockOut)
	}
}

func TestFinalizeSale_NegativeQuantity(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(2)
	_ = inv.Events()

	err := inv.FinalizeSale(-1)
	if err != inventory.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestFinalizeSale_EmitsStockOutEvent(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(2)
	_ = inv.Events()

	_ = inv.FinalizeSale(2)

	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	evt, ok := events[0].(inventory.StockOutEvent)
	if !ok {
		t.Fatalf("expected StockOutEvent, got %T", events[0])
	}
	if evt.InventoryID != 0 {
		t.Errorf("InventoryID = %d, want 0", evt.InventoryID)
	}
	if evt.Quantity != 2 {
		t.Errorf("Quantity = %d, want 2", evt.Quantity)
	}
}

func TestResetPrice(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	inv.ApplyDiscount(15)
	_ = inv.Events()

	inv.ResetPrice()

	if inv.FinalPrice != 1200000 {
		t.Errorf("FinalPrice = %d, want 1200000", inv.FinalPrice)
	}
}

func TestResetPrice_EmitsPriceUpdated(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	inv.ApplyDiscount(15)
	_ = inv.Events()

	inv.ResetPrice()

	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	evt, ok := events[0].(inventory.PriceUpdatedEvent)
	if !ok {
		t.Fatalf("expected PriceUpdatedEvent, got %T", events[0])
	}
	if evt.InventoryID != 0 {
		t.Errorf("InventoryID = %d, want 0", evt.InventoryID)
	}
	if evt.OldPrice != 1020000 {
		t.Errorf("OldPrice = %d, want 1020000", evt.OldPrice)
	}
	if evt.NewPrice != 1200000 {
		t.Errorf("NewPrice = %d, want 1200000", evt.NewPrice)
	}
}

func TestResetPrice_NoDiscountApplied(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	inv.ResetPrice()

	if inv.FinalPrice != 1200000 {
		t.Errorf("FinalPrice = %d, want 1200000", inv.FinalPrice)
	}
}

func TestRestoreStock(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(5)
	_ = inv.Events()

	inv.RestoreStock(3)

	if inv.AvailableStock != 48 {
		t.Errorf("AvailableStock = %d, want 48", inv.AvailableStock)
	}
	if inv.ReservedStock != 2 {
		t.Errorf("ReservedStock = %d, want 2", inv.ReservedStock)
	}
}

func TestRestoreStock_ZeroQuantity(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(5)
	_ = inv.Events()

	inv.RestoreStock(0)

	if inv.AvailableStock != 45 {
		t.Errorf("AvailableStock = %d, want 45", inv.AvailableStock)
	}
	if inv.ReservedStock != 5 {
		t.Errorf("ReservedStock = %d, want 5", inv.ReservedStock)
	}
}

func TestRestoreStock_NegativeQuantity(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()
	_ = inv.ReserveStock(5)
	_ = inv.Events()

	inv.RestoreStock(-1)

	if inv.AvailableStock != 45 {
		t.Errorf("AvailableStock should remain 45, got %d", inv.AvailableStock)
	}
	if inv.ReservedStock != 5 {
		t.Errorf("ReservedStock should remain 5, got %d", inv.ReservedStock)
	}
}

func TestEvents_ClearedAfterRetrieval(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	events := inv.Events()
	if len(events) == 0 {
		t.Fatal("expected events from creation")
	}

	events = inv.Events()
	if len(events) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events))
	}
}

func TestMultipleOperations_EventOrder(t *testing.T) {
	inv, _ := inventory.NewInventory(10, 1, 1200000, 50)
	_ = inv.Events()

	inv.ApplyDiscount(15)
	_ = inv.ReserveStock(2)
	_ = inv.FinalizeSale(2)

	events := inv.Events()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	if _, ok := events[0].(inventory.PriceUpdatedEvent); !ok {
		t.Errorf("events[0] expected PriceUpdatedEvent, got %T", events[0])
	}
	if _, ok := events[1].(inventory.ReservedEvent); !ok {
		t.Errorf("events[1] expected ReservedEvent, got %T", events[1])
	}
	if _, ok := events[2].(inventory.StockOutEvent); !ok {
		t.Errorf("events[2] expected StockOutEvent, got %T", events[2])
	}
}
