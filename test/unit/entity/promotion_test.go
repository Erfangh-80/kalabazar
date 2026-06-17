package entity_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestNewPromotion_Success(t *testing.T) {
	now := time.Now()
	p, err := entity.NewPromotion("promo-1", "Summer Sale", "desc", now, now.Add(72*time.Hour), false, 0, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ID != "promo-1" {
		t.Errorf("expected promo-1, got %s", p.ID)
	}
	if p.Title != "Summer Sale" {
		t.Errorf("expected Summer Sale, got %s", p.Title)
	}
	if p.IsActive {
		t.Error("expected IsActive to be false initially")
	}
}

func TestNewPromotion_InvalidID(t *testing.T) {
	now := time.Now()
	_, err := entity.NewPromotion("", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	if err != entity.ErrPromotionInvalidID {
		t.Errorf("expected ErrPromotionInvalidID, got %v", err)
	}
}

func TestNewPromotion_InvalidTitle(t *testing.T) {
	now := time.Now()
	_, err := entity.NewPromotion("promo-1", "", "", now, now.Add(24*time.Hour), false, 0, false)
	if err != entity.ErrPromotionInvalidTitle {
		t.Errorf("expected ErrPromotionInvalidTitle, got %v", err)
	}
}

func TestNewPromotion_StartAfterEnd(t *testing.T) {
	now := time.Now()
	_, err := entity.NewPromotion("promo-1", "Sale", "", now.Add(24*time.Hour), now, false, 0, false)
	if err != entity.ErrPromotionInvalidTimeRange {
		t.Errorf("expected ErrPromotionInvalidTimeRange, got %v", err)
	}
}

func TestNewPromotion_StartEqualToEnd(t *testing.T) {
	now := time.Now()
	_, err := entity.NewPromotion("promo-1", "Sale", "", now, now, false, 0, false)
	if err != entity.ErrPromotionInvalidTimeRange {
		t.Errorf("expected ErrPromotionInvalidTimeRange, got %v", err)
	}
}

func TestNewPromotion_RequiresApprovalDefaultStatus(t *testing.T) {
	now := time.Now()
	p, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), true, 0, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ApprovalStatus != entity.PromotionApprovalPending {
		t.Errorf("expected pending approval, got %s", p.ApprovalStatus)
	}
}

func TestNewPromotion_NoApprovalDefaultStatus(t *testing.T) {
	now := time.Now()
	p, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ApprovalStatus != entity.PromotionApprovalNone {
		t.Errorf("expected none approval, got %s", p.ApprovalStatus)
	}
}

func TestNewPromotion_EventEmitted(t *testing.T) {
	now := time.Now()
	p, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.PromotionCreated); !ok {
		t.Fatalf("expected PromotionCreated, got %T", events[0])
	}
}

func TestPromotion_Approve_Success(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), true, 0, false)
	p.Events()

	err := p.Approve()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ApprovalStatus != entity.PromotionApprovalApproved {
		t.Errorf("expected approved, got %s", p.ApprovalStatus)
	}
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.PromotionApproved); !ok {
		t.Fatalf("expected PromotionApproved, got %T", events[0])
	}
}

func TestPromotion_Approve_AlreadyApproved(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), true, 0, false)
	p.Approve()
	p.Events()

	err := p.Approve()
	if err != entity.ErrPromotionAlreadyApproved {
		t.Errorf("expected ErrPromotionAlreadyApproved, got %v", err)
	}
}

func TestPromotion_Approve_NotRequired(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)

	err := p.Approve()
	if err != entity.ErrPromotionApprovalNotRequired {
		t.Errorf("expected ErrPromotionApprovalNotRequired, got %v", err)
	}
}

func TestPromotion_Reject_Success(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), true, 0, false)
	p.Events()

	err := p.Reject()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ApprovalStatus != entity.PromotionApprovalRejected {
		t.Errorf("expected rejected, got %s", p.ApprovalStatus)
	}
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.PromotionRejected); !ok {
		t.Fatalf("expected PromotionRejected, got %T", events[0])
	}
}

