package payout_test

import (
	"testing"

	"stock-service-version-three/internal/domain/payout"
)

func TestNewPayout_Success(t *testing.T) {
	sellerID := int64(1)
	amount := int64(1836000)

	p, err := payout.NewPayout(sellerID, amount)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p == nil {
		t.Fatal("expected payout, got nil")
	}
	if p.ID() == 0 {
		t.Error("expected non-zero ID")
	}
	if p.SellerID() != sellerID {
		t.Errorf("expected SellerID %d, got %d", sellerID, p.SellerID())
	}
	if p.Amount() != amount {
		t.Errorf("expected Amount %d, got %d", amount, p.Amount())
	}
	if p.Status() != payout.PayoutStatusPending {
		t.Errorf("expected status %s, got %s", payout.PayoutStatusPending, p.Status())
	}
	if p.CreatedAt().IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestNewPayout_ZeroAmount(t *testing.T) {
	_, err := payout.NewPayout(1, 0)
	if err == nil {
		t.Fatal("expected error for zero amount, got nil")
	}
	if err != payout.ErrInvalidPayoutAmount {
		t.Errorf("expected ErrInvalidPayoutAmount, got %v", err)
	}
}

func TestNewPayout_NegativeAmount(t *testing.T) {
	_, err := payout.NewPayout(1, -100)
	if err == nil {
		t.Fatal("expected error for negative amount, got nil")
	}
	if err != payout.ErrInvalidPayoutAmount {
		t.Errorf("expected ErrInvalidPayoutAmount, got %v", err)
	}
}

func TestPayout_Execute_Success(t *testing.T) {
	p, err := payout.NewPayout(1, 1836000)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	event, err := p.Execute()
	if err != nil {
		t.Fatalf("expected no error on execute, got %v", err)
	}
	if p.Status() != payout.PayoutStatusExecuted {
		t.Errorf("expected status %s, got %s", payout.PayoutStatusExecuted, p.Status())
	}
	if event == nil {
		t.Fatal("expected event, got nil")
	}
	if event.PayoutID != p.ID() {
		t.Errorf("expected PayoutID %d, got %d", p.ID(), event.PayoutID)
	}
	if event.SellerID != p.SellerID() {
		t.Errorf("expected SellerID %d, got %d", p.SellerID(), event.SellerID)
	}
	if event.Amount != p.Amount() {
		t.Errorf("expected Amount %d, got %d", p.Amount(), event.Amount)
	}
}

func TestPayout_Execute_AlreadyExecuted(t *testing.T) {
	p, err := payout.NewPayout(1, 1836000)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = p.Execute()
	if err != nil {
		t.Fatalf("expected no error on first execute, got %v", err)
	}

	_, err = p.Execute()
	if err == nil {
		t.Fatal("expected error on second execute, got nil")
	}
	if err != payout.ErrPayoutAlreadyExecuted {
		t.Errorf("expected ErrPayoutAlreadyExecuted, got %v", err)
	}
}

func TestPayout_Execute_FailedStatus(t *testing.T) {
	p, err := payout.NewPayout(1, 1836000)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = p.Execute()
	if err != nil {
		t.Fatalf("expected no error on execute, got %v", err)
	}

	if _, err := p.Execute(); err != payout.ErrPayoutAlreadyExecuted {
		t.Errorf("expected ErrPayoutAlreadyExecuted after execute, got %v", err)
	}
}

func TestPayout_Getters(t *testing.T) {
	sellerID := int64(42)
	amount := int64(500000)

	p, err := payout.NewPayout(sellerID, amount)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.SellerID() != sellerID {
		t.Errorf("expected SellerID %d, got %d", sellerID, p.SellerID())
	}
	if p.Amount() != amount {
		t.Errorf("expected Amount %d, got %d", amount, p.Amount())
	}
	if p.Status() != payout.PayoutStatusPending {
		t.Errorf("expected status %s", p.Status())
	}
}
