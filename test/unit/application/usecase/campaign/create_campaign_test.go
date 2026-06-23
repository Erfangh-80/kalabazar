package campaign_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/application/campaign"
)

func TestCreateCampaign_Success(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewCreateCampaignUseCase(repo)

	now := time.Now()
	req := campaign.CreateCampaignRequest{
		Title:        "Test Campaign",
		DiscountType: "percentage",
		Value:        15,
		StartAt:      now,
		EndAt:        now.Add(24 * time.Hour),
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.CampaignID == 0 {
		t.Error("expected non-zero campaign ID")
	}
	if resp.Status != "INACTIVE" {
		t.Errorf("expected INACTIVE, got %s", resp.Status)
	}
	if resp.ApprovalStatus != "PENDING" {
		t.Errorf("expected PENDING, got %s", resp.ApprovalStatus)
	}

	if len(repo.campaigns) != 1 {
		t.Errorf("expected 1 campaign in repo, got %d", len(repo.campaigns))
	}
}
