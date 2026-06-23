package warehouse

type WarehouseCreatedEvent struct {
	WarehouseID int64
	Name        string
}

type WarehouseLinkedToStoreEvent struct {
	WarehouseID int64
	StoreID     int64
	LinkType    string
}
