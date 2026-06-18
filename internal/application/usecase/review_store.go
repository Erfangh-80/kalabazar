package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrInvalidReviewDecision = errors.New("review decision must be 'approve' or 'reject'")
)

// ReviewStoreInput contains the data needed to review a pending store.
type ReviewStoreInput struct {
	StoreID  string
	Decision string
}

// ReviewStoreOutput contains the result of a store review.
type ReviewStoreOutput struct {
	StoreID string
	Status  string
	Event   any
}

// ReviewStoreUseCase orchestrates the review of a pending store by an admin.
type ReviewStoreUseCase struct {
	repo entity.StoreRepository
}

// NewReviewStoreUseCase creates a new ReviewStoreUseCase.
func NewReviewStoreUseCase(repo entity.StoreRepository) *ReviewStoreUseCase {
	return &ReviewStoreUseCase{repo: repo}
}

// Execute reviews a pending store with the given decision (approve/reject).
func (uc *ReviewStoreUseCase) Execute(input ReviewStoreInput) (*ReviewStoreOutput, error) {
	store, err := uc.repo.FindByID(input.StoreID)
	if err != nil {
		return nil, err
	}

	switch input.Decision {
	case "approve":
		if err := store.Approve(); err != nil {
			return nil, err
		}
	case "reject":
		if err := store.Reject(); err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidReviewDecision
	}

	if err := uc.repo.Update(store); err != nil {
		return nil, err
	}

	events := store.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &ReviewStoreOutput{
		StoreID: store.ID,
		Status:  string(store.Status),
		Event:   domainEvent,
	}, nil
}
