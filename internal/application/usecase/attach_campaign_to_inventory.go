package usecase

import (
	"kalabazar-stock-service/internal/domain/entity"
)

// AttachCampaignToInventoryInput contains the data needed to attach a campaign to inventory.
type AttachCampaignToInventoryInput struct {
	InventoryID string
	PromotionID string
}

// AttachCampaignToInventoryOutput contains the result of attaching a campaign to inventory.
type AttachCampaignToInventoryOutput struct {
	InventoryID    string
	PromotionID    string
	ProductID      string
	InventoryEvent any
	PromotionEvent any
}

// AttachCampaignToInventoryUseCase orchestrates attaching a campaign to an inventory item.
type AttachCampaignToInventoryUseCase struct {
	inventoryRepo entity.InventoryRepository
	promotionRepo entity.PromotionRepository
}

// NewAttachCampaignToInventoryUseCase creates a new AttachCampaignToInventoryUseCase.
func NewAttachCampaignToInventoryUseCase(
	inventoryRepo entity.InventoryRepository,
	promotionRepo entity.PromotionRepository,
) *AttachCampaignToInventoryUseCase {
	return &AttachCampaignToInventoryUseCase{
		inventoryRepo: inventoryRepo,
		promotionRepo: promotionRepo,
	}
}

// Execute attaches a campaign to an inventory item, replacing any existing campaign.
func (uc *AttachCampaignToInventoryUseCase) Execute(input AttachCampaignToInventoryInput) (*AttachCampaignToInventoryOutput, error) {
	inv, err := uc.inventoryRepo.FindByID(input.InventoryID)
	if err != nil {
		return nil, err
	}

	promo, err := uc.promotionRepo.FindByID(input.PromotionID)
	if err != nil {
		return nil, err
	}

	if inv.PromotionID != nil {
		inv.PromotionID = nil
	}

	if err := inv.LinkPromotion(input.PromotionID); err != nil {
		return nil, err
	}

	if err := promo.LinkToProduct(inv.ProductID); err != nil {
		return nil, err
	}

	if err := uc.inventoryRepo.Update(inv); err != nil {
		return nil, err
	}

	if err := uc.promotionRepo.Update(promo); err != nil {
		return nil, err
	}

	invEvents := inv.Events()
	promoEvents := promo.Events()

	var invEvent any
	if len(invEvents) > 0 {
		invEvent = invEvents[0]
	}

	var promoEvent any
	if len(promoEvents) > 0 {
		promoEvent = promoEvents[0]
	}

	return &AttachCampaignToInventoryOutput{
		InventoryID:    inv.ID,
		PromotionID:    *inv.PromotionID,
		ProductID:      inv.ProductID,
		InventoryEvent: invEvent,
		PromotionEvent: promoEvent,
	}, nil
}
