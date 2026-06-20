package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockDefineCommissionInventoryRepo struct {
	items map[string]*entity.Inventory
}

func (m *mockDefineCommissionInventoryRepo) Save(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

func (m *mockDefineCommissionInventoryRepo) FindByID(id string) (*entity.Inventory, error) {
	inv, ok := m.items[id]
	if !ok {
		return nil, entity.ErrInventoryNotFound
	}
	return inv, nil
}

func (m *mockDefineCommissionInventoryRepo) FindByStoreID(storeID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockDefineCommissionInventoryRepo) FindByWarehouseID(warehouseID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockDefineCommissionInventoryRepo) FindByProductID(productID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockDefineCommissionInventoryRepo) FindByPromotionID(promotionID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockDefineCommissionInventoryRepo) Update(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

type mockDefineCommissionRepo struct {
	items map[string]*entity.Commission
}

func (m *mockDefineCommissionRepo) Save(c *entity.Commission) error {
	m.items[c.ID] = c
	return nil
}

func (m *mockDefineCommissionRepo) FindByID(id string) (*entity.Commission, error) {
	c, ok := m.items[id]
	if !ok {
		return nil, entity.ErrCommissionNotFound
	}
	return c, nil
}

func (m *mockDefineCommissionRepo) FindByInventoryID(inventoryID string) (*entity.Commission, error) {
	for _, c := range m.items {
		if c.InventoryID == inventoryID {
			return c, nil
		}
	}
	return nil, entity.ErrCommissionNotFound
}

func newCommissionInventory(id, productID string) *entity.Inventory {
	inv, _ := entity.NewInventory(id, "store-1", "wh-1", productID, 500000, 50, "fixed", "new", 1, nil, nil)
	inv.VendorSaleStatus = "active"
	inv.SystemSaleStatus = "active"
	inv.Events()
	return inv
}

func TestDefineCommissionRule_Success(t *testing.T) {
	inv := newCommissionInventory("inv-1", "prod-x200")

	invRepo := &mockDefineCommissionInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	commRepo := &mockDefineCommissionRepo{items: map[string]*entity.Commission{}}

	uc := usecase.NewDefineCommissionRuleUseCase(invRepo, commRepo)

	output, err := uc.Execute(usecase.DefineCommissionRuleInput{
		CommissionID: "comm-1",
		InventoryID:  "inv-1",
		RatePercent:  10,
		MinPrice:     100000,
		MaxPrice:     1000000,
		MinQty:       1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.CommissionID != "comm-1" {
		t.Errorf("expected comm-1, got %s", output.CommissionID)
	}

	// verify saved
	saved, err := commRepo.FindByID("comm-1")
	if err != nil {
		t.Fatalf("expected saved commission, got %v", err)
	}
	if saved.InventoryID != "inv-1" {
		t.Errorf("expected inv-1, got %s", saved.InventoryID)
	}
	if saved.RatePercent != 10 {
		t.Errorf("expected 10, got %f", saved.RatePercent)
	}

	// verify event
	if output.Event == nil {
		t.Fatal("expected event, got nil")
	}
	e, ok := output.Event.(event.CommissionRuleCreated)
	if !ok {
		t.Fatalf("expected CommissionRuleCreated, got %T", output.Event)
	}
	if e.CommissionID != "comm-1" {
		t.Errorf("expected comm-1, got %s", e.CommissionID)
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected inv-1, got %s", e.InventoryID)
	}
	if e.EventName() != "commission.rule_defined" {
		t.Errorf("expected commission.rule_defined, got %s", e.EventName())
	}
}

func TestDefineCommissionRule_InventoryNotFound(t *testing.T) {
	invRepo := &mockDefineCommissionInventoryRepo{items: map[string]*entity.Inventory{}}
	commRepo := &mockDefineCommissionRepo{items: map[string]*entity.Commission{}}

	uc := usecase.NewDefineCommissionRuleUseCase(invRepo, commRepo)

	_, err := uc.Execute(usecase.DefineCommissionRuleInput{
		CommissionID: "comm-1",
		InventoryID:  "inv-not-found",
		RatePercent:  10,
		MinPrice:     100000,
		MaxPrice:     1000000,
		MinQty:       1,
	})
	if err != usecase.ErrCommissionInventoryNotFound {
		t.Errorf("expected ErrCommissionInventoryNotFound, got %v", err)
	}
}

func TestDefineCommissionRule_DuplicateRule(t *testing.T) {
	inv := newCommissionInventory("inv-1", "prod-x200")
	existingRule, _ := entity.NewCommission("comm-existing", "inv-1", "retail", 5, 0, 1000000, 1)
	existingRule.Events()

	invRepo := &mockDefineCommissionInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	commRepo := &mockDefineCommissionRepo{items: map[string]*entity.Commission{"comm-existing": existingRule}}

	uc := usecase.NewDefineCommissionRuleUseCase(invRepo, commRepo)

	_, err := uc.Execute(usecase.DefineCommissionRuleInput{
		CommissionID: "comm-1",
		InventoryID:  "inv-1",
		RatePercent:  10,
		MinPrice:     100000,
		MaxPrice:     1000000,
		MinQty:       1,
	})
	if err != usecase.ErrCommissionRuleAlreadyDefined {
		t.Errorf("expected ErrCommissionRuleAlreadyDefined, got %v", err)
	}
}

func TestDefineCommissionRule_InvalidInput(t *testing.T) {
	inv := newCommissionInventory("inv-1", "prod-x200")
	invRepo := &mockDefineCommissionInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	commRepo := &mockDefineCommissionRepo{items: map[string]*entity.Commission{}}

	uc := usecase.NewDefineCommissionRuleUseCase(invRepo, commRepo)

	tests := []struct {
		name  string
		input usecase.DefineCommissionRuleInput
	}{
		{"zero rate", usecase.DefineCommissionRuleInput{CommissionID: "comm-1", InventoryID: "inv-1", RatePercent: 0, MinPrice: 100000, MaxPrice: 1000000, MinQty: 1}},
		{"negative rate", usecase.DefineCommissionRuleInput{CommissionID: "comm-1", InventoryID: "inv-1", RatePercent: -5, MinPrice: 100000, MaxPrice: 1000000, MinQty: 1}},
		{"rate over 100", usecase.DefineCommissionRuleInput{CommissionID: "comm-1", InventoryID: "inv-1", RatePercent: 150, MinPrice: 100000, MaxPrice: 1000000, MinQty: 1}},
		{"min price negative", usecase.DefineCommissionRuleInput{CommissionID: "comm-1", InventoryID: "inv-1", RatePercent: 10, MinPrice: -100, MaxPrice: 1000000, MinQty: 1}},
		{"max price negative", usecase.DefineCommissionRuleInput{CommissionID: "comm-1", InventoryID: "inv-1", RatePercent: 10, MinPrice: 0, MaxPrice: -100, MinQty: 1}},
		{"min above max", usecase.DefineCommissionRuleInput{CommissionID: "comm-1", InventoryID: "inv-1", RatePercent: 10, MinPrice: 500, MaxPrice: 100, MinQty: 1}},
		{"negative min qty", usecase.DefineCommissionRuleInput{CommissionID: "comm-1", InventoryID: "inv-1", RatePercent: 10, MinPrice: 0, MaxPrice: 100, MinQty: -1}},
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
