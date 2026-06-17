package event_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestReferencePriceEvents_Created(t *testing.T) {
	rp, _ := entity.NewReferencePrice("rp-1", "prod-1", 550000, "DigiKala")
	events := rp.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.ReferencePriceCreated)
	if !ok {
		t.Fatalf("expected ReferencePriceCreated, got %T", events[0])
	}
	if e.ReferencePriceID != "rp-1" {
		t.Errorf("expected ReferencePriceID rp-1, got %s", e.ReferencePriceID)
	}
	if e.ProductID != "prod-1" {
		t.Errorf("expected ProductID prod-1, got %s", e.ProductID)
	}
	if e.Price != 550000 {
		t.Errorf("expected Price 550000, got %f", e.Price)
	}
	if e.Source != "DigiKala" {
		t.Errorf("expected Source DigiKala, got %s", e.Source)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if e.EventName() != "pricing.reference_price_recorded" {
		t.Errorf("expected pricing.reference_price_recorded, got %s", e.EventName())
	}
}

func TestReferencePriceEvents_ClearedAfterEventsCall(t *testing.T) {
	rp, _ := entity.NewReferencePrice("rp-1", "prod-1", 550000, "DigiKala")
	rp.Events()
	remaining := rp.Events()
	if len(remaining) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(remaining))
	}
}
