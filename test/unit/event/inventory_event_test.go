package event_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func ip(i int) *int { return &i }

func TestInventoryEvents_ItemCreated(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryItemCreated)
	if !ok {
		t.Fatalf("expected InventoryItemCreated, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.ProductID != "prod-1" {
		t.Errorf("expected ProductID prod-1, got %s", e.ProductID)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if e.EventName() != "inventory.item_created" {
		t.Errorf("expected inventory.item_created, got %s", e.EventName())
	}
}

func TestInventoryEvents_StockUpdated(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	inv.UpdateStock(25)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryStockUpdated)
	if !ok {
		t.Fatalf("expected InventoryStockUpdated, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.NewQty != 25 {
		t.Errorf("expected NewQty 25, got %d", e.NewQty)
	}
	if e.EventName() != "inventory.stock_updated" {
		t.Errorf("expected inventory.stock_updated, got %s", e.EventName())
	}
}

func TestInventoryEvents_PriceUpdated(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	inv.UpdatePrice(200, 180)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryPriceUpdated)
	if !ok {
		t.Fatalf("expected InventoryPriceUpdated, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.BasePrice != 200 {
		t.Errorf("expected BasePrice 200, got %f", e.BasePrice)
	}
	if e.FinalPrice != 180 {
		t.Errorf("expected FinalPrice 180, got %f", e.FinalPrice)
	}
	if e.EventName() != "inventory.price_updated" {
		t.Errorf("expected inventory.price_updated, got %s", e.EventName())
	}
}

func TestInventoryEvents_VendorActivated(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.SetVendorStatus(entity.VendorSaleStatusInactive)
	inv.Events()
	inv.SetVendorStatus(entity.VendorSaleStatusActive)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryItemActivated)
	if !ok {
		t.Fatalf("expected InventoryItemActivated, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.EventName() != "inventory.item_activated" {
		t.Errorf("expected inventory.item_activated, got %s", e.EventName())
	}
}

func TestInventoryEvents_VendorDeactivated(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	inv.SetVendorStatus(entity.VendorSaleStatusInactive)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryItemDeactivated)
	if !ok {
		t.Fatalf("expected InventoryItemDeactivated, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.EventName() != "inventory.item_deactivated" {
		t.Errorf("expected inventory.item_deactivated, got %s", e.EventName())
	}
}

func TestInventoryEvents_SystemBlocked(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	inv.SetSystemStatus(entity.SystemSaleStatusInactive)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventorySystemBlocked)
	if !ok {
		t.Fatalf("expected InventorySystemBlocked, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.EventName() != "inventory.system_blocked" {
		t.Errorf("expected inventory.system_blocked, got %s", e.EventName())
	}
}

func TestInventoryEvents_SystemUnblocked(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.SetSystemStatus(entity.SystemSaleStatusInactive)
	inv.Events()
	inv.SetSystemStatus(entity.SystemSaleStatusActive)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventorySystemUnblocked)
	if !ok {
		t.Fatalf("expected InventorySystemUnblocked, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.EventName() != "inventory.system_unblocked" {
		t.Errorf("expected inventory.system_unblocked, got %s", e.EventName())
	}
}

func TestInventoryEvents_SaleScheduled(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	start := time.Now().Add(24 * time.Hour)
	end := time.Now().Add(72 * time.Hour)
	inv.SetSaleSchedule(&start, &end)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventorySaleScheduled)
	if !ok {
		t.Fatalf("expected InventorySaleScheduled, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.StartAt == nil || !e.StartAt.Equal(start) {
		t.Errorf("expected StartAt %v, got %v", start, e.StartAt)
	}
	if e.EndAt == nil || !e.EndAt.Equal(end) {
		t.Errorf("expected EndAt %v, got %v", end, e.EndAt)
	}
	if e.EventName() != "inventory.sale_scheduled" {
		t.Errorf("expected inventory.sale_scheduled, got %s", e.EventName())
	}
}

func TestInventoryEvents_PromotionLinked(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	inv.LinkPromotion("promo-1")
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryPromotionLinked)
	if !ok {
		t.Fatalf("expected InventoryPromotionLinked, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.EventName() != "inventory.promotion_linked" {
		t.Errorf("expected inventory.promotion_linked, got %s", e.EventName())
	}
}

func TestInventoryEvents_PromotionStatusChanged(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	inv.UpdatePromotionStatus(entity.CampaignApprovalApproved)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryPromotionStatusChanged)
	if !ok {
		t.Fatalf("expected InventoryPromotionStatusChanged, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.Status != "approved" {
		t.Errorf("expected status approved, got %s", e.Status)
	}
	if e.EventName() != "inventory.promotion_status_changed" {
		t.Errorf("expected inventory.promotion_status_changed, got %s", e.EventName())
	}
}

func TestInventoryEvents_ClearedAfterEventsCall(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "s-1", "wh-1", "prod-1", 100, 10,
		"retail", "new", 1, ip(100), nil)
	inv.Events()
	remaining := inv.Events()
	if len(remaining) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(remaining))
	}
}


