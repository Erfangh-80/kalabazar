package store_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/domain/store"
)

func TestNewStore(t *testing.T) {
	t.Parallel()

	s := store.NewStore(1, "SportLine", "09120001111")

	if s.ID != 0 {
		t.Errorf("expected ID 0, got %d", s.ID)
	}
	if s.SellerID != 1 {
		t.Errorf("expected SellerID 1, got %d", s.SellerID)
	}
	if s.Name != "SportLine" {
		t.Errorf("expected Name 'SportLine', got %s", s.Name)
	}
	if s.Phone != "09120001111" {
		t.Errorf("expected Phone '09120001111', got %s", s.Phone)
	}
	if s.Status != store.StoreStatusPENDING {
		t.Errorf("expected Status PENDING, got %s", s.Status)
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if s.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
	if !s.UpdatedAt.Equal(s.CreatedAt) {
		t.Error("expected UpdatedAt to equal CreatedAt for new store")
	}
}

func TestNewStore_EmitsCreatedEvent(t *testing.T) {
	t.Parallel()

	s := store.NewStore(1, "SportLine", "09120001111")
	events := s.Events()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e, ok := events[0].(store.StoreCreatedEvent)
	if !ok {
		t.Fatalf("expected StoreCreatedEvent, got %T", events[0])
	}

	if e.StoreID != 0 {
		t.Errorf("expected event StoreID 0, got %d", e.StoreID)
	}
	if e.SellerID != 1 {
		t.Errorf("expected event SellerID 1, got %d", e.SellerID)
	}
	if e.Name != "SportLine" {
		t.Errorf("expected event Name 'SportLine', got %s", e.Name)
	}
}

func TestActivate_PendingStore(t *testing.T) {
	t.Parallel()

	s := store.NewStore(1, "SportLine", "09120001111")
	_ = s.Events()

	before := time.Now()
	err := s.Activate()
	after := time.Now()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Status != store.StoreStatusACTIVE {
		t.Errorf("expected Status ACTIVE, got %s", s.Status)
	}
	if s.UpdatedAt.Before(before) || s.UpdatedAt.After(after) {
		t.Error("expected UpdatedAt to be updated after activation")
	}

	events := s.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	_, ok := events[0].(store.StoreActivatedEvent)
	if !ok {
		t.Fatalf("expected StoreActivatedEvent, got %T", events[0])
	}
}

func TestActivate_AlreadyActiveStore_ReturnsError(t *testing.T) {
	t.Parallel()

	s := store.NewStore(1, "SportLine", "09120001111")
	_ = s.Events()

	_ = s.Activate()
	_ = s.Events()

	err := s.Activate()
	if err != store.ErrStoreAlreadyActive {
		t.Errorf("expected ErrStoreAlreadyActive, got %v", err)
	}
}

func TestActivate_SuspendedStore_ReturnsError(t *testing.T) {
	t.Parallel()

	s := store.NewStore(1, "SportLine", "09120001111")
	_ = s.Events()

	s.Suspend()
	_ = s.Events()

	err := s.Activate()
	if err != store.ErrStoreSuspended {
		t.Errorf("expected ErrStoreSuspended, got %v", err)
	}
}

func TestSuspend(t *testing.T) {
	t.Parallel()

	s := store.NewStore(1, "SportLine", "09120001111")

	before := time.Now()
	s.Suspend()
	after := time.Now()

	if s.Status != store.StoreStatusSUSPENDED {
		t.Errorf("expected Status SUSPENDED, got %s", s.Status)
	}
	if s.UpdatedAt.Before(before) || s.UpdatedAt.After(after) {
		t.Error("expected UpdatedAt to be updated after suspend")
	}
}

func TestNewStoreAllowedCategory(t *testing.T) {
	t.Parallel()

	c := store.NewStoreAllowedCategory(1, 42)

	if c.StoreID != 1 {
		t.Errorf("expected StoreID 1, got %d", c.StoreID)
	}
	if c.CategoryID != 42 {
		t.Errorf("expected CategoryID 42, got %d", c.CategoryID)
	}
	if c.Status != store.AllowedCategoryStatusAPPROVED {
		t.Errorf("expected Status APPROVED, got %s", c.Status)
	}
	if c.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewStoreAllowedCategory_EmitsCategoryAllowedEvent(t *testing.T) {
	t.Parallel()

	c := store.NewStoreAllowedCategory(1, 42)
	events := c.Events()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e, ok := events[0].(store.StoreCategoryAllowedEvent)
	if !ok {
		t.Fatalf("expected StoreCategoryAllowedEvent, got %T", events[0])
	}
	if e.StoreID != 1 {
		t.Errorf("expected event StoreID 1, got %d", e.StoreID)
	}
	if e.CategoryID != 42 {
		t.Errorf("expected event CategoryID 42, got %d", e.CategoryID)
	}
}

func TestApproveCategory(t *testing.T) {
	t.Parallel()

	c := store.NewStoreAllowedCategory(1, 42)
	if c.Status != store.AllowedCategoryStatusAPPROVED {
		t.Fatalf("expected APPROVED after creation, got %s", c.Status)
	}
}

func TestEvents_ClearedAfterRead(t *testing.T) {
	t.Parallel()

	s := store.NewStore(1, "SportLine", "09120001111")

	first := s.Events()
	if len(first) != 1 {
		t.Fatalf("expected 1 event in first read, got %d", len(first))
	}

	second := s.Events()
	if len(second) != 0 {
		t.Errorf("expected 0 events after clearing, got %d", len(second))
	}
}
