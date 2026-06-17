package event

import "time"

// WarehouseCreated is emitted when a new warehouse is registered.
type WarehouseCreated struct {
	WarehouseID   string
	SellerID      string
	WarehouseName string
	Timestamp     time.Time
}

func (e WarehouseCreated) EventName() string { return "warehouse.created" }

// WarehouseUpdated is emitted when warehouse information is modified.
type WarehouseUpdated struct {
	WarehouseID string
	Timestamp   time.Time
}

func (e WarehouseUpdated) EventName() string { return "warehouse.updated" }

// WarehouseActivated is emitted when a warehouse transitions to active status.
type WarehouseActivated struct {
	WarehouseID string
	Timestamp   time.Time
}

func (e WarehouseActivated) EventName() string { return "warehouse.activated" }

// WarehouseDeactivated is emitted when a warehouse transitions to inactive status.
type WarehouseDeactivated struct {
	WarehouseID string
	Timestamp   time.Time
}

func (e WarehouseDeactivated) EventName() string { return "warehouse.deactivated" }
