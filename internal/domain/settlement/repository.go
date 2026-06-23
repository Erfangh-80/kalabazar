package settlement

type SettlementRepository interface {
	Save(settlement *Settlement) error
	FindByID(id int64) (*Settlement, error)
	FindBySellerID(sellerID int64) ([]*Settlement, error)
}
