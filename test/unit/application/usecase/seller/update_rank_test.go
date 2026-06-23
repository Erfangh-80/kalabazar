package seller_test

import (
	"testing"

	appSeller "stock-service-version-three/internal/application/seller"
	domainseller "stock-service-version-three/internal/domain/seller"
)

func TestUpdateRankUseCase_Success(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	uc := appSeller.NewUpdateRankUseCase(sellerRepo)

	seller := domainseller.NewSeller(1, "Test Store", "09123456789")
	sellerRepo.Save(seller)

	req := appSeller.UpdateRankRequest{
		SellerID: seller.ID,
		Score:    95.5,
		Rank:     "GOLD",
	}

	resp, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.SellerID != seller.ID {
		t.Errorf("expected SellerID %d, got %d", seller.ID, resp.SellerID)
	}
	if resp.Score != req.Score {
		t.Errorf("expected Score %f, got %f", req.Score, resp.Score)
	}
	if resp.Rank != req.Rank {
		t.Errorf("expected Rank %s, got %s", req.Rank, resp.Rank)
	}

	updated, _ := sellerRepo.FindByID(seller.ID)
	if updated.Score != req.Score {
		t.Errorf("expected stored Score %f, got %f", req.Score, updated.Score)
	}
	if updated.Rank != req.Rank {
		t.Errorf("expected stored Rank %s, got %s", req.Rank, updated.Rank)
	}
}

func TestUpdateRankUseCase_NotFound(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	uc := appSeller.NewUpdateRankUseCase(sellerRepo)

	req := appSeller.UpdateRankRequest{
		SellerID: 999,
		Score:    50.0,
		Rank:     "SILVER",
	}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatal("expected error for non-existent seller, got nil")
	}
	if err != domainseller.ErrSellerNotFound {
		t.Errorf("expected ErrSellerNotFound, got %v", err)
	}
}

func TestUpdateRankUseCase_OverwritePreviousRank(t *testing.T) {
	sellerRepo := newMockSellerRepo()
	uc := appSeller.NewUpdateRankUseCase(sellerRepo)

	seller := domainseller.NewSeller(1, "Test Store", "09123456789")
	sellerRepo.Save(seller)

	firstReq := appSeller.UpdateRankRequest{
		SellerID: seller.ID,
		Score:    10.0,
		Rank:     "BRONZE",
	}

	_, err := uc.Execute(firstReq)
	if err != nil {
		t.Fatalf("expected no error on first update, got %v", err)
	}

	secondReq := appSeller.UpdateRankRequest{
		SellerID: seller.ID,
		Score:    90.0,
		Rank:     "PLATINUM",
	}

	resp, err := uc.Execute(secondReq)
	if err != nil {
		t.Fatalf("expected no error on second update, got %v", err)
	}

	if resp.Score != 90.0 {
		t.Errorf("expected Score 90.0, got %f", resp.Score)
	}
	if resp.Rank != "PLATINUM" {
		t.Errorf("expected Rank PLATINUM, got %s", resp.Rank)
	}
}
