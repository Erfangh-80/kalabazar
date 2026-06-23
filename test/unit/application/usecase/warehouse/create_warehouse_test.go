package warehouse_test

import (
	"errors"
	"testing"

	appWarehouse "stock-service-version-three/internal/application/warehouse"
	domain "stock-service-version-three/internal/domain/warehouse"
)

type mockWarehouseRepository struct {
	warehouses map[int64]*domain.Warehouse
	nextID     int64
	saveErr    error
}

func (m *mockWarehouseRepository) Save(w *domain.Warehouse) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.nextID++
	w.ID = m.nextID
	m.warehouses[w.ID] = w
	return nil
}

func (m *mockWarehouseRepository) FindByID(id int64) (*domain.Warehouse, error) {
	w, ok := m.warehouses[id]
	if !ok {
		return nil, domain.ErrWarehouseNotFound
	}
	return w, nil
}

func (m *mockWarehouseRepository) Update(w *domain.Warehouse) error {
	m.warehouses[w.ID] = w
	return nil
}

func TestCreateWarehouse_Success(t *testing.T) {
	repo := &mockWarehouseRepository{warehouses: make(map[int64]*domain.Warehouse)}
	uc := appWarehouse.NewCreateWarehouseUseCase(repo)

	resp, err := uc.Execute(appWarehouse.CreateWarehouseRequest{Name: "Main Warehouse", Capacity: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.WarehouseID != 1 {
		t.Errorf("expected WarehouseID 1, got %d", resp.WarehouseID)
	}
	if resp.Name != "Main Warehouse" {
		t.Errorf("expected Name 'Main Warehouse', got %s", resp.Name)
	}
}

func TestCreateWarehouse_InvalidCapacity(t *testing.T) {
	repo := &mockWarehouseRepository{warehouses: make(map[int64]*domain.Warehouse)}
	uc := appWarehouse.NewCreateWarehouseUseCase(repo)

	_, err := uc.Execute(appWarehouse.CreateWarehouseRequest{Name: "Bad", Capacity: 0})
	if err != domain.ErrInvalidCapacity {
		t.Errorf("expected ErrInvalidCapacity, got %v", err)
	}
}

func TestCreateWarehouse_RepositoryError(t *testing.T) {
	expectedErr := errors.New("db error")
	repo := &mockWarehouseRepository{
		warehouses: make(map[int64]*domain.Warehouse),
		saveErr:    expectedErr,
	}
	uc := appWarehouse.NewCreateWarehouseUseCase(repo)

	_, err := uc.Execute(appWarehouse.CreateWarehouseRequest{Name: "Main", Capacity: 100})
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}
