package campaign_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/application/campaign"
	domaincampaign "stock-service-version-three/internal/domain/campaign"
)

func TestEndCampaign_Success(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewEndCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()
	c.Activate(now)
	repo.Save(c)

	resp, err := uc.Execute(campaign.EndCampaignRequest{
		CampaignID: c.ID,
		Now:        now.Add(48 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.CampaignID != c.ID {
		t.Errorf("expected CampaignID %d, got %d", c.ID, resp.CampaignID)
	}
	if resp.Status != "INACTIVE" {
		t.Errorf("expected INACTIVE, got %s", resp.Status)
	}
}

func TestEndCampaign_CampaignNotFound(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewEndCampaignUseCase(repo)

	_, err := uc.Execute(campaign.EndCampaignRequest{
		CampaignID: 999,
		Now:        time.Now(),
	})
	if err != domaincampaign.ErrCampaignNotFound {
		t.Errorf("expected ErrCampaignNotFound, got %v", err)
	}
}

func TestEndCampaign_NotExpired(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewEndCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	repo.Save(c)

	_, err := uc.Execute(campaign.EndCampaignRequest{
		CampaignID: c.ID,
		Now:        now,
	})
	if err != domaincampaign.ErrCampaignNotExpired {
		t.Errorf("expected ErrCampaignNotExpired, got %v", err)
	}
}
