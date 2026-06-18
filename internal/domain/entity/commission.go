package entity

import (
	"errors"
	"time"

	"kalabazar-stock-service/internal/domain/event"
)

var (
	ErrCommissionInvalidID          = errors.New("commission id cannot be empty")
	ErrCommissionInvalidProductID   = errors.New("product id cannot be empty")
	ErrCommissionInvalidRate        = errors.New("commission rate must be between 0 and 100")
	ErrCommissionInvalidPriceRange  = errors.New("price range is invalid")
	ErrCommissionInvalidMinQty      = errors.New("minimum quantity cannot be negative")
	ErrCommissionConditionsNotMet   = errors.New("commission conditions not met")
	ErrCommissionNotFound           = errors.New("commission rule not found")
)

// Commission represents a commission rule defined for a product or category.
type Commission struct {
	ID          string
	ProductID   string
	SalesModel  string
	RatePercent float64
	MinPrice    float64
	MaxPrice    float64
	MinQty      int
	CreatedAt   time.Time

	events []any
}

// NewCommission creates a new Commission rule with validation.
func NewCommission(id, productID, salesModel string, ratePercent float64, minPrice, maxPrice float64, minQty int) (*Commission, error) {
	now := time.Now()
	c := &Commission{
		ID:          id,
		ProductID:   productID,
		SalesModel:  salesModel,
		RatePercent: ratePercent,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		MinQty:      minQty,
		CreatedAt:   now,
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	c.events = append(c.events, event.CommissionRuleCreated{
		CommissionID: id,
		ProductID:    productID,
		RatePercent:  ratePercent,
		Timestamp:    now,
	})
	return c, nil
}

// validate checks all invariant business rules for the Commission entity.
func (c *Commission) validate() error {
	switch {
	case c.ID == "":
		return ErrCommissionInvalidID
	case c.ProductID == "":
		return ErrCommissionInvalidProductID
	case c.RatePercent <= 0 || c.RatePercent > 100:
		return ErrCommissionInvalidRate
	case c.MinPrice < 0 || c.MaxPrice < 0 || c.MinPrice > c.MaxPrice:
		return ErrCommissionInvalidPriceRange
	case c.MinQty < 0:
		return ErrCommissionInvalidMinQty
	default:
		return nil
	}
}

// Calculate computes the commission amount for a given sale.
// Returns an error if the sale does not meet the rule conditions.
func (c *Commission) Calculate(saleAmount float64, quantity int) (float64, error) {
	if saleAmount < c.MinPrice || saleAmount > c.MaxPrice {
		return 0, ErrCommissionConditionsNotMet
	}
	if quantity < c.MinQty {
		return 0, ErrCommissionConditionsNotMet
	}
	commissionAmount := saleAmount * (c.RatePercent / 100)
	c.events = append(c.events, event.CommissionCalculated{
		CommissionID:     c.ID,
		SaleAmount:       saleAmount,
		CommissionAmount: commissionAmount,
		Timestamp:        time.Now(),
	})
	return commissionAmount, nil
}

// Events returns and clears the domain events produced by the entity.
func (c *Commission) Events() []any {
	events := c.events
	c.events = nil
	return events
}

// CommissionRepository defines the persistence contract for Commission entities.
type CommissionRepository interface {
	Save(commission *Commission) error
	FindByID(id string) (*Commission, error)
	FindByProductID(productID string) (*Commission, error)
}
