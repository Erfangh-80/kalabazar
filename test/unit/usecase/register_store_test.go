package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
)

type mockStoreRepo struct {
	saved []*entity.Store
}

func (m *mockStoreRepo) Save(store *entity.Store) error {
	m.saved = append(m.saved, store)
	return nil
}

func (m *mockStoreRepo) FindByID(id string) (*entity.Store, error) {
	for _, s := range m.saved {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, entity.ErrStoreNotFound
}

func (m *mockStoreRepo) FindByUserID(userID string) ([]*entity.Store, error) {
	return nil, nil
}

func (m *mockStoreRepo) Update(store *entity.Store) error {
	return nil
}

func TestRegisterStore_Success(t *testing.T) {
	repo := &mockStoreRepo{}
	uc := usecase.NewRegisterStoreUseCase(repo)

	input := usecase.RegisterStoreInput{
		ID:        "store-1",
		UserID:    "user-1",
		StoreName: "ElectronicsShop",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ID != "store-1" {
		t.Errorf("expected store-1, got %s", output.ID)
	}
	if output.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", output.UserID)
	}
	if output.StoreName != "ElectronicsShop" {
		t.Errorf("expected ElectronicsShop, got %s", output.StoreName)
	}
	if !output.IsCommissionApplicable {
		t.Error("expected commission applicable to be true")
	}

	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved store, got %d", len(repo.saved))
	}
	if repo.saved[0].ID != "store-1" {
		t.Errorf("expected saved store ID store-1, got %s", repo.saved[0].ID)
	}
}

func TestRegisterStore_InvalidInput(t *testing.T) {
	repo := &mockStoreRepo{}
	uc := usecase.NewRegisterStoreUseCase(repo)

	tests := []struct {
		name  string
		input usecase.RegisterStoreInput
	}{
		{"empty id", usecase.RegisterStoreInput{ID: "", UserID: "user-1", StoreName: "Shop"}},
		{"empty user id", usecase.RegisterStoreInput{ID: "store-1", UserID: "", StoreName: "Shop"}},
		{"empty name", usecase.RegisterStoreInput{ID: "store-1", UserID: "user-1", StoreName: ""}},
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

func TestRegisterStore_EventEmitted(t *testing.T) {
	repo := &mockStoreRepo{}
	uc := usecase.NewRegisterStoreUseCase(repo)

	input := usecase.RegisterStoreInput{
		ID:        "store-1",
		UserID:    "user-1",
		StoreName: "ElectronicsShop",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Event == nil {
		t.Fatal("expected a domain event, got nil")
	}
	e, ok := output.Event.(interface{ EventName() string })
	if !ok {
		t.Fatal("expected event with EventName method")
	}
	if e.EventName() != "store.created" {
		t.Errorf("expected store.created event, got %s", e.EventName())
	}
}
