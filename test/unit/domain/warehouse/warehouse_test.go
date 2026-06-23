package warehouse_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/domain/warehouse"
)

func TestNewWarehouse_Success(t *testing.T) {
	before := time.Now()

	w, err := warehouse.NewWarehouse("Tehran Central Warehouse", 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected warehouse, got nil")
	}
	if w.Name != "Tehran Central Warehouse" {
		t.Errorf("expected name 'Tehran Central Warehouse', got '%s'", w.Name)
	}
	if w.Capacity != 10000 {
		t.Errorf("expected capacity 10000, got %d", w.Capacity)
	}
	if w.CreatedAt.Before(before) || w.CreatedAt.After(time.Now()) {
		t.Errorf("CreatedAt should be between before and now, got %v", w.CreatedAt)
	}

	events := w.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	createdEvent, ok := events[0].(warehouse.WarehouseCreatedEvent)
	if !ok {
		t.Fatalf("expected WarehouseCreatedEvent, got %T", events[0])
	}
	if createdEvent.WarehouseID != w.ID {
		t.Errorf("expected WarehouseID %d, got %d", w.ID, createdEvent.WarehouseID)
	}
	if createdEvent.Name != w.Name {
		t.Errorf("expected Name '%s', got '%s'", w.Name, createdEvent.Name)
	}
}

func TestNewWarehouse_InvalidCapacity(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
	}{
		{"zero capacity", 0},
		{"negative capacity", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := warehouse.NewWarehouse("Test Warehouse", tt.capacity)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err != warehouse.ErrInvalidCapacity {
				t.Errorf("expected ErrInvalidCapacity, got %v", err)
			}
			if w != nil {
				t.Errorf("expected nil warehouse, got %v", w)
			}
		})
	}
}

func TestStoreWarehouseLink_Success(t *testing.T) {
	before := time.Now()

	link, err := warehouse.NewStoreWarehouseLink(1, 2, "primary")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link == nil {
		t.Fatal("expected link, got nil")
	}
	if link.StoreID != 1 {
		t.Errorf("expected StoreID 1, got %d", link.StoreID)
	}
	if link.WarehouseID != 2 {
		t.Errorf("expected WarehouseID 2, got %d", link.WarehouseID)
	}
	if link.LinkType != "primary" {
		t.Errorf("expected LinkType 'primary', got '%s'", link.LinkType)
	}
	if link.CreatedAt.Before(before) || link.CreatedAt.After(time.Now()) {
		t.Errorf("CreatedAt should be between before and now, got %v", link.CreatedAt)
	}

	events := link.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	linkedEvent, ok := events[0].(warehouse.WarehouseLinkedToStoreEvent)
	if !ok {
		t.Fatalf("expected WarehouseLinkedToStoreEvent, got %T", events[0])
	}
	if linkedEvent.WarehouseID != link.WarehouseID {
		t.Errorf("expected WarehouseID %d, got %d", link.WarehouseID, linkedEvent.WarehouseID)
	}
	if linkedEvent.StoreID != link.StoreID {
		t.Errorf("expected StoreID %d, got %d", link.StoreID, linkedEvent.StoreID)
	}
	if linkedEvent.LinkType != link.LinkType {
		t.Errorf("expected LinkType '%s', got '%s'", link.LinkType, linkedEvent.LinkType)
	}
}

func TestEvents_Cleared(t *testing.T) {
	w, err := warehouse.NewWarehouse("Test", 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(w.Events()) != 1 {
		t.Fatal("expected 1 event after creation")
	}

	w.ClearEvents()
	if len(w.Events()) != 0 {
		t.Errorf("expected 0 events after ClearEvents, got %d", len(w.Events()))
	}
}
