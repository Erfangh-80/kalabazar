package event

import "time"

// InventoryItemCreated is emitted when a new inventory item is registered.
type InventoryItemCreated struct {
	InventoryID string
	ProductID   string
	Timestamp   time.Time
}

func (e InventoryItemCreated) EventName() string { return "inventory.item_created" }

// InventoryStockUpdated is emitted when the stock quantity changes.
type InventoryStockUpdated struct {
	InventoryID string
	NewQty      int
	Timestamp   time.Time
}

func (e InventoryStockUpdated) EventName() string { return "inventory.stock_updated" }

// InventoryItemActivated is emitted when the vendor activates the item for sale.
type InventoryItemActivated struct {
	InventoryID string
	Timestamp   time.Time
}

func (e InventoryItemActivated) EventName() string { return "inventory.item_activated" }

// InventoryItemDeactivated is emitted when the vendor deactivates the item.
type InventoryItemDeactivated struct {
	InventoryID string
	Timestamp   time.Time
}

func (e InventoryItemDeactivated) EventName() string { return "inventory.item_deactivated" }

// InventorySystemBlocked is emitted when the system blocks an item.
type InventorySystemBlocked struct {
	InventoryID string
	Timestamp   time.Time
}

func (e InventorySystemBlocked) EventName() string { return "inventory.system_blocked" }

// InventorySystemUnblocked is emitted when the system unblocks an item.
type InventorySystemUnblocked struct {
	InventoryID string
	Timestamp   time.Time
}

func (e InventorySystemUnblocked) EventName() string { return "inventory.system_unblocked" }

// InventorySaleScheduled is emitted when a sale schedule is set on an item.
type InventorySaleScheduled struct {
	InventoryID string
	StartAt     *time.Time
	EndAt       *time.Time
	Timestamp   time.Time
}

func (e InventorySaleScheduled) EventName() string { return "inventory.sale_scheduled" }

// InventoryPriceUpdated is emitted when the item price changes.
type InventoryPriceUpdated struct {
	InventoryID string
	BasePrice   float64
	FinalPrice  float64
	Timestamp   time.Time
}

func (e InventoryPriceUpdated) EventName() string { return "inventory.price_updated" }

// InventoryPromotionLinked is emitted when a promotion campaign is linked to an inventory item.
type InventoryPromotionLinked struct {
	InventoryID string
	PromotionID string
	Timestamp   time.Time
}

func (e InventoryPromotionLinked) EventName() string { return "inventory.promotion_linked" }

// InventoryPromotionStatusChanged is emitted when the promotion approval status changes.
type InventoryPromotionStatusChanged struct {
	InventoryID string
	Status      string
	Timestamp   time.Time
}

func (e InventoryPromotionStatusChanged) EventName() string { return "inventory.promotion_status_changed" }
