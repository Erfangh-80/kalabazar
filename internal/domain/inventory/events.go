package inventory

type InventoryCreatedEvent struct {
	InventoryID    int64
	ProductID      int64
	AvailableStock int
	FinalPrice     int64
}

type StockInEvent struct {
	InventoryID int64
	Quantity    int
}

type PriceUpdatedEvent struct {
	InventoryID int64
	OldPrice    int64
	NewPrice    int64
}

type ReservedEvent struct {
	InventoryID int64
	Quantity    int
}

type StockOutEvent struct {
	InventoryID int64
	Quantity    int
}
