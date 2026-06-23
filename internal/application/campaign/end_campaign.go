package campaign

import domaincampaign "stock-service-version-three/internal/domain/campaign"

type EndCampaignUseCase struct {
	repo domaincampaign.CampaignRepository
}

func NewEndCampaignUseCase(repo domaincampaign.CampaignRepository) *EndCampaignUseCase {
	return &EndCampaignUseCase{repo: repo}
}

func (uc *EndCampaignUseCase) Execute(req EndCampaignRequest) (*EndCampaignResponse, error) {
	c, err := uc.repo.FindByID(req.CampaignID)
	if err != nil {
		return nil, err
	}

	if _, err := c.End(req.Now); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(c); err != nil {
		return nil, err
	}

	return &EndCampaignResponse{
		CampaignID: c.ID,
		Status:     string(c.Status),
	}, nil
}
