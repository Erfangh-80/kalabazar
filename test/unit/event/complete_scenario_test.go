package event_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/domain/entity"
	"kalabazar-stock-service/internal/domain/event"
)

func TestStep1_RegisterStore_Event(t *testing.T) {
	phone := "09121234567"
	store, _ := entity.NewStore("store-1", "user-saeed", "Electronics Shop", &phone, nil, nil)
	events := store.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreCreated)
	if !ok {
		t.Fatalf("expected StoreCreated, got %T", events[0])
	}
	if e.StoreID != "store-1" {
		t.Errorf("expected StoreID store-1, got %s", e.StoreID)
	}
	if e.EventName() != "store.created" {
		t.Errorf("expected store.created, got %s", e.EventName())
	}
}

func TestStep1dot5_CategoryPermission_Event(t *testing.T) {
	sc, _ := entity.NewStoreCategory("store-1", "cat-7")
	events := sc.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.StoreCategoryAllowed)
	if !ok {
		t.Fatalf("expected StoreCategoryAllowed, got %T", events[0])
	}
	if e.StoreID != "store-1" {
		t.Errorf("expected StoreID store-1, got %s", e.StoreID)
	}
	if e.CategoryID != "cat-7" {
		t.Errorf("expected CategoryID cat-7, got %s", e.CategoryID)
	}
	if e.Status != "pending" {
		t.Errorf("expected status pending, got %s", e.Status)
	}
	if e.EventName() != "store.category_allowed" {
		t.Errorf("expected store.category_allowed, got %s", e.EventName())
	}

	sc.Approve()
	events = sc.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event after approve, got %d", len(events))
	}
	e2, ok := events[0].(event.StoreCategoryAllowed)
	if !ok {
		t.Fatalf("expected StoreCategoryAllowed, got %T", events[0])
	}
	if e2.Status != "approved" {
		t.Errorf("expected status approved, got %s", e2.Status)
	}
}

func TestStep2_CreateWarehouseAndLink_Event(t *testing.T) {
	addr := entity.Address{Street: "123 Main St", City: "Tehran", Country: "Iran"}
	wh, _ := entity.NewWarehouse("wh-1", "user-saeed", "Tehran Central Warehouse", addr, 10000, "public")
	events := wh.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	we, ok := events[0].(event.WarehouseCreated)
	if !ok {
		t.Fatalf("expected WarehouseCreated, got %T", events[0])
	}
	if we.WarehouseID != "wh-1" {
		t.Errorf("expected WarehouseID wh-1, got %s", we.WarehouseID)
	}
	if we.EventName() != "warehouse.created" {
		t.Errorf("expected warehouse.created, got %s", we.EventName())
	}

	wh.LinkToStore("store-1", "primary")
	events = wh.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	le, ok := events[0].(event.WarehouseLinkedToStore)
	if !ok {
		t.Fatalf("expected WarehouseLinkedToStore, got %T", events[0])
	}
	if le.StoreID != "store-1" {
		t.Errorf("expected StoreID store-1, got %s", le.StoreID)
	}
	if le.RelationType != "primary" {
		t.Errorf("expected RelationType primary, got %s", le.RelationType)
	}
	if le.EventName() != "warehouse.linked_to_store" {
		t.Errorf("expected warehouse.linked_to_store, got %s", le.EventName())
	}
}

func TestStep3_RegisterGoods_Event(t *testing.T) {
	attrs := map[string]string{"color": "black", "warranty_months": "18"}
	maxQty := 100
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, &maxQty, attrs)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryItemCreated)
	if !ok {
		t.Fatalf("expected InventoryItemCreated, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.ProductID != "prod-x200" {
		t.Errorf("expected ProductID prod-x200, got %s", e.ProductID)
	}
	if e.EventName() != "inventory.item_created" {
		t.Errorf("expected inventory.item_created, got %s", e.EventName())
	}
}

