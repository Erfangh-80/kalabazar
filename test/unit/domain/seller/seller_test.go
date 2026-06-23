package seller_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/domain/seller"
)

func TestNewSeller(t *testing.T) {
	s := seller.NewSeller(101, "SportLine", "09120001111")

	if s.UserID != 101 {
		t.Errorf("expected UserID 101, got %d", s.UserID)
	}
	if s.Name != "SportLine" {
		t.Errorf("expected Name SportLine, got %s", s.Name)
	}
	if s.Phone != "09120001111" {
		t.Errorf("expected Phone 09120001111, got %s", s.Phone)
	}
	if s.Status != seller.SellerStatusUnverified {
		t.Errorf("expected Status UNVERIFIED, got %s", s.Status)
	}
	if s.Score != 0.0 {
		t.Errorf("expected Score 0.0, got %f", s.Score)
	}
	if s.Rank != "" {
		t.Errorf("expected empty Rank, got %s", s.Rank)
	}
	if s.ID != 0 {
		t.Errorf("expected ID 0 for new seller, got %d", s.ID)
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if s.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestVerifyKYC(t *testing.T) {
	s := seller.NewSeller(101, "SportLine", "09120001111")

	err := s.VerifyKYC()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Status != seller.SellerStatusVerified {
		t.Errorf("expected Status VERIFIED, got %s", s.Status)
	}
}

func TestVerifyKYC_AlreadyVerified(t *testing.T) {
	s := seller.NewSeller(101, "SportLine", "09120001111")
	s.VerifyKYC()

	err := s.VerifyKYC()
	if err != seller.ErrSellerAlreadyVerified {
		t.Errorf("expected ErrSellerAlreadyVerified, got %v", err)
	}
}

func TestUpdateRank(t *testing.T) {
	s := seller.NewSeller(101, "SportLine", "09120001111")
	s.VerifyKYC()

	s.UpdateRank(4.7, "A")
	if s.Score != 4.7 {
		t.Errorf("expected Score 4.7, got %f", s.Score)
	}
	if s.Rank != "A" {
		t.Errorf("expected Rank A, got %s", s.Rank)
	}
}

func TestValidateSellerName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"SportLine", false},
		{"Ali", false},
		{"A", true},
		{"", true},
		{"ThisNameIsWayTooLongForASellerStoreName", true},
	}
	for _, tt := range tests {
		err := seller.ValidateSellerName(tt.name)
		if tt.wantErr && err == nil {
			t.Errorf("ValidateSellerName(%q) expected error, got nil", tt.name)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("ValidateSellerName(%q) unexpected error: %v", tt.name, err)
		}
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		phone   string
		wantErr bool
	}{
		{"09120001111", false},
		{"02112345678", false},
		{"", true},
		{"123", true},
		{"notaphone", true},
	}
	for _, tt := range tests {
		err := seller.ValidatePhone(tt.phone)
		if tt.wantErr && err == nil {
			t.Errorf("ValidatePhone(%q) expected error, got nil", tt.phone)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("ValidatePhone(%q) unexpected error: %v", tt.phone, err)
		}
	}
}

func TestSellerVerifiedEvent(t *testing.T) {
	e := seller.SellerVerifiedEvent{SellerID: 42}
	if e.EventName() != "seller.verified" {
		t.Errorf("expected seller.verified, got %s", e.EventName())
	}
	if e.SellerID != 42 {
		t.Errorf("expected SellerID 42, got %d", e.SellerID)
	}
}

func TestSellerRankUpdatedEvent(t *testing.T) {
	e := seller.SellerRankUpdatedEvent{SellerID: 42, Score: 4.7, Rank: "A"}
	if e.EventName() != "seller.rank.updated" {
		t.Errorf("expected seller.rank.updated, got %s", e.EventName())
	}
	if e.SellerID != 42 {
		t.Errorf("expected SellerID 42, got %d", e.SellerID)
	}
	if e.Score != 4.7 {
		t.Errorf("expected Score 4.7, got %f", e.Score)
	}
	if e.Rank != "A" {
		t.Errorf("expected Rank A, got %s", e.Rank)
	}
}

func TestEventInterface(t *testing.T) {
	var evt seller.Event
	evt = seller.SellerVerifiedEvent{SellerID: 1}
	if _, ok := evt.(seller.Event); !ok {
		t.Error("SellerVerifiedEvent does not implement Event interface")
	}
	evt = seller.SellerRankUpdatedEvent{SellerID: 1, Score: 4.0, Rank: "B"}
	if _, ok := evt.(seller.Event); !ok {
		t.Error("SellerRankUpdatedEvent does not implement Event interface")
	}
}

func TestCreatedAtUpdatedAt(t *testing.T) {
	before := time.Now()
	s := seller.NewSeller(101, "SportLine", "09120001111")
	after := time.Now()

	if s.CreatedAt.Before(before) || s.CreatedAt.After(after) {
		t.Error("CreatedAt should be between before and after timestamps")
	}
	if s.UpdatedAt.Before(before) || s.UpdatedAt.After(after) {
		t.Error("UpdatedAt should be between before and after timestamps")
	}

	time.Sleep(time.Millisecond)
	s.VerifyKYC()
	if !s.UpdatedAt.After(s.CreatedAt) {
		t.Error("UpdatedAt should be updated after VerifyKYC")
	}
}
