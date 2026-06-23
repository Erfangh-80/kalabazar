package payout

type PayoutRepository interface {
	Save(payout *Payout) error
	FindByID(id int64) (*Payout, error)
	FindBySellerID(sellerID int64) ([]*Payout, error)
	Update(payout *Payout) error
}
