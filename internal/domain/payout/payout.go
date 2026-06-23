package payout

import (
	"sync/atomic"
	"time"
)

type PayoutStatus string

const (
	PayoutStatusPending  PayoutStatus = "PENDING"
	PayoutStatusExecuted PayoutStatus = "EXECUTED"
	PayoutStatusFailed   PayoutStatus = "FAILED"
)

var payoutIDCounter int64

type Payout struct {
	id        int64
	sellerID  int64
	amount    int64
	status    PayoutStatus
	createdAt time.Time
}

func NewPayout(sellerID int64, amount int64) (*Payout, error) {
	if amount <= 0 {
		return nil, ErrInvalidPayoutAmount
	}

	id := atomic.AddInt64(&payoutIDCounter, 1)

	return &Payout{
		id:        id,
		sellerID:  sellerID,
		amount:    amount,
		status:    PayoutStatusPending,
		createdAt: time.Now(),
	}, nil
}

func (p *Payout) Execute() (*PayoutExecutedEvent, error) {
	if p.status == PayoutStatusExecuted {
		return nil, ErrPayoutAlreadyExecuted
	}

	p.status = PayoutStatusExecuted

	event := &PayoutExecutedEvent{
		PayoutID: p.id,
		SellerID: p.sellerID,
		Amount:   p.amount,
	}

	return event, nil
}

func (p *Payout) ID() int64              { return p.id }
func (p *Payout) SellerID() int64        { return p.sellerID }
func (p *Payout) Amount() int64          { return p.amount }
func (p *Payout) Status() PayoutStatus   { return p.status }
func (p *Payout) CreatedAt() time.Time   { return p.createdAt }
