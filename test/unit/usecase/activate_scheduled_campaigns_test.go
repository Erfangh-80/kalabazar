package usecase_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockActivateCampaignRepo struct {
	items map[string]*entity.Promotion
}

func (m *mockActivateCampaignRepo) Save(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

func (m *mockActivateCampaignRepo) FindByID(id string) (*entity.Promotion, error) {
	p, ok := m.items[id]
	if !ok {
		return nil, entity.ErrPromotionNotFound
	}
	return p, nil
}

func (m *mockActivateCampaignRepo) FindSchedulable(now time.Time) ([]*entity.Promotion, error) {
	var result []*entity.Promotion
	for _, p := range m.items {
		if p.IsActive {
			continue
		}
		if p.StartAt.After(now) {
			continue
		}
		if p.RequiresApproval && p.ApprovalStatus != entity.PromotionApprovalApproved {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

func (m *mockActivateCampaignRepo) FindExpired(now time.Time) ([]*entity.Promotion, error) {
	return nil, nil
}

func (m *mockActivateCampaignRepo) Update(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

func newInactiveCampaign(id, title string, startAt time.Time, requiresApproval bool, approvalStatus entity.ApprovalStatus) *entity.Promotion {
	end := startAt.Add(15 * 24 * time.Hour)
	p, _ := entity.NewPromotion(id, title, "", startAt, end, requiresApproval, 20, false)
	p.ApprovalStatus = approvalStatus
	p.Events()
	return p
}

func TestActivateScheduledCampaigns_Success(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(24 * time.Hour)

	promo1 := newInactiveCampaign("promo-1", "Nowruz Auction", past, false, entity.PromotionApprovalNone)
	promo2 := newInactiveCampaign("promo-2", "Summer Sale", past, true, entity.PromotionApprovalApproved)
	promo3 := newInactiveCampaign("promo-3", "Not Ready", future, false, entity.PromotionApprovalNone)

	repo := &mockActivateCampaignRepo{items: map[string]*entity.Promotion{
		"promo-1": promo1,
		"promo-2": promo2,
		"promo-3": promo3,
	}}

	uc := usecase.NewActivateScheduledCampaignsUseCase(repo)

	output, err := uc.Execute(usecase.ActivateScheduledCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.ActivatedCampaigns) != 2 {
		t.Fatalf("expected 2 activated campaigns, got %d", len(output.ActivatedCampaigns))
	}

	ids := make(map[string]bool)
	for _, c := range output.ActivatedCampaigns {
		ids[c.CampaignID] = true
	}
	if !ids["promo-1"] {
		t.Error("expected promo-1 to be activated")
	}
	if !ids["promo-2"] {
		t.Error("expected promo-2 to be activated")
	}
	if ids["promo-3"] {
		t.Error("expected promo-3 to NOT be activated")
	}

	if !promo1.IsActive {
		t.Error("expected promo-1 to be active")
	}
	if !promo2.IsActive {
		t.Error("expected promo-2 to be active")
	}
	if promo3.IsActive {
		t.Error("expected promo-3 to remain inactive")
	}
}

func TestActivateScheduledCampaigns_Events(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)

	promo1 := newInactiveCampaign("promo-1", "Campaign 1", past, false, entity.PromotionApprovalNone)
	promo2 := newInactiveCampaign("promo-2", "Campaign 2", past, false, entity.PromotionApprovalNone)

	repo := &mockActivateCampaignRepo{items: map[string]*entity.Promotion{
		"promo-1": promo1,
		"promo-2": promo2,
	}}

	uc := usecase.NewActivateScheduledCampaignsUseCase(repo)

	output, err := uc.Execute(usecase.ActivateScheduledCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, c := range output.ActivatedCampaigns {
		if c.Event == nil {
			t.Fatalf("expected event for campaign %s, got nil", c.CampaignID)
		}
		e, ok := c.Event.(event.PromotionActivated)
		if !ok {
			t.Fatalf("expected PromotionActivated for %s, got %T", c.CampaignID, c.Event)
		}
		if e.PromotionID != c.CampaignID {
			t.Errorf("expected promotion id %s, got %s", c.CampaignID, e.PromotionID)
		}
		if e.EventName() != "promotion.campaign_activated" {
			t.Errorf("expected promotion.campaign_activated, got %s", e.EventName())
		}
	}
}

func TestActivateScheduledCampaigns_RequiresApproval(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)

	promo1 := newInactiveCampaign("promo-1", "Approved", past, true, entity.PromotionApprovalApproved)
	promo2 := newInactiveCampaign("promo-2", "Pending", past, true, entity.PromotionApprovalPending)
	promo3 := newInactiveCampaign("promo-3", "Rejected", past, true, entity.PromotionApprovalRejected)

	repo := &mockActivateCampaignRepo{items: map[string]*entity.Promotion{
		"promo-1": promo1,
		"promo-2": promo2,
		"promo-3": promo3,
	}}

	uc := usecase.NewActivateScheduledCampaignsUseCase(repo)

	output, err := uc.Execute(usecase.ActivateScheduledCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.ActivatedCampaigns) != 1 {
		t.Fatalf("expected 1 activated campaign, got %d", len(output.ActivatedCampaigns))
	}
	if output.ActivatedCampaigns[0].CampaignID != "promo-1" {
		t.Errorf("expected promo-1 to be activated, got %s", output.ActivatedCampaigns[0].CampaignID)
	}
}

func TestActivateScheduledCampaigns_NoSchedulable(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)

	promo := newInactiveCampaign("promo-1", "Future Campaign", future, false, entity.PromotionApprovalNone)

	repo := &mockActivateCampaignRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewActivateScheduledCampaignsUseCase(repo)

	output, err := uc.Execute(usecase.ActivateScheduledCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.ActivatedCampaigns) != 0 {
		t.Errorf("expected 0 activated campaigns, got %d", len(output.ActivatedCampaigns))
	}
}
