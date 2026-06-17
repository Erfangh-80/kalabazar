package event_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestStoreCategoryEvents_CategoryAllowed(t *testing.T) {
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	events := sc.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreCategoryAllowed)
	if !ok {
		t.Fatalf("expected StoreCategoryAllowed, got %T", events[0])
	}
	if e.StoreID != "store-1" {
		t.Errorf("expected StoreID store-1, got %s", e.StoreID)
	}
	if e.CategoryID != "cat-7" {
		t.Errorf("expected CategoryID cat-7, got %s", e.CategoryID)
	}
	if e.Status != "pending" {
		t.Errorf("expected status pending, got %s", e.Status)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if e.EventName() != "store.category_allowed" {
		t.Errorf("expected store.category_allowed, got %s", e.EventName())
	}
}

func TestStoreCategoryEvents_Approved(t *testing.T) {
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	sc.Events()
	sc.Approve()
	events := sc.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreCategoryAllowed)
	if !ok {
		t.Fatalf("expected StoreCategoryAllowed, got %T", events[0])
	}
	if e.Status != "approved" {
		t.Errorf("expected status approved, got %s", e.Status)
	}
	if e.EventName() != "store.category_allowed" {
		t.Errorf("expected store.category_allowed, got %s", e.EventName())
	}
}
