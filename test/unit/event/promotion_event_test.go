package event_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestPromotionEvents_Created(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Summer Sale", "desc", now, now.Add(72*time.Hour), true, 20, true)
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.PromotionCreated)
	if !ok {
		t.Fatalf("expected PromotionCreated, got %T", events[0])
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.Title != "Summer Sale" {
		t.Errorf("expected Title Summer Sale, got %s", e.Title)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if e.EventName() != "promotion.campaign_created" {
		t.Errorf("expected promotion.campaign_created, got %s", e.EventName())
	}
}

func TestPromotionEvents_LinkedToProduct(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(72*time.Hour), false, 10, false)
	p.Events()
	p.LinkToProduct("prod-1")
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.PromotionCampaignLinkedToProduct)
	if !ok {
		t.Fatalf("expected PromotionCampaignLinkedToProduct, got %T", events[0])
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.ProductID != "prod-1" {
		t.Errorf("expected ProductID prod-1, got %s", e.ProductID)
	}
	if e.EventName() != "promotion.campaign_linked_to_product" {
		t.Errorf("expected promotion.campaign_linked_to_product, got %s", e.EventName())
	}
}

func TestPromotionEvents_Approved(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(72*time.Hour), true, 10, false)
	p.Events()
	p.Approve()
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.PromotionApproved)
	if !ok {
		t.Fatalf("expected PromotionApproved, got %T", events[0])
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.EventName() != "promotion.campaign_approved" {
		t.Errorf("expected promotion.campaign_approved, got %s", e.EventName())
	}
}

func TestPromotionEvents_Rejected(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(72*time.Hour), true, 10, false)
	p.Events()
	p.Reject()
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.PromotionRejected)
	if !ok {
		t.Fatalf("expected PromotionRejected, got %T", events[0])
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.EventName() != "promotion.campaign_rejected" {
		t.Errorf("expected promotion.campaign_rejected, got %s", e.EventName())
	}
}

func TestPromotionEvents_Activated(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(72*time.Hour), false, 10, false)
	p.Events()
	p.Activate()
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.PromotionActivated)
	if !ok {
		t.Fatalf("expected PromotionActivated, got %T", events[0])
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.EventName() != "promotion.campaign_activated" {
		t.Errorf("expected promotion.campaign_activated, got %s", e.EventName())
	}
}

func TestPromotionEvents_Deactivated(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(72*time.Hour), false, 10, false)
	p.Activate()
	p.Events()
	p.Deactivate()
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.PromotionDeactivated)
	if !ok {
		t.Fatalf("expected PromotionDeactivated, got %T", events[0])
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.EventName() != "promotion.campaign_deactivated" {
		t.Errorf("expected promotion.campaign_deactivated, got %s", e.EventName())
	}
}

func TestPromotionEvents_ClearedAfterEventsCall(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(72*time.Hour), false, 10, false)
	p.Events()
	remaining := p.Events()
	if len(remaining) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(remaining))
	}
}
