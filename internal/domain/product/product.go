package product

import "time"

type ProductStatus string

const (
	PENDING_REVIEW ProductStatus = "PENDING_REVIEW"
	ACTIVE         ProductStatus = "ACTIVE"
	REJECTED       ProductStatus = "REJECTED"
)

type Product struct {
	ID         int64
	StoreID    int64
	Title      string
	CategoryID int64
	Brand      string
	Status     ProductStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time

	events []any
}

func NewProduct(storeID int64, title string, categoryID int64, brand string) *Product {
	now := time.Now()
	p := &Product{
		StoreID:    storeID,
		Title:      title,
		CategoryID: categoryID,
		Brand:      brand,
		Status:     PENDING_REVIEW,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	p.emit(ProductCreatedEvent{
		ProductID: p.ID,
		StoreID:   p.StoreID,
		Title:     p.Title,
	})
	return p
}

func (p *Product) Approve() error {
	if p.Status == ACTIVE {
		return ErrProductAlreadyApproved
	}
	p.Status = ACTIVE
	p.UpdatedAt = time.Now()
	p.emit(ProductApprovedEvent{
		ProductID: p.ID,
	})
	return nil
}

func (p *Product) Reject() error {
	if p.Status == ACTIVE {
		return ErrProductAlreadyApproved
	}
	p.Status = REJECTED
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) emit(event any) {
	p.events = append(p.events, event)
}

func (p *Product) Events() []any {
	return p.events
}

func (p *Product) ClearEvents() {
	p.events = nil
}
