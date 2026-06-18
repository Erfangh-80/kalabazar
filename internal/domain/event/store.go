package event

import "time"

// StoreCategoryAllowed is emitted when a category is requested/approved for a store.
type StoreCategoryAllowed struct {
	StoreID    string
	CategoryID string
	Status     string
	Timestamp  time.Time
}

func (e StoreCategoryAllowed) EventName() string { return "store.category_allowed" }

// StoreCreated is emitted when a new store is registered.
type StoreCreated struct {
	StoreID   string
	UserID    string
	StoreName string
	Timestamp time.Time
}

func (e StoreCreated) EventName() string { return "store.created" }

// StoreUpdated is emitted when store information is modified.
type StoreUpdated struct {
	StoreID   string
	Timestamp time.Time
}

func (e StoreUpdated) EventName() string { return "store.updated" }

// StoreActivated is emitted when a store transitions to active status.
type StoreActivated struct {
	StoreID   string
	Timestamp time.Time
}

func (e StoreActivated) EventName() string { return "store.activated" }

// StoreDeactivated is emitted when a store transitions to inactive status.
type StoreDeactivated struct {
	StoreID   string
	Timestamp time.Time
}

func (e StoreDeactivated) EventName() string { return "store.deactivated" }

// StoreRejected is emitted when a pending store is rejected by admin.
type StoreRejected struct {
	StoreID   string
	Timestamp time.Time
}

func (e StoreRejected) EventName() string { return "store.rejected" }