func TestPromotion_Reject_AlreadyRejected(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), true, 0, false)
	p.Reject()
	p.Events()

	err := p.Reject()
	if err != entity.ErrPromotionAlreadyRejected {
		t.Errorf("expected ErrPromotionAlreadyRejected, got %v", err)
	}
}

func TestPromotion_Reject_NotRequired(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)

	err := p.Reject()
	if err != entity.ErrPromotionApprovalNotRequired {
		t.Errorf("expected ErrPromotionApprovalNotRequired, got %v", err)
	}
}

func TestPromotion_Activate_Success(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	p.Events()

	err := p.Activate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !p.IsActive {
		t.Error("expected IsActive to be true")
	}
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.PromotionActivated); !ok {
		t.Fatalf("expected PromotionActivated, got %T", events[0])
	}
}

func TestPromotion_Activate_NotApproved(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), true, 0, false)

	err := p.Activate()
	if err != entity.ErrPromotionNotApproved {
		t.Errorf("expected ErrPromotionNotApproved, got %v", err)
	}
}

func TestPromotion_Activate_Rejected(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), true, 0, false)
	p.Reject()

	err := p.Activate()
	if err != entity.ErrPromotionNotApproved {
		t.Errorf("expected ErrPromotionNotApproved, got %v", err)
	}
}

func TestPromotion_Activate_AlreadyActive(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	p.Activate()
	p.Events()

	err := p.Activate()
	if err != entity.ErrPromotionAlreadyActive {
		t.Errorf("expected ErrPromotionAlreadyActive, got %v", err)
	}
}

func TestPromotion_Deactivate_Success(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	p.Activate()
	p.Events()

	err := p.Deactivate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.IsActive {
		t.Error("expected IsActive to be false")
	}
	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(event.PromotionDeactivated); !ok {
		t.Fatalf("expected PromotionDeactivated, got %T", events[0])
	}
}

func TestPromotion_Deactivate_AlreadyInactive(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	p.Events()

	err := p.Deactivate()
	if err != entity.ErrPromotionAlreadyInactive {
		t.Errorf("expected ErrPromotionAlreadyInactive, got %v", err)
	}
}

func TestPromotion_Events_ClearedAfterCall(t *testing.T) {
	now := time.Now()
	p, _ := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	events1 := p.Events()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}
	events2 := p.Events()
	if len(events2) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events2))
	}
}

func TestNewPromotion_DiscountPercent(t *testing.T) {
	now := time.Now()
	p, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 25, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.DiscountPercent != 25 {
		t.Errorf("expected discount 25, got %f", p.DiscountPercent)
	}
}

func TestNewPromotion_InvalidDiscountPercent_Negative(t *testing.T) {
	now := time.Now()
	_, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, -10, false)
	if err != entity.ErrPromotionInvalidDiscountPct {
		t.Errorf("expected ErrPromotionInvalidDiscountPct, got %v", err)
	}
}

func TestNewPromotion_InvalidDiscountPercent_Over100(t *testing.T) {
	now := time.Now()
	_, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 150, false)
	if err != entity.ErrPromotionInvalidDiscountPct {
		t.Errorf("expected ErrPromotionInvalidDiscountPct, got %v", err)
	}
}

func TestNewPromotion_IsCountdown(t *testing.T) {
	now := time.Now()
	p, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !p.IsCountdown {
		t.Error("expected IsCountdown to be true")
	}
}

func TestNewPromotion_NotCountdown(t *testing.T) {
	now := time.Now()
	p, err := entity.NewPromotion("promo-1", "Sale", "", now, now.Add(24*time.Hour), false, 0, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.IsCountdown {
		t.Error("expected IsCountdown to be false")
	}
}
