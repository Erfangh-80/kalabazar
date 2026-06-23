package settlement

type CreateSettlementRequest struct {
	SellerID   int64
	GrossSales int64
	Commission int64
}

type CreateSettlementResponse struct {
	SettlementID int64
	GrossSales   int64
	Commission   int64
	NetAmount    int64
}
