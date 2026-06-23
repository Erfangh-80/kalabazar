package inventory

type CreateInventoryRequest struct {
	ProductID   int64
	WarehouseID int64
	BasePrice   int64
	Stock       int
}

type CreateInventoryResponse struct {
	InventoryID   int64
	AvailableStock int
	FinalPrice    int64
}

type UpdatePriceRequest struct {
	InventoryID       int64
	DiscountPercentage int
}

type UpdatePriceResponse struct {
	InventoryID int64
	OldPrice    int64
	NewPrice    int64
}

type HandleOrderPaidRequest struct {
	InventoryID int64
	Quantity    int
}

type HandleOrderPaidResponse struct {
	AvailableStock int
	ReservedStock  int
}

type HandleOrderDeliveredRequest struct {
	InventoryID int64
	Quantity    int
}

type HandleOrderDeliveredResponse struct {
	StockOut int
}

type ResetPriceRequest struct {
	InventoryID int64
}

type ResetPriceResponse struct {
	InventoryID int64
	FinalPrice  int64
}

type RecordReferencePriceRequest struct {
	ProductID int64
	Price     int64
	Source    string
}

type RecordReferencePriceResponse struct {
	ProductID int64
	Price     int64
}
