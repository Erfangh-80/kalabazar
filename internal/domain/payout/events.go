package payout

type PayoutExecutedEvent struct {
	PayoutID int64
	SellerID int64
	Amount   int64
}
