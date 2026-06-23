package dto

type CalculateCommissionRequest struct {
	SellerID    int64   `json:"seller_id"`
	SalesAmount int64   `json:"sales_amount"`
	Rate        float64 `json:"rate"`
}
type CalculateCommissionResponse struct {
	CommissionID int64 `json:"commission_id"`
	Amount       int64 `json:"amount"`
	SalesAmount  int64 `json:"sales_amount"`
}
