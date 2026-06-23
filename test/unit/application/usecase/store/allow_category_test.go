package store_test

import (
	"context"
	"testing"

	appStore "stock-service-version-three/internal/application/store"
	domainstore "stock-service-version-three/internal/domain/store"
)

func TestAllowCategory_Success(t *testing.T) {
	storeRepo := newMockStoreRepo()
	uc := appStore.NewAllowCategoryUseCase(storeRepo)

	store := domainstore.NewStore(1, "Test Store", "09123456789")
	storeRepo.Save(context.Background(), store)

	req := appStore.AllowCategoryRequest{
		StoreID:    store.ID,
		CategoryID: 42,
	}

	resp, err := uc.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StoreID != store.ID {
		t.Errorf("expected store ID %d, got %d", store.ID, resp.StoreID)
	}
	if resp.CategoryID != 42 {
		t.Errorf("expected category ID 42, got %d", resp.CategoryID)
	}
	if len(resp.Events) == 0 {
		t.Error("expected at least one event (StoreCategoryAllowedEvent)")
	}
}

func TestAllowCategory_StoreNotFound(t *testing.T) {
	storeRepo := newMockStoreRepo()
	uc := appStore.NewAllowCategoryUseCase(storeRepo)

	req := appStore.AllowCategoryRequest{
		StoreID:    999,
		CategoryID: 42,
	}

	_, err := uc.Execute(context.Background(), req)
	if err != domainstore.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}
