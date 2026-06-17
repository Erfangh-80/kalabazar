package entity_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestNewInventory_Success(t *testing.T) {
	inv, err := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 150000, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.ID != "inv-1" {
		t.Errorf("expected inv-1, got %s", inv.ID)
	}
	if inv.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", inv.StoreID)
	}
	if inv.WarehouseID != "wh-1" {
		t.Errorf("expected wh-1, got %s", inv.WarehouseID)
	}
	if inv.ProductID != "prod-1" {
		t.Errorf("expected prod-1, got %s", inv.ProductID)
	}
	if inv.BasePrice != 150000 {
		t.Errorf("expected 150000, got %f", inv.BasePrice)
	}
	if inv.InstantQty != 10 {
		t.Errorf("expected 10, got %d", inv.InstantQty)
	}
}

func TestNewInventory_Defaults(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100000, 0)
	if inv.VendorSaleStatus != entity.VendorSaleStatusActive {
		t.Errorf("expected vendor status active, got %s", inv.VendorSaleStatus)
	}
	if inv.SystemSaleStatus != entity.SystemSaleStatusActive {
		t.Errorf("expected system status active, got %s", inv.SystemSaleStatus)
	}
	if inv.CampaignApprovalStatus != entity.CampaignApprovalPending {
		t.Errorf("expected campaign approval pending, got %s", inv.CampaignApprovalStatus)
	}
	if inv.FinalPrice != 100000 {
		t.Errorf("expected FinalPrice to equal BasePrice, got %f", inv.FinalPrice)
	}
}

func TestNewInventory_InvalidID(t *testing.T) {
	_, err := entity.NewInventory("", "store-1", "wh-1", "prod-1", 100, 1)
	if err != entity.ErrInventoryInvalidID {
		t.Errorf("expected ErrInventoryInvalidID, got %v", err)
	}
}

func TestNewInventory_InvalidStoreID(t *testing.T) {
	_, err := entity.NewInventory("inv-1", "", "wh-1", "prod-1", 100, 1)
	if err != entity.ErrInventoryInvalidStoreID {
		t.Errorf("expected ErrInventoryInvalidStoreID, got %v", err)
	}
}

func TestNewInventory_InvalidWarehouseID(t *testing.T) {
	_, err := entity.NewInventory("inv-1", "store-1", "", "prod-1", 100, 1)
	if err != entity.ErrInventoryInvalidWarehouseID {
		t.Errorf("expected ErrInventoryInvalidWarehouseID, got %v", err)
	}
}

func TestNewInventory_InvalidProductID(t *testing.T) {
	_, err := entity.NewInventory("inv-1", "store-1", "wh-1", "", 100, 1)
	if err != entity.ErrInventoryInvalidProductID {
		t.Errorf("expected ErrInventoryInvalidProductID, got %v", err)
	}
}

func TestNewInventory_InvalidBasePrice(t *testing.T) {
	_, err := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 0, 1)
	if err != entity.ErrInventoryInvalidBasePrice {
		t.Errorf("expected ErrInventoryInvalidBasePrice, got %v", err)
	}
}

func TestNewInventory_InvalidStock(t *testing.T) {
	_, err := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, -5)
	if err != entity.ErrInventoryInvalidStock {
		t.Errorf("expected ErrInventoryInvalidStock, got %v", err)
	}
}

func TestNewInventory_EventEmitted(t *testing.T) {
	inv, err := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventoryItemCreated); !ok {
		t.Fatalf("expected InventoryItemCreated, got %T", events[0])
	}
}

func TestInventory_UpdateStock_Success(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 10)
	inv.Events()

	err := inv.UpdateStock(25)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.InstantQty != 25 {
		t.Errorf("expected 25, got %d", inv.InstantQty)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventoryStockUpdated); !ok {
		t.Fatalf("expected InventoryStockUpdated, got %T", events[0])
	}
}

func TestInventory_UpdateStock_ToZero(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 10)
	inv.Events()

	err := inv.UpdateStock(0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.InstantQty != 0 {
		t.Errorf("expected 0, got %d", inv.InstantQty)
	}
}

