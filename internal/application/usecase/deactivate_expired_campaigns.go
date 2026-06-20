package usecase

import (
	"time"

	"kalabazar-stock-service/internal/domain/entity"
)

type DeactivateExpiredCampaignsInput struct{}

type DeactivatedCampaignInfo struct {
	CampaignID      string
	CampaignEvent   any
	InventoryEvents []any
}

type DeactivateExpiredCampaignsOutput struct {
	DeactivatedCampaigns []DeactivatedCampaignInfo
}

type DeactivateExpiredCampaignsUseCase struct {
	promotionRepo  entity.PromotionRepository
	inventoryRepo  entity.InventoryRepository
}

func NewDeactivateExpiredCampaignsUseCase(promotionRepo entity.PromotionRepository, inventoryRepo entity.InventoryRepository) *DeactivateExpiredCampaignsUseCase {
	return &DeactivateExpiredCampaignsUseCase{
		promotionRepo: promotionRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (uc *DeactivateExpiredCampaignsUseCase) Execute(_ DeactivateExpiredCampaignsInput) (*DeactivateExpiredCampaignsOutput, error) {
	campaigns, err := uc.promotionRepo.FindExpired(time.Now())
	if err != nil {
		return nil, err
	}

	var result DeactivateExpiredCampaignsOutput

	for _, campaign := range campaigns {
		if err := campaign.Deactivate(); err != nil {
			continue
		}

		promoEvents := campaign.Events()
		info := DeactivatedCampaignInfo{
			CampaignID:    campaign.ID,
			CampaignEvent: promoEvents[len(promoEvents)-1],
		}

		linkedInventories, err := uc.inventoryRepo.FindByPromotionID(campaign.ID)
		if err != nil {
			return nil, err
		}

		for _, inv := range linkedInventories {
			if err := inv.ResetPrice(); err != nil {
				return nil, err
			}

			invEvents := inv.Events()
			info.InventoryEvents = append(info.InventoryEvents, invEvents...)

			if campaign.ExpireSaleWithPromotion {
				if err := inv.SetVendorStatus(entity.VendorSaleStatusInactive); err != nil {
					return nil, err
				}
				vendorEvents := inv.Events()
				info.InventoryEvents = append(info.InventoryEvents, vendorEvents...)
			}

			if err := uc.inventoryRepo.Update(inv); err != nil {
				return nil, err
			}
		}

		if err := uc.promotionRepo.Update(campaign); err != nil {
			return nil, err
		}

		result.DeactivatedCampaigns = append(result.DeactivatedCampaigns, info)
	}

	return &result, nil
}
