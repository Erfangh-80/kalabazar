package dto

type RegisterSellerRequest struct {
	UserID    int64  `json:"user_id"`
	StoreName string `json:"store_name"`
	Phone     string `json:"phone"`
}
type RegisterSellerResponse struct {
	SellerID int64  `json:"seller_id"`
	StoreID  int64  `json:"store_id"`
	Message  string `json:"message,omitempty"`
}
type VerifyKYCRequest struct {
	SellerID  int64  `json:"seller_id"`
	KYCStatus string `json:"kyc_status"`
}
type VerifyKYCResponse struct {
	SellerID int64  `json:"seller_id"`
	Status   string `json:"status"`
}
type UpdateRankRequest struct {
	SellerID int64   `json:"seller_id"`
	Score    float64 `json:"score"`
	Rank     string  `json:"rank"`
}
type UpdateRankResponse struct {
	SellerID int64   `json:"seller_id"`
	Score    float64 `json:"score"`
	Rank     string  `json:"rank"`
}
