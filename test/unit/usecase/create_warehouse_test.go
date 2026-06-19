package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockWarehouseRepo struct {
	saved []*entity.Warehouse
}

func (m *mockWarehouseRepo) Save(warehouse *entity.Warehouse) error {
	m.saved = append(m.saved, warehouse)
	return nil
}

func (m *mockWarehouseRepo) FindByID(id string) (*entity.Warehouse, error) {
	return nil, nil
}

func (m *mockWarehouseRepo) FindBySellerID(sellerID string) ([]*entity.Warehouse, error) {
	return nil, nil
}

func (m *mockWarehouseRepo) Update(warehouse *entity.Warehouse) error {
	return nil
}

func validAddress() entity.Address {
	return entity.Address{
		Street:  "Main St",
		City:    "Tehran",
		Country: "Iran",
	}
}

func TestCreateWarehouse_Success(t *testing.T) {
	repo := &mockWarehouseRepo{}
	uc := usecase.NewCreateWarehouseUseCase(repo)

	input := usecase.CreateWarehouseInput{
		ID:            "wh-1",
		SellerID:      "seller-1",
		Name:          "Tehran Central Warehouse",
		Address:       validAddress(),
		TotalCapacity: 1000,
		AccessType:    "public",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ID != "wh-1" {
		t.Errorf("expected wh-1, got %s", output.ID)
	}
	if output.SellerID != "seller-1" {
		t.Errorf("expected seller-1, got %s", output.SellerID)
	}
	if output.Name != "Tehran Central Warehouse" {
		t.Errorf("expected Tehran Central Warehouse, got %s", output.Name)
	}
	if output.AccessType != "public" {
		t.Errorf("expected public, got %s", output.AccessType)
	}
	if output.TotalCapacity != 1000 {
		t.Errorf("expected 1000, got %d", output.TotalCapacity)
	}
	if output.Status != string(entity.WarehouseStatusActive) {
		t.Errorf("expected active, got %s", output.Status)
	}

	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved warehouse, got %d", len(repo.saved))
	}
	if repo.saved[0].ID != "wh-1" {
		t.Errorf("expected saved warehouse ID wh-1, got %s", repo.saved[0].ID)
	}
}

func TestCreateWarehouse_EventEmitted(t *testing.T) {
	repo := &mockWarehouseRepo{}
	uc := usecase.NewCreateWarehouseUseCase(repo)

	input := usecase.CreateWarehouseInput{
		ID:            "wh-1",
		SellerID:      "seller-1",
		Name:          "Tehran Central Warehouse",
		Address:       validAddress(),
		TotalCapacity: 1000,
		AccessType:    "public",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Event == nil {
		t.Fatal("expected a domain event, got nil")
	}
	e, ok := output.Event.(event.WarehouseCreated)
	if !ok {
		t.Fatalf("expected WarehouseCreated, got %T", output.Event)
	}
	if e.WarehouseID != "wh-1" {
		t.Errorf("expected wh-1, got %s", e.WarehouseID)
	}
	if e.EventName() != "warehouse.created" {
		t.Errorf("expected warehouse.created, got %s", e.EventName())
	}
}

func TestCreateWarehouse_InvalidInput(t *testing.T) {
	repo := &mockWarehouseRepo{}
	uc := usecase.NewCreateWarehouseUseCase(repo)

	tests := []struct {
		name  string
		input usecase.CreateWarehouseInput
	}{
		{"empty id", usecase.CreateWarehouseInput{ID: "", SellerID: "seller-1", Name: "Warehouse", Address: validAddress(), TotalCapacity: 100, AccessType: "public"}},
		{"empty seller id", usecase.CreateWarehouseInput{ID: "wh-1", SellerID: "", Name: "Warehouse", Address: validAddress(), TotalCapacity: 100, AccessType: "public"}},
		{"empty name", usecase.CreateWarehouseInput{ID: "wh-1", SellerID: "seller-1", Name: "", Address: validAddress(), TotalCapacity: 100, AccessType: "public"}},
		{"invalid access type", usecase.CreateWarehouseInput{ID: "wh-1", SellerID: "seller-1", Name: "Warehouse", Address: validAddress(), TotalCapacity: 100, AccessType: "invalid"}},
		{"zero capacity", usecase.CreateWarehouseInput{ID: "wh-1", SellerID: "seller-1", Name: "Warehouse", Address: validAddress(), TotalCapacity: 0, AccessType: "public"}},
		{"invalid address", usecase.CreateWarehouseInput{ID: "wh-1", SellerID: "seller-1", Name: "Warehouse", Address: entity.Address{}, TotalCapacity: 100, AccessType: "public"}},
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

func TestCreateWarehouse_PrivateAccess(t *testing.T) {
	repo := &mockWarehouseRepo{}
	uc := usecase.NewCreateWarehouseUseCase(repo)

	input := usecase.CreateWarehouseInput{
		ID:            "wh-2",
		SellerID:      "seller-1",
		Name:          "Internal Storage",
		Address:       validAddress(),
		TotalCapacity: 500,
		AccessType:    "private",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.AccessType != "private" {
		t.Errorf("expected private, got %s", output.AccessType)
	}
}
