package settlement

type SettlementCreatedEvent struct {
	SettlementID int64
	SellerID     int64
	GrossSales   int64
	Commission   int64
	NetAmount    int64
}
