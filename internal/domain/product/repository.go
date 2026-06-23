package product

type ProductRepository interface {
	Save(product *Product) error
	FindByID(id int64) (*Product, error)
	FindByStoreID(storeID int64) ([]*Product, error)
	Update(product *Product) error
}
