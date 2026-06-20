package usecase

import (
	"time"

	"kalabazar-stock-service/internal/domain/entity"
)

// CreateCampaignInput contains the data needed to create a new campaign.
type CreateCampaignInput struct {
	ID                     string
	Title                  string
	Description            string
	StartAt                time.Time
	EndAt                  time.Time
	RequiresApproval       bool
	DiscountPercent        float64
	IsCountdown            bool
	ExpireSaleWithPromotion bool
}

// CreateCampaignOutput contains the result of creating a campaign.
type CreateCampaignOutput struct {
	ID                     string
	Title                  string
	Description            string
	StartAt                time.Time
	EndAt                  time.Time
	RequiresApproval       bool
	ApprovalStatus         string
	IsActive               bool
	IsCountdown            bool
	DiscountPercent        float64
	ExpireSaleWithPromotion bool
	Event                  any
	CreatedAt              time.Time
}

// CreateCampaignUseCase orchestrates creating a new discount campaign.
type CreateCampaignUseCase struct {
	repo entity.PromotionRepository
}

// NewCreateCampaignUseCase creates a new CreateCampaignUseCase.
func NewCreateCampaignUseCase(repo entity.PromotionRepository) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{repo: repo}
}

// Execute creates a new campaign with the given input.
func (uc *CreateCampaignUseCase) Execute(input CreateCampaignInput) (*CreateCampaignOutput, error) {
	campaign, err := entity.NewPromotion(
		input.ID, input.Title, input.Description,
		input.StartAt, input.EndAt,
		input.RequiresApproval, input.DiscountPercent, input.IsCountdown,
	)
	if err != nil {
		return nil, err
	}

	campaign.ExpireSaleWithPromotion = input.ExpireSaleWithPromotion

	if err := uc.repo.Save(campaign); err != nil {
		return nil, err
	}

	events := campaign.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &CreateCampaignOutput{
		ID:                     campaign.ID,
		Title:                  campaign.Title,
		Description:            campaign.Description,
		StartAt:                campaign.StartAt,
		EndAt:                  campaign.EndAt,
		RequiresApproval:       campaign.RequiresApproval,
		ApprovalStatus:         string(campaign.ApprovalStatus),
		IsActive:               campaign.IsActive,
		IsCountdown:            campaign.IsCountdown,
		DiscountPercent:        campaign.DiscountPercent,
		ExpireSaleWithPromotion: campaign.ExpireSaleWithPromotion,
		Event:                  domainEvent,
		CreatedAt:              campaign.CreatedAt,
	}, nil
}
