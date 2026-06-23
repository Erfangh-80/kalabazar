package store

import (
	"context"

	domainstore "stock-service-version-three/internal/domain/store"
)

type AllowCategoryUseCase struct {
	storeRepo domainstore.StoreRepository
}

func NewAllowCategoryUseCase(storeRepo domainstore.StoreRepository) *AllowCategoryUseCase {
	return &AllowCategoryUseCase{
		storeRepo: storeRepo,
	}
}

func (uc *AllowCategoryUseCase) Execute(ctx context.Context, req AllowCategoryRequest) (*AllowCategoryResponse, error) {
	if _, err := uc.storeRepo.FindByID(ctx, req.StoreID); err != nil {
		return nil, err
	}

	cat := domainstore.NewStoreAllowedCategory(req.StoreID, req.CategoryID)
	events := cat.Events()

	return &AllowCategoryResponse{
		StoreID:    req.StoreID,
		CategoryID: req.CategoryID,
		Events:     events,
	}, nil
}
