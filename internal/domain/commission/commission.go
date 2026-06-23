package commission

import (
	"sync/atomic"
	"time"
)

var commissionIDCounter int64

type Commission struct {
	id          int64
	sellerID    int64
	rate        float64
	salesAmount int64
	amount      int64
	createdAt   time.Time
}

func NewCommission(sellerID int64, rate float64, salesAmount int64) (*Commission, *CommissionCalculatedEvent, error) {
	if rate <= 0 {
		return nil, nil, ErrInvalidCommissionRate
	}

	amount := int64(float64(salesAmount) * rate)
	id := atomic.AddInt64(&commissionIDCounter, 1)

	c := &Commission{
		id:          id,
		sellerID:    sellerID,
		rate:        rate,
		salesAmount: salesAmount,
		amount:      amount,
		createdAt:   time.Now(),
	}

	event := &CommissionCalculatedEvent{
		CommissionID: c.id,
		SellerID:     c.sellerID,
		Amount:       c.amount,
		SalesAmount:  c.salesAmount,
	}

	return c, event, nil
}

func (c *Commission) ID() int64              { return c.id }
func (c *Commission) SellerID() int64        { return c.sellerID }
func (c *Commission) Rate() float64          { return c.rate }
func (c *Commission) SalesAmount() int64     { return c.salesAmount }
func (c *Commission) Amount() int64          { return c.amount }
func (c *Commission) CreatedAt() time.Time   { return c.createdAt }
