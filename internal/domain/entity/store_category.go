package entity

import (
	"errors"
	"time"

	"kalabazar-stock-service/internal/domain/event"
)

var (
	ErrStoreCategoryInvalidStoreID    = errors.New("store id cannot be empty")
	ErrStoreCategoryInvalidCategoryID = errors.New("category id cannot be empty")
	ErrStoreCategoryAlreadyApproved   = errors.New("store category is already approved")
)

// StoreCategoryStatus represents the approval status of a store category permission.
type StoreCategoryStatus string

const (
	StoreCategoryStatusPending  StoreCategoryStatus = "pending"
	StoreCategoryStatusApproved StoreCategoryStatus = "approved"
)

// StoreCategory represents a category permission request for a store.
type StoreCategory struct {
	StoreID    string
	CategoryID string
	Status     StoreCategoryStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time

	events []any
}

// NewStoreCategory creates a new StoreCategory with pending status.
func NewStoreCategory(storeID, categoryID string) (*StoreCategory, error) {
	if storeID == "" {
		return nil, ErrStoreCategoryInvalidStoreID
	}
	if categoryID == "" {
		return nil, ErrStoreCategoryInvalidCategoryID
	}

	now := time.Now()
	sc := &StoreCategory{
		StoreID:    storeID,
		CategoryID: categoryID,
		Status:     StoreCategoryStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	sc.events = append(sc.events, event.StoreCategoryAllowed{
		StoreID:    storeID,
		CategoryID: categoryID,
		Status:     string(StoreCategoryStatusPending),
		Timestamp:  now,
	})
	return sc, nil
}

// Approve transitions the category permission to approved status.
func (sc *StoreCategory) Approve() error {
	if sc.Status == StoreCategoryStatusApproved {
		return ErrStoreCategoryAlreadyApproved
	}
	sc.Status = StoreCategoryStatusApproved
	sc.UpdatedAt = time.Now()
	sc.events = append(sc.events, event.StoreCategoryAllowed{
		StoreID:    sc.StoreID,
		CategoryID: sc.CategoryID,
		Status:     string(StoreCategoryStatusApproved),
		Timestamp:  sc.UpdatedAt,
	})
	return nil
}

// Events returns and clears the domain events produced by the entity.
func (sc *StoreCategory) Events() []any {
	events := sc.events
	sc.events = nil
	return events
}
