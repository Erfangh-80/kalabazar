package usecase

import (
	"time"

	"kalabazar-stock-service/internal/domain/entity"
)

type ActivateScheduledCampaignsInput struct{}

type ActivatedCampaignInfo struct {
	CampaignID string
	Event      any
}

type ActivateScheduledCampaignsOutput struct {
	ActivatedCampaigns []ActivatedCampaignInfo
}

type ActivateScheduledCampaignsUseCase struct {
	promotionRepo entity.PromotionRepository
}

func NewActivateScheduledCampaignsUseCase(promotionRepo entity.PromotionRepository) *ActivateScheduledCampaignsUseCase {
	return &ActivateScheduledCampaignsUseCase{
		promotionRepo: promotionRepo,
	}
}

func (uc *ActivateScheduledCampaignsUseCase) Execute(_ ActivateScheduledCampaignsInput) (*ActivateScheduledCampaignsOutput, error) {
	campaigns, err := uc.promotionRepo.FindSchedulable(time.Now())
	if err != nil {
		return nil, err
	}

	var result ActivateScheduledCampaignsOutput

	for _, campaign := range campaigns {
		if err := campaign.Activate(); err != nil {
			continue
		}
		if err := uc.promotionRepo.Update(campaign); err != nil {
			return nil, err
		}
		events := campaign.Events()
		result.ActivatedCampaigns = append(result.ActivatedCampaigns, ActivatedCampaignInfo{
			CampaignID: campaign.ID,
			Event:      events[len(events)-1],
		})
	}

	return &result, nil
}
