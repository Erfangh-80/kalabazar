package entity

import (
	"errors"
	"time"

	"kalabazar-stock-service/internal/domain/event"
)

var (
	ErrInventoryInvalidID                = errors.New("inventory id cannot be empty")
	ErrInventoryInvalidStoreID           = errors.New("store id cannot be empty")
	ErrInventoryInvalidWarehouseID       = errors.New("warehouse id cannot be empty")
	ErrInventoryInvalidProductID         = errors.New("product id cannot be empty")
	ErrInventoryInvalidBasePrice         = errors.New("base price must be greater than zero")
	ErrInventoryInvalidStock             = errors.New("stock quantity cannot be negative")
	ErrInventoryInvalidVendorStatus      = errors.New("invalid vendor sale status")
	ErrInventoryInvalidSystemStatus      = errors.New("invalid system sale status")
	ErrInventoryInvalidPrice             = errors.New("final price cannot be negative")
	ErrInventoryInvalidTimeRange         = errors.New("end time must be after start time")
	ErrInventoryNotFound                 = errors.New("inventory item not found")
	ErrInventoryAlreadyLinkedToPromotion = errors.New("inventory already linked to a promotion")
	ErrInventoryInvalidPromotionID       = errors.New("promotion id cannot be empty")
	ErrInventoryInvalidPromotionStatus   = errors.New("invalid promotion status")
	ErrInventoryInvalidSaleModel         = errors.New("sale model cannot be empty")
	ErrInventoryInvalidCondition         = errors.New("condition cannot be empty")
	ErrInventoryInvalidMinOrderQty       = errors.New("minimum order quantity must be greater than zero")
	ErrInventoryInvalidMaxOrderQty       = errors.New("maximum order quantity must be greater than zero")
	ErrInventoryBelowMinOrder            = errors.New("quantity is below minimum order quantity")
	ErrInventoryAboveMaxOrder            = errors.New("quantity exceeds maximum order quantity")
	ErrInventoryInsufficientStock        = errors.New("insufficient stock")
)

// VendorSaleStatus represents the sale status controlled by the seller.
type VendorSaleStatus string

const (
	VendorSaleStatusActive    VendorSaleStatus = "active"
	VendorSaleStatusInactive  VendorSaleStatus = "inactive"
	VendorSaleStatusScheduled VendorSaleStatus = "scheduled"
	VendorSaleStatusDraft     VendorSaleStatus = "draft"
)

// SystemSaleStatus represents the sale status controlled by the system administrator.
type SystemSaleStatus string

const (
	SystemSaleStatusActive   SystemSaleStatus = "active"
	SystemSaleStatusInactive SystemSaleStatus = "inactive"
)

// CampaignApprovalStatus represents the approval state of a campaign linked to the item.
type CampaignApprovalStatus string

const (
	CampaignApprovalPending  CampaignApprovalStatus = "pending"
	CampaignApprovalApproved CampaignApprovalStatus = "approved"
	CampaignApprovalRejected CampaignApprovalStatus = "rejected"
)

// validVendorStatuses contains all valid vendor sale status values.
var validVendorStatuses = map[VendorSaleStatus]bool{
	VendorSaleStatusActive:    true,
	VendorSaleStatusInactive:  true,
	VendorSaleStatusScheduled: true,
	VendorSaleStatusDraft:     true,
}

// validSystemStatuses contains all valid system sale status values.
var validSystemStatuses = map[SystemSaleStatus]bool{
	SystemSaleStatusActive:   true,
	SystemSaleStatusInactive: true,
}

// Inventory represents a product item registered for sale in a specific warehouse.
type Inventory struct {
	ID                     string
	StoreID                string
	WarehouseID            string
	ProductID              string
	PromotionID            *string
	SaleModel              string
	Condition              string
	VendorSaleStatus       VendorSaleStatus
	SystemSaleStatus       SystemSaleStatus
	CampaignApprovalStatus CampaignApprovalStatus
	InstantQty             int
	MinOrderQty            int
	MaxOrderQty            *int
	BasePrice              float64
	FinalPrice             float64
	Attributes             map[string]string
	StartAt                *time.Time
	EndAt                  *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time

	events []any
}

