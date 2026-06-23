package store

type ApproveStoreRequest struct {
	StoreID  int64
	Decision string
}

type ApproveStoreResponse struct {
	StoreID int64
	Status  string
	Events  []any
}

type AllowCategoryRequest struct {
	StoreID    int64
	CategoryID int64
}

type AllowCategoryResponse struct {
	StoreID    int64
	CategoryID int64
	Events     []any
}
