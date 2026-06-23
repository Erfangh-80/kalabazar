package commission

type CommissionRepository interface {
	Save(commission *Commission) error
	FindByID(id int64) (*Commission, error)
	FindBySellerID(sellerID int64) ([]*Commission, error)
}
