package usecase_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockReviewCampaignInventoryRepo struct {
	items map[string]*entity.Inventory
}

func (m *mockReviewCampaignInventoryRepo) Save(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

func (m *mockReviewCampaignInventoryRepo) FindByID(id string) (*entity.Inventory, error) {
	inv, ok := m.items[id]
	if !ok {
		return nil, entity.ErrInventoryNotFound
	}
	return inv, nil
}

func (m *mockReviewCampaignInventoryRepo) FindByStoreID(storeID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockReviewCampaignInventoryRepo) FindByWarehouseID(warehouseID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockReviewCampaignInventoryRepo) FindByProductID(productID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockReviewCampaignInventoryRepo) Update(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

type mockReviewCampaignPromotionRepo struct {
	items map[string]*entity.Promotion
}

func (m *mockReviewCampaignPromotionRepo) Save(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

func (m *mockReviewCampaignPromotionRepo) FindByID(id string) (*entity.Promotion, error) {
	p, ok := m.items[id]
	if !ok {
		return nil, entity.ErrPromotionNotFound
	}
	return p, nil
}

func (m *mockReviewCampaignPromotionRepo) Update(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

func newInventoryWithPromotion(invID, productID, promoID string, basePrice float64) *entity.Inventory {
	inv, _ := entity.NewInventory(invID, "store-1", "wh-1", productID, basePrice, 50, "fixed", "new", 1, nil, nil)
	inv.PromotionID = &promoID
	inv.Events()
	return inv
}

func newCampaignRequiringApproval(id, title string, discountPercent float64) *entity.Promotion {
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(15 * 24 * time.Hour)
	p, _ := entity.NewPromotion(id, title, "", start, end, true, discountPercent, false)
	p.Events()
	return p
}

func TestReviewCampaign_Approve(t *testing.T) {
	inv := newInventoryWithPromotion("inv-1", "prod-1", "promo-1", 500000)
	promo := newCampaignRequiringApproval("promo-1", "Nowruz Auction", 20)

	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockReviewCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "inv-1",
		Decision:    "approve",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.InventoryID != "inv-1" {
		t.Errorf("expected inv-1, got %s", output.InventoryID)
	}
	if output.PromotionID != "promo-1" {
		t.Errorf("expected promo-1, got %s", output.PromotionID)
	}
	if output.ApprovalStatus != string(entity.CampaignApprovalApproved) {
		t.Errorf("expected approved, got %s", output.ApprovalStatus)
	}
	if output.FinalPrice != 400000 {
		t.Errorf("expected 400000 (500000 - 20%%), got %f", output.FinalPrice)
	}

	if inv.CampaignApprovalStatus != entity.CampaignApprovalApproved {
		t.Errorf("expected inventory campaign status approved, got %s", inv.CampaignApprovalStatus)
	}
	if inv.FinalPrice != 400000 {
		t.Errorf("expected inventory final price 400000, got %f", inv.FinalPrice)
	}
	if promo.ApprovalStatus != entity.PromotionApprovalApproved {
		t.Errorf("expected promotion approval status approved, got %s", promo.ApprovalStatus)
	}
}

func TestReviewCampaign_Approve_Events(t *testing.T) {
	inv := newInventoryWithPromotion("inv-1", "prod-1", "promo-1", 500000)
	promo := newCampaignRequiringApproval("promo-1", "Nowruz Auction", 20)

	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockReviewCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "inv-1",
		Decision:    "approve",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.PromotionEvent == nil {
		t.Fatal("expected promotion event, got nil")
	}
	promEvt, ok := output.PromotionEvent.(event.PromotionApproved)
	if !ok {
		t.Fatalf("expected PromotionApproved, got %T", output.PromotionEvent)
	}
	if promEvt.EventName() != "promotion.campaign_approved" {
		t.Errorf("expected promotion.campaign_approved, got %s", promEvt.EventName())
	}

	foundStatusChanged := false
	foundPriceUpdated := false
	for _, e := range output.InventoryEvents {
		switch e.(type) {
		case event.InventoryPromotionStatusChanged:
			foundStatusChanged = true
		case event.InventoryPriceUpdated:
			foundPriceUpdated = true
		}
	}
	if !foundStatusChanged {
		t.Error("expected InventoryPromotionStatusChanged event")
	}
	if !foundPriceUpdated {
		t.Error("expected InventoryPriceUpdated event")
	}
}

func TestReviewCampaign_Reject(t *testing.T) {
	inv := newInventoryWithPromotion("inv-1", "prod-1", "promo-1", 500000)
	promo := newCampaignRequiringApproval("promo-1", "Nowruz Auction", 20)

	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockReviewCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "inv-1",
		Decision:    "reject",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ApprovalStatus != string(entity.CampaignApprovalRejected) {
		t.Errorf("expected rejected, got %s", output.ApprovalStatus)
	}
	if output.FinalPrice != 500000 {
		t.Errorf("expected final price unchanged at 500000, got %f", output.FinalPrice)
	}

	if inv.CampaignApprovalStatus != entity.CampaignApprovalRejected {
		t.Errorf("expected inventory campaign status rejected, got %s", inv.CampaignApprovalStatus)
	}
	if inv.FinalPrice != 500000 {
		t.Errorf("expected inventory final price unchanged at 500000, got %f", inv.FinalPrice)
	}
	if promo.ApprovalStatus != entity.PromotionApprovalRejected {
		t.Errorf("expected promotion approval status rejected, got %s", promo.ApprovalStatus)
	}
}

func TestReviewCampaign_Reject_Events(t *testing.T) {
	inv := newInventoryWithPromotion("inv-1", "prod-1", "promo-1", 500000)
	promo := newCampaignRequiringApproval("promo-1", "Nowruz Auction", 20)

	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockReviewCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "inv-1",
		Decision:    "reject",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.PromotionEvent == nil {
		t.Fatal("expected promotion event, got nil")
	}
	promEvt, ok := output.PromotionEvent.(event.PromotionRejected)
	if !ok {
		t.Fatalf("expected PromotionRejected, got %T", output.PromotionEvent)
	}
	if promEvt.EventName() != "promotion.campaign_rejected" {
		t.Errorf("expected promotion.campaign_rejected, got %s", promEvt.EventName())
	}

	foundStatusChanged := false
	for _, e := range output.InventoryEvents {
		if _, ok := e.(event.InventoryPromotionStatusChanged); ok {
			foundStatusChanged = true
		}
	}
	if !foundStatusChanged {
		t.Error("expected InventoryPromotionStatusChanged event")
	}
	if len(output.InventoryEvents) != 1 {
		t.Errorf("expected exactly 1 inventory event for reject, got %d", len(output.InventoryEvents))
	}
}

func TestReviewCampaign_InvalidDecision(t *testing.T) {
	inv := newInventoryWithPromotion("inv-1", "prod-1", "promo-1", 500000)
	promo := newCampaignRequiringApproval("promo-1", "Nowruz Auction", 20)

	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockReviewCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "inv-1",
		Decision:    "invalid",
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrInvalidReviewCampaignDecision {
		t.Errorf("expected ErrInvalidReviewCampaignDecision, got %v", err)
	}
}

func TestReviewCampaign_InventoryNotFound(t *testing.T) {
	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{}}
	promoRepo := &mockReviewCampaignPromotionRepo{items: map[string]*entity.Promotion{}}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "nonexistent",
		Decision:    "approve",
	}

	_, err := uc.Execute(input)
	if err != entity.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}

func TestReviewCampaign_NoLinkedPromotion(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 500000, 50, "fixed", "new", 1, nil, nil)
	inv.Events()

	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockReviewCampaignPromotionRepo{}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "inv-1",
		Decision:    "approve",
	}

	_, err := uc.Execute(input)
	if err != usecase.ErrInventoryHasNoLinkedPromotion {
		t.Errorf("expected ErrInventoryHasNoLinkedPromotion, got %v", err)
	}
}

func TestReviewCampaign_PromotionNotFound(t *testing.T) {
	inv := newInventoryWithPromotion("inv-1", "prod-1", "promo-1", 500000)

	invRepo := &mockReviewCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockReviewCampaignPromotionRepo{items: map[string]*entity.Promotion{}}

	uc := usecase.NewReviewCampaignUseCase(invRepo, promoRepo)

	input := usecase.ReviewCampaignInput{
		InventoryID: "inv-1",
		Decision:    "approve",
	}

	_, err := uc.Execute(input)
	if err != entity.ErrPromotionNotFound {
		t.Errorf("expected ErrPromotionNotFound, got %v", err)
	}
}
