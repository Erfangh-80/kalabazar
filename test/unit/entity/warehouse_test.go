package entity_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func validAddress() entity.Address {
	return entity.Address{
		Street: "123 Main St", City: "Tehran", Country: "Iran",
	}
}

func TestNewWarehouse_Success(t *testing.T) {
	addr := validAddress()
	w, err := entity.NewWarehouse("wh-1", "usr-1", "Main Warehouse", addr, 1000, "public")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.ID != "wh-1" {
		t.Errorf("expected wh-1, got %s", w.ID)
	}
	if w.SellerID != "usr-1" {
		t.Errorf("expected usr-1, got %s", w.SellerID)
	}
	if w.Name != "Main Warehouse" {
		t.Errorf("expected Main Warehouse, got %s", w.Name)
	}
	if w.Address.Street != "123 Main St" {
		t.Errorf("expected 123 Main St, got %s", w.Address.Street)
	}
	if w.TotalCapacity != 1000 {
		t.Errorf("expected 1000, got %d", w.TotalCapacity)
	}
	if w.UsedCapacity != 0 {
		t.Errorf("expected 0, got %d", w.UsedCapacity)
	}
	if w.Status != entity.WarehouseStatusActive {
		t.Errorf("expected active, got %s", w.Status)
	}
}

func TestNewWarehouse_InvalidID(t *testing.T) {
	_, err := entity.NewWarehouse("", "usr-1", "WH", validAddress(), 100, "public")
	if err != entity.ErrWarehouseInvalidID {
		t.Errorf("expected ErrWarehouseInvalidID, got %v", err)
	}
}

func TestNewWarehouse_InvalidSellerID(t *testing.T) {
	_, err := entity.NewWarehouse("wh-1", "", "WH", validAddress(), 100, "public")
	if err != entity.ErrWarehouseInvalidSellerID {
		t.Errorf("expected ErrWarehouseInvalidSellerID, got %v", err)
	}
}

func TestNewWarehouse_InvalidName(t *testing.T) {
	_, err := entity.NewWarehouse("wh-1", "usr-1", "", validAddress(), 100, "public")
	if err != entity.ErrWarehouseInvalidName {
		t.Errorf("expected ErrWarehouseInvalidName, got %v", err)
	}
}

func TestNewWarehouse_NameTooLong(t *testing.T) {
	name := make([]byte, 256)
	for i := range name {
		name[i] = 'a'
	}
	_, err := entity.NewWarehouse("wh-1", "usr-1", string(name), validAddress(), 100, "public")
	if err != entity.ErrWarehouseNameTooLong {
		t.Errorf("expected ErrWarehouseNameTooLong, got %v", err)
	}
}

func TestNewWarehouse_InvalidAddress(t *testing.T) {
	addr := entity.Address{Street: "", City: "", Country: ""}
	_, err := entity.NewWarehouse("wh-1", "usr-1", "WH", addr, 100, "public")
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestNewWarehouse_InvalidCapacity(t *testing.T) {
	_, err := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 0, "public")
	if err != entity.ErrWarehouseInvalidCapacity {
		t.Errorf("expected ErrWarehouseInvalidCapacity, got %v", err)
	}
}

func TestNewWarehouse_NegativeCapacity(t *testing.T) {
	_, err := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), -10, "public")
	if err != entity.ErrWarehouseInvalidCapacity {
		t.Errorf("expected ErrWarehouseInvalidCapacity, got %v", err)
	}
}

func TestNewWarehouse_EventEmitted(t *testing.T) {
	w, err := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := w.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.WarehouseCreated)
	if !ok {
		t.Fatalf("expected WarehouseCreated, got %T", events[0])
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected wh-1, got %s", e.WarehouseID)
	}
}

