package commission

import (
	domaincommission "stock-service-version-three/internal/domain/commission"
)

type CalculateCommissionUseCase struct {
	repo domaincommission.CommissionRepository
}

func NewCalculateCommissionUseCase(repo domaincommission.CommissionRepository) *CalculateCommissionUseCase {
	return &CalculateCommissionUseCase{repo: repo}
}

func (uc *CalculateCommissionUseCase) Execute(req CalculateCommissionRequest) (*CalculateCommissionResponse, error) {
	c, _, err := domaincommission.NewCommission(req.SellerID, req.Rate, req.SalesAmount)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(c); err != nil {
		return nil, err
	}

	return &CalculateCommissionResponse{
		CommissionID: c.ID(),
		Amount:       c.Amount(),
		SalesAmount:  c.SalesAmount(),
	}, nil
}
