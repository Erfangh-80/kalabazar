package seller_test

import (
	"context"
	"testing"

	appSeller "stock-service-version-three/internal/application/seller"
	domainseller "stock-service-version-three/internal/domain/seller"
)

func TestVerifyKYC_Success(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	uc := appSeller.NewVerifyKYCUseCase(sellerRepo)

	seller := domainseller.NewSeller(1, "Test Store", "09123456789")
	sellerRepo.Save(seller)

	req := appSeller.VerifyKYCRequest{
		SellerID:  seller.ID,
		KYCStatus: "approved",
	}

	resp, err := uc.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.SellerID != seller.ID {
		t.Errorf("expected seller ID %d, got %d", seller.ID, resp.SellerID)
	}
	if resp.Status != string(domainseller.SellerStatusVerified) {
		t.Errorf("expected status VERIFIED, got '%s'", resp.Status)
	}
}

func TestVerifyKYC_NotApproved(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	uc := appSeller.NewVerifyKYCUseCase(sellerRepo)

	seller := domainseller.NewSeller(1, "Test Store", "09123456789")
	sellerRepo.Save(seller)

	req := appSeller.VerifyKYCRequest{
		SellerID:  seller.ID,
		KYCStatus: "rejected",
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainseller.ErrInvalidKYCStatus {
		t.Errorf("expected ErrInvalidKYCStatus, got %v", err)
	}
}

func TestVerifyKYC_SellerNotFound(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	uc := appSeller.NewVerifyKYCUseCase(sellerRepo)

	req := appSeller.VerifyKYCRequest{
		SellerID:  999,
		KYCStatus: "approved",
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainseller.ErrSellerNotFound {
		t.Errorf("expected ErrSellerNotFound, got %v", err)
	}
}

func TestVerifyKYC_AlreadyVerified(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	uc := appSeller.NewVerifyKYCUseCase(sellerRepo)

	seller := domainseller.NewSeller(1, "Test Store", "09123456789")
	sellerRepo.Save(seller)
	seller.VerifyKYC()
	sellerRepo.Update(seller)

	req := appSeller.VerifyKYCRequest{
		SellerID:  seller.ID,
		KYCStatus: "approved",
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainseller.ErrSellerAlreadyVerified {
		t.Errorf("expected ErrSellerAlreadyVerified, got %v", err)
	}
}
