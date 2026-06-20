package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockReferencePriceRepo struct {
	saved []*entity.ReferencePrice
}

func (m *mockReferencePriceRepo) Save(rp *entity.ReferencePrice) error {
	m.saved = append(m.saved, rp)
	return nil
}

func (m *mockReferencePriceRepo) FindByID(id string) (*entity.ReferencePrice, error) {
	return nil, nil
}

func (m *mockReferencePriceRepo) FindByProductID(productID string) ([]*entity.ReferencePrice, error) {
	return nil, nil
}

func TestRecordReferencePrice_Success(t *testing.T) {
	repo := &mockReferencePriceRepo{}
	uc := usecase.NewRecordReferencePriceUseCase(repo)

	input := usecase.RecordReferencePriceInput{
		ID:        "rp-1",
		ProductID: "prod-1",
		Price:     550000,
		Source:    "Digikala",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ID != "rp-1" {
		t.Errorf("expected rp-1, got %s", output.ID)
	}
	if output.ProductID != "prod-1" {
		t.Errorf("expected prod-1, got %s", output.ProductID)
	}
	if output.Price != 550000 {
		t.Errorf("expected 550000, got %f", output.Price)
	}
	if output.Source != "Digikala" {
		t.Errorf("expected Digikala, got %s", output.Source)
	}

	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved reference price, got %d", len(repo.saved))
	}
	if repo.saved[0].ID != "rp-1" {
		t.Errorf("expected saved ID rp-1, got %s", repo.saved[0].ID)
	}
}

func TestRecordReferencePrice_EventEmitted(t *testing.T) {
	repo := &mockReferencePriceRepo{}
	uc := usecase.NewRecordReferencePriceUseCase(repo)

	input := usecase.RecordReferencePriceInput{
		ID:        "rp-1",
		ProductID: "prod-1",
		Price:     550000,
		Source:    "Digikala",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Event == nil {
		t.Fatal("expected a domain event, got nil")
	}
	e, ok := output.Event.(event.ReferencePriceCreated)
	if !ok {
		t.Fatalf("expected ReferencePriceCreated, got %T", output.Event)
	}
	if e.ReferencePriceID != "rp-1" {
		t.Errorf("expected rp-1, got %s", e.ReferencePriceID)
	}
	if e.EventName() != "pricing.reference_price_recorded" {
		t.Errorf("expected pricing.reference_price_recorded, got %s", e.EventName())
	}
}

func TestRecordReferencePrice_InvalidInput(t *testing.T) {
	repo := &mockReferencePriceRepo{}
	uc := usecase.NewRecordReferencePriceUseCase(repo)

	tests := []struct {
		name  string
		input usecase.RecordReferencePriceInput
	}{
		{"empty id", usecase.RecordReferencePriceInput{ID: "", ProductID: "prod-1", Price: 100, Source: "Source"}},
		{"empty product id", usecase.RecordReferencePriceInput{ID: "rp-1", ProductID: "", Price: 100, Source: "Source"}},
		{"zero price", usecase.RecordReferencePriceInput{ID: "rp-1", ProductID: "prod-1", Price: 0, Source: "Source"}},
		{"negative price", usecase.RecordReferencePriceInput{ID: "rp-1", ProductID: "prod-1", Price: -100, Source: "Source"}},
		{"empty source", usecase.RecordReferencePriceInput{ID: "rp-1", ProductID: "prod-1", Price: 100, Source: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(tt.input)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
