package event

import "time"

// PromotionCreated is emitted when a new promotion campaign is created.
type PromotionCreated struct {
	PromotionID string
	Title       string
	Timestamp   time.Time
}

func (e PromotionCreated) EventName() string { return "promotion.created" }

// PromotionApproved is emitted when a campaign is approved by an administrator.
type PromotionApproved struct {
	PromotionID string
	Timestamp   time.Time
}

func (e PromotionApproved) EventName() string { return "promotion.approved" }

// PromotionRejected is emitted when a campaign is rejected by an administrator.
type PromotionRejected struct {
	PromotionID string
	Timestamp   time.Time
}

func (e PromotionRejected) EventName() string { return "promotion.rejected" }

// PromotionActivated is emitted when a campaign is activated based on its time frame.
type PromotionActivated struct {
	PromotionID string
	Timestamp   time.Time
}

func (e PromotionActivated) EventName() string { return "promotion.campaign_activated" }

// PromotionDeactivated is emitted when a campaign is deactivated.
type PromotionDeactivated struct {
	PromotionID string
	Timestamp   time.Time
}

func (e PromotionDeactivated) EventName() string { return "promotion.campaign_deactivated" }
