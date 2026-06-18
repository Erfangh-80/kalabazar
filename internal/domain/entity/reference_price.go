package entity

import (
	"errors"
	"time"

	"kalabazar-stock-service/internal/domain/event"
)

var (
	ErrReferencePriceInvalidID        = errors.New("reference price id cannot be empty")
	ErrReferencePriceInvalidProductID = errors.New("product id cannot be empty")
	ErrReferencePriceInvalidPrice     = errors.New("reference price must be greater than zero")
	ErrReferencePriceInvalidSource    = errors.New("reference price source cannot be empty")
	ErrReferencePriceNotFound         = errors.New("reference price not found")
	ErrInvalidBasePrice               = errors.New("base price must be greater than zero")
	ErrInvalidDiscountPercent         = errors.New("discount percentage must be between 0 and 100")
)

// ReferencePrice represents a market price observation for a product from a specific source.
type ReferencePrice struct {
	ID        string
	ProductID string
	Price     float64
	Source    string
	CreatedAt time.Time

	events []any
}

// NewReferencePrice creates a new ReferencePrice with the given values.
func NewReferencePrice(id, productID string, price float64, source string) (*ReferencePrice, error) {
	now := time.Now()
	rp := &ReferencePrice{
		ID:        id,
		ProductID: productID,
		Price:     price,
		Source:    source,
		CreatedAt: now,
	}
	if err := rp.validate(); err != nil {
		return nil, err
	}
	rp.events = append(rp.events, event.ReferencePriceCreated{
		ReferencePriceID: id,
		ProductID:        productID,
		Price:            price,
		Source:           source,
		Timestamp:        now,
	})
	return rp, nil
}

// validate checks all invariant business rules for the ReferencePrice entity.
func (rp *ReferencePrice) validate() error {
	switch {
	case rp.ID == "":
		return ErrReferencePriceInvalidID
	case rp.ProductID == "":
		return ErrReferencePriceInvalidProductID
	case rp.Price <= 0:
		return ErrReferencePriceInvalidPrice
	case rp.Source == "":
		return ErrReferencePriceInvalidSource
	default:
		return nil
	}
}

// Events returns and clears the domain events produced by the entity.
func (rp *ReferencePrice) Events() []any {
	events := rp.events
	rp.events = nil
	return events
}

// CalculateFinalPrice computes the final price after applying a campaign discount.
// basePrice is the product's base price, discountPercent is the discount percentage (0–100).
func CalculateFinalPrice(basePrice, discountPercent float64) (float64, error) {
	if basePrice <= 0 {
		return 0, ErrInvalidBasePrice
	}
	if discountPercent < 0 || discountPercent > 100 {
		return 0, ErrInvalidDiscountPercent
	}
	discount := basePrice * (discountPercent / 100)
	return basePrice - discount, nil
}

// ReferencePriceRepository defines the persistence contract for ReferencePrice entities.
type ReferencePriceRepository interface {
	Save(referencePrice *ReferencePrice) error
	FindByID(id string) (*ReferencePrice, error)
	FindByProductID(productID string) ([]*ReferencePrice, error)
}
