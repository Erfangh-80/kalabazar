package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrInvalidCategoryReviewDecision = errors.New("review decision must be 'approve' or 'reject'")
	ErrCategoryAccessNotFound        = errors.New("category access request not found")
)

// ReviewCategoryAccessInput contains the data needed to review a category access request.
type ReviewCategoryAccessInput struct {
	StoreID     string
	CategoryID  string
	Decision    string
	SupportNote string
}

// ReviewCategoryAccessOutput contains the result of a category access review.
type ReviewCategoryAccessOutput struct {
	StoreID     string
	CategoryID  string
	Status      string
	SupportNote string
	Event       any
}

// ReviewCategoryAccessUseCase orchestrates the review of a category access request.
type ReviewCategoryAccessUseCase struct {
	repo entity.StoreCategoryRepository
}

// NewReviewCategoryAccessUseCase creates a new ReviewCategoryAccessUseCase.
func NewReviewCategoryAccessUseCase(repo entity.StoreCategoryRepository) *ReviewCategoryAccessUseCase {
	return &ReviewCategoryAccessUseCase{repo: repo}
}

// Execute reviews a pending category access request with the given decision.
func (uc *ReviewCategoryAccessUseCase) Execute(input ReviewCategoryAccessInput) (*ReviewCategoryAccessOutput, error) {
	sc, err := uc.repo.FindByStoreIDAndCategoryID(input.StoreID, input.CategoryID)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, ErrCategoryAccessNotFound
	}

	switch input.Decision {
	case "approve":
		if err := sc.Approve(); err != nil {
			return nil, err
		}
	case "reject":
		if err := sc.Reject(input.SupportNote); err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidCategoryReviewDecision
	}

	if err := uc.repo.Update(sc); err != nil {
		return nil, err
	}

	events := sc.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &ReviewCategoryAccessOutput{
		StoreID:     sc.StoreID,
		CategoryID:  sc.CategoryID,
		Status:      string(sc.Status),
		SupportNote: sc.SupportNote,
		Event:       domainEvent,
	}, nil
}
