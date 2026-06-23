package store_test

import (
	"context"
	"testing"

	appStore "stock-service-version-three/internal/application/store"
	domainstore "stock-service-version-three/internal/domain/store"
)

type mockStoreRepo struct {
	stores map[int64]*domainstore.Store
	nextID int64
}

func newMockStoreRepo() *mockStoreRepo {
	return &mockStoreRepo{stores: make(map[int64]*domainstore.Store)}
}

func (m *mockStoreRepo) Save(ctx context.Context, store *domainstore.Store) error {
	m.nextID++
	store.ID = m.nextID
	m.stores[store.ID] = store
	return nil
}

func (m *mockStoreRepo) FindByID(ctx context.Context, id int64) (*domainstore.Store, error) {
	s, ok := m.stores[id]
	if !ok {
		return nil, domainstore.ErrStoreNotFound
	}
	return s, nil
}

func (m *mockStoreRepo) FindBySellerID(ctx context.Context, sellerID int64) (*domainstore.Store, error) {
	for _, s := range m.stores {
		if s.SellerID == sellerID {
			return s, nil
		}
	}
	return nil, domainstore.ErrStoreNotFound
}

func (m *mockStoreRepo) Update(ctx context.Context, store *domainstore.Store) error {
	_, ok := m.stores[store.ID]
	if !ok {
		return domainstore.ErrStoreNotFound
	}
	m.stores[store.ID] = store
	return nil
}

func TestApproveStore_Success(t *testing.T) {
	storeRepo := newMockStoreRepo()
	uc := appStore.NewApproveStoreUseCase(storeRepo)

	store := domainstore.NewStore(1, "Test Store", "09123456789")
	storeRepo.Save(context.Background(), store)

	req := appStore.ApproveStoreRequest{
		StoreID:  store.ID,
		Decision: "approved",
	}

	resp, err := uc.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StoreID != store.ID {
		t.Errorf("expected store ID %d, got %d", store.ID, resp.StoreID)
	}
	if resp.Status != string(domainstore.StoreStatusACTIVE) {
		t.Errorf("expected status ACTIVE, got '%s'", resp.Status)
	}
	if len(resp.Events) == 0 {
		t.Error("expected at least one event (StoreActivatedEvent)")
	}
}

func TestApproveStore_NotApproved(t *testing.T) {
	storeRepo := newMockStoreRepo()
	uc := appStore.NewApproveStoreUseCase(storeRepo)

	store := domainstore.NewStore(1, "Test Store", "09123456789")
	storeRepo.Save(context.Background(), store)

	req := appStore.ApproveStoreRequest{
		StoreID:  store.ID,
		Decision: "rejected",
	}

	resp, err := uc.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != string(domainstore.StoreStatusPENDING) {
		t.Errorf("expected status PENDING, got '%s'", resp.Status)
	}
}

func TestApproveStore_NotFound(t *testing.T) {
	storeRepo := newMockStoreRepo()
	uc := appStore.NewApproveStoreUseCase(storeRepo)

	req := appStore.ApproveStoreRequest{
		StoreID:  999,
		Decision: "approved",
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainstore.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}

func TestApproveStore_AlreadyActive(t *testing.T) {
	storeRepo := newMockStoreRepo()
	uc := appStore.NewApproveStoreUseCase(storeRepo)

	store := domainstore.NewStore(1, "Test Store", "09123456789")
	storeRepo.Save(context.Background(), store)
	store.Activate()
	storeRepo.Update(context.Background(), store)

	req := appStore.ApproveStoreRequest{
		StoreID:  store.ID,
		Decision: "approved",
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainstore.ErrStoreAlreadyActive {
		t.Errorf("expected ErrStoreAlreadyActive, got %v", err)
	}
}
