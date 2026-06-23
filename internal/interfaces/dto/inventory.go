package dto

type CreateInventoryRequest struct {
	ProductID   int64 `json:"product_id"`
	WarehouseID int64 `json:"warehouse_id"`
	BasePrice   int64 `json:"base_price"`
	Stock       int   `json:"stock"`
}
type CreateInventoryResponse struct {
	InventoryID    int64 `json:"inventory_id"`
	AvailableStock int   `json:"available_stock"`
	FinalPrice     int64 `json:"final_price"`
}
type UpdatePriceRequest struct {
	InventoryID        int64 `json:"inventory_id"`
	DiscountPercentage int   `json:"discount_percentage"`
}
type UpdatePriceResponse struct {
	InventoryID int64 `json:"inventory_id"`
	OldPrice    int64 `json:"old_price"`
	NewPrice    int64 `json:"new_price"`
}
type HandleOrderPaidRequest struct {
	InventoryID int64 `json:"inventory_id"`
	Quantity    int   `json:"quantity"`
}
type HandleOrderPaidResponse struct {
	AvailableStock int `json:"available_stock"`
	ReservedStock  int `json:"reserved_stock"`
}
type HandleOrderDeliveredRequest struct {
	InventoryID int64 `json:"inventory_id"`
	Quantity    int   `json:"quantity"`
}
type HandleOrderDeliveredResponse struct {
	StockOut int `json:"stock_out"`
}
type ResetPriceRequest struct {
	InventoryID int64 `json:"inventory_id"`
}
type ResetPriceResponse struct {
	InventoryID int64 `json:"inventory_id"`
	FinalPrice  int64 `json:"final_price"`
}
type RecordReferencePriceRequest struct {
	ProductID int64  `json:"product_id"`
	Price     int64  `json:"price"`
	Source    string `json:"source"`
}
type RecordReferencePriceResponse struct {
	ProductID int64 `json:"product_id"`
	Price     int64 `json:"price"`
}
