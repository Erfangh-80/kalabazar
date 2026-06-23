package dto

type CreateProductRequest struct {
	StoreID    int64  `json:"store_id"`
	Title      string `json:"title"`
	CategoryID int64  `json:"category_id"`
	Brand      string `json:"brand"`
}
type CreateProductResponse struct {
	ProductID int64  `json:"product_id"`
	Status    string `json:"status"`
}
type ApproveProductRequest struct {
	ProductID int64  `json:"product_id"`
	Decision  string `json:"decision"`
}
type ApproveProductResponse struct {
	ProductID int64  `json:"product_id"`
	Status    string `json:"status"`
}
