package payout

import (
	domainpayout "stock-service-version-three/internal/domain/payout"
)

type ExecutePayoutUseCase struct {
	repo domainpayout.PayoutRepository
}

func NewExecutePayoutUseCase(repo domainpayout.PayoutRepository) *ExecutePayoutUseCase {
	return &ExecutePayoutUseCase{repo: repo}
}

func (uc *ExecutePayoutUseCase) Execute(req ExecutePayoutRequest) (*ExecutePayoutResponse, error) {
	p, err := domainpayout.NewPayout(req.SellerID, req.Amount)
	if err != nil {
		return nil, err
	}

	if _, err := p.Execute(); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return &ExecutePayoutResponse{
		PayoutID: p.ID(),
		Amount:   p.Amount(),
		Status:   p.Status(),
	}, nil
}
