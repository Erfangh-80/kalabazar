package usecase

import (
	"errors"

	"kalabazar-stock-service/internal/domain/entity"
)

var (
	ErrCategoryAccessAlreadyRequested = errors.New("category access already requested for this store and category")
)

// RequestCategoryAccessInput contains the data needed to request category access.
type RequestCategoryAccessInput struct {
	StoreID    string
	CategoryID string
}

// RequestCategoryAccessOutput contains the result of a category access request.
type RequestCategoryAccessOutput struct {
	StoreID    string
	CategoryID string
	Status     string
	Event      any
}

// RequestCategoryAccessUseCase orchestrates a request for category access.
type RequestCategoryAccessUseCase struct {
	repo entity.StoreCategoryRepository
}

// NewRequestCategoryAccessUseCase creates a new RequestCategoryAccessUseCase.
func NewRequestCategoryAccessUseCase(repo entity.StoreCategoryRepository) *RequestCategoryAccessUseCase {
	return &RequestCategoryAccessUseCase{repo: repo}
}

// Execute requests access to a category for a store.
func (uc *RequestCategoryAccessUseCase) Execute(input RequestCategoryAccessInput) (*RequestCategoryAccessOutput, error) {
	existing, err := uc.repo.FindByStoreIDAndCategoryID(input.StoreID, input.CategoryID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCategoryAccessAlreadyRequested
	}

	sc, err := entity.NewStoreCategory(input.StoreID, input.CategoryID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(sc); err != nil {
		return nil, err
	}

	events := sc.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &RequestCategoryAccessOutput{
		StoreID:    sc.StoreID,
		CategoryID: sc.CategoryID,
		Status:     string(sc.Status),
		Event:      domainEvent,
	}, nil
}
