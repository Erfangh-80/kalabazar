package campaign

import (
	"sync/atomic"
	"time"
)

var nextID int64

func nextCampaignID() int64 {
	return atomic.AddInt64(&nextID, 1)
}

type CampaignStatus string

const (
	CampaignStatusInactive CampaignStatus = "INACTIVE"
	CampaignStatusActive   CampaignStatus = "ACTIVE"
)

type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "PENDING"
	ApprovalStatusApproved ApprovalStatus = "APPROVED"
	ApprovalStatusRejected ApprovalStatus = "REJECTED"
)

type Campaign struct {
	ID             int64
	Title          string
	DiscountType   string
	Value          float64
	StartAt        time.Time
	EndAt          time.Time
	Status         CampaignStatus
	ApprovalStatus ApprovalStatus
	CreatedAt      time.Time

	events []interface{}
}

func NewCampaign(title, discountType string, value float64, startAt, endAt time.Time) *Campaign {
	c := &Campaign{
		ID:             nextCampaignID(),
		Title:          title,
		DiscountType:   discountType,
		Value:          value,
		StartAt:        startAt,
		EndAt:          endAt,
		Status:         CampaignStatusInactive,
		ApprovalStatus: ApprovalStatusPending,
		CreatedAt:      time.Now(),
	}
	c.events = append(c.events, CampaignCreatedEvent{
		CampaignID: c.ID,
		Title:      c.Title,
	})
	return c
}

func (c *Campaign) Events() []interface{} {
	events := c.events
	c.events = nil
	return events
}

func (c *Campaign) Approve() ([]interface{}, error) {
	if c.ApprovalStatus == ApprovalStatusApproved {
		return nil, ErrCampaignAlreadyApproved
	}
	c.ApprovalStatus = ApprovalStatusApproved
	evt := CampaignApprovedEvent{CampaignID: c.ID}
	c.events = append(c.events, evt)
	return []interface{}{evt}, nil
}

func (c *Campaign) Activate(now time.Time) ([]interface{}, error) {
	if c.ApprovalStatus != ApprovalStatusApproved {
		return nil, ErrCampaignNotApproved
	}
	if c.Status == CampaignStatusActive {
		return nil, ErrCampaignAlreadyActive
	}
	if now.Before(c.StartAt) {
		return nil, ErrCampaignNotStarted
	}
	c.Status = CampaignStatusActive
	evt := CampaignActivatedEvent{CampaignID: c.ID}
	c.events = append(c.events, evt)
	return []interface{}{evt}, nil
}

func (c *Campaign) End(now time.Time) ([]interface{}, error) {
	if !now.After(c.EndAt) {
		return nil, ErrCampaignNotExpired
	}
	c.Status = CampaignStatusInactive
	evt := CampaignEndedEvent{CampaignID: c.ID}
	c.events = append(c.events, evt)
	return []interface{}{evt}, nil
}

func (c *Campaign) IsActive() bool {
	return c.Status == CampaignStatusActive
}

func (c *Campaign) IsExpired() bool {
	return time.Now().After(c.EndAt)
}

func (c *Campaign) IsExpiredAt(t time.Time) bool {
	return t.After(c.EndAt)
}

func (c *Campaign) LinkToInventory(inventoryID int64) []interface{} {
	evt := CampaignLinkedToInventoryEvent{
		CampaignID:  c.ID,
		InventoryID: inventoryID,
	}
	c.events = append(c.events, evt)
	return []interface{}{evt}
}