func TestInventory_UpdateStock_Negative(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 10)
	err := inv.UpdateStock(-1)
	if err != entity.ErrInventoryInvalidStock {
		t.Errorf("expected ErrInventoryInvalidStock, got %v", err)
	}
}

func TestInventory_CanBeSold_AllConditionsMet(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	if !inv.CanBeSold() {
		t.Error("expected CanBeSold to be true")
	}
}

func TestInventory_CanBeSold_VendorInactive(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.SetVendorStatus(entity.VendorSaleStatusInactive)
	if inv.CanBeSold() {
		t.Error("expected CanBeSold to be false when vendor inactive")
	}
}

func TestInventory_CanBeSold_SystemInactive(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.SetSystemStatus(entity.SystemSaleStatusInactive)
	if inv.CanBeSold() {
		t.Error("expected CanBeSold to be false when system inactive")
	}
}

func TestInventory_CanBeSold_ZeroStock(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 0)
	if inv.CanBeSold() {
		t.Error("expected CanBeSold to be false when stock is zero")
	}
}

func TestInventory_CanBeSold_DraftStatus(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.SetVendorStatus(entity.VendorSaleStatusDraft)
	if inv.CanBeSold() {
		t.Error("expected CanBeSold to be false when vendor status is draft")
	}
}

func TestInventory_CanBeSold_ScheduledStatus(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.SetVendorStatus(entity.VendorSaleStatusScheduled)
	if inv.CanBeSold() {
		t.Error("expected CanBeSold to be false when vendor status is scheduled")
	}
}

func TestInventory_SetVendorStatus_Activate(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.SetVendorStatus(entity.VendorSaleStatusInactive)
	inv.Events()

	err := inv.SetVendorStatus(entity.VendorSaleStatusActive)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.VendorSaleStatus != entity.VendorSaleStatusActive {
		t.Errorf("expected active, got %s", inv.VendorSaleStatus)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventoryItemActivated); !ok {
		t.Fatalf("expected InventoryItemActivated, got %T", events[0])
	}
}

func TestInventory_SetVendorStatus_Deactivate(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetVendorStatus(entity.VendorSaleStatusInactive)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.VendorSaleStatus != entity.VendorSaleStatusInactive {
		t.Errorf("expected inactive, got %s", inv.VendorSaleStatus)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventoryItemDeactivated); !ok {
		t.Fatalf("expected InventoryItemDeactivated, got %T", events[0])
	}
}

func TestInventory_SetVendorStatus_Scheduled(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetVendorStatus(entity.VendorSaleStatusScheduled)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.VendorSaleStatus != entity.VendorSaleStatusScheduled {
		t.Errorf("expected scheduled, got %s", inv.VendorSaleStatus)
	}
}

func TestInventory_SetVendorStatus_Draft(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetVendorStatus(entity.VendorSaleStatusDraft)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.VendorSaleStatus != entity.VendorSaleStatusDraft {
		t.Errorf("expected draft, got %s", inv.VendorSaleStatus)
	}
}

func TestInventory_SetVendorStatus_SameStatus(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetVendorStatus(entity.VendorSaleStatusActive)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := inv.Events()
	if len(events) != 0 {
		t.Errorf("expected 0 events for same status, got %d", len(events))
	}
}

func TestInventory_SetVendorStatus_Invalid(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	err := inv.SetVendorStatus("invalid")
	if err != entity.ErrInventoryInvalidVendorStatus {
		t.Errorf("expected ErrInventoryInvalidVendorStatus, got %v", err)
	}
}

func TestInventory_SetSystemStatus_Unblocked(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.SetSystemStatus(entity.SystemSaleStatusInactive)
	inv.Events()

	err := inv.SetSystemStatus(entity.SystemSaleStatusActive)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.SystemSaleStatus != entity.SystemSaleStatusActive {
		t.Errorf("expected active, got %s", inv.SystemSaleStatus)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventorySystemUnblocked); !ok {
		t.Fatalf("expected InventorySystemUnblocked, got %T", events[0])
	}
}

