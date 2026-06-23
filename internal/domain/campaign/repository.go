package campaign

type CampaignRepository interface {
	Save(campaign *Campaign) error
	FindByID(id int64) (*Campaign, error)
	Update(campaign *Campaign) error
}
