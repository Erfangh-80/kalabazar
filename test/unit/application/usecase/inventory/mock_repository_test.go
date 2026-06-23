package inventory_test

import (
	"sync"

	domain "stock-service-version-three/internal/domain/inventory"
)

type mockInventoryRepository struct {
	mu     sync.Mutex
	items  map[int64]*domain.Inventory
	nextID int64
}

func newMockInventoryRepository() *mockInventoryRepository {
	return &mockInventoryRepository{
		items:  make(map[int64]*domain.Inventory),
		nextID: 1,
	}
}

func (r *mockInventoryRepository) Save(inv *domain.Inventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv.ID = r.nextID
	r.nextID++
	r.items[inv.ID] = inv
	return nil
}

func (r *mockInventoryRepository) FindByID(id int64) (*domain.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.items[id]
	if !ok {
		return nil, domain.ErrInventoryNotFound
	}
	return inv, nil
}

func (r *mockInventoryRepository) FindByProductID(productID int64) ([]*domain.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Inventory
	for _, inv := range r.items {
		if inv.ProductID == productID {
			result = append(result, inv)
		}
	}
	return result, nil
}

func (r *mockInventoryRepository) FindByWarehouseID(warehouseID int64) ([]*domain.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Inventory
	for _, inv := range r.items {
		if inv.WarehouseID == warehouseID {
			result = append(result, inv)
		}
	}
	return result, nil
}

func (r *mockInventoryRepository) Update(inv *domain.Inventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[inv.ID]; !ok {
		return domain.ErrInventoryNotFound
	}
	r.items[inv.ID] = inv
	return nil
}
