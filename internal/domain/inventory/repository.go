package inventory

type InventoryRepository interface {
	Save(inventory *Inventory) error
	FindByID(id int64) (*Inventory, error)
	FindByProductID(productID int64) ([]*Inventory, error)
	FindByWarehouseID(warehouseID int64) ([]*Inventory, error)
	Update(inventory *Inventory) error
}
