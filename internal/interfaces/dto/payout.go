package dto

type ExecutePayoutRequest struct {
	SellerID int64 `json:"seller_id"`
	Amount   int64 `json:"amount"`
}
type ExecutePayoutResponse struct {
	PayoutID int64  `json:"payout_id"`
	Amount   int64  `json:"amount"`
	Status   string `json:"status"`
}
