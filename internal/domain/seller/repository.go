package seller

type SellerRepository interface {
	Save(seller *Seller) error
	FindByID(id int64) (*Seller, error)
	FindByUserID(userID int64) (*Seller, error)
	Update(seller *Seller) error
}
