package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
)

type mockStoreCategoryRepo struct {
	records []*entity.StoreCategory
}

func (m *mockStoreCategoryRepo) Save(sc *entity.StoreCategory) error {
	m.records = append(m.records, sc)
	return nil
}

func (m *mockStoreCategoryRepo) FindByStoreIDAndCategoryID(storeID, categoryID string) (*entity.StoreCategory, error) {
	for _, sc := range m.records {
		if sc.StoreID == storeID && sc.CategoryID == categoryID {
			return sc, nil
		}
	}
	return nil, nil
}

func (m *mockStoreCategoryRepo) Update(sc *entity.StoreCategory) error {
	return nil
}

func TestRequestCategoryAccess_Success(t *testing.T) {
	repo := &mockStoreCategoryRepo{}
	uc := usecase.NewRequestCategoryAccessUseCase(repo)

	input := usecase.RequestCategoryAccessInput{
		StoreID:    "store-1",
		CategoryID: "cat-7",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", output.StoreID)
	}
	if output.CategoryID != "cat-7" {
		t.Errorf("expected cat-7, got %s", output.CategoryID)
	}
	if output.Status != string(entity.StoreCategoryStatusPending) {
		t.Errorf("expected pending, got %s", output.Status)
	}
	if output.Event == nil {
		t.Fatal("expected a domain event, got nil")
	}
	e, ok := output.Event.(interface{ EventName() string })
	if !ok {
		t.Fatal("expected event with EventName method")
	}
	if e.EventName() != "store.category_allowed" {
		t.Errorf("expected store.category_allowed event, got %s", e.EventName())
	}

	if len(repo.records) != 1 {
		t.Fatalf("expected 1 saved record, got %d", len(repo.records))
	}
	if repo.records[0].StoreID != "store-1" {
		t.Errorf("expected saved store ID store-1, got %s", repo.records[0].StoreID)
	}
}

func TestRequestCategoryAccess_InvalidInput(t *testing.T) {
	repo := &mockStoreCategoryRepo{}
	uc := usecase.NewRequestCategoryAccessUseCase(repo)

	tests := []struct {
		name  string
		input usecase.RequestCategoryAccessInput
	}{
		{"empty store id", usecase.RequestCategoryAccessInput{StoreID: "", CategoryID: "cat-7"}},
		{"empty category id", usecase.RequestCategoryAccessInput{StoreID: "store-1", CategoryID: ""}},
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

func TestRequestCategoryAccess_AlreadyRequested(t *testing.T) {
	repo := &mockStoreCategoryRepo{}
	uc := usecase.NewRequestCategoryAccessUseCase(repo)

	input := usecase.RequestCategoryAccessInput{
		StoreID:    "store-1",
		CategoryID: "cat-7",
	}

	_, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error on first request, got %v", err)
	}

	_, err = uc.Execute(input)
	if err != usecase.ErrCategoryAccessAlreadyRequested {
		t.Errorf("expected ErrCategoryAccessAlreadyRequested, got %v", err)
	}
}
