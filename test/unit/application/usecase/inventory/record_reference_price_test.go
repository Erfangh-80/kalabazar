package inventory_test

import (
	"testing"

	appInventory "stock-service-version-three/internal/application/inventory"
)

func TestRecordReferencePrice_Success(t *testing.T) {
	uc := appInventory.NewRecordReferencePriceUseCase()

	req := appInventory.RecordReferencePriceRequest{
		ProductID: 10,
		Price:     150000,
		Source:    "competitor_a",
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ProductID != 10 {
		t.Errorf("ProductID = %d, want 10", resp.ProductID)
	}
	if resp.Price != 150000 {
		t.Errorf("Price = %d, want 150000", resp.Price)
	}
}

func TestRecordReferencePrice_ZeroPrice(t *testing.T) {
	uc := appInventory.NewRecordReferencePriceUseCase()

	req := appInventory.RecordReferencePriceRequest{
		ProductID: 10,
		Price:     0,
		Source:    "competitor_a",
	}

	_, err := uc.Execute(req)
	if err != appInventory.ErrInvalidReferencePrice {
		t.Errorf("expected ErrInvalidReferencePrice, got %v", err)
	}
}

func TestRecordReferencePrice_NegativePrice(t *testing.T) {
	uc := appInventory.NewRecordReferencePriceUseCase()

	req := appInventory.RecordReferencePriceRequest{
		ProductID: 10,
		Price:     -1000,
		Source:    "competitor_a",
	}

	_, err := uc.Execute(req)
	if err != appInventory.ErrInvalidReferencePrice {
		t.Errorf("expected ErrInvalidReferencePrice, got %v", err)
	}
}

func TestRecordReferencePrice_EmptySource(t *testing.T) {
	uc := appInventory.NewRecordReferencePriceUseCase()

	req := appInventory.RecordReferencePriceRequest{
		ProductID: 10,
		Price:     150000,
		Source:    "",
	}

	_, err := uc.Execute(req)
	if err != appInventory.ErrInvalidReferenceSource {
		t.Errorf("expected ErrInvalidReferenceSource, got %v", err)
	}
}
