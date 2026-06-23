package warehouse_test

import (
	"errors"
	"testing"

	appWarehouse "stock-service-version-three/internal/application/warehouse"
	domain "stock-service-version-three/internal/domain/warehouse"
)

type mockLinkRepository struct {
	links   []*domain.StoreWarehouseLink
	saveErr error
}

func (m *mockLinkRepository) SaveLink(link *domain.StoreWarehouseLink) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	link.ID = int64(len(m.links) + 1)
	m.links = append(m.links, link)
	return nil
}

func TestLinkWarehouseToStore_Success(t *testing.T) {
	repo := &mockLinkRepository{}
	uc := appWarehouse.NewLinkWarehouseToStoreUseCase(repo)

	resp, err := uc.Execute(appWarehouse.LinkWarehouseRequest{
		StoreID:     10,
		WarehouseID: 20,
		Type:        "primary",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StoreID != 10 {
		t.Errorf("expected StoreID 10, got %d", resp.StoreID)
	}
	if resp.WarehouseID != 20 {
		t.Errorf("expected WarehouseID 20, got %d", resp.WarehouseID)
	}
}

func TestLinkWarehouseToStore_RepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := &mockLinkRepository{saveErr: expectedErr}
	uc := appWarehouse.NewLinkWarehouseToStoreUseCase(repo)

	_, err := uc.Execute(appWarehouse.LinkWarehouseRequest{
		StoreID:     10,
		WarehouseID: 20,
		Type:        "primary",
	})
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}
