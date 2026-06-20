package usecase_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockAttachCampaignInventoryRepo struct {
	items map[string]*entity.Inventory
}

func (m *mockAttachCampaignInventoryRepo) Save(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

func (m *mockAttachCampaignInventoryRepo) FindByID(id string) (*entity.Inventory, error) {
	inv, ok := m.items[id]
	if !ok {
		return nil, entity.ErrInventoryNotFound
	}
	return inv, nil
}

func (m *mockAttachCampaignInventoryRepo) FindByStoreID(storeID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockAttachCampaignInventoryRepo) FindByWarehouseID(warehouseID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockAttachCampaignInventoryRepo) FindByProductID(productID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockAttachCampaignInventoryRepo) Update(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

type mockAttachCampaignPromotionRepo struct {
	items map[string]*entity.Promotion
}

func (m *mockAttachCampaignPromotionRepo) Save(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

func (m *mockAttachCampaignPromotionRepo) FindByID(id string) (*entity.Promotion, error) {
	p, ok := m.items[id]
	if !ok {
		return nil, entity.ErrPromotionNotFound
	}
	return p, nil
}

func (m *mockAttachCampaignPromotionRepo) FindSchedulable(now time.Time) ([]*entity.Promotion, error) {
	return nil, nil
}

func (m *mockAttachCampaignPromotionRepo) Update(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

func newInventoryForCampaign(id, productID, storeID, warehouseID string) *entity.Inventory {
	inv, _ := entity.NewInventory(id, storeID, warehouseID, productID, 500000, 50, "fixed", "new", 1, nil, nil)
	inv.Events()
	return inv
}

func newCampaign(id, title string) *entity.Promotion {
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(15 * 24 * time.Hour)
	p, _ := entity.NewPromotion(id, title, "", start, end, false, 20, false)
	p.Events()
	return p
}

func TestAttachCampaignToInventory_Success(t *testing.T) {
	inv := newInventoryForCampaign("inv-1", "prod-1", "store-1", "wh-1")
	promo := newCampaign("promo-1", "Nowruz Auction")

	invRepo := &mockAttachCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockAttachCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewAttachCampaignToInventoryUseCase(invRepo, promoRepo)

	input := usecase.AttachCampaignToInventoryInput{
		InventoryID: "inv-1",
		PromotionID: "promo-1",
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
	if output.ProductID != "prod-1" {
		t.Errorf("expected prod-1, got %s", output.ProductID)
	}

	// Check inventory was linked
	if inv.PromotionID == nil || *inv.PromotionID != "promo-1" {
		t.Errorf("expected inventory PromotionID to be promo-1, got %v", inv.PromotionID)
	}

	// Check promotion was linked
	found := false
	for _, pid := range promo.LinkedProductIDs {
		if pid == "prod-1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected prod-1 in promotion's LinkedProductIDs")
	}
}

func TestAttachCampaignToInventory_EventsEmitted(t *testing.T) {
	inv := newInventoryForCampaign("inv-1", "prod-1", "store-1", "wh-1")
	promo := newCampaign("promo-1", "Nowruz Auction")

	invRepo := &mockAttachCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockAttachCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewAttachCampaignToInventoryUseCase(invRepo, promoRepo)

	input := usecase.AttachCampaignToInventoryInput{
		InventoryID: "inv-1",
		PromotionID: "promo-1",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.InventoryEvent == nil {
		t.Fatal("expected inventory event, got nil")
	}
	if output.PromotionEvent == nil {
		t.Fatal("expected promotion event, got nil")
	}

	invEvent, ok := output.InventoryEvent.(event.InventoryPromotionLinked)
	if !ok {
		t.Fatalf("expected InventoryPromotionLinked, got %T", output.InventoryEvent)
	}
	if invEvent.InventoryID != "inv-1" {
		t.Errorf("expected inv-1, got %s", invEvent.InventoryID)
	}
	if invEvent.PromotionID != "promo-1" {
		t.Errorf("expected promo-1, got %s", invEvent.PromotionID)
	}
	if invEvent.EventName() != "inventory.promotion_linked" {
		t.Errorf("expected inventory.promotion_linked, got %s", invEvent.EventName())
	}

	promoEvent, ok := output.PromotionEvent.(event.PromotionCampaignLinkedToProduct)
	if !ok {
		t.Fatalf("expected PromotionCampaignLinkedToProduct, got %T", output.PromotionEvent)
	}
	if promoEvent.PromotionID != "promo-1" {
		t.Errorf("expected promo-1, got %s", promoEvent.PromotionID)
	}
	if promoEvent.ProductID != "prod-1" {
		t.Errorf("expected prod-1, got %s", promoEvent.ProductID)
	}
	if promoEvent.EventName() != "promotion.campaign_linked_to_product" {
		t.Errorf("expected promotion.campaign_linked_to_product, got %s", promoEvent.EventName())
	}
}

func TestAttachCampaignToInventory_ReplaceExisting(t *testing.T) {
	inv := newInventoryForCampaign("inv-1", "prod-1", "store-1", "wh-1")
	oldPromoID := "promo-0"
	inv.PromotionID = &oldPromoID
	inv.Events()

	promo := newCampaign("promo-1", "Nowruz Auction")

	invRepo := &mockAttachCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockAttachCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": promo}}

	uc := usecase.NewAttachCampaignToInventoryUseCase(invRepo, promoRepo)

	input := usecase.AttachCampaignToInventoryInput{
		InventoryID: "inv-1",
		PromotionID: "promo-1",
	}

	output, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.PromotionID != "promo-1" {
		t.Errorf("expected promo-1, got %s", output.PromotionID)
	}
	if inv.PromotionID == nil || *inv.PromotionID != "promo-1" {
		t.Errorf("expected inventory PromotionID to be promo-1, got %v", inv.PromotionID)
	}
}

func TestAttachCampaignToInventory_InventoryNotFound(t *testing.T) {
	invRepo := &mockAttachCampaignInventoryRepo{items: map[string]*entity.Inventory{}}
	promoRepo := &mockAttachCampaignPromotionRepo{items: map[string]*entity.Promotion{"promo-1": newCampaign("promo-1", "Sale")}}

	uc := usecase.NewAttachCampaignToInventoryUseCase(invRepo, promoRepo)

	input := usecase.AttachCampaignToInventoryInput{
		InventoryID: "nonexistent",
		PromotionID: "promo-1",
	}

	_, err := uc.Execute(input)
	if err != entity.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}

func TestAttachCampaignToInventory_PromotionNotFound(t *testing.T) {
	inv := newInventoryForCampaign("inv-1", "prod-1", "store-1", "wh-1")
	invRepo := &mockAttachCampaignInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}
	promoRepo := &mockAttachCampaignPromotionRepo{items: map[string]*entity.Promotion{}}

	uc := usecase.NewAttachCampaignToInventoryUseCase(invRepo, promoRepo)

	input := usecase.AttachCampaignToInventoryInput{
		InventoryID: "inv-1",
		PromotionID: "promo-1",
	}

	_, err := uc.Execute(input)
	if err != entity.ErrPromotionNotFound {
		t.Errorf("expected ErrPromotionNotFound, got %v", err)
	}
}
