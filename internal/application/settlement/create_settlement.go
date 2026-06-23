package settlement

import (
	domainsettlement "stock-service-version-three/internal/domain/settlement"
)

type CreateSettlementUseCase struct {
	repo domainsettlement.SettlementRepository
}

func NewCreateSettlementUseCase(repo domainsettlement.SettlementRepository) *CreateSettlementUseCase {
	return &CreateSettlementUseCase{repo: repo}
}

func (uc *CreateSettlementUseCase) Execute(req CreateSettlementRequest) (*CreateSettlementResponse, error) {
	s, _, err := domainsettlement.NewSettlement(req.SellerID, req.GrossSales, req.Commission)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(s); err != nil {
		return nil, err
	}

	return &CreateSettlementResponse{
		SettlementID: s.ID(),
		GrossSales:   s.GrossSales(),
		Commission:   s.Commission(),
		NetAmount:    s.NetAmount(),
	}, nil
}
