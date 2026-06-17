package entity_test

import (
	"testing"
	"time"

	"kalabazar-stock-service/internal/domain/entity"
)

func ptr(i int) *int { return &i }

func TestStep1_RegisterStore(t *testing.T) {
	phone := "09121234567"
	store, err := entity.NewStore("store-1", "user-saeed", "Electronics Shop", &phone, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Status != entity.StoreStatusActive {
		t.Errorf("expected active status, got %s", store.Status)
	}
	store.Events()
}

func TestStep1dot5_CategoryPermission(t *testing.T) {
	catPerm, err := entity.NewStoreCategory("store-1", "cat-7")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if catPerm.Status != entity.StoreCategoryStatusPending {
		t.Errorf("expected pending, got %s", catPerm.Status)
	}
	catPerm.Events()

	err = catPerm.Approve()
	if err != nil {
		t.Fatalf("approve failed: %v", err)
	}
	if catPerm.Status != entity.StoreCategoryStatusApproved {
		t.Errorf("expected approved, got %s", catPerm.Status)
	}
	catPerm.Events()
}

func TestStep2_CreateWarehouseAndLink(t *testing.T) {
	addr := entity.Address{
		Street: "123 Main St", City: "Tehran", Country: "Iran",
	}

	wh, err := entity.NewWarehouse("wh-1", "user-saeed", "Tehran Central Warehouse", addr, 10000, "public")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	wh.Events()

	err = wh.LinkToStore("store-1", "primary")
	if err != nil {
		t.Fatalf("link to store failed: %v", err)
	}
	if wh.StoreID != "store-1" {
		t.Errorf("expected store-1, got %s", wh.StoreID)
	}
	if wh.RelationType != "primary" {
		t.Errorf("expected primary, got %s", wh.RelationType)
	}
	wh.Events()
}

func TestStep3_RegisterGoods(t *testing.T) {
	attrs := map[string]string{"color": "black", "warranty_months": "18"}
	maxQty := 100
	inv, err := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, &maxQty, attrs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if inv.InstantQty != 50 {
		t.Errorf("expected stock 50, got %d", inv.InstantQty)
	}
	if inv.FinalPrice != 500000 {
		t.Errorf("expected final_price 500000, got %f", inv.FinalPrice)
	}
	if inv.VendorSaleStatus != entity.VendorSaleStatusActive {
		t.Errorf("expected vendor_sale_status active, got %s", inv.VendorSaleStatus)
	}
	if inv.SystemSaleStatus != entity.SystemSaleStatusActive {
		t.Errorf("expected system_sale_status active, got %s", inv.SystemSaleStatus)
	}
	if inv.SaleModel != "retail" {
		t.Errorf("expected retail, got %s", inv.SaleModel)
	}
	if inv.Condition != "new" {
		t.Errorf("expected new, got %s", inv.Condition)
	}
	if inv.MinOrderQty != 1 {
		t.Errorf("expected min_order_qty 1, got %d", inv.MinOrderQty)
	}
	if inv.MaxOrderQty == nil || *inv.MaxOrderQty != 100 {
		t.Errorf("expected max_order_qty 100, got %v", inv.MaxOrderQty)
	}
	if inv.Attributes["color"] != "black" {
		t.Errorf("expected color black, got %s", inv.Attributes["color"])
	}
	if inv.PromotionID != nil {
		t.Errorf("expected promotion_id nil, got %v", inv.PromotionID)
	}
	if inv.CampaignApprovalStatus != entity.CampaignApprovalPending {
		t.Errorf("expected promotion_status pending, got %s", inv.CampaignApprovalStatus)
	}
	inv.Events()
}

func TestStep4_RecordReferencePrice(t *testing.T) {
	rp, err := entity.NewReferencePrice("rp-1", "prod-x200", 550000, "DigiKala")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	rp.Events()
}

func TestStep5_CreateCampaignAndLink(t *testing.T) {
	now := time.Now()

	campaign, err := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if campaign.IsActive {
		t.Error("expected campaign to be inactive initially")
	}
	if campaign.DiscountPercent != 20 {
		t.Errorf("expected discount 20, got %f", campaign.DiscountPercent)
	}
	if !campaign.IsCountdown {
		t.Error("expected is_countdown to be true")
	}
	campaign.Events()

	err = campaign.LinkToProduct("prod-x200")
	if err != nil {
		t.Fatalf("link to product failed: %v", err)
	}
	campaign.Events()

	attrs := map[string]string{"color": "black"}
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, ptr(100), attrs)
	inv.Events()

	err = inv.LinkPromotion("promo-1")
	if err != nil {
		t.Fatalf("link promotion to inventory failed: %v", err)
	}
	if inv.PromotionID == nil || *inv.PromotionID != "promo-1" {
		t.Errorf("expected promotion promo-1, got %v", inv.PromotionID)
	}
	inv.Events()
}

