package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrCommissionInventoryNotFound   = errors.New("inventory not found")
	ErrCommissionRuleAlreadyDefined  = errors.New("commission rule already defined for this inventory")
)

type DefineCommissionRuleInput struct {
	CommissionID string
	InventoryID  string
	RatePercent  float64
	MinPrice     float64
	MaxPrice     float64
	MinQty       int
}

type DefineCommissionRuleOutput struct {
	CommissionID string
	Event        any
}

type DefineCommissionRuleUseCase struct {
	inventoryRepo  entity.InventoryRepository
	commissionRepo entity.CommissionRepository
}

func NewDefineCommissionRuleUseCase(inventoryRepo entity.InventoryRepository, commissionRepo entity.CommissionRepository) *DefineCommissionRuleUseCase {
	return &DefineCommissionRuleUseCase{
		inventoryRepo:  inventoryRepo,
		commissionRepo: commissionRepo,
	}
}

func (uc *DefineCommissionRuleUseCase) Execute(input DefineCommissionRuleInput) (*DefineCommissionRuleOutput, error) {
	if _, err := uc.inventoryRepo.FindByID(input.InventoryID); err != nil {
		return nil, ErrCommissionInventoryNotFound
	}

	existing, _ := uc.commissionRepo.FindByInventoryID(input.InventoryID)
	if existing != nil {
		return nil, ErrCommissionRuleAlreadyDefined
	}

	commission, err := entity.NewCommission(input.CommissionID, input.InventoryID, "retail", input.RatePercent, input.MinPrice, input.MaxPrice, input.MinQty)
	if err != nil {
		return nil, err
	}

	if err := uc.commissionRepo.Save(commission); err != nil {
		return nil, err
	}

	events := commission.Events()

	return &DefineCommissionRuleOutput{
		CommissionID: commission.ID,
		Event:        events[0],
	}, nil
}