func TestStep4_RecordReferencePrice_Event(t *testing.T) {
	rp, _ := entity.NewReferencePrice("rp-1", "prod-x200", 550000, "DigiKala")
	events := rp.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.ReferencePriceCreated)
	if !ok {
		t.Fatalf("expected ReferencePriceCreated, got %T", events[0])
	}
	if e.ReferencePriceID != "rp-1" {
		t.Errorf("expected ReferencePriceID rp-1, got %s", e.ReferencePriceID)
	}
	if e.ProductID != "prod-x200" {
		t.Errorf("expected ProductID prod-x200, got %s", e.ProductID)
	}
	if e.Price != 550000 {
		t.Errorf("expected Price 550000, got %f", e.Price)
	}
	if e.Source != "DigiKala" {
		t.Errorf("expected Source DigiKala, got %s", e.Source)
	}
	if e.EventName() != "pricing.reference_price_recorded" {
		t.Errorf("expected pricing.reference_price_recorded, got %s", e.EventName())
	}
}

func TestStep5_CreateCampaignAndLink_Event(t *testing.T) {
	now := time.Now()

	campaign, _ := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	events := campaign.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ce, ok := events[0].(event.PromotionCreated)
	if !ok {
		t.Fatalf("expected PromotionCreated, got %T", events[0])
	}
	if ce.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", ce.PromotionID)
	}
	if ce.Title != "Nowruz Auction" {
		t.Errorf("expected Title Nowruz Auction, got %s", ce.Title)
	}
	if ce.EventName() != "promotion.campaign_created" {
		t.Errorf("expected promotion.campaign_created, got %s", ce.EventName())
	}

	campaign.LinkToProduct("prod-x200")
	events = campaign.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	lpe, ok := events[0].(event.PromotionCampaignLinkedToProduct)
	if !ok {
		t.Fatalf("expected PromotionCampaignLinkedToProduct, got %T", events[0])
	}
	if lpe.ProductID != "prod-x200" {
		t.Errorf("expected ProductID prod-x200, got %s", lpe.ProductID)
	}
	if lpe.EventName() != "promotion.campaign_linked_to_product" {
		t.Errorf("expected promotion.campaign_linked_to_product, got %s", lpe.EventName())
	}

	maxQty := 100
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, &maxQty, nil)
	inv.Events()

	inv.LinkPromotion("promo-1")
	events = inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	pie, ok := events[0].(event.InventoryPromotionLinked)
	if !ok {
		t.Fatalf("expected InventoryPromotionLinked, got %T", events[0])
	}
	if pie.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", pie.PromotionID)
	}
	if pie.EventName() != "inventory.promotion_linked" {
		t.Errorf("expected inventory.promotion_linked, got %s", pie.EventName())
	}
}

func TestStep6_CampaignApproval_Event(t *testing.T) {
	now := time.Now()
	campaign, _ := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	campaign.Events()

	campaign.Approve()
	events := campaign.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	pae, ok := events[0].(event.PromotionApproved)
	if !ok {
		t.Fatalf("expected PromotionApproved, got %T", events[0])
	}
	if pae.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", pae.PromotionID)
	}
	if pae.EventName() != "promotion.campaign_approved" {
		t.Errorf("expected promotion.campaign_approved, got %s", pae.EventName())
	}

	maxQty := 100
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, &maxQty, nil)
	inv.Events()

	inv.UpdatePromotionStatus(entity.CampaignApprovalApproved)
	events = inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	psce, ok := events[0].(event.InventoryPromotionStatusChanged)
	if !ok {
		t.Fatalf("expected InventoryPromotionStatusChanged, got %T", events[0])
	}
	if psce.Status != "approved" {
		t.Errorf("expected status approved, got %s", psce.Status)
	}
	if psce.EventName() != "inventory.promotion_status_changed" {
		t.Errorf("expected inventory.promotion_status_changed, got %s", psce.EventName())
	}
}

func TestStep7_CalculateFinalPrice_Event(t *testing.T) {
	maxQty := 100
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, &maxQty, nil)
	inv.Events()

	inv.UpdatePrice(500000, 400000)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryPriceUpdated)
	if !ok {
		t.Fatalf("expected InventoryPriceUpdated, got %T", events[0])
	}
	if e.BasePrice != 500000 {
		t.Errorf("expected BasePrice 500000, got %f", e.BasePrice)
	}
	if e.FinalPrice != 400000 {
		t.Errorf("expected FinalPrice 400000, got %f", e.FinalPrice)
	}
	if e.EventName() != "inventory.price_updated" {
		t.Errorf("expected inventory.price_updated, got %s", e.EventName())
	}
}

