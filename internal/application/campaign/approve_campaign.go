package campaign

import domaincampaign "stock-service-version-three/internal/domain/campaign"

type ApproveCampaignUseCase struct {
	repo domaincampaign.CampaignRepository
}

func NewApproveCampaignUseCase(repo domaincampaign.CampaignRepository) *ApproveCampaignUseCase {
	return &ApproveCampaignUseCase{repo: repo}
}

func (uc *ApproveCampaignUseCase) Execute(req ApproveCampaignRequest) (*ApproveCampaignResponse, error) {
	c, err := uc.repo.FindByID(req.CampaignID)
	if err != nil {
		return nil, err
	}

	if req.Decision == "approved" {
		if _, err := c.Approve(); err != nil {
			return nil, err
		}
	} else {
		c.ApprovalStatus = domaincampaign.ApprovalStatusRejected
	}

	if err := uc.repo.Update(c); err != nil {
		return nil, err
	}

	return &ApproveCampaignResponse{
		CampaignID:     c.ID,
		ApprovalStatus: string(c.ApprovalStatus),
	}, nil
}
