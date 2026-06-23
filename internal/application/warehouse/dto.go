package warehouse

type CreateWarehouseRequest struct {
	Name     string
	Capacity int
}

type CreateWarehouseResponse struct {
	WarehouseID int64
	Name        string
}

type LinkWarehouseRequest struct {
	StoreID     int64
	WarehouseID int64
	Type        string
}

type LinkWarehouseResponse struct {
	StoreID     int64
	WarehouseID int64
}
