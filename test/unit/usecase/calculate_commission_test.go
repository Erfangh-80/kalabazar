package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockCalculateCommissionRepo struct {
	items map[string]*entity.Commission
}

func (m *mockCalculateCommissionRepo) Save(c *entity.Commission) error {
	m.items[c.ID] = c
	return nil
}

func (m *mockCalculateCommissionRepo) FindByID(id string) (*entity.Commission, error) {
	c, ok := m.items[id]
	if !ok {
		return nil, entity.ErrCommissionNotFound
	}
	return c, nil
}

func (m *mockCalculateCommissionRepo) FindByInventoryID(inventoryID string) (*entity.Commission, error) {
	for _, c := range m.items {
		if c.InventoryID == inventoryID {
			return c, nil
		}
	}
	return nil, entity.ErrCommissionNotFound
}

func TestCalculateCommission_Success(t *testing.T) {
	rule, _ := entity.NewCommission("comm-1", "inv-1", "retail", 10, 100000, 1000000, 1)
	rule.Events()

	repo := &mockCalculateCommissionRepo{items: map[string]*entity.Commission{"comm-1": rule}}

	uc := usecase.NewCalculateCommissionUseCase(repo)

	output, err := uc.Execute(usecase.CalculateCommissionInput{
		InventoryID: "inv-1",
		SaleAmount:  400000,
		Quantity:    1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.CommissionID != "comm-1" {
		t.Errorf("expected comm-1, got %s", output.CommissionID)
	}
	if output.CommissionAmount != 40000 {
		t.Errorf("expected 40000, got %f", output.CommissionAmount)
	}

	// verify event
	if output.Event == nil {
		t.Fatal("expected event, got nil")
	}
	e, ok := output.Event.(event.CommissionCalculated)
	if !ok {
		t.Fatalf("expected CommissionCalculated, got %T", output.Event)
	}
	if e.CommissionID != "comm-1" {
		t.Errorf("expected comm-1, got %s", e.CommissionID)
	}
	if e.SaleAmount != 400000 {
		t.Errorf("expected 400000, got %f", e.SaleAmount)
	}
	if e.CommissionAmount != 40000 {
		t.Errorf("expected 40000, got %f", e.CommissionAmount)
	}
	if e.EventName() != "commission.calculated" {
		t.Errorf("expected commission.calculated, got %s", e.EventName())
	}
}

func TestCalculateCommission_RuleNotFound(t *testing.T) {
	repo := &mockCalculateCommissionRepo{items: map[string]*entity.Commission{}}
	uc := usecase.NewCalculateCommissionUseCase(repo)

	_, err := uc.Execute(usecase.CalculateCommissionInput{
		InventoryID: "inv-not-found",
		SaleAmount:  400000,
		Quantity:    1,
	})
	if err != usecase.ErrCommissionRuleNotFound {
		t.Errorf("expected ErrCommissionRuleNotFound, got %v", err)
	}
}

func TestCalculateCommission_ConditionsNotMet(t *testing.T) {
	rule, _ := entity.NewCommission("comm-1", "inv-1", "retail", 10, 100000, 1000000, 1)
	rule.Events()

	repo := &mockCalculateCommissionRepo{items: map[string]*entity.Commission{"comm-1": rule}}
	uc := usecase.NewCalculateCommissionUseCase(repo)

	// price below min
	_, err := uc.Execute(usecase.CalculateCommissionInput{
		InventoryID: "inv-1",
		SaleAmount:  50000,
		Quantity:    1,
	})
	if err != usecase.ErrCommissionConditionsNotMet {
		t.Errorf("expected ErrCommissionConditionsNotMet, got %v", err)
	}

	// price above max
	_, err = uc.Execute(usecase.CalculateCommissionInput{
		InventoryID: "inv-1",
		SaleAmount:  2000000,
		Quantity:    1,
	})
	if err != usecase.ErrCommissionConditionsNotMet {
		t.Errorf("expected ErrCommissionConditionsNotMet, got %v", err)
	}

	// qty below min
	rule2, _ := entity.NewCommission("comm-2", "inv-2", "retail", 10, 0, 1000000, 5)
	rule2.Events()
	repo2 := &mockCalculateCommissionRepo{items: map[string]*entity.Commission{"comm-2": rule2}}
	uc2 := usecase.NewCalculateCommissionUseCase(repo2)

	_, err = uc2.Execute(usecase.CalculateCommissionInput{
		InventoryID: "inv-2",
		SaleAmount:  400000,
		Quantity:    1,
	})
	if err != usecase.ErrCommissionConditionsNotMet {
		t.Errorf("expected ErrCommissionConditionsNotMet, got %v", err)
	}
}
