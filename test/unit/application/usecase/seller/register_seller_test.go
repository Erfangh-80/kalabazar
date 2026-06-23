package seller_test

import (
	"context"
	"testing"

	appSeller "stock-service-version-three/internal/application/seller"
	domainseller "stock-service-version-three/internal/domain/seller"
	domainstore "stock-service-version-three/internal/domain/store"
)

type mockSellerRepo struct {
	sellers map[int64]*domainseller.Seller
	nextID  int64
}

func newMockSellerRepo() *mockSellerRepo {
	return &mockSellerRepo{sellers: make(map[int64]*domainseller.Seller)}
}

func (m *mockSellerRepo) Save(seller *domainseller.Seller) error {
	m.nextID++
	seller.ID = m.nextID
	m.sellers[seller.ID] = seller
	return nil
}

func (m *mockSellerRepo) FindByID(id int64) (*domainseller.Seller, error) {
	s, ok := m.sellers[id]
	if !ok {
		return nil, domainseller.ErrSellerNotFound
	}
	return s, nil
}

func (m *mockSellerRepo) FindByUserID(userID int64) (*domainseller.Seller, error) {
	for _, s := range m.sellers {
		if s.UserID == userID {
			return s, nil
		}
	}
	return nil, domainseller.ErrSellerNotFound
}

func (m *mockSellerRepo) Update(seller *domainseller.Seller) error {
	_, ok := m.sellers[seller.ID]
	if !ok {
		return domainseller.ErrSellerNotFound
	}
	m.sellers[seller.ID] = seller
	return nil
}

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

func TestRegisterSeller_Success(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	storeRepo := newMockStoreRepo()
	uc := appSeller.NewRegisterSellerUseCase(sellerRepo, storeRepo)

	req := appSeller.RegisterSellerRequest{
		UserID:    1,
		StoreName: "My Store",
		Phone:     "09123456789",
	}

	resp, err := uc.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.SellerID == 0 {
		t.Error("expected non-zero seller ID")
	}
	if resp.StoreID == 0 {
		t.Error("expected non-zero store ID")
	}

	s, err := sellerRepo.FindByID(resp.SellerID)
	if err != nil {
		t.Fatalf("seller not found: %v", err)
	}
	if s.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", s.UserID)
	}
	if s.Name != "My Store" {
		t.Errorf("expected name 'My Store', got '%s'", s.Name)
	}
	if s.Phone != "09123456789" {
		t.Errorf("expected phone '09123456789', got '%s'", s.Phone)
	}
	if s.Status != domainseller.SellerStatusUnverified {
		t.Errorf("expected status UNVERIFIED, got '%s'", s.Status)
	}

	st, err := storeRepo.FindByID(context.Background(), resp.StoreID)
	if err != nil {
		t.Fatalf("store not found: %v", err)
	}
	if st.SellerID != resp.SellerID {
		t.Errorf("expected seller ID %d, got %d", resp.SellerID, st.SellerID)
	}
	if st.Name != "My Store" {
		t.Errorf("expected name 'My Store', got '%s'", st.Name)
	}
	if st.Status != domainstore.StoreStatusPENDING {
		t.Errorf("expected status PENDING, got '%s'", st.Status)
	}

	if len(resp.Events) == 0 {
		t.Error("expected at least one event (StoreCreatedEvent)")
	}
}

func TestRegisterSeller_InvalidStoreName(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	storeRepo := newMockStoreRepo()
	uc := appSeller.NewRegisterSellerUseCase(sellerRepo, storeRepo)

	req := appSeller.RegisterSellerRequest{
		UserID:    1,
		StoreName: "",
		Phone:     "09123456789",
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainseller.ErrInvalidSellerName {
		t.Errorf("expected ErrInvalidSellerName, got %v", err)
	}
}

func TestRegisterSeller_InvalidPhone(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	storeRepo := newMockStoreRepo()
	uc := appSeller.NewRegisterSellerUseCase(sellerRepo, storeRepo)

	req := appSeller.RegisterSellerRequest{
		UserID:    1,
		StoreName: "My Store",
		Phone:     "123",
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainseller.ErrInvalidPhone {
		t.Errorf("expected ErrInvalidPhone, got %v", err)
	}
}
