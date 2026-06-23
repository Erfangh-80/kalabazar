package seller

type Event interface {
	EventName() string
}

type SellerVerifiedEvent struct {
	SellerID int64
}

func (e SellerVerifiedEvent) EventName() string {
	return "seller.verified"
}

type SellerRankUpdatedEvent struct {
	SellerID int64
	Score    float64
	Rank     string
}

func (e SellerRankUpdatedEvent) EventName() string {
	return "seller.rank.updated"
}
