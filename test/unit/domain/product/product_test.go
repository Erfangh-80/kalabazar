package product_test

import (
	"testing"
	"time"

	"stock-service-version-three/internal/domain/product"
)

func TestNewProduct(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")

	if p.StoreID != 1 {
		t.Errorf("expected StoreID 1, got %d", p.StoreID)
	}
	if p.Title != "Runner 3000" {
		t.Errorf("expected Title 'Runner 3000', got %s", p.Title)
	}
	if p.CategoryID != 12 {
		t.Errorf("expected CategoryID 12, got %d", p.CategoryID)
	}
	if p.Brand != "SportLine" {
		t.Errorf("expected Brand 'SportLine', got %s", p.Brand)
	}
	if p.Status != product.PENDING_REVIEW {
		t.Errorf("expected Status %s, got %s", product.PENDING_REVIEW, p.Status)
	}
	if p.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set, got zero")
	}
	if p.CreatedAt != p.UpdatedAt {
		t.Errorf("expected CreatedAt == UpdatedAt for new product")
	}
	if p.ID != 0 {
		t.Errorf("expected ID 0 for new product, got %d", p.ID)
	}
}

func TestNewProduct_EmitsCreatedEvent(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")

	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e, ok := events[0].(product.ProductCreatedEvent)
	if !ok {
		t.Fatalf("expected ProductCreatedEvent, got %T", events[0])
	}
	if e.ProductID != p.ID {
		t.Errorf("expected ProductID %d, got %d", p.ID, e.ProductID)
	}
	if e.StoreID != 1 {
		t.Errorf("expected StoreID 1, got %d", e.StoreID)
	}
	if e.Title != "Runner 3000" {
		t.Errorf("expected Title 'Runner 3000', got %s", e.Title)
	}
}

func TestProductApprove(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")
	p.ID = 10

	err := p.Approve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Status != product.ACTIVE {
		t.Errorf("expected Status %s, got %s", product.ACTIVE, p.Status)
	}
}

func TestProductApprove_EmitsApprovedEvent(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")
	p.ID = 10
	p.ClearEvents()

	err := p.Approve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := p.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	e, ok := events[0].(product.ProductApprovedEvent)
	if !ok {
		t.Fatalf("expected ProductApprovedEvent, got %T", events[0])
	}
	if e.ProductID != 10 {
		t.Errorf("expected ProductID 10, got %d", e.ProductID)
	}
}

func TestProductApprove_AlreadyApproved(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")
	p.ID = 10

	err := p.Approve()
	if err != nil {
		t.Fatalf("unexpected error on first approve: %v", err)
	}

	err = p.Approve()
	if err != product.ErrProductAlreadyApproved {
		t.Errorf("expected ErrProductAlreadyApproved, got %v", err)
	}
}

func TestProductReject(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")
	p.ID = 10

	err := p.Reject()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Status != product.REJECTED {
		t.Errorf("expected Status %s, got %s", product.REJECTED, p.Status)
	}
}

func TestProductReject_AlreadyApproved(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")
	p.ID = 10

	p.Approve()

	err := p.Reject()
	if err != product.ErrProductAlreadyApproved {
		t.Errorf("expected ErrProductAlreadyApproved, got %v", err)
	}
	if p.Status != product.ACTIVE {
		t.Errorf("expected Status to remain %s, got %s", product.ACTIVE, p.Status)
	}
}

func TestProductEvents_Cleared(t *testing.T) {
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")

	p.ClearEvents()
	if len(p.Events()) != 0 {
		t.Error("expected no events after ClearEvents")
	}
}

func TestProductTimestamps(t *testing.T) {
	now := time.Now()
	p := product.NewProduct(1, "Runner 3000", 12, "SportLine")

	if p.CreatedAt.Before(now) || p.CreatedAt.After(time.Now()) {
		t.Error("expected CreatedAt to be set to current time")
	}

	prev := p.UpdatedAt
	time.Sleep(time.Millisecond)
	p.Approve()

	if !p.UpdatedAt.After(prev) {
		t.Error("expected UpdatedAt to be updated on approve")
	}
}
