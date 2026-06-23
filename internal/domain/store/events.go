package store

type StoreCreatedEvent struct {
	StoreID  int64
	SellerID int64
	Name     string
}

type StoreActivatedEvent struct {
	StoreID int64
}

type StoreCategoryAllowedEvent struct {
	StoreID    int64
	CategoryID int64
}
