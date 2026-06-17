package event

import "time"

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
