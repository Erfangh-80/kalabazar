package event_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestWarehouseEvents_Created(t *testing.T) {
	addr := entity.Address{Street: "St", City: "C", Country: "Iran"}
	wh, _ := entity.NewWarehouse("wh-1", "seller-1", "Warehouse", addr, 1000, "public")
	events := wh.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.WarehouseCreated)
	if !ok {
		t.Fatalf("expected WarehouseCreated, got %T", events[0])
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected WarehouseID wh-1, got %s", e.WarehouseID)
	}
	if e.SellerID != "seller-1" {
		t.Errorf("expected SellerID seller-1, got %s", e.SellerID)
	}
	if e.WarehouseName != "Warehouse" {
		t.Errorf("expected WarehouseName Warehouse, got %s", e.WarehouseName)
	}
	if e.EventName() != "warehouse.created" {
		t.Errorf("expected warehouse.created, got %s", e.EventName())
	}
}

func TestWarehouseEvents_LinkedToStore(t *testing.T) {
	addr := entity.Address{Street: "St", City: "C", Country: "Iran"}
	wh, _ := entity.NewWarehouse("wh-1", "seller-1", "Warehouse", addr, 1000, "public")
	wh.Events()

	if err := wh.LinkToStore("store-1", "primary"); err != nil {
		t.Fatal(err)
	}
	events := wh.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.WarehouseLinkedToStore)
	if !ok {
		t.Fatalf("expected WarehouseLinkedToStore, got %T", events[0])
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected WarehouseID wh-1, got %s", e.WarehouseID)
	}
	if e.StoreID != "store-1" {
		t.Errorf("expected StoreID store-1, got %s", e.StoreID)
	}
	if e.RelationType != "primary" {
		t.Errorf("expected RelationType primary, got %s", e.RelationType)
	}
	if e.EventName() != "warehouse.linked_to_store" {
		t.Errorf("expected warehouse.linked_to_store, got %s", e.EventName())
	}
}

func TestWarehouseEvents_Updated(t *testing.T) {
	addr := entity.Address{Street: "St", City: "C", Country: "Iran"}
	wh, _ := entity.NewWarehouse("wh-1", "seller-1", "Warehouse", addr, 1000, "public")
	wh.Events()

	newAddr := entity.Address{Street: "New St", City: "C", Country: "Iran"}
	if err := wh.UpdateInfo("New Name", newAddr); err != nil {
		t.Fatal(err)
	}
	events := wh.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.WarehouseUpdated)
	if !ok {
		t.Fatalf("expected WarehouseUpdated, got %T", events[0])
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected WarehouseID wh-1, got %s", e.WarehouseID)
	}
	if e.EventName() != "warehouse.updated" {
		t.Errorf("expected warehouse.updated, got %s", e.EventName())
	}
}

func TestWarehouseEvents_Activated(t *testing.T) {
	addr := entity.Address{Street: "St", City: "C", Country: "Iran"}
	wh, _ := entity.NewWarehouse("wh-1", "seller-1", "Warehouse", addr, 1000, "public")
	wh.Deactivate()
	wh.Events()

	if err := wh.Activate(); err != nil {
		t.Fatal(err)
	}
	events := wh.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.WarehouseActivated)
	if !ok {
		t.Fatalf("expected WarehouseActivated, got %T", events[0])
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected WarehouseID wh-1, got %s", e.WarehouseID)
	}
	if e.EventName() != "warehouse.activated" {
		t.Errorf("expected warehouse.activated, got %s", e.EventName())
	}
}

func TestWarehouseEvents_Deactivated(t *testing.T) {
	addr := entity.Address{Street: "St", City: "C", Country: "Iran"}
	wh, _ := entity.NewWarehouse("wh-1", "seller-1", "Warehouse", addr, 1000, "public")
	wh.Events()

	if err := wh.Deactivate(); err != nil {
		t.Fatal(err)
	}
	events := wh.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.WarehouseDeactivated)
	if !ok {
		t.Fatalf("expected WarehouseDeactivated, got %T", events[0])
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected WarehouseID wh-1, got %s", e.WarehouseID)
	}
	if e.EventName() != "warehouse.deactivated" {
		t.Errorf("expected warehouse.deactivated, got %s", e.EventName())
	}
}

func TestWarehouseEvents_ClearedAfterEventsCall(t *testing.T) {
	addr := entity.Address{Street: "St", City: "C", Country: "Iran"}
	wh, _ := entity.NewWarehouse("wh-1", "seller-1", "Warehouse", addr, 1000, "public")
	wh.Events()
	remaining := wh.Events()
	if len(remaining) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(remaining))
	}
}
