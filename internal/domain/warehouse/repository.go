package warehouse

type WarehouseRepository interface {
	Save(warehouse *Warehouse) error
	FindByID(id int64) (*Warehouse, error)
	Update(warehouse *Warehouse) error
}
