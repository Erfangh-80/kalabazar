package settlement_test

import (
	"testing"

	"stock-service-version-three/internal/domain/settlement"
)

func TestNewSettlement_Success(t *testing.T) {
	sellerID := int64(1)
	grossSales := int64(2040000)
	commission := int64(204000)

	sett, event, err := settlement.NewSettlement(sellerID, grossSales, commission)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sett == nil {
		t.Fatal("expected settlement, got nil")
	}
	if sett.ID() == 0 {
		t.Error("expected non-zero ID")
	}
	if sett.SellerID() != sellerID {
		t.Errorf("expected SellerID %d, got %d", sellerID, sett.SellerID())
	}
	if sett.GrossSales() != grossSales {
		t.Errorf("expected GrossSales %d, got %d", grossSales, sett.GrossSales())
	}
	if sett.Commission() != commission {
		t.Errorf("expected Commission %d, got %d", commission, sett.Commission())
	}
	expectedNet := grossSales - commission
	if sett.NetAmount() != expectedNet {
		t.Errorf("expected NetAmount %d, got %d", expectedNet, sett.NetAmount())
	}
	if sett.CreatedAt().IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if event == nil {
		t.Fatal("expected event, got nil")
	}
	if event.SettlementID != sett.ID() {
		t.Errorf("expected SettlementID %d, got %d", sett.ID(), event.SettlementID)
	}
	if event.SellerID != sellerID {
		t.Errorf("expected SellerID %d, got %d", sellerID, event.SellerID)
	}
	if event.GrossSales != grossSales {
		t.Errorf("expected GrossSales %d, got %d", grossSales, event.GrossSales)
	}
	if event.Commission != commission {
		t.Errorf("expected Commission %d, got %d", commission, event.Commission)
	}
	if event.NetAmount != expectedNet {
		t.Errorf("expected NetAmount %d, got %d", expectedNet, event.NetAmount)
	}
}

func TestNewSettlement_CommissionExceedsGrossSales(t *testing.T) {
	_, _, err := settlement.NewSettlement(1, 1000, 2000)
	if err == nil {
		t.Fatal("expected error when commission exceeds gross sales, got nil")
	}
	if err != settlement.ErrInvalidSettlementAmount {
		t.Errorf("expected ErrInvalidSettlementAmount, got %v", err)
	}
}

func TestNewSettlement_CommissionEqualsGrossSales(t *testing.T) {
	sett, event, err := settlement.NewSettlement(1, 1000, 1000)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sett.NetAmount() != 0 {
		t.Errorf("expected NetAmount 0, got %d", sett.NetAmount())
	}
	if event.NetAmount != 0 {
		t.Errorf("expected event NetAmount 0, got %d", event.NetAmount)
	}
}

func TestNewSettlement_ZeroValues(t *testing.T) {
	sett, event, err := settlement.NewSettlement(1, 0, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sett.NetAmount() != 0 {
		t.Errorf("expected NetAmount 0, got %d", sett.NetAmount())
	}
	if sett.GrossSales() != 0 {
		t.Errorf("expected GrossSales 0, got %d", sett.GrossSales())
	}
	if event.NetAmount != 0 {
		t.Errorf("expected event NetAmount 0, got %d", event.NetAmount)
	}
}

func TestNewSettlement_NegativeValues(t *testing.T) {
	_, _, err := settlement.NewSettlement(1, -1000, 100)
	if err == nil {
		t.Fatal("expected error for negative gross sales, got nil")
	}
	if err != settlement.ErrInvalidSettlementAmount {
		t.Errorf("expected ErrInvalidSettlementAmount, got %v", err)
	}
}

func TestSettlement_Getters(t *testing.T) {
	sellerID := int64(42)
	grossSales := int64(500000)
	commission := int64(50000)

	sett, _, err := settlement.NewSettlement(sellerID, grossSales, commission)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sett.SellerID() != sellerID {
		t.Errorf("expected SellerID %d, got %d", sellerID, sett.SellerID())
	}
	if sett.GrossSales() != grossSales {
		t.Errorf("expected GrossSales %d, got %d", grossSales, sett.GrossSales())
	}
	if sett.Commission() != commission {
		t.Errorf("expected Commission %d, got %d", commission, sett.Commission())
	}
	if sett.NetAmount() != grossSales-commission {
		t.Errorf("expected NetAmount %d, got %d", grossSales-commission, sett.NetAmount())
	}
	if sett.CreatedAt().IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}
