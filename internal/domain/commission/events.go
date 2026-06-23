package commission

type CommissionCalculatedEvent struct {
	CommissionID int64
	SellerID     int64
	Amount       int64
	SalesAmount  int64
}
