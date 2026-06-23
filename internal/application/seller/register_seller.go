package seller

import (
	"context"

	domainseller "stock-service-version-three/internal/domain/seller"
	domainstore "stock-service-version-three/internal/domain/store"
)

type RegisterSellerUseCase struct {
	sellerRepo domainseller.SellerRepository
	storeRepo  domainstore.StoreRepository
}

func NewRegisterSellerUseCase(sellerRepo domainseller.SellerRepository, storeRepo domainstore.StoreRepository) *RegisterSellerUseCase {
	return &RegisterSellerUseCase{
		sellerRepo: sellerRepo,
		storeRepo:  storeRepo,
	}
}

func (uc *RegisterSellerUseCase) Execute(ctx context.Context, req RegisterSellerRequest) (*RegisterSellerResponse, error) {
	if err := domainseller.ValidateSellerName(req.StoreName); err != nil {
		return nil, err
	}
	if err := domainseller.ValidatePhone(req.Phone); err != nil {
		return nil, err
	}

	s := domainseller.NewSeller(req.UserID, req.StoreName, req.Phone)
	if err := uc.sellerRepo.Save(s); err != nil {
		return nil, err
	}

	st := domainstore.NewStore(s.ID, req.StoreName, req.Phone)
	if err := uc.storeRepo.Save(ctx, st); err != nil {
		return nil, err
	}

	events := st.Events()

	return &RegisterSellerResponse{
		SellerID: s.ID,
		StoreID:  st.ID,
		Events:   events,
	}, nil
}
