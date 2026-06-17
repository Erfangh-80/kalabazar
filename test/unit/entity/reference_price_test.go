package entity_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestNewReferencePrice_Success(t *testing.T) {
	rp, err := entity.NewReferencePrice("rp-1", "prod-1", 150.0, "manual")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rp.ID != "rp-1" {
		t.Errorf("expected rp-1, got %s", rp.ID)
	}
	if rp.ProductID != "prod-1" {
		t.Errorf("expected prod-1, got %s", rp.ProductID)
	}
	if rp.Price != 150.0 {
		t.Errorf("expected 150.0, got %f", rp.Price)
	}
	if rp.Source != "manual" {
		t.Errorf("expected manual, got %s", rp.Source)
	}
}

func TestNewReferencePrice_InvalidID(t *testing.T) {
	_, err := entity.NewReferencePrice("", "prod-1", 100, "source")
	if err != entity.ErrReferencePriceInvalidID {
		t.Errorf("expected ErrReferencePriceInvalidID, got %v", err)
	}
}

func TestNewReferencePrice_InvalidProductID(t *testing.T) {
	_, err := entity.NewReferencePrice("rp-1", "", 100, "source")
	if err != entity.ErrReferencePriceInvalidProductID {
		t.Errorf("expected ErrReferencePriceInvalidProductID, got %v", err)
	}
}

func TestNewReferencePrice_ZeroPrice(t *testing.T) {
	_, err := entity.NewReferencePrice("rp-1", "prod-1", 0, "source")
	if err != entity.ErrReferencePriceInvalidPrice {
		t.Errorf("expected ErrReferencePriceInvalidPrice, got %v", err)
	}
}

func TestNewReferencePrice_NegativePrice(t *testing.T) {
	_, err := entity.NewReferencePrice("rp-1", "prod-1", -10, "source")
	if err != entity.ErrReferencePriceInvalidPrice {
		t.Errorf("expected ErrReferencePriceInvalidPrice, got %v", err)
	}
}

func TestNewReferencePrice_InvalidSource(t *testing.T) {
	_, err := entity.NewReferencePrice("rp-1", "prod-1", 100, "")
	if err != entity.ErrReferencePriceInvalidSource {
		t.Errorf("expected ErrReferencePriceInvalidSource, got %v", err)
	}
}

func TestNewReferencePrice_EventEmitted(t *testing.T) {
	rp, err := entity.NewReferencePrice("rp-1", "prod-1", 200, "market")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := rp.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.ReferencePriceCreated)
	if !ok {
		t.Fatalf("expected ReferencePriceCreated, got %T", events[0])
	}
	if e.ReferencePriceID != "rp-1" {
		t.Errorf("expected rp-1, got %s", e.ReferencePriceID)
	}
	if e.Price != 200 {
		t.Errorf("expected 200, got %f", e.Price)
	}
}

func TestReferencePrice_Events_ClearedAfterCall(t *testing.T) {
	rp, _ := entity.NewReferencePrice("rp-1", "prod-1", 100, "test")
	events1 := rp.Events()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}
	events2 := rp.Events()
	if len(events2) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events2))
	}
}

func TestCalculateFinalPrice_NoDiscount(t *testing.T) {
	price, err := entity.CalculateFinalPrice(100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 100 {
		t.Errorf("expected 100, got %f", price)
	}
}

func TestCalculateFinalPrice_WithDiscount(t *testing.T) {
	price, err := entity.CalculateFinalPrice(200, 25)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 150 {
		t.Errorf("expected 150, got %f", price)
	}
}

func TestCalculateFinalPrice_FullDiscount(t *testing.T) {
	price, err := entity.CalculateFinalPrice(100, 100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 0 {
		t.Errorf("expected 0, got %f", price)
	}
}

func TestCalculateFinalPrice_ZeroBasePrice(t *testing.T) {
	_, err := entity.CalculateFinalPrice(0, 10)
	if err != entity.ErrInvalidBasePrice {
		t.Errorf("expected ErrInvalidBasePrice, got %v", err)
	}
}

func TestCalculateFinalPrice_NegativeBasePrice(t *testing.T) {
	_, err := entity.CalculateFinalPrice(-50, 10)
	if err != entity.ErrInvalidBasePrice {
		t.Errorf("expected ErrInvalidBasePrice, got %v", err)
	}
}

func TestCalculateFinalPrice_NegativeDiscount(t *testing.T) {
	_, err := entity.CalculateFinalPrice(100, -10)
	if err != entity.ErrInvalidDiscountPercent {
		t.Errorf("expected ErrInvalidDiscountPercent, got %v", err)
	}
}

func TestCalculateFinalPrice_DiscountOver100(t *testing.T) {
	_, err := entity.CalculateFinalPrice(100, 150)
	if err != entity.ErrInvalidDiscountPercent {
		t.Errorf("expected ErrInvalidDiscountPercent, got %v", err)
	}
}

func TestCalculateFinalPrice_FractionalDiscount(t *testing.T) {
	price, err := entity.CalculateFinalPrice(200, 12.5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 175 {
		t.Errorf("expected 175, got %f", price)
	}
}