func TestStep8_ActivateCampaign_Event(t *testing.T) {
	now := time.Now()
	campaign, _ := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	campaign.Approve()
	campaign.Events()

	campaign.Activate()
	events := campaign.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.PromotionActivated)
	if !ok {
		t.Fatalf("expected PromotionActivated, got %T", events[0])
	}
	if e.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", e.PromotionID)
	}
	if e.EventName() != "promotion.campaign_activated" {
		t.Errorf("expected promotion.campaign_activated, got %s", e.EventName())
	}
}

func TestStep9_Sales_Event(t *testing.T) {
	maxQty := 100
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, &maxQty, nil)
	inv.Events()

	inv.UpdateStock(49)
	events := inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e, ok := events[0].(event.InventoryStockUpdated)
	if !ok {
		t.Fatalf("expected InventoryStockUpdated, got %T", events[0])
	}
	if e.InventoryID != "inv-1" {
		t.Errorf("expected InventoryID inv-1, got %s", e.InventoryID)
	}
	if e.NewQty != 49 {
		t.Errorf("expected NewQty 49, got %d", e.NewQty)
	}
	if e.EventName() != "inventory.stock_updated" {
		t.Errorf("expected inventory.stock_updated, got %s", e.EventName())
	}
}

func TestStep10_CalculateCommission_Event(t *testing.T) {
	comm, _ := entity.NewCommission("comm-1", "prod-x200", "retail", 10, 100000, 1000000, 1)
	events := comm.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	rce, ok := events[0].(event.CommissionRuleCreated)
	if !ok {
		t.Fatalf("expected CommissionRuleCreated, got %T", events[0])
	}
	if rce.CommissionID != "comm-1" {
		t.Errorf("expected CommissionID comm-1, got %s", rce.CommissionID)
	}
	if rce.RatePercent != 10 {
		t.Errorf("expected RatePercent 10, got %f", rce.RatePercent)
	}
	if rce.EventName() != "commission.rule.created" {
		t.Errorf("expected commission.rule.created, got %s", rce.EventName())
	}

	comm.Calculate(400000, 1)
	events = comm.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	cce, ok := events[0].(event.CommissionCalculated)
	if !ok {
		t.Fatalf("expected CommissionCalculated, got %T", events[0])
	}
	if cce.SaleAmount != 400000 {
		t.Errorf("expected SaleAmount 400000, got %f", cce.SaleAmount)
	}
	if cce.CommissionAmount != 40000 {
		t.Errorf("expected CommissionAmount 40000, got %f", cce.CommissionAmount)
	}
	if cce.EventName() != "commission.calculated" {
		t.Errorf("expected commission.calculated, got %s", cce.EventName())
	}
}

func TestStep11_EndCampaign_Event(t *testing.T) {
	now := time.Now()
	campaign, _ := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	campaign.Approve()
	campaign.Activate()
	campaign.Events()

	maxQty := 100
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, &maxQty, nil)
	inv.UpdatePrice(500000, 400000)
	inv.Events()

	campaign.Deactivate()
	events := campaign.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	pde, ok := events[0].(event.PromotionDeactivated)
	if !ok {
		t.Fatalf("expected PromotionDeactivated, got %T", events[0])
	}
	if pde.PromotionID != "promo-1" {
		t.Errorf("expected PromotionID promo-1, got %s", pde.PromotionID)
	}
	if pde.EventName() != "promotion.campaign_deactivated" {
		t.Errorf("expected promotion.campaign_deactivated, got %s", pde.EventName())
	}

	inv.ResetPrice()
	events = inv.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ipue, ok := events[0].(event.InventoryPriceUpdated)
	if !ok {
		t.Fatalf("expected InventoryPriceUpdated, got %T", events[0])
	}
	if ipue.FinalPrice != 500000 {
		t.Errorf("expected FinalPrice 500000, got %f", ipue.FinalPrice)
	}
	if ipue.EventName() != "inventory.price_updated" {
		t.Errorf("expected inventory.price_updated, got %s", ipue.EventName())
	}
}
