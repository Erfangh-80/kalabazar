package entity_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestNewCommission_Success(t *testing.T) {
	c, err := entity.NewCommission("comm-1", "prod-1", "retail", 10.0, 50000, 500000, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.ID != "comm-1" {
		t.Errorf("expected comm-1, got %s", c.ID)
	}
	if c.ProductID != "prod-1" {
		t.Errorf("expected prod-1, got %s", c.ProductID)
	}
	if c.RatePercent != 10.0 {
		t.Errorf("expected 10.0, got %f", c.RatePercent)
	}
	if c.MinPrice != 50000 {
		t.Errorf("expected 50000, got %f", c.MinPrice)
	}
	if c.MaxPrice != 500000 {
		t.Errorf("expected 500000, got %f", c.MaxPrice)
	}
	if c.MinQty != 1 {
		t.Errorf("expected 1, got %d", c.MinQty)
	}
}

func TestNewCommission_InvalidID(t *testing.T) {
	_, err := entity.NewCommission("", "prod-1", "retail", 10, 0, 100, 1)
	if err != entity.ErrCommissionInvalidID {
		t.Errorf("expected ErrCommissionInvalidID, got %v", err)
	}
}

func TestNewCommission_InvalidProductID(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "", "retail", 10, 0, 100, 1)
	if err != entity.ErrCommissionInvalidProductID {
		t.Errorf("expected ErrCommissionInvalidProductID, got %v", err)
	}
}

func TestNewCommission_ZeroRate(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "prod-1", "retail", 0, 0, 100, 1)
	if err != entity.ErrCommissionInvalidRate {
		t.Errorf("expected ErrCommissionInvalidRate, got %v", err)
	}
}

func TestNewCommission_NegativeRate(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "prod-1", "retail", -5, 0, 100, 1)
	if err != entity.ErrCommissionInvalidRate {
		t.Errorf("expected ErrCommissionInvalidRate, got %v", err)
	}
}

func TestNewCommission_RateOver100(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "prod-1", "retail", 150, 0, 100, 1)
	if err != entity.ErrCommissionInvalidRate {
		t.Errorf("expected ErrCommissionInvalidRate, got %v", err)
	}
}

func TestNewCommission_MinPriceNegative(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "prod-1", "retail", 10, -100, 100, 1)
	if err != entity.ErrCommissionInvalidPriceRange {
		t.Errorf("expected ErrCommissionInvalidPriceRange, got %v", err)
	}
}

func TestNewCommission_MaxPriceNegative(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "prod-1", "retail", 10, 0, -100, 1)
	if err != entity.ErrCommissionInvalidPriceRange {
		t.Errorf("expected ErrCommissionInvalidPriceRange, got %v", err)
	}
}

func TestNewCommission_MinAboveMax(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "prod-1", "retail", 10, 500, 100, 1)
	if err != entity.ErrCommissionInvalidPriceRange {
		t.Errorf("expected ErrCommissionInvalidPriceRange, got %v", err)
	}
}

func TestNewCommission_InvalidMinQty(t *testing.T) {
	_, err := entity.NewCommission("comm-1", "prod-1", "retail", 10, 0, 100, -1)
	if err != entity.ErrCommissionInvalidMinQty {
		t.Errorf("expected ErrCommissionInvalidMinQty, got %v", err)
	}
}

func TestNewCommission_EventEmitted(t *testing.T) {
	c, err := entity.NewCommission("comm-1", "prod-1", "retail", 10, 0, 100, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := c.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.CommissionRuleCreated)
	if !ok {
		t.Fatalf("expected CommissionRuleCreated, got %T", events[0])
	}
	if e.CommissionID != "comm-1" {
		t.Errorf("expected comm-1, got %s", e.CommissionID)
	}
}

func TestCommission_Calculate_Success(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 50000, 500000, 1)
	amount, err := c.Calculate(200000, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if amount != 20000 {
		t.Errorf("expected 20000, got %f", amount)
	}
}

func TestCommission_Calculate_PriceBelowMin(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 50000, 500000, 1)
	_, err := c.Calculate(10000, 1)
	if err != entity.ErrCommissionConditionsNotMet {
		t.Errorf("expected ErrCommissionConditionsNotMet, got %v", err)
	}
}

func TestCommission_Calculate_PriceAboveMax(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 50000, 500000, 1)
	_, err := c.Calculate(600000, 1)
	if err != entity.ErrCommissionConditionsNotMet {
		t.Errorf("expected ErrCommissionConditionsNotMet, got %v", err)
	}
}

func TestCommission_Calculate_QtyBelowMin(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 50000, 500000, 3)
	_, err := c.Calculate(200000, 1)
	if err != entity.ErrCommissionConditionsNotMet {
		t.Errorf("expected ErrCommissionConditionsNotMet, got %v", err)
	}
}

func TestCommission_Calculate_PriceAtMinBoundary(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 50000, 500000, 1)
	amount, err := c.Calculate(50000, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if amount != 5000 {
		t.Errorf("expected 5000, got %f", amount)
	}
}

func TestCommission_Calculate_PriceAtMaxBoundary(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 50000, 500000, 1)
	amount, err := c.Calculate(500000, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if amount != 50000 {
		t.Errorf("expected 50000, got %f", amount)
	}
}

func TestCommission_Calculate_QtyAtMinBoundary(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 0, 1000, 5)
	amount, err := c.Calculate(100, 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if amount != 10 {
		t.Errorf("expected 10, got %f", amount)
	}
}

func TestCommission_Events_ClearedAfterCall(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 0, 100, 1)
	events1 := c.Events()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}
	events2 := c.Events()
	if len(events2) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events2))
	}
}
