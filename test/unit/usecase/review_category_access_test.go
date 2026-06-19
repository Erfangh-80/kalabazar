package usecase_test

import (
	"testing"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockStoreCategoryRepoForReview struct {
	records []*entity.StoreCategory
}

func (m *mockStoreCategoryRepoForReview) Save(sc *entity.StoreCategory) error {
	m.records = append(m.records, sc)
	return nil
}

func (m *mockStoreCategoryRepoForReview) FindByStoreIDAndCategoryID(storeID, categoryID string) (*entity.StoreCategory, error) {
	for _, sc := range m.records {
		if sc.StoreID == storeID && sc.CategoryID == categoryID {
			return sc, nil
		}
	}
	return nil, nil
}

func (m *mockStoreCategoryRepoForReview) Update(sc *entity.StoreCategory) error {
	for i, r := range m.records {
		if r.StoreID == sc.StoreID && r.CategoryID == sc.CategoryID {
			m.records[i] = sc
			return nil
		}
	}
	return nil
}

func newPendingCategoryAccess() *entity.StoreCategory {
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	sc.Events()
	return sc
}

func TestReviewCategoryAccess_Approve(t *testing.T) {
	sc := newPendingCategoryAccess()
	repo := &mockStoreCategoryRepoForReview{records: []*entity.StoreCategory{sc}}
	uc := usecase.NewReviewCategoryAccessUseCase(repo)

	input := usecase.ReviewCategoryAccessInput{
		StoreID:    "store-1",
		CategoryID: "cat-7",
		Decision:   "approve",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", output.StoreID)
	}
	if output.CategoryID != "cat-7" {
		t.Errorf("expected cat-7, got %s", output.CategoryID)
	}
	if output.Status != string(entity.StoreCategoryStatusApproved) {
		t.Errorf("expected approved, got %s", output.Status)
	}
	if output.Event == nil {
		t.Fatal("expected event, got nil")
	}
	e, ok := output.Event.(event.StoreCategoryAllowed)
	if !ok {
		t.Fatalf("expected StoreCategoryAllowed, got %T", output.Event)
	}
	if e.Status != string(entity.StoreCategoryStatusApproved) {
		t.Errorf("expected approved status in event, got %s", e.Status)
	}
}

func TestReviewCategoryAccess_Reject(t *testing.T) {
	sc := newPendingCategoryAccess()
	repo := &mockStoreCategoryRepoForReview{records: []*entity.StoreCategory{sc}}
	uc := usecase.NewReviewCategoryAccessUseCase(repo)

	input := usecase.ReviewCategoryAccessInput{
		StoreID:     "store-1",
		CategoryID:  "cat-7",
		Decision:    "reject",
		SupportNote: "Complete your documents",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Status != string(entity.StoreCategoryStatusRejected) {
		t.Errorf("expected rejected, got %s", output.Status)
	}
	if output.SupportNote != "Complete your documents" {
		t.Errorf("expected 'Complete your documents', got %s", output.SupportNote)
	}
	if output.Event == nil {
		t.Fatal("expected event, got nil")
	}
	e, ok := output.Event.(event.StoreCategoryRejected)
	if !ok {
		t.Fatalf("expected StoreCategoryRejected, got %T", output.Event)
	}
	if e.SupportNote != "Complete your documents" {
		t.Errorf("expected support note in event, got %s", e.SupportNote)
	}
}

func TestReviewCategoryAccess_RejectWithoutNote(t *testing.T) {
	sc := newPendingCategoryAccess()
	repo := &mockStoreCategoryRepoForReview{records: []*entity.StoreCategory{sc}}
	uc := usecase.NewReviewCategoryAccessUseCase(repo)

	input := usecase.ReviewCategoryAccessInput{
		StoreID:    "store-1",
		CategoryID: "cat-7",
		Decision:   "reject",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Status != string(entity.StoreCategoryStatusRejected) {
		t.Errorf("expected rejected, got %s", output.Status)
	}
	if output.SupportNote != "" {
		t.Errorf("expected empty support note, got %s", output.SupportNote)
	}
}

func TestReviewCategoryAccess_NotFound(t *testing.T) {
	repo := &mockStoreCategoryRepoForReview{}
	uc := usecase.NewReviewCategoryAccessUseCase(repo)

	input := usecase.ReviewCategoryAccessInput{
		StoreID:    "nonexistent",
		CategoryID: "cat-7",
		Decision:   "approve",
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrCategoryAccessNotFound {
		t.Errorf("expected ErrCategoryAccessNotFound, got %v", err)
	}
}

func TestReviewCategoryAccess_InvalidDecision(t *testing.T) {
	sc := newPendingCategoryAccess()
	repo := &mockStoreCategoryRepoForReview{records: []*entity.StoreCategory{sc}}
	uc := usecase.NewReviewCategoryAccessUseCase(repo)

	input := usecase.ReviewCategoryAccessInput{
		StoreID:    "store-1",
		CategoryID: "cat-7",
		Decision:   "invalid",
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrInvalidCategoryReviewDecision {
		t.Errorf("expected ErrInvalidCategoryReviewDecision, got %v", err)
	}
}

func TestReviewCategoryAccess_AlreadyApproved(t *testing.T) {
	sc := newPendingCategoryAccess()
	sc.Approve()
	sc.Events()
	repo := &mockStoreCategoryRepoForReview{records: []*entity.StoreCategory{sc}}
	uc := usecase.NewReviewCategoryAccessUseCase(repo)

	t.Run("reject after approve", func(t *testing.T) {
		input := usecase.ReviewCategoryAccessInput{
			StoreID:    "store-1",
			CategoryID: "cat-7",
			Decision:   "reject",
		}
		_, err := uc.Execute(input)
		if err != entity.ErrStoreCategoryAlreadyDecided {
			t.Errorf("expected ErrStoreCategoryAlreadyDecided, got %v", err)
		}
	})

	t.Run("approve again", func(t *testing.T) {
		input := usecase.ReviewCategoryAccessInput{
			StoreID:    "store-1",
			CategoryID: "cat-7",
			Decision:   "approve",
		}
		_, err := uc.Execute(input)
		if err != entity.ErrStoreCategoryAlreadyApproved {
			t.Errorf("expected ErrStoreCategoryAlreadyApproved, got %v", err)
		}
	})
}

func TestReviewCategoryAccess_AlreadyRejected(t *testing.T) {
	sc := newPendingCategoryAccess()
	sc.Reject("")
	sc.Events()
	repo := &mockStoreCategoryRepoForReview{records: []*entity.StoreCategory{sc}}
	uc := usecase.NewReviewCategoryAccessUseCase(repo)

	t.Run("approve after reject", func(t *testing.T) {
		input := usecase.ReviewCategoryAccessInput{
			StoreID:    "store-1",
			CategoryID: "cat-7",
			Decision:   "approve",
		}
		_, err := uc.Execute(input)
		if err != entity.ErrStoreCategoryAlreadyDecided {
			t.Errorf("expected ErrStoreCategoryAlreadyDecided, got %v", err)
		}
	})

	t.Run("reject again", func(t *testing.T) {
		input := usecase.ReviewCategoryAccessInput{
			StoreID:    "store-1",
			CategoryID: "cat-7",
			Decision:   "reject",
		}
		_, err := uc.Execute(input)
		if err != entity.ErrStoreCategoryAlreadyRejected {
			t.Errorf("expected ErrStoreCategoryAlreadyRejected, got %v", err)
		}
	})
}
