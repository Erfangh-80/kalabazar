package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrInvalidReviewCampaignDecision  = errors.New("review decision must be 'approve' or 'reject'")
	ErrInventoryHasNoLinkedPromotion  = errors.New("inventory item has no linked promotion")
)

// ReviewCampaignInput contains the data needed to review a campaign for an inventory item.
type ReviewCampaignInput struct {
	InventoryID string
	Decision    string
}

// ReviewCampaignOutput contains the result of reviewing a campaign.
type ReviewCampaignOutput struct {
	InventoryID    string
	PromotionID    string
	ApprovalStatus string
	FinalPrice     float64
	InventoryEvents []any
	PromotionEvent  any
}

// ReviewCampaignUseCase orchestrates the review of a campaign for an inventory item.
type ReviewCampaignUseCase struct {
	inventoryRepo entity.InventoryRepository
	promotionRepo entity.PromotionRepository
}

// NewReviewCampaignUseCase creates a new ReviewCampaignUseCase.
func NewReviewCampaignUseCase(
	inventoryRepo entity.InventoryRepository,
	promotionRepo entity.PromotionRepository,
) *ReviewCampaignUseCase {
	return &ReviewCampaignUseCase{
		inventoryRepo: inventoryRepo,
		promotionRepo: promotionRepo,
	}
}

// Execute reviews a campaign for an inventory item with the given decision.
func (uc *ReviewCampaignUseCase) Execute(input ReviewCampaignInput) (*ReviewCampaignOutput, error) {
	inv, err := uc.inventoryRepo.FindByID(input.InventoryID)
	if err != nil {
		return nil, err
	}

	if inv.PromotionID == nil {
		return nil, ErrInventoryHasNoLinkedPromotion
	}

	promo, err := uc.promotionRepo.FindByID(*inv.PromotionID)
	if err != nil {
		return nil, err
	}

	switch input.Decision {
	case "approve":
		if err := promo.Approve(); err != nil {
			return nil, err
		}
		inv.UpdatePromotionStatus(entity.CampaignApprovalApproved)

		newFinalPrice, calcErr := entity.CalculateFinalPrice(inv.BasePrice, promo.DiscountPercent)
		if calcErr != nil {
			return nil, calcErr
		}
		if err := inv.UpdatePrice(inv.BasePrice, newFinalPrice); err != nil {
			return nil, err
		}

	case "reject":
		if err := promo.Reject(); err != nil {
			return nil, err
		}
		inv.UpdatePromotionStatus(entity.CampaignApprovalRejected)

	default:
		return nil, ErrInvalidReviewCampaignDecision
	}

	if err := uc.inventoryRepo.Update(inv); err != nil {
		return nil, err
	}
	if err := uc.promotionRepo.Update(promo); err != nil {
		return nil, err
	}

	promoEvents := promo.Events()
	var promoEvent any
	if len(promoEvents) > 0 {
		promoEvent = promoEvents[0]
	}

	return &ReviewCampaignOutput{
		InventoryID:     inv.ID,
		PromotionID:     *inv.PromotionID,
		ApprovalStatus:  string(inv.CampaignApprovalStatus),
		FinalPrice:      inv.FinalPrice,
		InventoryEvents: inv.Events(),
		PromotionEvent:  promoEvent,
	}, nil
}
