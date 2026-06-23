package seller

import (
	domainseller "stock-service-version-three/internal/domain/seller"
)

type UpdateRankUseCase struct {
	repo domainseller.SellerRepository
}

func NewUpdateRankUseCase(repo domainseller.SellerRepository) *UpdateRankUseCase {
	return &UpdateRankUseCase{repo: repo}
}

func (uc *UpdateRankUseCase) Execute(req UpdateRankRequest) (*UpdateRankResponse, error) {
	s, err := uc.repo.FindByID(req.SellerID)
	if err != nil {
		return nil, err
	}

	s.UpdateRank(req.Score, req.Rank)

	if err := uc.repo.Update(s); err != nil {
		return nil, err
	}

	return &UpdateRankResponse{
		SellerID: s.ID,
		Score:    s.Score,
		Rank:     s.Rank,
	}, nil
}
