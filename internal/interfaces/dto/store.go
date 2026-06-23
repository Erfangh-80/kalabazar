package dto

type ApproveStoreRequest struct {
	StoreID  int64  `json:"store_id"`
	Decision string `json:"decision"`
}
type ApproveStoreResponse struct {
	StoreID int64  `json:"store_id"`
	Status  string `json:"status"`
}
type AllowCategoryRequest struct {
	StoreID    int64 `json:"store_id"`
	CategoryID int64 `json:"category_id"`
}
type AllowCategoryResponse struct {
	StoreID    int64 `json:"store_id"`
	CategoryID int64 `json:"category_id"`
}
