package event_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestCommissionEvents_RuleCreated(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 100000, 1000000, 1)
	events := c.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.CommissionRuleCreated)
	if !ok {
		t.Fatalf("expected CommissionRuleCreated, got %T", events[0])
	}
	if e.CommissionID != "comm-1" {
		t.Errorf("expected CommissionID comm-1, got %s", e.CommissionID)
	}
	if e.ProductID != "prod-1" {
		t.Errorf("expected ProductID prod-1, got %s", e.ProductID)
	}
	if e.RatePercent != 10 {
		t.Errorf("expected RatePercent 10, got %f", e.RatePercent)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if e.EventName() != "commission.rule.created" {
		t.Errorf("expected commission.rule.created, got %s", e.EventName())
	}
}

func TestCommissionEvents_Calculated(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 100000, 1000000, 1)
	c.Events()
	amount, _ := c.Calculate(400000, 1)
	events := c.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.CommissionCalculated)
	if !ok {
		t.Fatalf("expected CommissionCalculated, got %T", events[0])
	}
	if e.CommissionID != "comm-1" {
		t.Errorf("expected CommissionID comm-1, got %s", e.CommissionID)
	}
	if e.SaleAmount != 400000 {
		t.Errorf("expected SaleAmount 400000, got %f", e.SaleAmount)
	}
	if e.CommissionAmount != amount {
		t.Errorf("expected CommissionAmount %f, got %f", amount, e.CommissionAmount)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if e.EventName() != "commission.calculated" {
		t.Errorf("expected commission.calculated, got %s", e.EventName())
	}
}

func TestCommissionEvents_ClearedAfterEventsCall(t *testing.T) {
	c, _ := entity.NewCommission("comm-1", "prod-1", "retail", 10, 100000, 1000000, 1)
	c.Events()
	remaining := c.Events()
	if len(remaining) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(remaining))
	}
}
