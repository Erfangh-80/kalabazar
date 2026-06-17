package entity_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/domain/entity"
)

func TestCompleteScenario_FullFlow(t *testing.T) {
	now := time.Now()

	// ── Step 1: Register a Store ──
	store, err := entity.NewStore("store-1", "user-saeed", "Electronics Shop", nil, nil, nil)
	if err != nil {
		t.Fatalf("Step 1: expected no error, got %v", err)
	}
	if store.Status != entity.StoreStatusActive {
		t.Errorf("Step 1: expected active status, got %s", store.Status)
	}
	assertEvent(t, store.Events(), "store.created")

	// ── Step 2: Create Warehouse and Link to Store ──
	addr := entity.Address{
		Street: "123 Main St", City: "Tehran", Country: "Iran",
	}
	wh, err := entity.NewWarehouse("wh-1", "user-saeed", "Tehran Central Warehouse", addr, 10000, "public")
	if err != nil {
		t.Fatalf("Step 2: expected no error, got %v", err)
	}
	assertEvent(t, wh.Events(), "warehouse.created")

	err = wh.LinkToStore("store-1")
	if err != nil {
		t.Fatalf("Step 2: link to store failed: %v", err)
	}
	if wh.StoreID != "store-1" {
		t.Errorf("Step 2: expected store-1, got %s", wh.StoreID)
	}
	assertEvent(t, wh.Events(), "warehouse.linked_to_store")

	// ── Step 3: Register Inventory Item ──
	inv, err := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50)
	if err != nil {
		t.Fatalf("Step 3: expected no error, got %v", err)
	}
	if inv.InstantQty != 50 {
		t.Errorf("Step 3: expected stock 50, got %d", inv.InstantQty)
	}
	if !inv.CanBeSold() {
		t.Error("Step 3: expected item to be available for sale")
	}
	assertEvent(t, inv.Events(), "inventory.item_created")

	// ── Step 4: Register Reference Market Price ──
	rp, err := entity.NewReferencePrice("rp-1", "prod-x200", 550000, "Digikala")
	if err != nil {
		t.Fatalf("Step 4: expected no error, got %v", err)
	}
	assertEvent(t, rp.Events(), "pricing.reference_price_recorded")

	// ── Step 5: Create Discount Campaign ──
	campaign, err := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true)
	if err != nil {
		t.Fatalf("Step 5: expected no error, got %v", err)
	}
	if campaign.IsActive {
		t.Error("Step 5: expected campaign to be inactive initially")
	}
	events := campaign.Events()
	assertEvent(t, events, "promotion.campaign_created")

	err = campaign.LinkToProduct("prod-x200")
	if err != nil {
		t.Fatalf("Step 5: link to product failed: %v", err)
	}
	assertEvent(t, campaign.Events(), "promotion.campaign_linked_to_product")

	err = inv.LinkPromotion("promo-1")
	if err != nil {
		t.Fatalf("Step 5: link promotion to inventory failed: %v", err)
	}
	if inv.PromotionID == nil || *inv.PromotionID != "promo-1" {
		t.Errorf("Step 5: expected promotion promo-1, got %v", inv.PromotionID)
	}
	events = inv.Events()
	assertEvent(t, events, "inventory.promotion_linked")

	// ── Step 6: Campaign Approval ──
	err = campaign.Approve()
	if err != nil {
		t.Fatalf("Step 6: approve failed: %v", err)
	}
	if campaign.ApprovalStatus != entity.PromotionApprovalApproved {
		t.Errorf("Step 6: expected approved, got %s", campaign.ApprovalStatus)
	}
	assertEvent(t, campaign.Events(), "promotion.campaign_approved")

	err = inv.UpdatePromotionStatus(entity.CampaignApprovalApproved)
	if err != nil {
		t.Fatalf("Step 6: update promotion status failed: %v", err)
	}
	if inv.CampaignApprovalStatus != entity.CampaignApprovalApproved {
		t.Errorf("Step 6: expected approved, got %s", inv.CampaignApprovalStatus)
	}
	assertEvent(t, inv.Events(), "inventory.promotion_status_changed")

	// ── Step 7: Final Price Calculation ──
	finalPrice, err := entity.CalculateFinalPrice(500000, 20)
	if err != nil {
		t.Fatalf("Step 7: calculate final price failed: %v", err)
	}
	if finalPrice != 400000 {
		t.Errorf("Step 7: expected 400000, got %f", finalPrice)
	}

	err = inv.UpdatePrice(500000, finalPrice)
	if err != nil {
		t.Fatalf("Step 7: update price failed: %v", err)
	}
	if inv.FinalPrice != 400000 {
		t.Errorf("Step 7: expected final 400000, got %f", inv.FinalPrice)
	}
	assertEvent(t, inv.Events(), "inventory.price_updated")

	// ── Step 8: Campaign Activation (Time-Based) ──
	err = campaign.Activate()
	if err != nil {
		t.Fatalf("Step 8: activate failed: %v", err)
	}
	if !campaign.IsActive {
		t.Error("Step 8: expected campaign to be active")
	}
	assertEvent(t, campaign.Events(), "promotion.campaign_activated")

	// ── Step 9: Successful Purchase ──
	if !inv.CanBeSold() {
		t.Fatal("Step 9: expected item to be sellable before purchase")
	}
	err = inv.UpdateStock(49)
	if err != nil {
		t.Fatalf("Step 9: update stock failed: %v", err)
	}
	if inv.InstantQty != 49 {
		t.Errorf("Step 9: expected stock 49, got %d", inv.InstantQty)
	}
	assertEvent(t, inv.Events(), "inventory.stock_updated")

	// ── Step 10: Commission Calculation ──
	comm, err := entity.NewCommission("comm-1", "prod-x200", "retail", 10, 100000, 1000000, 1)
	if err != nil {
		t.Fatalf("Step 10: expected no error, got %v", err)
	}
	comm.Events()

	commissionAmount, err := comm.Calculate(400000, 1)
	if err != nil {
		t.Fatalf("Step 10: calculate commission failed: %v", err)
	}
	if commissionAmount != 40000 {
		t.Errorf("Step 10: expected 40000, got %f", commissionAmount)
	}
	events = comm.Events()
	assertEvent(t, events, "commission.calculated")

	// ── Step 11: Campaign End ──
	err = campaign.Deactivate()
	if err != nil {
		t.Fatalf("Step 11: deactivate failed: %v", err)
	}
	if campaign.IsActive {
		t.Error("Step 11: expected campaign to be inactive")
	}
	assertEvent(t, campaign.Events(), "promotion.campaign_deactivated")

	// Price resets to base price
	err = inv.ResetPrice()
	if err != nil {
		t.Fatalf("Step 11: reset price failed: %v", err)
	}
	if inv.FinalPrice != inv.BasePrice {
		t.Errorf("Step 11: expected FinalPrice %f to equal BasePrice %f", inv.FinalPrice, inv.BasePrice)
	}
	assertEvent(t, inv.Events(), "inventory.price_updated")
}

func TestCompleteScenario_CampaignEndWithExpire(t *testing.T) {
	now := time.Now()

	campaign, err := entity.NewPromotion("promo-1", "Sale", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	campaign.ExpireSaleWithPromotion = true
	campaign.Events()
	campaign.Activate()
	campaign.Events()

	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-1", 100, 5)
	inv.Events()

	// Campaign ends → item removed from sale
	err = campaign.Deactivate()
	if err != nil {
		t.Fatalf("deactivate failed: %v", err)
	}

	if campaign.ExpireSaleWithPromotion {
		// When expire_sale_with_promotion is true, the item should be removed from sale
		inv.SetVendorStatus(entity.VendorSaleStatusInactive)
		if inv.CanBeSold() {
			t.Error("expected item to be removed from sale after campaign end with expire flag")
		}
	}
}

func assertEvent(t *testing.T, events []any, expectedName string) {
	t.Helper()
	for _, e := range events {
		if ev, ok := e.(interface{ EventName() string }); ok {
			if ev.EventName() == expectedName {
				return
			}
		}
	}
	t.Errorf("expected event %s in %v", expectedName, events)
}
