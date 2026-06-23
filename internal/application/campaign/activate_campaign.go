package campaign

import domaincampaign "stock-service-version-three/internal/domain/campaign"

type ActivateCampaignUseCase struct {
	repo domaincampaign.CampaignRepository
}

func NewActivateCampaignUseCase(repo domaincampaign.CampaignRepository) *ActivateCampaignUseCase {
	return &ActivateCampaignUseCase{repo: repo}
}

func (uc *ActivateCampaignUseCase) Execute(req ActivateCampaignRequest) (*ActivateCampaignResponse, error) {
	c, err := uc.repo.FindByID(req.CampaignID)
	if err != nil {
		return nil, err
	}

	if _, err := c.Activate(req.Now); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(c); err != nil {
		return nil, err
	}

	return &ActivateCampaignResponse{
		CampaignID: c.ID,
		Status:     string(c.Status),
	}, nil
}
