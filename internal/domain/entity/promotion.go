package entity

import (
	"errors"
	"time"

	"kalabazar-stock-service/internal/domain/event"
)

var (
	ErrPromotionInvalidID          = errors.New("promotion id cannot be empty")
	ErrPromotionInvalidTitle       = errors.New("promotion title cannot be empty")
	ErrPromotionInvalidTimeRange   = errors.New("start time must be before end time")
	ErrPromotionAlreadyApproved    = errors.New("promotion is already approved")
	ErrPromotionAlreadyRejected    = errors.New("promotion is already rejected")
	ErrPromotionApprovalNotRequired = errors.New("promotion does not require approval")
	ErrPromotionNotApproved        = errors.New("promotion is not approved")
	ErrPromotionAlreadyActive      = errors.New("promotion is already active")
	ErrPromotionAlreadyInactive    = errors.New("promotion is already inactive")
	ErrPromotionNotFound           = errors.New("promotion not found")
)

// ApprovalStatus represents the approval state of a promotion campaign.
type ApprovalStatus string

const (
	PromotionApprovalNone     ApprovalStatus = "none"
	PromotionApprovalPending  ApprovalStatus = "pending"
	PromotionApprovalApproved ApprovalStatus = "approved"
	PromotionApprovalRejected ApprovalStatus = "rejected"
)

// Promotion represents a discount campaign with a time frame and optional approval.
type Promotion struct {
	ID               string
	Title            string
	Description      string
	StartAt          time.Time
	EndAt            time.Time
	RequiresApproval bool
	ApprovalStatus   ApprovalStatus
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time

	events []any
}

// NewPromotion creates a new Promotion campaign.
func NewPromotion(id, title, description string, startAt, endAt time.Time, requiresApproval bool) (*Promotion, error) {
	if id == "" {
		return nil, ErrPromotionInvalidID
	}
	if title == "" {
		return nil, ErrPromotionInvalidTitle
	}
	if !startAt.Before(endAt) {
		return nil, ErrPromotionInvalidTimeRange
	}

	approvalStatus := PromotionApprovalNone
	if requiresApproval {
		approvalStatus = PromotionApprovalPending
	}

	now := time.Now()
	p := &Promotion{
		ID:               id,
		Title:            title,
		Description:      description,
		StartAt:          startAt,
		EndAt:            endAt,
		RequiresApproval: requiresApproval,
		ApprovalStatus:   approvalStatus,
		IsActive:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	p.events = append(p.events, event.PromotionCreated{
		PromotionID: id,
		Title:       title,
		Timestamp:   now,
	})
	return p, nil
}

// Approve transitions the campaign to approved status.
func (p *Promotion) Approve() error {
	if !p.RequiresApproval {
		return ErrPromotionApprovalNotRequired
	}
	if p.ApprovalStatus == PromotionApprovalApproved {
		return ErrPromotionAlreadyApproved
	}
	if p.ApprovalStatus == PromotionApprovalRejected {
		return ErrPromotionAlreadyRejected
	}
	p.ApprovalStatus = PromotionApprovalApproved
	p.UpdatedAt = time.Now()
	p.events = append(p.events, event.PromotionApproved{
		PromotionID: p.ID,
		Timestamp:   p.UpdatedAt,
	})
	return nil
}

// Reject transitions the campaign to rejected status.
func (p *Promotion) Reject() error {
	if !p.RequiresApproval {
		return ErrPromotionApprovalNotRequired
	}
	if p.ApprovalStatus == PromotionApprovalRejected {
		return ErrPromotionAlreadyRejected
	}
	if p.ApprovalStatus == PromotionApprovalApproved {
		return ErrPromotionAlreadyApproved
	}
	p.ApprovalStatus = PromotionApprovalRejected
	p.UpdatedAt = time.Now()
	p.events = append(p.events, event.PromotionRejected{
		PromotionID: p.ID,
		Timestamp:   p.UpdatedAt,
	})
	return nil
}

// Activate enables the campaign so discounts apply to linked products.
func (p *Promotion) Activate() error {
	if p.IsActive {
		return ErrPromotionAlreadyActive
	}
	if p.RequiresApproval && p.ApprovalStatus != PromotionApprovalApproved {
		return ErrPromotionNotApproved
	}
	p.IsActive = true
	p.UpdatedAt = time.Now()
	p.events = append(p.events, event.PromotionActivated{
		PromotionID: p.ID,
		Timestamp:   p.UpdatedAt,
	})
	return nil
}

// Deactivate disables the campaign.
func (p *Promotion) Deactivate() error {
	if !p.IsActive {
		return ErrPromotionAlreadyInactive
	}
	p.IsActive = false
	p.UpdatedAt = time.Now()
	p.events = append(p.events, event.PromotionDeactivated{
		PromotionID: p.ID,
		Timestamp:   p.UpdatedAt,
	})
	return nil
}

// Events returns and clears the domain events produced by the entity.
func (p *Promotion) Events() []any {
	events := p.events
	p.events = nil
	return events
}

// PromotionRepository defines the persistence contract for Promotion entities.
type PromotionRepository interface {
	Save(promotion *Promotion) error
	FindByID(id string) (*Promotion, error)
	Update(promotion *Promotion) error
}
