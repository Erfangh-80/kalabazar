package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrCommissionRuleNotFound      = errors.New("commission rule not found for the given inventory")
	ErrCommissionConditionsNotMet  = errors.New("commission conditions not met")
)

type CalculateCommissionInput struct {
	InventoryID string
	SaleAmount  float64
	Quantity    int
}

type CalculateCommissionOutput struct {
	CommissionID     string
	CommissionAmount float64
	Event            any
}

type CalculateCommissionUseCase struct {
	commissionRepo entity.CommissionRepository
}

func NewCalculateCommissionUseCase(commissionRepo entity.CommissionRepository) *CalculateCommissionUseCase {
	return &CalculateCommissionUseCase{
		commissionRepo: commissionRepo,
	}
}

func (uc *CalculateCommissionUseCase) Execute(input CalculateCommissionInput) (*CalculateCommissionOutput, error) {
	rule, err := uc.commissionRepo.FindByInventoryID(input.InventoryID)
	if err != nil {
		return nil, ErrCommissionRuleNotFound
	}

	amount, err := rule.Calculate(input.SaleAmount, input.Quantity)
	if err != nil {
		return nil, ErrCommissionConditionsNotMet
	}

	events := rule.Events()

	return &CalculateCommissionOutput{
		CommissionID:     rule.ID,
		CommissionAmount: amount,
		Event:            events[0],
	}, nil
}
