package campaign

type CampaignCreatedEvent struct {
	CampaignID int64
	Title      string
}

type CampaignLinkedToInventoryEvent struct {
	CampaignID  int64
	InventoryID int64
}

type CampaignApprovedEvent struct {
	CampaignID int64
}

type CampaignActivatedEvent struct {
	CampaignID int64
}

type CampaignEndedEvent struct {
	CampaignID int64
}
