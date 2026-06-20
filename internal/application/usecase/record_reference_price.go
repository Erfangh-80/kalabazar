package usecase

import (
	"kalabazar-stock-service/internal/domain/entity"
)

// RecordReferencePriceInput contains the data needed to record a reference price.
type RecordReferencePriceInput struct {
	ID        string
	ProductID string
	Price     float64
	Source    string
}

// RecordReferencePriceOutput contains the result of recording a reference price.
type RecordReferencePriceOutput struct {
	ID        string
	ProductID string
	Price     float64
	Source    string
	Event     any
}

// RecordReferencePriceUseCase orchestrates recording a reference price observation.
type RecordReferencePriceUseCase struct {
	repo entity.ReferencePriceRepository
}

// NewRecordReferencePriceUseCase creates a new RecordReferencePriceUseCase.
func NewRecordReferencePriceUseCase(repo entity.ReferencePriceRepository) *RecordReferencePriceUseCase {
	return &RecordReferencePriceUseCase{repo: repo}
}

// Execute records a new reference price observation.
func (uc *RecordReferencePriceUseCase) Execute(input RecordReferencePriceInput) (*RecordReferencePriceOutput, error) {
	rp, err := entity.NewReferencePrice(input.ID, input.ProductID, input.Price, input.Source)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(rp); err != nil {
		return nil, err
	}

	events := rp.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &RecordReferencePriceOutput{
		ID:        rp.ID,
		ProductID: rp.ProductID,
		Price:     rp.Price,
		Source:    rp.Source,
		Event:     domainEvent,
	}, nil
}
