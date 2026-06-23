package commission

type CalculateCommissionRequest struct {
	SellerID    int64
	SalesAmount int64
	Rate        float64
}

type CalculateCommissionResponse struct {
	CommissionID int64
	Amount       int64
	SalesAmount  int64
}
