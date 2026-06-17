package event_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestStoreEvents_Created(t *testing.T) {
	phone := "09121234567"
	store, _ := entity.NewStore("s-1", "user-42", "Electronics Shop", &phone, nil, nil)
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreCreated)
	if !ok {
		t.Fatalf("expected StoreCreated, got %T", events[0])
	}
	if e.StoreID != "s-1" {
		t.Errorf("expected StoreID s-1, got %s", e.StoreID)
	}
	if e.StoreName != "Electronics Shop" {
		t.Errorf("expected StoreName Electronics Shop, got %s", e.StoreName)
	}
	if e.UserID != "user-42" {
		t.Errorf("expected UserID user-42, got %s", e.UserID)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if e.EventName() != "store.created" {
		t.Errorf("expected store.created, got %s", e.EventName())
	}
}

func TestStoreEvents_Updated(t *testing.T) {
	store, _ := entity.NewStore("s-1", "user-42", "Shop", nil, nil, nil)
	store.Events()
	newAddr := entity.Address{Street: "New St", City: "City", Country: "C"}
	if err := store.UpdateInfo("New Name", nil, &newAddr, nil); err != nil {
		t.Fatal(err)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreUpdated)
	if !ok {
		t.Fatalf("expected StoreUpdated, got %T", events[0])
	}
	if e.StoreID != "s-1" {
		t.Errorf("expected StoreID s-1, got %s", e.StoreID)
	}
	if e.EventName() != "store.updated" {
		t.Errorf("expected store.updated, got %s", e.EventName())
	}
}

func TestStoreEvents_Activated(t *testing.T) {
	store, _ := entity.NewStore("s-1", "user-42", "Shop", nil, nil, nil)
	store.Deactivate()
	store.Events()
	if err := store.Activate(); err != nil {
		t.Fatal(err)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreActivated)
	if !ok {
		t.Fatalf("expected StoreActivated, got %T", events[0])
	}
	if e.StoreID != "s-1" {
		t.Errorf("expected StoreID s-1, got %s", e.StoreID)
	}
	if e.EventName() != "store.activated" {
		t.Errorf("expected store.activated, got %s", e.EventName())
	}
}

func TestStoreEvents_Deactivated(t *testing.T) {
	store, _ := entity.NewStore("s-1", "user-42", "Shop", nil, nil, nil)
	store.Events()
	if err := store.Deactivate(); err != nil {
		t.Fatal(err)
	}
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreDeactivated)
	if !ok {
		t.Fatalf("expected StoreDeactivated, got %T", events[0])
	}
	if e.StoreID != "s-1" {
		t.Errorf("expected StoreID s-1, got %s", e.StoreID)
	}
	if e.EventName() != "store.deactivated" {
		t.Errorf("expected store.deactivated, got %s", e.EventName())
	}
}

func TestStoreEvents_ClearedAfterEventsCall(t *testing.T) {
	store, _ := entity.NewStore("s-1", "user-42", "Shop", nil, nil, nil)
	store.Events()
	remaining := store.Events()
	if len(remaining) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(remaining))
	}
}
