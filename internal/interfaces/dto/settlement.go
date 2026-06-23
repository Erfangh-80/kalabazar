package dto

type CreateSettlementRequest struct {
	SellerID   int64 `json:"seller_id"`
	GrossSales int64 `json:"gross_sales"`
	Commission int64 `json:"commission"`
}
type CreateSettlementResponse struct {
	SettlementID int64 `json:"settlement_id"`
	GrossSales   int64 `json:"gross_sales"`
	Commission   int64 `json:"commission"`
	NetAmount    int64 `json:"net_amount"`
}
