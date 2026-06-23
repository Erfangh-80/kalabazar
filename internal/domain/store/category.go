package store

import "time"

type AllowedCategoryStatus string

const (
	AllowedCategoryStatusAPPROVED AllowedCategoryStatus = "APPROVED"
	AllowedCategoryStatusPENDING  AllowedCategoryStatus = "PENDING"
	AllowedCategoryStatusREJECTED AllowedCategoryStatus = "REJECTED"
)

type StoreAllowedCategory struct {
	ID         int64
	StoreID    int64
	CategoryID int64
	Status     AllowedCategoryStatus
	CreatedAt  time.Time
	events     []any
}

func NewStoreAllowedCategory(storeID, categoryID int64) *StoreAllowedCategory {
	c := &StoreAllowedCategory{
		StoreID:    storeID,
		CategoryID: categoryID,
		Status:     AllowedCategoryStatusAPPROVED,
		CreatedAt:  time.Now(),
	}
	c.emit(StoreCategoryAllowedEvent{
		StoreID:    storeID,
		CategoryID: categoryID,
	})
	return c
}

func (c *StoreAllowedCategory) Events() []any {
	evts := c.events
	c.events = nil
	return evts
}

func (c *StoreAllowedCategory) emit(event any) {
	c.events = append(c.events, event)
}
