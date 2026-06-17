package event

import "time"

// ReferencePriceCreated is emitted when a new reference price is registered.
type ReferencePriceCreated struct {
	ReferencePriceID string
	ProductID        string
	Price            float64
	Source           string
	Timestamp        time.Time
}

func (e ReferencePriceCreated) EventName() string { return "pricing.reference_price_recorded" }