func TestStep6_CampaignApproval(t *testing.T) {
	now := time.Now()
	campaign, _ := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	campaign.Events()

	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, ptr(100), nil)
	inv.Events()

	err := campaign.Approve()
	if err != nil {
		t.Fatalf("approve failed: %v", err)
	}
	if campaign.ApprovalStatus != entity.PromotionApprovalApproved {
		t.Errorf("expected approved, got %s", campaign.ApprovalStatus)
	}
	campaign.Events()

	err = inv.UpdatePromotionStatus(entity.CampaignApprovalApproved)
	if err != nil {
		t.Fatalf("update promotion status failed: %v", err)
	}
	if inv.CampaignApprovalStatus != entity.CampaignApprovalApproved {
		t.Errorf("expected approved, got %s", inv.CampaignApprovalStatus)
	}
	inv.Events()
}

func TestStep7_CalculateFinalPrice(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, ptr(100), nil)
	inv.Events()

	finalPrice, err := entity.CalculateFinalPrice(500000, 20)
	if err != nil {
		t.Fatalf("calculate final price failed: %v", err)
	}
	if finalPrice != 400000 {
		t.Errorf("expected 400000, got %f", finalPrice)
	}

	err = inv.UpdatePrice(500000, finalPrice)
	if err != nil {
		t.Fatalf("update price failed: %v", err)
	}
	if inv.FinalPrice != 400000 {
		t.Errorf("expected 400000, got %f", inv.FinalPrice)
	}
	inv.Events()
}

func TestStep8_ActivateCampaign(t *testing.T) {
	now := time.Now()
	campaign, _ := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	campaign.Approve()
	campaign.Events()

	err := campaign.Activate()
	if err != nil {
		t.Fatalf("activate failed: %v", err)
	}
	if !campaign.IsActive {
		t.Error("expected campaign to be active")
	}
	campaign.Events()
}

func TestStep9_Sales(t *testing.T) {
	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, ptr(100), nil)
	inv.Events()

	if !inv.CanBeSold() {
		t.Fatal("expected item to be sellable before purchase")
	}

	err := inv.ValidatePurchase(1)
	if err != nil {
		t.Fatalf("validate purchase failed: %v", err)
	}

	err = inv.UpdateStock(49)
	if err != nil {
		t.Fatalf("update stock failed: %v", err)
	}
	if inv.InstantQty != 49 {
		t.Errorf("expected stock 49, got %d", inv.InstantQty)
	}
	inv.Events()
}

func TestStep10_CalculateCommission(t *testing.T) {
	comm, err := entity.NewCommission("comm-1", "prod-x200", "retail", 10, 100000, 1000000, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	comm.Events()

	commissionAmount, err := comm.Calculate(400000, 1)
	if err != nil {
		t.Fatalf("calculate commission failed: %v", err)
	}
	if commissionAmount != 40000 {
		t.Errorf("expected 40000, got %f", commissionAmount)
	}
	comm.Events()
}

func TestStep11_EndCampaign(t *testing.T) {
	now := time.Now()
	campaign, _ := entity.NewPromotion("promo-1", "Nowruz Auction", "",
		now.Add(24*time.Hour), now.Add(72*time.Hour), true, 20, true)
	campaign.Approve()
	campaign.Activate()
	campaign.Events()

	inv, _ := entity.NewInventory("inv-1", "store-1", "wh-1", "prod-x200", 500000, 50,
		"retail", "new", 1, ptr(100), nil)
	inv.UpdatePrice(500000, 400000)
	inv.Events()

	err := campaign.Deactivate()
	if err != nil {
		t.Fatalf("deactivate failed: %v", err)
	}
	if campaign.IsActive {
		t.Error("expected campaign to be inactive")
	}
	campaign.Events()

	err = inv.ResetPrice()
	if err != nil {
		t.Fatalf("reset price failed: %v", err)
	}
	if inv.FinalPrice != inv.BasePrice {
		t.Errorf("expected FinalPrice %f to equal BasePrice %f", inv.FinalPrice, inv.BasePrice)
	}
	inv.Events()
}
