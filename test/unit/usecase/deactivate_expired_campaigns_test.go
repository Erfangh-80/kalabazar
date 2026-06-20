package usecase_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/application/usecase"
	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

type mockDeactivatePromotionRepo struct {
	items map[string]*entity.Promotion
}

func (m *mockDeactivatePromotionRepo) Save(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

func (m *mockDeactivatePromotionRepo) FindByID(id string) (*entity.Promotion, error) {
	p, ok := m.items[id]
	if !ok {
		return nil, entity.ErrPromotionNotFound
	}
	return p, nil
}

func (m *mockDeactivatePromotionRepo) FindSchedulable(now time.Time) ([]*entity.Promotion, error) {
	return nil, nil
}

func (m *mockDeactivatePromotionRepo) FindExpired(now time.Time) ([]*entity.Promotion, error) {
	var result []*entity.Promotion
	for _, p := range m.items {
		if !p.IsActive {
			continue
		}
		if p.EndAt.After(now) {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

func (m *mockDeactivatePromotionRepo) Update(p *entity.Promotion) error {
	m.items[p.ID] = p
	return nil
}

type mockDeactivateInventoryRepo struct {
	items map[string]*entity.Inventory
}

func (m *mockDeactivateInventoryRepo) Save(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

func (m *mockDeactivateInventoryRepo) FindByID(id string) (*entity.Inventory, error) {
	inv, ok := m.items[id]
	if !ok {
		return nil, entity.ErrInventoryNotFound
	}
	return inv, nil
}

func (m *mockDeactivateInventoryRepo) FindByStoreID(storeID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockDeactivateInventoryRepo) FindByWarehouseID(warehouseID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockDeactivateInventoryRepo) FindByProductID(productID string) ([]*entity.Inventory, error) {
	return nil, nil
}

func (m *mockDeactivateInventoryRepo) FindByPromotionID(promotionID string) ([]*entity.Inventory, error) {
	var result []*entity.Inventory
	for _, inv := range m.items {
		if inv.PromotionID != nil && *inv.PromotionID == promotionID {
			result = append(result, inv)
		}
	}
	return result, nil
}

func (m *mockDeactivateInventoryRepo) Update(inv *entity.Inventory) error {
	m.items[inv.ID] = inv
	return nil
}

func newActiveCampaign(id string, startAt, endAt time.Time, expireWithSale bool) *entity.Promotion {
	p, _ := entity.NewPromotion(id, "Campaign "+id, "", startAt, endAt, false, 20, false)
	p.IsActive = true
	p.ExpireSaleWithPromotion = expireWithSale
	p.Events()
	return p
}

func newInventoryLinkedToCampaign(id, productID, campaignID string, basePrice, finalPrice float64) *entity.Inventory {
	inv, _ := entity.NewInventory(id, "store-1", "wh-1", productID, basePrice, 50, "fixed", "new", 1, nil, nil)
	inv.FinalPrice = finalPrice
	inv.PromotionID = &campaignID
	inv.VendorSaleStatus = entity.VendorSaleStatusActive
	inv.SystemSaleStatus = entity.SystemSaleStatusActive
	inv.Events()
	return inv
}

func TestDeactivateExpiredCampaigns_Success(t *testing.T) {
	now := time.Now()
	start := now.Add(-48 * time.Hour)
	end := now.Add(-1 * time.Hour) // expired

	campaign := newActiveCampaign("promo-1", start, end, false)
	inv := newInventoryLinkedToCampaign("inv-1", "prod-x200", "promo-1", 500000, 400000)

	promoRepo := &mockDeactivatePromotionRepo{items: map[string]*entity.Promotion{"promo-1": campaign}}
	invRepo := &mockDeactivateInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}

	uc := usecase.NewDeactivateExpiredCampaignsUseCase(promoRepo, invRepo)

	output, err := uc.Execute(usecase.DeactivateExpiredCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.DeactivatedCampaigns) != 1 {
		t.Fatalf("expected 1 deactivated campaign, got %d", len(output.DeactivatedCampaigns))
	}

	info := output.DeactivatedCampaigns[0]
	if info.CampaignID != "promo-1" {
		t.Errorf("expected promo-1, got %s", info.CampaignID)
	}
	if campaign.IsActive {
		t.Error("expected campaign to be deactivated")
	}

	// verify price reset
	savedInv, _ := invRepo.FindByID("inv-1")
	if savedInv.FinalPrice != savedInv.BasePrice {
		t.Errorf("expected final_price %f to be reset to base_price %f", savedInv.FinalPrice, savedInv.BasePrice)
	}
	if savedInv.VendorSaleStatus != entity.VendorSaleStatusActive {
		t.Errorf("expected vendor_sale_status to remain active, got %s", savedInv.VendorSaleStatus)
	}

	// verify events
	if info.CampaignEvent == nil {
		t.Fatal("expected campaign event, got nil")
	}
	_, ok := info.CampaignEvent.(event.PromotionDeactivated)
	if !ok {
		t.Fatalf("expected PromotionDeactivated, got %T", info.CampaignEvent)
	}

	if len(info.InventoryEvents) != 1 {
		t.Fatalf("expected 1 inventory event, got %d", len(info.InventoryEvents))
	}
	_, ok = info.InventoryEvents[0].(event.InventoryPriceUpdated)
	if !ok {
		t.Fatalf("expected InventoryPriceUpdated, got %T", info.InventoryEvents[0])
	}
}

func TestDeactivateExpiredCampaigns_ExpireSaleWithPromotion(t *testing.T) {
	now := time.Now()
	start := now.Add(-48 * time.Hour)
	end := now.Add(-1 * time.Hour) // expired

	campaign := newActiveCampaign("promo-1", start, end, true) // expire_sale_with_promotion = true
	inv := newInventoryLinkedToCampaign("inv-1", "prod-x200", "promo-1", 500000, 400000)

	promoRepo := &mockDeactivatePromotionRepo{items: map[string]*entity.Promotion{"promo-1": campaign}}
	invRepo := &mockDeactivateInventoryRepo{items: map[string]*entity.Inventory{"inv-1": inv}}

	uc := usecase.NewDeactivateExpiredCampaignsUseCase(promoRepo, invRepo)

	output, err := uc.Execute(usecase.DeactivateExpiredCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.DeactivatedCampaigns) != 1 {
		t.Fatalf("expected 1 deactivated campaign, got %d", len(output.DeactivatedCampaigns))
	}

	savedInv, _ := invRepo.FindByID("inv-1")
	if savedInv.VendorSaleStatus != entity.VendorSaleStatusInactive {
		t.Errorf("expected vendor_sale_status to be inactive, got %s", savedInv.VendorSaleStatus)
	}

	info := output.DeactivatedCampaigns[0]
	if len(info.InventoryEvents) != 2 {
		t.Fatalf("expected 2 inventory events (price_updated + item_deactivated), got %d", len(info.InventoryEvents))
	}
}

func TestDeactivateExpiredCampaigns_NoExpired(t *testing.T) {
	now := time.Now()
	start := now.Add(24 * time.Hour)
	end := now.Add(72 * time.Hour) // future

	campaign := newActiveCampaign("promo-1", start, end, false)

	promoRepo := &mockDeactivatePromotionRepo{items: map[string]*entity.Promotion{"promo-1": campaign}}
	invRepo := &mockDeactivateInventoryRepo{items: map[string]*entity.Inventory{}}

	uc := usecase.NewDeactivateExpiredCampaignsUseCase(promoRepo, invRepo)

	output, err := uc.Execute(usecase.DeactivateExpiredCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.DeactivatedCampaigns) != 0 {
		t.Errorf("expected 0 deactivated campaigns, got %d", len(output.DeactivatedCampaigns))
	}
}

func TestDeactivateExpiredCampaigns_MultipleCampaigns(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(72 * time.Hour)

	camp1 := newActiveCampaign("promo-1", now.Add(-48*time.Hour), past, false)
	camp2 := newActiveCampaign("promo-2", now.Add(-24*time.Hour), past, false)
	camp3 := newActiveCampaign("promo-3", now.Add(24*time.Hour), future, false) // not expired

	inv1 := newInventoryLinkedToCampaign("inv-1", "prod-a", "promo-1", 100000, 80000)
	inv2 := newInventoryLinkedToCampaign("inv-2", "prod-b", "promo-2", 200000, 160000)

	promoRepo := &mockDeactivatePromotionRepo{items: map[string]*entity.Promotion{
		"promo-1": camp1,
		"promo-2": camp2,
		"promo-3": camp3,
	}}
	invRepo := &mockDeactivateInventoryRepo{items: map[string]*entity.Inventory{
		"inv-1": inv1,
		"inv-2": inv2,
	}}

	uc := usecase.NewDeactivateExpiredCampaignsUseCase(promoRepo, invRepo)

	output, err := uc.Execute(usecase.DeactivateExpiredCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.DeactivatedCampaigns) != 2 {
		t.Fatalf("expected 2 deactivated campaigns, got %d", len(output.DeactivatedCampaigns))
	}

	ids := make(map[string]bool)
	for _, c := range output.DeactivatedCampaigns {
		ids[c.CampaignID] = true
	}
	if !ids["promo-1"] {
		t.Error("expected promo-1 to be deactivated")
	}
	if !ids["promo-2"] {
		t.Error("expected promo-2 to be deactivated")
	}
	if ids["promo-3"] {
		t.Error("expected promo-3 to NOT be deactivated")
	}

	if !camp3.IsActive {
		t.Error("expected promo-3 to remain active")
	}
}

func TestDeactivateExpiredCampaigns_NoLinkedInventory(t *testing.T) {
	now := time.Now()
	end := now.Add(-1 * time.Hour)

	campaign := newActiveCampaign("promo-1", now.Add(-48*time.Hour), end, false)

	promoRepo := &mockDeactivatePromotionRepo{items: map[string]*entity.Promotion{"promo-1": campaign}}
	invRepo := &mockDeactivateInventoryRepo{items: map[string]*entity.Inventory{}}

	uc := usecase.NewDeactivateExpiredCampaignsUseCase(promoRepo, invRepo)

	output, err := uc.Execute(usecase.DeactivateExpiredCampaignsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.DeactivatedCampaigns) != 1 {
		t.Fatalf("expected 1 deactivated campaign, got %d", len(output.DeactivatedCampaigns))
	}
	if len(output.DeactivatedCampaigns[0].InventoryEvents) != 0 {
		t.Errorf("expected 0 inventory events, got %d", len(output.DeactivatedCampaigns[0].InventoryEvents))
	}
}
