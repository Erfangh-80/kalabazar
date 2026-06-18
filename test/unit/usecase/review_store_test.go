package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
)

type mockStoreRepoForReview struct {
	stores map[string]*entity.Store
}

func (m *mockStoreRepoForReview) Save(store *entity.Store) error {
	m.stores[store.ID] = store
	return nil
}

func (m *mockStoreRepoForReview) FindByID(id string) (*entity.Store, error) {
	s, ok := m.stores[id]
	if !ok {
		return nil, entity.ErrStoreNotFound
	}
	return s, nil
}

func (m *mockStoreRepoForReview) FindByUserID(userID string) ([]*entity.Store, error) {
	return nil, nil
}

func (m *mockStoreRepoForReview) Update(store *entity.Store) error {
	m.stores[store.ID] = store
	return nil
}

func newPendingStore() *entity.Store {
	s, _ := entity.NewStore("store-1", "user-1", "Electronics Shop", nil, nil, nil)
	return s
}

func TestReviewStore_Approve(t *testing.T) {
	store := newPendingStore()
	store.Events()

	repo := &mockStoreRepoForReview{stores: map[string]*entity.Store{"store-1": store}}
	uc := usecase.NewReviewStoreUseCase(repo)

	input := usecase.ReviewStoreInput{StoreID: "store-1", Decision: "approve"}
	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Status != string(entity.StoreStatusActive) {
		t.Errorf("expected active, got %s", output.Status)
	}
	if output.Event == nil {
		t.Fatal("expected event, got nil")
	}
	e, ok := output.Event.(interface{ EventName() string })
	if !ok {
		t.Fatal("expected event with EventName method")
	}
	if e.EventName() != "store.activated" {
		t.Errorf("expected store.activated, got %s", e.EventName())
	}
}

func TestReviewStore_Reject(t *testing.T) {
	store := newPendingStore()
	store.Events()

	repo := &mockStoreRepoForReview{stores: map[string]*entity.Store{"store-1": store}}
	uc := usecase.NewReviewStoreUseCase(repo)

	input := usecase.ReviewStoreInput{StoreID: "store-1", Decision: "reject"}
	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Status != string(entity.StoreStatusRejected) {
		t.Errorf("expected rejected, got %s", output.Status)
	}
	if output.Event == nil {
		t.Fatal("expected event, got nil")
	}
	e, ok := output.Event.(interface{ EventName() string })
	if !ok {
		t.Fatal("expected event with EventName method")
	}
	if e.EventName() != "store.rejected" {
		t.Errorf("expected store.rejected, got %s", e.EventName())
	}
}

func TestReviewStore_StoreNotFound(t *testing.T) {
	repo := &mockStoreRepoForReview{stores: map[string]*entity.Store{}}
	uc := usecase.NewReviewStoreUseCase(repo)

	input := usecase.ReviewStoreInput{StoreID: "nonexistent", Decision: "approve"}
	_, err := uc.Execute(input)
	if err != entity.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}

func TestReviewStore_InvalidDecision(t *testing.T) {
	store := newPendingStore()
	store.Events()

	repo := &mockStoreRepoForReview{stores: map[string]*entity.Store{"store-1": store}}
	uc := usecase.NewReviewStoreUseCase(repo)

	input := usecase.ReviewStoreInput{StoreID: "store-1", Decision: "invalid"}
	_, err := uc.Execute(input)
	if err != usecase.ErrInvalidReviewDecision {
		t.Errorf("expected ErrInvalidReviewDecision, got %v", err)
	}
}
