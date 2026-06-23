package seller

type RegisterSellerRequest struct {
	UserID    int64
	StoreName string
	Phone     string
}

type RegisterSellerResponse struct {
	SellerID int64
	StoreID  int64
	Events   []interface{}
}

type VerifyKYCRequest struct {
	SellerID  int64
	KYCStatus string
}

type VerifyKYCResponse struct {
	SellerID int64
	Status   string
}

type UpdateRankRequest struct {
	SellerID int64
	Score    float64
	Rank     string
}

type UpdateRankResponse struct {
	SellerID int64
	Score    float64
	Rank     string
}
