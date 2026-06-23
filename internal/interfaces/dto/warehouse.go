package dto

type CreateWarehouseRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}
type CreateWarehouseResponse struct {
	WarehouseID int64  `json:"warehouse_id"`
	Name        string `json:"name"`
}
type LinkWarehouseRequest struct {
	StoreID     int64  `json:"store_id"`
	WarehouseID int64  `json:"warehouse_id"`
	Type        string `json:"type"`
}
type LinkWarehouseResponse struct {
	StoreID     int64 `json:"store_id"`
	WarehouseID int64 `json:"warehouse_id"`
}
