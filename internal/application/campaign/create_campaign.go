package campaign

import (
	domaincampaign "stock-service-version-three/internal/domain/campaign"
)

type CreateCampaignUseCase struct {
	repo domaincampaign.CampaignRepository
}

func NewCreateCampaignUseCase(repo domaincampaign.CampaignRepository) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{repo: repo}
}

func (uc *CreateCampaignUseCase) Execute(req CreateCampaignRequest) (*CreateCampaignResponse, error) {
	c := domaincampaign.NewCampaign(req.Title, req.DiscountType, req.Value, req.StartAt, req.EndAt)

	if err := uc.repo.Save(c); err != nil {
		return nil, err
	}

	return &CreateCampaignResponse{
		CampaignID:     c.ID,
		Status:         string(c.Status),
		ApprovalStatus: string(c.ApprovalStatus),
	}, nil
}
