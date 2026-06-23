package seller

import (
	"context"

	domainseller "stock-service-version-three/internal/domain/seller"
)

type VerifyKYCUseCase struct {
	sellerRepo domainseller.SellerRepository
}

func NewVerifyKYCUseCase(sellerRepo domainseller.SellerRepository) *VerifyKYCUseCase {
	return &VerifyKYCUseCase{
		sellerRepo: sellerRepo,
	}
}

func (uc *VerifyKYCUseCase) Execute(ctx context.Context, req VerifyKYCRequest) (*VerifyKYCResponse, error) {
	if req.KYCStatus != "approved" {
		return nil, domainseller.ErrInvalidKYCStatus
	}

	s, err := uc.sellerRepo.FindByID(req.SellerID)
	if err != nil {
		return nil, err
	}

	if err := s.VerifyKYC(); err != nil {
		return nil, err
	}

	if err := uc.sellerRepo.Update(s); err != nil {
		return nil, err
	}

	return &VerifyKYCResponse{
		SellerID: s.ID,
		Status:   string(s.Status),
	}, nil
}
