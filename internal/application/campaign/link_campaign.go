package campaign

import domaincampaign "stock-service-version-three/internal/domain/campaign"

type LinkCampaignUseCase struct {
	repo domaincampaign.CampaignRepository
}

func NewLinkCampaignUseCase(repo domaincampaign.CampaignRepository) *LinkCampaignUseCase {
	return &LinkCampaignUseCase{repo: repo}
}

func (uc *LinkCampaignUseCase) Execute(req LinkCampaignRequest) (*LinkCampaignResponse, error) {
	c, err := uc.repo.FindByID(req.CampaignID)
	if err != nil {
		return nil, err
	}

	_ = c.LinkToInventory(req.InventoryID)

	if err := uc.repo.Update(c); err != nil {
		return nil, err
	}

	return &LinkCampaignResponse{
		CampaignID:  req.CampaignID,
		InventoryID: req.InventoryID,
	}, nil
}
