package usecase

import (
	"time"

	"kalabazar-stock-service/internal/domain/entity"
)

// RegisterStoreInput contains the data needed to register a new store.
type RegisterStoreInput struct {
	ID        string
	UserID    string
	StoreName string
}

// RegisterStoreOutput contains the result of a store registration.
type RegisterStoreOutput struct {
	ID                     string
	UserID                 string
	StoreName              string
	IsCommissionApplicable bool
	Event                  any
	CreatedAt              time.Time
}

// RegisterStoreUseCase orchestrates the registration of a new store.
type RegisterStoreUseCase struct {
	repo entity.StoreRepository
}

// NewRegisterStoreUseCase creates a new RegisterStoreUseCase.
func NewRegisterStoreUseCase(repo entity.StoreRepository) *RegisterStoreUseCase {
	return &RegisterStoreUseCase{repo: repo}
}

// Execute registers a new store with the given input.
func (uc *RegisterStoreUseCase) Execute(input RegisterStoreInput) (*RegisterStoreOutput, error) {
	store, err := entity.NewStore(input.ID, input.UserID, input.StoreName, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(store); err != nil {
		return nil, err
	}

	events := store.Events()
	var domainEvent any
	if len(events) > 0 {
		domainEvent = events[0]
	}

	return &RegisterStoreOutput{
		ID:                     store.ID,
		UserID:                 store.UserID,
		StoreName:              store.StoreName,
		IsCommissionApplicable: store.IsCommissionApplicable,
		Event:                  domainEvent,
		CreatedAt:              store.CreatedAt,
	}, nil
}
