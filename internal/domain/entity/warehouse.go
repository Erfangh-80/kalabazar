package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"kalabazar-stock-service/internal/domain/event"
)

var (
	ErrWarehouseInvalidName       = errors.New("warehouse name cannot be empty")
	ErrWarehouseNameTooLong       = errors.New("warehouse name cannot exceed 255 characters")
	ErrWarehouseInvalidSellerID   = errors.New("seller id cannot be empty")
	ErrWarehouseInvalidID         = errors.New("warehouse id cannot be empty")
	ErrWarehouseInvalidCapacity   = errors.New("total capacity must be greater than zero")
	ErrWarehouseCapacityExceeded  = errors.New("cannot exceed warehouse storage capacity")
	ErrWarehouseInvalidUsedAmount = errors.New("used capacity cannot be negative")
	ErrWarehouseInactive          = errors.New("warehouse is inactive")
	ErrWarehouseAlreadyActive     = errors.New("warehouse is already active")
	ErrWarehouseAlreadyInactive   = errors.New("warehouse is already inactive")
	ErrWarehouseNotFound          = errors.New("warehouse not found")
)

// WarehouseStatus represents the operational status of a warehouse.
type WarehouseStatus string

const (
	WarehouseStatusActive   WarehouseStatus = "active"
	WarehouseStatusInactive WarehouseStatus = "inactive"
)

// Address represents a physical location.
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// Validate checks that the address represents a real physical location.
func (a Address) Validate() error {
	if a.Street == "" {
		return errors.New("street cannot be empty")
	}
	if a.City == "" {
		return errors.New("city cannot be empty")
	}
	if a.Country == "" {
		return errors.New("country cannot be empty")
	}
	return nil
}

// Warehouse represents a physical storage location owned by a seller.
type Warehouse struct {
	ID           string
	SellerID     string
	Name         string
	Address      Address
	TotalCapacity   int
	UsedCapacity    int
	Status       WarehouseStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time

	events []any
}

// NewWarehouse creates a new Warehouse with active status.
func NewWarehouse(id, sellerID, name string, address Address, totalCapacity int) (*Warehouse, error) {
	if id == "" {
		return nil, ErrWarehouseInvalidID
	}
	if sellerID == "" {
		return nil, ErrWarehouseInvalidSellerID
	}
	if name == "" {
		return nil, ErrWarehouseInvalidName
	}
	if utf8.RuneCountInString(name) > 255 {
		return nil, ErrWarehouseNameTooLong
	}
	if err := address.Validate(); err != nil {
		return nil, err
	}
	if totalCapacity <= 0 {
		return nil, ErrWarehouseInvalidCapacity
	}

	now := time.Now()
	w := &Warehouse{
		ID:           id,
		SellerID:     sellerID,
		Name:         name,
		Address:      address,
		TotalCapacity:   totalCapacity,
		UsedCapacity:    0,
		Status:       WarehouseStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	w.events = append(w.events, event.WarehouseCreated{
		WarehouseID:   id,
		SellerID:      sellerID,
		WarehouseName: name,
		Timestamp:     now,
	})
	return w, nil
}

// UpdateInfo updates the warehouse's mutable information fields.
func (w *Warehouse) UpdateInfo(name string, address Address) error {
	if name == "" {
		return ErrWarehouseInvalidName
	}
	if utf8.RuneCountInString(name) > 255 {
		return ErrWarehouseNameTooLong
	}
	if err := address.Validate(); err != nil {
		return err
	}

	w.Name = name
	w.Address = address
	w.UpdatedAt = time.Now()
	w.events = append(w.events, event.WarehouseUpdated{
		WarehouseID: w.ID,
		Timestamp:   w.UpdatedAt,
	})
	return nil
}

// Activate transitions the warehouse to active status.
func (w *Warehouse) Activate() error {
	if w.Status == WarehouseStatusActive {
		return ErrWarehouseAlreadyActive
	}
	w.Status = WarehouseStatusActive
	w.UpdatedAt = time.Now()
	w.events = append(w.events, event.WarehouseActivated{
		WarehouseID: w.ID,
		Timestamp:   w.UpdatedAt,
	})
	return nil
}

// Deactivate transitions the warehouse to inactive status.
func (w *Warehouse) Deactivate() error {
	if w.Status == WarehouseStatusInactive {
		return ErrWarehouseAlreadyInactive
	}
	w.Status = WarehouseStatusInactive
	w.UpdatedAt = time.Now()
	w.events = append(w.events, event.WarehouseDeactivated{
		WarehouseID: w.ID,
		Timestamp:   w.UpdatedAt,
	})
	return nil
}

// IncreaseUsedCapacity adds the given amount to used capacity.
// Returns an error if the warehouse is inactive or if capacity would be exceeded.
func (w *Warehouse) IncreaseUsedCapacity(amount int) error {
	if amount < 0 {
		return ErrWarehouseInvalidUsedAmount
	}
	if w.Status != WarehouseStatusActive {
		return ErrWarehouseInactive
	}
	if w.UsedCapacity+amount > w.TotalCapacity {
		return ErrWarehouseCapacityExceeded
	}
	w.UsedCapacity += amount
	w.UpdatedAt = time.Now()
	return nil
}

// DecreaseUsedCapacity subtracts the given amount from used capacity.
func (w *Warehouse) DecreaseUsedCapacity(amount int) error {
	if amount < 0 {
		return ErrWarehouseInvalidUsedAmount
	}
	if w.UsedCapacity-amount < 0 {
		return ErrWarehouseInvalidUsedAmount
	}
	w.UsedCapacity -= amount
	w.UpdatedAt = time.Now()
	return nil
}

// IsAtCapacity returns true if the warehouse has reached its storage limit.
func (w *Warehouse) IsAtCapacity() bool {
	return w.UsedCapacity >= w.TotalCapacity
}

// AvailableCapacity returns the remaining storage space.
func (w *Warehouse) AvailableCapacity() int {
	return w.TotalCapacity - w.UsedCapacity
}

// Events returns and clears the domain events produced by the entity.
func (w *Warehouse) Events() []any {
	events := w.events
	w.events = nil
	return events
}

// WarehouseRepository defines the persistence contract for Warehouse entities.
type WarehouseRepository interface {
	Save(warehouse *Warehouse) error
	FindByID(id string) (*Warehouse, error)
	FindBySellerID(sellerID string) ([]*Warehouse, error)
	Update(warehouse *Warehouse) error
}
