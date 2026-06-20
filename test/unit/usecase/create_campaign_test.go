package usecase_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockPromotionRepo struct {
	saved []*entity.Promotion
}

func (m *mockPromotionRepo) Save(p *entity.Promotion) error {
	m.saved = append(m.saved, p)
	return nil
}

func (m *mockPromotionRepo) FindByID(id string) (*entity.Promotion, error) {
	return nil, nil
}

func (m *mockPromotionRepo) FindSchedulable(now time.Time) ([]*entity.Promotion, error) {
	return nil, nil
}

func (m *mockPromotionRepo) Update(p *entity.Promotion) error {
	return nil
}

func TestCreateCampaign_Success(t *testing.T) {
	repo := &mockPromotionRepo{}
	uc := usecase.NewCreateCampaignUseCase(repo)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(15 * 24 * time.Hour)

	input := usecase.CreateCampaignInput{
		ID:                     "promo-1",
		Title:                  "Nowruz Auction",
		Description:            "Spring sale campaign",
		StartAt:                start,
		EndAt:                  end,
		RequiresApproval:       true,
		DiscountPercent:        20,
		IsCountdown:            true,
		ExpireSaleWithPromotion: true,
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ID != "promo-1" {
		t.Errorf("expected promo-1, got %s", output.ID)
	}
	if output.Title != "Nowruz Auction" {
		t.Errorf("expected 'Nowruz Auction', got %s", output.Title)
	}
	if output.Description != "Spring sale campaign" {
		t.Errorf("expected 'Spring sale campaign', got %s", output.Description)
	}
	if output.DiscountPercent != 20 {
		t.Errorf("expected 20, got %f", output.DiscountPercent)
	}
	if output.IsActive {
		t.Error("expected output IsActive to be false (starts inactive)")
	}
	if output.RequiresApproval != true {
		t.Error("expected requires approval")
	}
	if output.ApprovalStatus != string(entity.PromotionApprovalPending) {
		t.Errorf("expected pending approval status, got %s", output.ApprovalStatus)
	}
	if !output.IsCountdown {
		t.Error("expected is_countdown to be true")
	}
	if !output.ExpireSaleWithPromotion {
		t.Error("expected expire_sale_with_promotion to be true")
	}

	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved promotion, got %d", len(repo.saved))
	}
	p := repo.saved[0]
	if p.ID != "promo-1" {
		t.Errorf("expected saved ID promo-1, got %s", p.ID)
	}
	if !p.ExpireSaleWithPromotion {
		t.Error("expected saved ExpireSaleWithPromotion to be true")
	}
}

func TestCreateCampaign_WithoutApproval(t *testing.T) {
	repo := &mockPromotionRepo{}
	uc := usecase.NewCreateCampaignUseCase(repo)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(15 * 24 * time.Hour)

	input := usecase.CreateCampaignInput{
		ID:               "promo-2",
		Title:            "Flash Sale",
		StartAt:          start,
		EndAt:            end,
		RequiresApproval: false,
		DiscountPercent:  10,
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ApprovalStatus != string(entity.PromotionApprovalNone) {
		t.Errorf("expected none approval status, got %s", output.ApprovalStatus)
	}
}

func TestCreateCampaign_EventEmitted(t *testing.T) {
	repo := &mockPromotionRepo{}
	uc := usecase.NewCreateCampaignUseCase(repo)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(15 * 24 * time.Hour)

	input := usecase.CreateCampaignInput{
		ID:               "promo-1",
		Title:            "Nowruz Auction",
		StartAt:          start,
		EndAt:            end,
		RequiresApproval: false,
		DiscountPercent:  0,
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Event == nil {
		t.Fatal("expected a domain event, got nil")
	}
	e, ok := output.Event.(event.PromotionCreated)
	if !ok {
		t.Fatalf("expected PromotionCreated, got %T", output.Event)
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected promo-1, got %s", e.PromotionID)
	}
	if e.EventName() != "promotion.campaign_created" {
		t.Errorf("expected promotion.campaign_created, got %s", e.EventName())
	}
}

func TestCreateCampaign_InvalidInput(t *testing.T) {
	repo := &mockPromotionRepo{}
	uc := usecase.NewCreateCampaignUseCase(repo)

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(15 * 24 * time.Hour)

	tests := []struct {
		name  string
		input usecase.CreateCampaignInput
	}{
		{"empty id", usecase.CreateCampaignInput{ID: "", Title: "Sale", StartAt: start, EndAt: end, DiscountPercent: 0}},
		{"empty title", usecase.CreateCampaignInput{ID: "promo-1", Title: "", StartAt: start, EndAt: end, DiscountPercent: 0}},
		{"start after end", usecase.CreateCampaignInput{ID: "promo-1", Title: "Sale", StartAt: end, EndAt: start, DiscountPercent: 0}},
		{"negative discount", usecase.CreateCampaignInput{ID: "promo-1", Title: "Sale", StartAt: start, EndAt: end, DiscountPercent: -10}},
		{"discount over 100", usecase.CreateCampaignInput{ID: "promo-1", Title: "Sale", StartAt: start, EndAt: end, DiscountPercent: 150}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(tt.input)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
