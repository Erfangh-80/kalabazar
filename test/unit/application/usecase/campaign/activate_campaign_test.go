package campaign_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/application/campaign"
	domaincampaign "stock-service-version-three/internal/domain/campaign"
)

func TestActivateCampaign_Success(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewActivateCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()
	repo.Save(c)

	resp, err := uc.Execute(campaign.ActivateCampaignRequest{
		CampaignID: c.ID,
		Now:        now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.CampaignID != c.ID {
		t.Errorf("expected CampaignID %d, got %d", c.ID, resp.CampaignID)
	}
	if resp.Status != "ACTIVE" {
		t.Errorf("expected ACTIVE, got %s", resp.Status)
	}
}

func TestActivateCampaign_CampaignNotFound(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewActivateCampaignUseCase(repo)

	_, err := uc.Execute(campaign.ActivateCampaignRequest{
		CampaignID: 999,
		Now:        time.Now(),
	})
	if err != domaincampaign.ErrCampaignNotFound {
		t.Errorf("expected ErrCampaignNotFound, got %v", err)
	}
}

func TestActivateCampaign_NotApproved(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewActivateCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	repo.Save(c)

	_, err := uc.Execute(campaign.ActivateCampaignRequest{
		CampaignID: c.ID,
		Now:        now,
	})
	if err != domaincampaign.ErrCampaignNotApproved {
		t.Errorf("expected ErrCampaignNotApproved, got %v", err)
	}
}

func TestActivateCampaign_AlreadyActive(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewActivateCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()
	c.Activate(now)
	repo.Save(c)

	_, err := uc.Execute(campaign.ActivateCampaignRequest{
		CampaignID: c.ID,
		Now:        now,
	})
	if err != domaincampaign.ErrCampaignAlreadyActive {
		t.Errorf("expected ErrCampaignAlreadyActive, got %v", err)
	}
}

func TestActivateCampaign_NotStarted(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewActivateCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now.Add(1*time.Hour), now.Add(24*time.Hour))
	c.Approve()
	repo.Save(c)

	_, err := uc.Execute(campaign.ActivateCampaignRequest{
		CampaignID: c.ID,
		Now:        now,
	})
	if err != domaincampaign.ErrCampaignNotStarted {
		t.Errorf("expected ErrCampaignNotStarted, got %v", err)
	}
}
