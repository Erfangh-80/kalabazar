package product

type CreateProductRequest struct {
	StoreID    int64
	Title      string
	CategoryID int64
	Brand      string
}

type CreateProductResponse struct {
	ProductID int64
	Status    string
}

type ApproveProductRequest struct {
	ProductID int64
	Decision  string
}

type ApproveProductResponse struct {
	ProductID int64
	Status    string
}
