package store

import (
	"context"

	domainstore "stock-service-version-three/internal/domain/store"
)

type ApproveStoreUseCase struct {
	storeRepo domainstore.StoreRepository
}

func NewApproveStoreUseCase(storeRepo domainstore.StoreRepository) *ApproveStoreUseCase {
	return &ApproveStoreUseCase{
		storeRepo: storeRepo,
	}
}

func (uc *ApproveStoreUseCase) Execute(ctx context.Context, req ApproveStoreRequest) (*ApproveStoreResponse, error) {
	s, err := uc.storeRepo.FindByID(ctx, req.StoreID)
	if err != nil {
		return nil, err
	}

	var events []any

	if req.Decision == "approved" {
		if err := s.Activate(); err != nil {
			return nil, err
		}
		if err := uc.storeRepo.Update(ctx, s); err != nil {
			return nil, err
		}
		events = s.Events()
	}

	return &ApproveStoreResponse{
		StoreID: s.ID,
		Status:  string(s.Status),
		Events:  events,
	}, nil
}
