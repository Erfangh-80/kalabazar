package product

type ProductCreatedEvent struct {
	ProductID int64
	StoreID   int64
	Title     string
}

type ProductApprovedEvent struct {
	ProductID int64
}
