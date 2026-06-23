package payout

import (
	domainpayout "stock-service-version-three/internal/domain/payout"
)

type ExecutePayoutRequest struct {
	SellerID int64
	Amount   int64
}

type ExecutePayoutResponse struct {
	PayoutID int64
	Amount   int64
	Status   domainpayout.PayoutStatus
}
