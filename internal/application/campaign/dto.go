package campaign

import "time"

type CreateCampaignRequest struct {
	Title        string
	DiscountType string
	Value        float64
	StartAt      time.Time
	EndAt        time.Time
}

type CreateCampaignResponse struct {
	CampaignID     int64
	Status         string
	ApprovalStatus string
}

type LinkCampaignRequest struct {
	CampaignID  int64
	InventoryID int64
}

type LinkCampaignResponse struct {
	CampaignID  int64
	InventoryID int64
}

type ApproveCampaignRequest struct {
	CampaignID int64
	Decision   string
}

type ApproveCampaignResponse struct {
	CampaignID     int64
	ApprovalStatus string
}

type ActivateCampaignRequest struct {
	CampaignID int64
	Now        time.Time
}

type ActivateCampaignResponse struct {
	CampaignID int64
	Status     string
}

type EndCampaignRequest struct {
	CampaignID int64
	Now        time.Time
}

type EndCampaignResponse struct {
	CampaignID int64
	Status     string
}