func TestWarehouse_UpdateInfo_Success(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	newAddr := entity.Address{
		Street: "456 New St", City: "Shiraz", Country: "Iran",
	}
	err := w.UpdateInfo("Updated WH", newAddr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.Name != "Updated WH" {
		t.Errorf("expected 'Updated WH', got %s", w.Name)
	}
	if w.Address.Street != "456 New St" {
		t.Errorf("expected 456 New St, got %s", w.Address.Street)
	}
}

func TestWarehouse_UpdateInfo_InvalidName(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.UpdateInfo("", validAddress())
	if err != entity.ErrWarehouseInvalidName {
		t.Errorf("expected ErrWarehouseInvalidName, got %v", err)
	}
}

func TestWarehouse_UpdateInfo_InvalidAddress(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.UpdateInfo("New Name", entity.Address{})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestWarehouse_Activate_Success(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	w.Status = entity.WarehouseStatusInactive
	w.Events()

	err := w.Activate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.Status != entity.WarehouseStatusActive {
		t.Errorf("expected active, got %s", w.Status)
	}
	events := w.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.WarehouseActivated); !ok {
		t.Fatalf("expected WarehouseActivated, got %T", events[0])
	}
}

func TestWarehouse_Activate_AlreadyActive(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	w.Events()

	err := w.Activate()
	if err != entity.ErrWarehouseAlreadyActive {
		t.Errorf("expected ErrWarehouseAlreadyActive, got %v", err)
	}
}

func TestWarehouse_Deactivate_Success(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	w.Events()

	err := w.Deactivate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.Status != entity.WarehouseStatusInactive {
		t.Errorf("expected inactive, got %s", w.Status)
	}
	events := w.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.WarehouseDeactivated); !ok {
		t.Fatalf("expected WarehouseDeactivated, got %T", events[0])
	}
}

func TestWarehouse_Deactivate_AlreadyInactive(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	w.Deactivate()
	w.Events()

	err := w.Deactivate()
	if err != entity.ErrWarehouseAlreadyInactive {
		t.Errorf("expected ErrWarehouseAlreadyInactive, got %v", err)
	}
}

func TestWarehouse_IncreaseUsedCapacity_Success(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.IncreaseUsedCapacity(30)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.UsedCapacity != 30 {
		t.Errorf("expected 30, got %d", w.UsedCapacity)
	}
}

func TestWarehouse_IncreaseUsedCapacity_ExactCapacity(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.IncreaseUsedCapacity(100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !w.IsAtCapacity() {
		t.Error("expected warehouse to be at capacity")
	}
}

func TestWarehouse_IncreaseUsedCapacity_Exceeded(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.IncreaseUsedCapacity(101)
	if err != entity.ErrWarehouseCapacityExceeded {
		t.Errorf("expected ErrWarehouseCapacityExceeded, got %v", err)
	}
}

func TestWarehouse_IncreaseUsedCapacity_Inactive(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	w.Deactivate()
	err := w.IncreaseUsedCapacity(10)
	if err != entity.ErrWarehouseInactive {
		t.Errorf("expected ErrWarehouseInactive, got %v", err)
	}
}

func TestWarehouse_IncreaseUsedCapacity_Negative(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.IncreaseUsedCapacity(-10)
	if err != entity.ErrWarehouseInvalidUsedAmount {
		t.Errorf("expected ErrWarehouseInvalidUsedAmount, got %v", err)
	}
}

func TestWarehouse_DecreaseUsedCapacity_Success(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	w.IncreaseUsedCapacity(50)
	err := w.DecreaseUsedCapacity(20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.UsedCapacity != 30 {
		t.Errorf("expected 30, got %d", w.UsedCapacity)
	}
}

func TestWarehouse_DecreaseUsedCapacity_Negative(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.DecreaseUsedCapacity(-10)
	if err != entity.ErrWarehouseInvalidUsedAmount {
		t.Errorf("expected ErrWarehouseInvalidUsedAmount, got %v", err)
	}
}

func TestWarehouse_DecreaseUsedCapacity_BelowZero(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	err := w.DecreaseUsedCapacity(10)
	if err != entity.ErrWarehouseInvalidUsedAmount {
		t.Errorf("expected ErrWarehouseInvalidUsedAmount, got %v", err)
	}
}

func TestWarehouse_IsAtCapacity(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	if w.IsAtCapacity() {
		t.Error("expected not at capacity initially")
	}
	w.IncreaseUsedCapacity(100)
	if !w.IsAtCapacity() {
		t.Error("expected at capacity")
	}
}

func TestWarehouse_AvailableCapacity(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	if avail := w.AvailableCapacity(); avail != 100 {
		t.Errorf("expected 100, got %d", avail)
	}
	w.IncreaseUsedCapacity(30)
	if avail := w.AvailableCapacity(); avail != 70 {
		t.Errorf("expected 70, got %d", avail)
	}
}

func TestWarehouse_Events_ClearedAfterCall(t *testing.T) {
	w, _ := entity.NewWarehouse("wh-1", "usr-1", "WH", validAddress(), 100, "public")
	events1 := w.Events()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}
	events2 := w.Events()
	if len(events2) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events2))
	}
}
