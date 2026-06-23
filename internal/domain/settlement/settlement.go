package settlement

import (
	"sync/atomic"
	"time"
)

var settlementIDCounter int64

type Settlement struct {
	id         int64
	sellerID   int64
	grossSales int64
	commission int64
	netAmount  int64
	createdAt  time.Time
}

func NewSettlement(sellerID int64, grossSales int64, commission int64) (*Settlement, *SettlementCreatedEvent, error) {
	if grossSales < 0 || commission < 0 {
		return nil, nil, ErrInvalidSettlementAmount
	}
	if commission > grossSales {
		return nil, nil, ErrInvalidSettlementAmount
	}

	netAmount := grossSales - commission
	id := atomic.AddInt64(&settlementIDCounter, 1)

	s := &Settlement{
		id:         id,
		sellerID:   sellerID,
		grossSales: grossSales,
		commission: commission,
		netAmount:  netAmount,
		createdAt:  time.Now(),
	}

	event := &SettlementCreatedEvent{
		SettlementID: s.id,
		SellerID:     s.sellerID,
		GrossSales:   s.grossSales,
		Commission:   s.commission,
		NetAmount:    s.netAmount,
	}

	return s, event, nil
}

func (s *Settlement) ID() int64            { return s.id }
func (s *Settlement) SellerID() int64      { return s.sellerID }
func (s *Settlement) GrossSales() int64    { return s.grossSales }
func (s *Settlement) Commission() int64    { return s.commission }
func (s *Settlement) NetAmount() int64     { return s.netAmount }
func (s *Settlement) CreatedAt() time.Time { return s.createdAt }
