package campaign

import "time"

type CampaignInventoryLink struct {
	ID          int64
	CampaignID  int64
	InventoryID int64
	CreatedAt   time.Time
}

func NewCampaignInventoryLink(campaignID, inventoryID int64) *CampaignInventoryLink {
	return &CampaignInventoryLink{
		CampaignID:  campaignID,
		InventoryID: inventoryID,
		CreatedAt:   time.Now(),
	}
}
