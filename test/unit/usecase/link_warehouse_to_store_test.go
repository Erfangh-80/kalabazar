package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockWarehouseRepoForLink struct {
	warehouses map[string]*entity.Warehouse
}

func (m *mockWarehouseRepoForLink) Save(warehouse *entity.Warehouse) error {
	m.warehouses[warehouse.ID] = warehouse
	return nil
}

func (m *mockWarehouseRepoForLink) FindByID(id string) (*entity.Warehouse, error) {
	w, ok := m.warehouses[id]
	if !ok {
		return nil, entity.ErrWarehouseNotFound
	}
	return w, nil
}

func (m *mockWarehouseRepoForLink) FindBySellerID(sellerID string) ([]*entity.Warehouse, error) {
	return nil, nil
}

func (m *mockWarehouseRepoForLink) Update(warehouse *entity.Warehouse) error {
	m.warehouses[warehouse.ID] = warehouse
	return nil
}

func newWarehouseForLink() *entity.Warehouse {
	w, _ := entity.NewWarehouse("wh-1", "seller-1", "Tehran Warehouse", validAddress(), 1000, "public")
	w.Events()
	return w
}

func TestLinkWarehouseToStore_Success(t *testing.T) {
	w := newWarehouseForLink()
	repo := &mockWarehouseRepoForLink{warehouses: map[string]*entity.Warehouse{"wh-1": w}}
	uc := usecase.NewLinkWarehouseToStoreUseCase(repo)

	input := usecase.LinkWarehouseToStoreInput{
		WarehouseID:  "wh-1",
		StoreID:      "store-1",
		RelationType: "primary",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.WarehouseID != "wh-1" {
		t.Errorf("expected wh-1, got %s", output.WarehouseID)
	}
	if output.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", output.StoreID)
	}
	if output.RelationType != "primary" {
		t.Errorf("expected primary, got %s", output.RelationType)
	}

	if w.StoreID != "store-1" {
		t.Errorf("expected warehouse StoreID store-1, got %s", w.StoreID)
	}
}

func TestLinkWarehouseToStore_Secondary(t *testing.T) {
	w := newWarehouseForLink()
	repo := &mockWarehouseRepoForLink{warehouses: map[string]*entity.Warehouse{"wh-1": w}}
	uc := usecase.NewLinkWarehouseToStoreUseCase(repo)

	input := usecase.LinkWarehouseToStoreInput{
		WarehouseID:  "wh-1",
		StoreID:      "store-1",
		RelationType: "secondary",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.RelationType != "secondary" {
		t.Errorf("expected secondary, got %s", output.RelationType)
	}
}

func TestLinkWarehouseToStore_EventEmitted(t *testing.T) {
	w := newWarehouseForLink()
	repo := &mockWarehouseRepoForLink{warehouses: map[string]*entity.Warehouse{"wh-1": w}}
	uc := usecase.NewLinkWarehouseToStoreUseCase(repo)

	input := usecase.LinkWarehouseToStoreInput{
		WarehouseID:  "wh-1",
		StoreID:      "store-1",
		RelationType: "primary",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Event == nil {
		t.Fatal("expected a domain event, got nil")
	}
	e, ok := output.Event.(event.WarehouseLinkedToStore)
	if !ok {
		t.Fatalf("expected WarehouseLinkedToStore, got %T", output.Event)
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected wh-1, got %s", e.WarehouseID)
	}
	if e.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", e.StoreID)
	}
	if e.RelationType != "primary" {
		t.Errorf("expected primary, got %s", e.RelationType)
	}
	if e.EventName() != "warehouse.linked_to_store" {
		t.Errorf("expected warehouse.linked_to_store, got %s", e.EventName())
	}
}

func TestLinkWarehouseToStore_WarehouseNotFound(t *testing.T) {
	repo := &mockWarehouseRepoForLink{warehouses: map[string]*entity.Warehouse{}}
	uc := usecase.NewLinkWarehouseToStoreUseCase(repo)

	input := usecase.LinkWarehouseToStoreInput{
		WarehouseID:  "nonexistent",
		StoreID:      "store-1",
		RelationType: "primary",
	}

	_, err := uc.Execute(input)
	if err != entity.ErrWarehouseNotFound {
		t.Errorf("expected ErrWarehouseNotFound, got %v", err)
	}
}

func TestLinkWarehouseToStore_InvalidRelationType(t *testing.T) {
	w := newWarehouseForLink()
	repo := &mockWarehouseRepoForLink{warehouses: map[string]*entity.Warehouse{"wh-1": w}}
	uc := usecase.NewLinkWarehouseToStoreUseCase(repo)

	tests := []struct {
		name         string
		relationType string
	}{
		{"empty", ""},
		{"invalid value", "tertiary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := usecase.LinkWarehouseToStoreInput{
				WarehouseID:  "wh-1",
				StoreID:      "store-1",
				RelationType: tt.relationType,
			}
			_, err := uc.Execute(input)
			if err != usecase.ErrInvalidRelationType {
				t.Errorf("expected ErrInvalidRelationType, got %v", err)
			}
		})
	}
}

func TestLinkWarehouseToStore_AlreadyLinked(t *testing.T) {
	w := newWarehouseForLink()
	w.LinkToStore("store-1", "primary")
	w.Events()
	repo := &mockWarehouseRepoForLink{warehouses: map[string]*entity.Warehouse{"wh-1": w}}
	uc := usecase.NewLinkWarehouseToStoreUseCase(repo)

	input := usecase.LinkWarehouseToStoreInput{
		WarehouseID:  "wh-1",
		StoreID:      "store-1",
		RelationType: "primary",
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrWarehouseAlreadyLinkedToStore {
		t.Errorf("expected ErrWarehouseAlreadyLinkedToStore, got %v", err)
	}
}