func TestInventory_SetSystemStatus_Blocked(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetSystemStatus(entity.SystemSaleStatusInactive)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.SystemSaleStatus != entity.SystemSaleStatusInactive {
		t.Errorf("expected inactive, got %s", inv.SystemSaleStatus)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventorySystemBlocked); !ok {
		t.Fatalf("expected InventorySystemBlocked, got %T", events[0])
	}
}

func TestInventory_SetSystemStatus_SameStatus(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetSystemStatus(entity.SystemSaleStatusActive)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := inv.Events()
	if len(events) != 0 {
		t.Errorf("expected 0 events for same status, got %d", len(events))
	}
}

func TestInventory_SetSystemStatus_Invalid(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	err := inv.SetSystemStatus("invalid")
	if err != entity.ErrInventoryInvalidSystemStatus {
		t.Errorf("expected ErrInventoryInvalidSystemStatus, got %v", err)
	}
}

func TestInventory_SetSaleSchedule_Full(t *testing.T) {
	now := time.Now()
	start := now.Add(24 * time.Hour)
	end := now.Add(72 * time.Hour)
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetSaleSchedule(&start, &end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.StartAt == nil || !inv.StartAt.Equal(start) {
		t.Errorf("expected start %v, got %v", start, inv.StartAt)
	}
	if inv.EndAt == nil || !inv.EndAt.Equal(end) {
		t.Errorf("expected end %v, got %v", end, inv.EndAt)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventorySaleScheduled); !ok {
		t.Fatalf("expected InventorySaleScheduled, got %T", events[0])
	}
}

func TestInventory_SetSaleSchedule_StartOnly(t *testing.T) {
	now := time.Now()
	start := now.Add(24 * time.Hour)
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.SetSaleSchedule(&start, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.StartAt == nil || !inv.StartAt.Equal(start) {
		t.Errorf("expected start %v, got %v", start, inv.StartAt)
	}
	if inv.EndAt != nil {
		t.Errorf("expected nil end, got %v", inv.EndAt)
	}
}

func TestInventory_SetSaleSchedule_EndBeforeStart(t *testing.T) {
	now := time.Now()
	start := now.Add(72 * time.Hour)
	end := now.Add(24 * time.Hour)
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)

	err := inv.SetSaleSchedule(&start, &end)
	if err != entity.ErrInventoryInvalidTimeRange {
		t.Errorf("expected ErrInventoryInvalidTimeRange, got %v", err)
	}
}

func TestInventory_SetSaleSchedule_NilStartWithEnd(t *testing.T) {
	now := time.Now()
	end := now.Add(24 * time.Hour)
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)

	err := inv.SetSaleSchedule(nil, &end)
	if err != entity.ErrInventoryInvalidTimeRange {
		t.Errorf("expected ErrInventoryInvalidTimeRange, got %v", err)
	}
}

func TestInventory_UpdatePrice_Success(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	err := inv.UpdatePrice(200, 180)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.BasePrice != 200 {
		t.Errorf("expected base 200, got %f", inv.BasePrice)
	}
	if inv.FinalPrice != 180 {
		t.Errorf("expected final 180, got %f", inv.FinalPrice)
	}
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.InventoryPriceUpdated); !ok {
		t.Fatalf("expected InventoryPriceUpdated, got %T", events[0])
	}
}

func TestInventory_UpdatePrice_InvalidBasePrice(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	err := inv.UpdatePrice(0, 0)
	if err != entity.ErrInventoryInvalidBasePrice {
		t.Errorf("expected ErrInventoryInvalidBasePrice, got %v", err)
	}
}

func TestInventory_UpdatePrice_NegativeFinalPrice(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	err := inv.UpdatePrice(100, -10)
	if err != entity.ErrInventoryInvalidPrice {
		t.Errorf("expected ErrInventoryInvalidPrice, got %v", err)
	}
}

func TestInventory_Events_ClearedAfterCall(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	events1 := inv.Events()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}
	events2 := inv.Events()
	if len(events2) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events2))
	}
}
