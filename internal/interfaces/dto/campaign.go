package dto

type CreateCampaignRequest struct {
	Title        string  `json:"title"`
	DiscountType string  `json:"discount_type"`
	Value        float64 `json:"value"`
	StartAt      string  `json:"start_at"`
	EndAt        string  `json:"end_at"`
}
type CreateCampaignResponse struct {
	CampaignID     int64  `json:"campaign_id"`
	Status         string `json:"status"`
	ApprovalStatus string `json:"approval_status"`
}
type LinkCampaignRequest struct {
	CampaignID  int64 `json:"campaign_id"`
	InventoryID int64 `json:"inventory_id"`
}
type LinkCampaignResponse struct {
	CampaignID  int64 `json:"campaign_id"`
	InventoryID int64 `json:"inventory_id"`
}
type ApproveCampaignRequest struct {
	CampaignID int64  `json:"campaign_id"`
	Decision   string `json:"decision"`
}
type ApproveCampaignResponse struct {
	CampaignID     int64  `json:"campaign_id"`
	ApprovalStatus string `json:"approval_status"`
}
type ActivateCampaignRequest struct {
	CampaignID int64  `json:"campaign_id"`
	Now        string `json:"now"`
}
type ActivateCampaignResponse struct {
	CampaignID int64  `json:"campaign_id"`
	Status     string `json:"status"`
}
type EndCampaignRequest struct {
	CampaignID int64  `json:"campaign_id"`
	Now        string `json:"now"`
}
type EndCampaignResponse struct {
	CampaignID int64  `json:"campaign_id"`
	Status     string `json:"status"`
}
