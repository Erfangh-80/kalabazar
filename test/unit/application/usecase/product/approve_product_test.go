package product_test

import (
	"errors"
	"testing"

	appProduct "stock-service-version-three/internal/application/product"
	domain "stock-service-version-three/internal/domain/product"
)

func TestApproveProduct_Approved(t *testing.T) {
	repo := &mockProductRepository{products: make(map[int64]*domain.Product)}
	uc := appProduct.NewCreateProductUseCase(repo)
	auc := appProduct.NewApproveProductUseCase(repo)

	createResp, err := uc.Execute(appProduct.CreateProductRequest{
		StoreID: 1, Title: "Runner 3000", CategoryID: 12, Brand: "SportLine",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := auc.Execute(appProduct.ApproveProductRequest{
		ProductID: createResp.ProductID,
		Decision:  "approved",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ProductID != createResp.ProductID {
		t.Errorf("expected ProductID %d, got %d", createResp.ProductID, resp.ProductID)
	}
	if resp.Status != string(domain.ACTIVE) {
		t.Errorf("expected Status %s, got %s", domain.ACTIVE, resp.Status)
	}
}

func TestApproveProduct_Rejected(t *testing.T) {
	repo := &mockProductRepository{products: make(map[int64]*domain.Product)}
	uc := appProduct.NewCreateProductUseCase(repo)
	auc := appProduct.NewApproveProductUseCase(repo)

	createResp, err := uc.Execute(appProduct.CreateProductRequest{
		StoreID: 1, Title: "Runner 3000", CategoryID: 12, Brand: "SportLine",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := auc.Execute(appProduct.ApproveProductRequest{
		ProductID: createResp.ProductID,
		Decision:  "rejected",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ProductID != createResp.ProductID {
		t.Errorf("expected ProductID %d, got %d", createResp.ProductID, resp.ProductID)
	}
	if resp.Status != string(domain.REJECTED) {
		t.Errorf("expected Status %s, got %s", domain.REJECTED, resp.Status)
	}
}

func TestApproveProduct_NotFound(t *testing.T) {
	repo := &mockProductRepository{products: make(map[int64]*domain.Product)}
	auc := appProduct.NewApproveProductUseCase(repo)

	_, err := auc.Execute(appProduct.ApproveProductRequest{
		ProductID: 999,
		Decision:  "approved",
	})
	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestApproveProduct_InvalidDecision(t *testing.T) {
	repo := &mockProductRepository{products: make(map[int64]*domain.Product)}
	uc := appProduct.NewCreateProductUseCase(repo)
	auc := appProduct.NewApproveProductUseCase(repo)

	createResp, err := uc.Execute(appProduct.CreateProductRequest{
		StoreID: 1, Title: "Runner 3000", CategoryID: 12, Brand: "SportLine",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = auc.Execute(appProduct.ApproveProductRequest{
		ProductID: createResp.ProductID,
		Decision:  "invalid",
	})
	if err == nil {
		t.Fatal("expected error for invalid decision, got nil")
	}
}

func TestApproveProduct_AlreadyApproved(t *testing.T) {
	repo := &mockProductRepository{products: make(map[int64]*domain.Product)}
	uc := appProduct.NewCreateProductUseCase(repo)
	auc := appProduct.NewApproveProductUseCase(repo)

	createResp, err := uc.Execute(appProduct.CreateProductRequest{
		StoreID: 1, Title: "Runner 3000", CategoryID: 12, Brand: "SportLine",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = auc.Execute(appProduct.ApproveProductRequest{
		ProductID: createResp.ProductID,
		Decision:  "approved",
	})
	if err != nil {
		t.Fatalf("unexpected error on first approve: %v", err)
	}

	_, err = auc.Execute(appProduct.ApproveProductRequest{
		ProductID: createResp.ProductID,
		Decision:  "approved",
	})
	if !errors.Is(err, domain.ErrProductAlreadyApproved) {
		t.Errorf("expected ErrProductAlreadyApproved, got %v", err)
	}
}
