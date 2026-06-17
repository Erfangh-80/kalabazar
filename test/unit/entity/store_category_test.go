package entity_test

import (
	"testing"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestNewStoreCategory_Success(t *testing.T) {
	sc, err := entity.NewStoreCategory("store-1", "cat-7")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sc.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", sc.StoreID)
	}
	if sc.CategoryID != "cat-7" {
		t.Errorf("expected cat-7, got %s", sc.CategoryID)
	}
	if sc.Status != entity.StoreCategoryStatusPending {
		t.Errorf("expected pending, got %s", sc.Status)
	}
}

func TestNewStoreCategory_InvalidStoreID(t *testing.T) {
	_, err := entity.NewStoreCategory("", "cat-7")
	if err != entity.ErrStoreCategoryInvalidStoreID {
		t.Errorf("expected ErrStoreCategoryInvalidStoreID, got %v", err)
	}
}

func TestNewStoreCategory_InvalidCategoryID(t *testing.T) {
	_, err := entity.NewStoreCategory("store-1", "")
	if err != entity.ErrStoreCategoryInvalidCategoryID {
		t.Errorf("expected ErrStoreCategoryInvalidCategoryID, got %v", err)
	}
}

func TestNewStoreCategory_EventEmitted(t *testing.T) {
	sc, err := entity.NewStoreCategory("store-1", "cat-7")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := sc.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreCategoryAllowed)
	if !ok {
		t.Fatalf("expected StoreCategoryAllowed, got %T", events[0])
	}
	if e.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", e.StoreID)
	}
}

func TestStoreCategory_Approve_Success(t *testing.T) {
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	sc.Events()

	err := sc.Approve()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sc.Status != entity.StoreCategoryStatusApproved {
		t.Errorf("expected approved, got %s", sc.Status)
	}
	events := sc.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event after approve, got %d", len(events))
	}
}

func TestStoreCategory_Approve_AlreadyApproved(t *testing.T) {
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	sc.Approve()
	sc.Events()

	err := sc.Approve()
	if err != entity.ErrStoreCategoryAlreadyApproved {
		t.Errorf("expected ErrStoreCategoryAlreadyApproved, got %v", err)
	}
}

func TestStoreCategory_Events_ClearedAfterCall(t *testing.T) {
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	events1 := sc.Events()
	if len(events1) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events1))
	}
	events2 := sc.Events()
	if len(events2) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(events2))
	}
}
