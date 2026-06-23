package campaign_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/application/campaign"
	domaincampaign "stock-service-version-three/internal/domain/campaign"
)

func TestApproveCampaign_Success(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewApproveCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	repo.Save(c)

	resp, err := uc.Execute(campaign.ApproveCampaignRequest{
		CampaignID: c.ID,
		Decision:   "approved",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.CampaignID != c.ID {
		t.Errorf("expected CampaignID %d, got %d", c.ID, resp.CampaignID)
	}
	if resp.ApprovalStatus != "APPROVED" {
		t.Errorf("expected APPROVED, got %s", resp.ApprovalStatus)
	}
}

func TestApproveCampaign_CampaignNotFound(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewApproveCampaignUseCase(repo)

	_, err := uc.Execute(campaign.ApproveCampaignRequest{
		CampaignID: 999,
		Decision:   "approved",
	})
	if err != domaincampaign.ErrCampaignNotFound {
		t.Errorf("expected ErrCampaignNotFound, got %v", err)
	}
}

func TestApproveCampaign_AlreadyApproved(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewApproveCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	c.Approve()
	repo.Save(c)

	_, err := uc.Execute(campaign.ApproveCampaignRequest{
		CampaignID: c.ID,
		Decision:   "approved",
	})
	if err != domaincampaign.ErrCampaignAlreadyApproved {
		t.Errorf("expected ErrCampaignAlreadyApproved, got %v", err)
	}
}

func TestApproveCampaign_Rejected(t *testing.T) {
	repo := newMockRepository()
	uc := campaign.NewApproveCampaignUseCase(repo)

	now := time.Now()
	c := domaincampaign.NewCampaign("Test", "percentage", 10, now, now.Add(24*time.Hour))
	repo.Save(c)

	resp, err := uc.Execute(campaign.ApproveCampaignRequest{
		CampaignID: c.ID,
		Decision:   "rejected",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ApprovalStatus != "REJECTED" {
		t.Errorf("expected REJECTED, got %s", resp.ApprovalStatus)
	}
}