// NewInventory creates a new Inventory item with default statuses.
func NewInventory(id, storeID, warehouseID, productID string, basePrice float64, instantQty int,
	saleModel, condition string, minOrderQty int, maxOrderQty *int, attributes map[string]string) (*Inventory, error) {
	if attributes == nil {
		attributes = make(map[string]string)
	}

	now := time.Now()
	inv := &Inventory{
		ID:                     id,
		StoreID:                storeID,
		WarehouseID:            warehouseID,
		ProductID:              productID,
		SaleModel:              saleModel,
		Condition:              condition,
		MinOrderQty:            minOrderQty,
		MaxOrderQty:            maxOrderQty,
		Attributes:             attributes,
		VendorSaleStatus:       VendorSaleStatusActive,
		SystemSaleStatus:       SystemSaleStatusActive,
		CampaignApprovalStatus: CampaignApprovalPending,
		InstantQty:             instantQty,
		BasePrice:              basePrice,
		FinalPrice:             basePrice,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if err := inv.validate(); err != nil {
		return nil, err
	}
	inv.events = append(inv.events, event.InventoryItemCreated{
		InventoryID: id,
		ProductID:   productID,
		Timestamp:   now,
	})
	return inv, nil
}

// UpdateStock changes the current stock quantity.
func (inv *Inventory) UpdateStock(qty int) error {
	oldQty := inv.InstantQty
	inv.InstantQty = qty
	if err := inv.validate(); err != nil {
		inv.InstantQty = oldQty
		return err
	}
	inv.UpdatedAt = time.Now()
	inv.events = append(inv.events, event.InventoryStockUpdated{
		InventoryID: inv.ID,
		NewQty:      qty,
		Timestamp:   inv.UpdatedAt,
	})
	return nil
}

// SetVendorStatus changes the vendor-controlled sale status.
func (inv *Inventory) SetVendorStatus(status VendorSaleStatus) error {
	if inv.VendorSaleStatus == status {
		return nil
	}

	oldStatus := inv.VendorSaleStatus
	inv.VendorSaleStatus = status
	if err := inv.validate(); err != nil {
		inv.VendorSaleStatus = oldStatus
		return err
	}
	inv.UpdatedAt = time.Now()

	if status == VendorSaleStatusActive && oldStatus != VendorSaleStatusActive {
		inv.events = append(inv.events, event.InventoryItemActivated{
			InventoryID: inv.ID,
			Timestamp:   inv.UpdatedAt,
		})
	} else if oldStatus == VendorSaleStatusActive && status != VendorSaleStatusActive {
		inv.events = append(inv.events, event.InventoryItemDeactivated{
			InventoryID: inv.ID,
			Timestamp:   inv.UpdatedAt,
		})
	}
	return nil
}

// SetSystemStatus changes the system-controlled sale status.
func (inv *Inventory) SetSystemStatus(status SystemSaleStatus) error {
	if inv.SystemSaleStatus == status {
		return nil
	}

	oldStatus := inv.SystemSaleStatus
	inv.SystemSaleStatus = status
	if err := inv.validate(); err != nil {
		inv.SystemSaleStatus = oldStatus
		return err
	}
	inv.UpdatedAt = time.Now()

	if status == SystemSaleStatusInactive {
		inv.events = append(inv.events, event.InventorySystemBlocked{
			InventoryID: inv.ID,
			Timestamp:   inv.UpdatedAt,
		})
	} else {
		inv.events = append(inv.events, event.InventorySystemUnblocked{
			InventoryID: inv.ID,
			Timestamp:   inv.UpdatedAt,
		})
	}
	return nil
}

// SetSaleSchedule sets the time window for the item's sale.
func (inv *Inventory) SetSaleSchedule(startAt, endAt *time.Time) error {
	oldStart, oldEnd := inv.StartAt, inv.EndAt
	inv.StartAt = startAt
	inv.EndAt = endAt
	if err := inv.validate(); err != nil {
		inv.StartAt, inv.EndAt = oldStart, oldEnd
		return err
	}
	inv.UpdatedAt = time.Now()

	inv.events = append(inv.events, event.InventorySaleScheduled{
		InventoryID: inv.ID,
		StartAt:     startAt,
		EndAt:       endAt,
		Timestamp:   inv.UpdatedAt,
	})
	return nil
}

// UpdatePrice changes the base and final price of the item.
func (inv *Inventory) UpdatePrice(basePrice, finalPrice float64) error {
	oldBase, oldFinal := inv.BasePrice, inv.FinalPrice
	inv.BasePrice = basePrice
	inv.FinalPrice = finalPrice
	if err := inv.validate(); err != nil {
		inv.BasePrice, inv.FinalPrice = oldBase, oldFinal
		return err
	}
	inv.UpdatedAt = time.Now()

	inv.events = append(inv.events, event.InventoryPriceUpdated{
		InventoryID: inv.ID,
		BasePrice:   basePrice,
		FinalPrice:  finalPrice,
		Timestamp:   inv.UpdatedAt,
	})
	return nil
}

// CanBeSold checks whether the item is available for sale.
// Both vendor and system status must be active and stock must be above zero.
func (inv *Inventory) CanBeSold() bool {
	return inv.VendorSaleStatus == VendorSaleStatusActive &&
		inv.SystemSaleStatus == SystemSaleStatusActive &&
		inv.InstantQty > 0
}

// validate checks all invariant business rules for the Inventory entity.
func (inv *Inventory) validate() error {
	switch {
	case inv.ID == "":
		return ErrInventoryInvalidID
	case inv.StoreID == "":
		return ErrInventoryInvalidStoreID
	case inv.WarehouseID == "":
		return ErrInventoryInvalidWarehouseID
	case inv.ProductID == "":
		return ErrInventoryInvalidProductID
	case inv.BasePrice <= 0:
		return ErrInventoryInvalidBasePrice
	case inv.InstantQty < 0:
		return ErrInventoryInvalidStock
	case inv.SaleModel == "":
		return ErrInventoryInvalidSaleModel
	case inv.Condition == "":
		return ErrInventoryInvalidCondition
	case inv.MinOrderQty <= 0:
		return ErrInventoryInvalidMinOrderQty
	case inv.MaxOrderQty != nil && *inv.MaxOrderQty <= 0:
		return ErrInventoryInvalidMaxOrderQty
	case inv.MaxOrderQty != nil && *inv.MaxOrderQty < inv.MinOrderQty:
		return ErrInventoryInvalidMaxOrderQty
	case !validVendorStatuses[inv.VendorSaleStatus]:
		return ErrInventoryInvalidVendorStatus
	case !validSystemStatuses[inv.SystemSaleStatus]:
		return ErrInventoryInvalidSystemStatus
	case inv.FinalPrice < 0:
		return ErrInventoryInvalidPrice
	case inv.EndAt != nil && (inv.StartAt == nil || !inv.StartAt.Before(*inv.EndAt)):
		return ErrInventoryInvalidTimeRange
	default:
		return nil
	}
}

// ValidatePurchase checks whether the given quantity can be purchased.
// Returns ErrInventoryBelowMinOrder, ErrInventoryAboveMaxOrder, or ErrInventoryInsufficientStock.
func (inv *Inventory) ValidatePurchase(qty int) error {
	return inv.validateQuantity(qty)
}

func (inv *Inventory) validateQuantity(qty int) error {
	if qty < inv.MinOrderQty {
		return ErrInventoryBelowMinOrder
	}
	if inv.MaxOrderQty != nil && qty > *inv.MaxOrderQty {
		return ErrInventoryAboveMaxOrder
	}
	if qty > inv.InstantQty {
		return ErrInventoryInsufficientStock
	}
	return nil
}

// LinkPromotion links a promotion campaign to this inventory item.
func (inv *Inventory) LinkPromotion(promotionID string) error {
	if promotionID == "" {
		return ErrInventoryInvalidPromotionID
	}
	if inv.PromotionID != nil {
		return ErrInventoryAlreadyLinkedToPromotion
	}
	inv.PromotionID = &promotionID
	inv.UpdatedAt = time.Now()
	inv.events = append(inv.events, event.InventoryPromotionLinked{
		InventoryID: inv.ID,
		PromotionID: promotionID,
		Timestamp:   inv.UpdatedAt,
	})
	return nil
}

// UpdatePromotionStatus changes the campaign approval status on this item.
func (inv *Inventory) UpdatePromotionStatus(status CampaignApprovalStatus) error {
	if status != CampaignApprovalPending && status != CampaignApprovalApproved && status != CampaignApprovalRejected {
		return ErrInventoryInvalidPromotionStatus
	}
	inv.CampaignApprovalStatus = status
	inv.UpdatedAt = time.Now()
	inv.events = append(inv.events, event.InventoryPromotionStatusChanged{
		InventoryID: inv.ID,
		Status:      string(status),
		Timestamp:   inv.UpdatedAt,
	})
	return nil
}

// ResetPrice resets the final price to match the base price.
func (inv *Inventory) ResetPrice() error {
	inv.FinalPrice = inv.BasePrice
	inv.UpdatedAt = time.Now()
	inv.events = append(inv.events, event.InventoryPriceUpdated{
		InventoryID: inv.ID,
		BasePrice:   inv.BasePrice,
		FinalPrice:  inv.FinalPrice,
		Timestamp:   inv.UpdatedAt,
	})
	return nil
}

// Events returns and clears the domain events produced by the entity.
func (inv *Inventory) Events() []any {
	events := inv.events
	inv.events = nil
	return events
}

// InventoryRepository defines the persistence contract for Inventory entities.
type InventoryRepository interface {
	Save(inventory *Inventory) error
	FindByID(id string) (*Inventory, error)
	FindByStoreID(storeID string) ([]*Inventory, error)
	FindByWarehouseID(warehouseID string) ([]*Inventory, error)
	FindByProductID(productID string) ([]*Inventory, error)
	Update(inventory *Inventory) error
}
