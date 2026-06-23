package commission_test

import (
	"testing"

	"stock-service-version-three/internal/domain/commission"
)

func TestNewCommission_Success(t *testing.T) {
	sellerID := int64(1)
	rate := 0.10
	salesAmount := int64(2040000)

	comm, event, err := commission.NewCommission(sellerID, rate, salesAmount)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comm == nil {
		t.Fatal("expected commission, got nil")
	}
	if comm.ID() == 0 {
		t.Error("expected non-zero ID")
	}
	if comm.SellerID() != sellerID {
		t.Errorf("expected sellerID %d, got %d", sellerID, comm.SellerID())
	}
	if comm.Rate() != rate {
		t.Errorf("expected rate %f, got %f", rate, comm.Rate())
	}
	if comm.SalesAmount() != salesAmount {
		t.Errorf("expected salesAmount %d, got %d", salesAmount, comm.SalesAmount())
	}
	expectedAmount := int64(float64(salesAmount) * rate)
	if comm.Amount() != expectedAmount {
		t.Errorf("expected amount %d, got %d", expectedAmount, comm.Amount())
	}
	if comm.CreatedAt().IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if event == nil {
		t.Fatal("expected event, got nil")
	}
	if event.CommissionID != comm.ID() {
		t.Errorf("expected CommissionID %d, got %d", comm.ID(), event.CommissionID)
	}
	if event.SellerID != sellerID {
		t.Errorf("expected SellerID %d, got %d", sellerID, event.SellerID)
	}
	if event.Amount != expectedAmount {
		t.Errorf("expected Amount %d, got %d", expectedAmount, event.Amount)
	}
	if event.SalesAmount != salesAmount {
		t.Errorf("expected SalesAmount %d, got %d", salesAmount, event.SalesAmount)
	}
}

func TestNewCommission_InvalidRateZero(t *testing.T) {
	_, _, err := commission.NewCommission(1, 0, 2040000)
	if err == nil {
		t.Fatal("expected error for zero rate, got nil")
	}
	if err != commission.ErrInvalidCommissionRate {
		t.Errorf("expected ErrInvalidCommissionRate, got %v", err)
	}
}

func TestNewCommission_InvalidRateNegative(t *testing.T) {
	_, _, err := commission.NewCommission(1, -0.05, 2040000)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
	if err != commission.ErrInvalidCommissionRate {
		t.Errorf("expected ErrInvalidCommissionRate, got %v", err)
	}
}

func TestNewCommission_ZeroSalesAmount(t *testing.T) {
	comm, event, err := commission.NewCommission(1, 0.10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comm.Amount() != 0 {
		t.Errorf("expected amount 0, got %d", comm.Amount())
	}
	if event.Amount != 0 {
		t.Errorf("expected event amount 0, got %d", event.Amount)
	}
}

func TestCommission_Getters(t *testing.T) {
	sellerID := int64(42)
	rate := 0.15
	salesAmount := int64(1000000)

	comm, _, err := commission.NewCommission(sellerID, rate, salesAmount)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comm.SellerID() != sellerID {
		t.Errorf("expected SellerID %d, got %d", sellerID, comm.SellerID())
	}
	if comm.Rate() != rate {
		t.Errorf("expected Rate %f, got %f", rate, comm.Rate())
	}
	if comm.SalesAmount() != salesAmount {
		t.Errorf("expected SalesAmount %d, got %d", salesAmount, comm.SalesAmount())
	}
	if comm.Amount() != int64(float64(salesAmount)*rate) {
		t.Errorf("expected Amount %d, got %d", int64(float64(salesAmount)*rate), comm.Amount())
	}
	if comm.CreatedAt().IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}
