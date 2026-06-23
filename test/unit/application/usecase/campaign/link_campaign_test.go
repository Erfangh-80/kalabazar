package campaign_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/application/campaign"
	domaincampaign "stock-service-version-three/internal/domain/campaign"
)

func TestLinkCampaign_Success(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewLinkCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	repo.Save(c)

	resp, err := uc.Execute(campaign.LinkCampaignRequest{
		CampaignID:  c.ID,
		InventoryID: 42,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.CampaignID != c.ID {
		t.Errorf("expected CampaignID %d, got %d", c.ID, resp.CampaignID)
	}
	if resp.InventoryID != 42 {
		t.Errorf("expected InventoryID 42, got %d", resp.InventoryID)
	}
}

func TestLinkCampaign_CampaignNotFound(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewLinkCampaignUseCase(repo)

	_, err := uc.Execute(campaign.LinkCampaignRequest{
		CampaignID:  999,
		InventoryID: 42,
	})
	if err != domaincampaign.ErrCampaignNotFound {
		t.Errorf("expected ErrCampaignNotFound, got %v", err)
	}
}
