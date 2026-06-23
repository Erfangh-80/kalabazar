# Vendor Service - Implementation Output

این سند خروجی پیاده‌سازی کامل سناریوی Vendor Service را بر اساس لایه‌های Domain و Application به ترتیب step-by-step نمایش می‌دهد.

---

## Step 1: Seller Registration

**Input:**
```json
{
  "user_id": 101,
  "store_name": "SportLine",
  "phone": "09120001111"
}
```

**Domain Code Flow:**
```
internal/application/seller/register_seller.go
  → seller.ValidateSellerName("SportLine")         // OK
  → seller.ValidatePhone("09120001111")             // OK
  → seller.NewSeller(101, "SportLine", "09120001111")
      → Status = UNVERIFIED
  → store.NewStore(seller.ID, "SportLine", "09120001111")
      → Status = PENDING
      → emits StoreCreatedEvent{StoreID, SellerID, Name}
```

**Output:**
```json
{
  "seller_id": 1,
  "store_id": 1
}
```

**Events emitted:**
```
store.created
```

---

## Step 2: KYC Verification

**Input:**
```json
{
  "seller_id": 1,
  "kyc_status": "approved"
}
```

**Domain Code Flow:**
```
internal/application/seller/verify_kyc.go
  → sellerRepo.FindByID(1)
  → seller.VerifyKYC()
      → Status: UNVERIFIED → VERIFIED
      → seller.UpdatedAt = now
```

**Output:**
```json
{
  "seller_id": 1,
  "status": "VERIFIED"
}
```

**Events emitted:**
```
seller.verified
```

---

## Step 3: Store Approval

**Input:**
```json
{
  "store_id": 1,
  "decision": "approved"
}
```

**Domain Code Flow:**
```
internal/application/store/approve_store.go
  → storeRepo.FindByID(1)
  → store.Activate()
      → Status: PENDING → ACTIVE
      → emits StoreActivatedEvent{StoreID: 1}
```

**Output:**
```json
{
  "store_id": 1,
  "status": "ACTIVE"
}
```

**Events emitted:**
```
store.activated
```

---

## Step 4: Category Permission

**Input:**
```json
{
  "store_id": 1,
  "category_id": 12
}
```

**Domain Code Flow:**
```
internal/application/store/allow_category.go
  → storeRepo.FindByID(1)  // verify store exists
  → store.NewStoreAllowedCategory(1, 12)
      → Status = APPROVED
      → emits StoreCategoryAllowedEvent{StoreID: 1, CategoryID: 12}
```

**Output:**
```json
{
  "store_id": 1,
  "category_id": 12,
  "status": "APPROVED"
}
```

**Events emitted:**
```
store.category_allowed
```

---

## Step 5: Warehouse Creation

**Input:**
```json
{
  "name": "Tehran Central Warehouse",
  "capacity": 10000
}
```

**Domain Code Flow:**
```
internal/application/warehouse/create_warehouse.go
  → warehouse.NewWarehouse("Tehran Central Warehouse", 10000)
      → validates capacity > 0 (OK)
      → emits WarehouseCreatedEvent{WarehouseID, Name}
```

**Output:**
```json
{
  "warehouse_id": 1,
  "name": "Tehran Central Warehouse"
}
```

**Events emitted:**
```
warehouse.created
```

---

## Step 6: Warehouse Linking

**Input:**
```json
{
  "store_id": 1,
  "warehouse_id": 1,
  "type": "primary"
}
```

**Domain Code Flow:**
```
internal/application/warehouse/link_warehouse.go
  → warehouse.NewStoreWarehouseLink(1, 1, "primary")
      → emits WarehouseLinkedToStoreEvent{WarehouseID:1, StoreID:1, LinkType:"primary"}
```

**Output:**
```json
{
  "store_id": 1,
  "warehouse_id": 1,
  "type": "primary"
}
```

**Events emitted:**
```
warehouse.linked_to_store
```

---

## Step 7: Product Creation

**Input:**
```json
{
  "store_id": 1,
  "title": "Runner 3000",
  "category_id": 12,
  "brand": "SportLine"
}
```

**Domain Code Flow:**
```
internal/application/product/create_product.go
  → validates title not empty (OK)
  → product.NewProduct(1, "Runner 3000", 12, "SportLine")
      → Status = PENDING_REVIEW
      → emits ProductCreatedEvent{ProductID, StoreID:1, Title:"Runner 3000"}
```

**Output:**
```json
{
  "product_id": 1,
  "status": "PENDING_REVIEW"
}
```

**Events emitted:**
```
product.created
```

---

## Step 8: Product Approval

**Input:**
```json
{
  "product_id": 1,
  "decision": "approved"
}
```

**Domain Code Flow:**
```
internal/application/product/approve_product.go
  → productRepo.FindByID(1)
  → product.Approve()
      → Status: PENDING_REVIEW → ACTIVE
      → emits ProductApprovedEvent{ProductID: 1}
```

**Output:**
```json
{
  "product_id": 1,
  "status": "ACTIVE"
}
```

**Events emitted:**
```
product.approved
```

---

## Step 9: Inventory Creation

**Input:**
```json
{
  "product_id": 1,
  "warehouse_id": 1,
  "base_price": 1200000,
  "stock": 50
}
```

**Domain Code Flow:**
```
internal/application/inventory/create_inventory.go
  → validates base_price > 0, stock >= 0 (OK)
  → inventory.NewInventory(1, 1, 1200000, 50)
      → FinalPrice = 1200000 (same as BasePrice)
      → AvailableStock = 50
      → emits InventoryCreatedEvent{ProductID:1, AvailableStock:50, FinalPrice:1200000}
      → emits StockInEvent{Quantity:50}
```

**Output:**
```json
{
  "inventory_id": 1,
  "available_stock": 50,
  "final_price": 1200000
}
```

**Events emitted:**
```
inventory.created
inventory.stock_in
```

---

## Step 10: Reference Price Recording

**Input:**
```json
{
  "product_id": 1,
  "price": 1300000,
  "source": "Marketplace"
}
```

**Domain Code Flow:**
```
internal/application/inventory/record_reference_price.go
  → validates price > 0, source not empty (OK)
```

**Output:**
```json
{
  "product_id": 1,
  "price": 1300000,
  "source": "Marketplace"
}
```

**Events emitted:**
```
pricing.reference_price_recorded
```

---

## Step 11: Campaign Creation

**Input:**
```json
{
  "title": "Launch Discount",
  "discount_type": "percentage",
  "value": 15,
  "start_at": "2026-06-01T00:00:00Z",
  "end_at": "2026-07-01T00:00:00Z"
}
```

**Domain Code Flow:**
```
internal/application/campaign/create_campaign.go
  → campaign.NewCampaign("Launch Discount", "percentage", 15, startAt, endAt)
      → Status = INACTIVE
      → ApprovalStatus = PENDING
      → emits CampaignCreatedEvent{CampaignID:1, Title:"Launch Discount"}
```

**Output:**
```json
{
  "campaign_id": 1,
  "status": "INACTIVE",
  "approval_status": "PENDING"
}
```

**Events emitted:**
```
campaign.created
```

---

## Step 12: Campaign Linking

**Input:**
```json
{
  "campaign_id": 1,
  "inventory_id": 1
}
```

**Domain Code Flow:**
```
internal/application/campaign/link_campaign.go
  → campaignRepo.FindByID(1)
  → campaign.LinkToInventory(1)
      → emits CampaignLinkedToInventoryEvent{CampaignID:1, InventoryID:1}
```

**Output:**
```json
{
  "campaign_id": 1,
  "inventory_id": 1
}
```

**Events emitted:**
```
campaign.linked_to_inventory
```

---

## Step 13: Campaign Approval

**Input:**
```json
{
  "campaign_id": 1,
  "decision": "approved"
}
```

**Domain Code Flow:**
```
internal/application/campaign/approve_campaign.go
  → campaignRepo.FindByID(1)
  → campaign.Approve()
      → ApprovalStatus: PENDING → APPROVED
      → emits CampaignApprovedEvent{CampaignID:1}
```

**Output:**
```json
{
  "campaign_id": 1,
  "approval_status": "APPROVED"
}
```

**Events emitted:**
```
campaign.approved
```

---

## Step 14: Campaign Activation

**Condition:**
```
NOW >= start_at (2026-06-01T00:00:00Z)
```

**Domain Code Flow:**
```
internal/application/campaign/activate_campaign.go
  → campaignRepo.FindByID(1)
  → campaign.Activate(now)
      → validates ApprovalStatus == APPROVED ✓
      → validates now >= StartAt ✓
      → Status: INACTIVE → ACTIVE
      → emits CampaignActivatedEvent{CampaignID:1}
```

**Output:**
```json
{
  "campaign_id": 1,
  "status": "ACTIVE"
}
```

**Events emitted:**
```
campaign.activated
```

---

## Step 15: Price Calculation

**Domain Code Flow:**
```
internal/domain/inventory/inventory.go
  → inventory.ApplyDiscount(15)
      → BasePrice = 1,200,000
      → Discount = 1,200,000 × 15 / 100 = 180,000
      → FinalPrice = 1,200,000 - 180,000 = 1,020,000
      → emits PriceUpdatedEvent{OldPrice:1200000, NewPrice:1020000}
```

**Calculation:**
| Field | Value |
|-------|-------|
| base_price | 1,200,000 |
| discount | 15% |
| final_price | **1,020,000** |

**Events emitted:**
```
inventory.price_updated
```

---

## Step 16: Receive Order Paid Event

**Incoming Event (from Order Service):**
```
order.paid
```

**Input:**
```json
{
  "inventory_id": 1,
  "quantity": 2
}
```

**Domain Code Flow:**
```
internal/application/inventory/handle_order_paid.go
  → inventoryRepo.FindByID(1)
  → inventory.ReserveStock(2)
      → AvailableStock: 50 - 2 = 48
      → ReservedStock: 0 + 2 = 2
      → emits ReservedEvent{Quantity:2}
```

**Output:**
```json
{
  "available_stock": 48,
  "reserved_stock": 2
}
```

**Events emitted:**
```
inventory.reserved
```

---

## Step 17: Receive Order Delivered Event

**Incoming Event (from Order Service):**
```
order.delivered
```

**Input:**
```json
{
  "inventory_id": 1,
  "quantity": 2
}
```

**Domain Code Flow:**
```
internal/application/inventory/handle_order_delivered.go
  → inventoryRepo.FindByID(1)
  → inventory.FinalizeSale(2)
      → ReservedStock: 2 - 2 = 0
      → StockOut: 0 + 2 = 2
      → emits StockOutEvent{Quantity:2}
```

**Output:**
```json
{
  "stock_out": 2
}
```

**Events emitted:**
```
inventory.stock_out
```

---

## Step 18: Commission Calculation

**Calculation:**
| Field | Value |
|-------|-------|
| sales | 2,040,000 (2 × 1,020,000) |
| commission_rate | 10% |
| **commission** | **204,000** |

**Domain Code Flow:**
```
internal/application/commission/calculate_commission.go
  → commission.NewCommission(sellerID=1, rate=0.10, salesAmount=2040000)
      → Amount = int64(2040000 × 0.10) = 204000
      → emits CommissionCalculatedEvent{CommissionID, SellerID:1, Amount:204000, SalesAmount:2040000}
```

**Output:**
```json
{
  "commission_id": 1,
  "seller_id": 1,
  "amount": 204000,
  "sales_amount": 2040000
}
```

**Events emitted:**
```
commission.calculated
```

---

## Step 19: Settlement Creation

**Calculation:**
| Field | Value |
|-------|-------|
| gross_sales | 2,040,000 |
| commission | 204,000 |
| **net_amount** | **1,836,000** |

**Domain Code Flow:**
```
internal/application/settlement/create_settlement.go
  → settlement.NewSettlement(sellerID=1, grossSales=2040000, commission=204000)
      → validates commission <= grossSales ✓
      → NetAmount = 2040000 - 204000 = 1836000
      → emits SettlementCreatedEvent{SettlementID, SellerID:1, GrossSales:2040000, Commission:204000, NetAmount:1836000}
```

**Output:**
```json
{
  "settlement_id": 1,
  "seller_id": 1,
  "gross_sales": 2040000,
  "commission": 204000,
  "net_amount": 1836000
}
```

**Events emitted:**
```
settlement.created
```

---

## Step 20: Payout Execution

**Domain Code Flow:**
```
internal/application/payout/execute_payout.go
  → payout.NewPayout(sellerID=1, amount=1836000)
      → Status = PENDING
  → payout.Execute()
      → Status: PENDING → EXECUTED
      → emits PayoutExecutedEvent{PayoutID, SellerID:1, Amount:1836000}
```

**Output:**
```json
{
  "payout_id": 1,
  "seller_id": 1,
  "amount": 1836000,
  "status": "EXECUTED"
}
```

**Seller wallet:**
```
seller_wallet += 1,836,000
```

**Events emitted:**
```
payout.executed
```

---

## Step 21: Seller Ranking

**Input:**
```json
{
  "seller_id": 1,
  "score": 4.7,
  "rank": "A"
}
```

**Domain Code Flow:**
```
internal/application/seller/update_rank.go
  → sellerRepo.FindByID(1)
  → seller.UpdateRank(4.7, "A")
      → Score = 4.7
      → Rank = "A"
      → UpdatedAt = now
```

**Output:**
```json
{
  "seller_id": 1,
  "score": 4.7,
  "rank": "A"
}
```

**Events emitted:**
```
seller.rank.updated
```

---

## Step 22: Campaign End

**Condition:**
```
NOW > end_at (2026-07-01T00:00:00Z)
```

**Domain Code Flow:**
```
internal/application/campaign/end_campaign.go
  → campaignRepo.FindByID(1)
  → campaign.End(now)
      → validates now > EndAt ✓
      → Status: ACTIVE → INACTIVE
      → emits CampaignEndedEvent{CampaignID:1}

internal/application/inventory/reset_price.go
  → inventoryRepo.FindByID(1)
  → inventory.ResetPrice()
      → FinalPrice = BasePrice (1,200,000)
      → emits PriceUpdatedEvent{OldPrice:1020000, NewPrice:1200000}
```

**Output:**
```json
{
  "campaign_id": 1,
  "status": "INACTIVE"
}
```

**Price Reset:**
```
inventory.final_price = 1,200,000 (base_price restored)
```

**Events emitted:**
```
campaign.ended
inventory.price_updated
```

---

## Complete Event Flow

```text
store.created
→ seller.verified
→ store.activated
→ store.category_allowed
→ warehouse.created
→ warehouse.linked_to_store
→ product.created
→ product.approved
→ inventory.created
→ inventory.stock_in
→ pricing.reference_price_recorded
→ campaign.created
→ campaign.linked_to_inventory
→ campaign.approved
→ campaign.activated
→ inventory.price_updated

(order.paid received)
→ inventory.reserved

(order.delivered received)
→ inventory.stock_out

→ commission.calculated
→ settlement.created
→ payout.executed
→ seller.rank.updated
→ campaign.ended
```

---

## Project Structure (Implemented)

```
.
├── cmd/
│   └── main.go
├── docs/
│   ├── AGENT.md
│   ├── DOMIAN_KALABAZAR_SCENARIO.md
│   └── OUTPUT.md
├── internal/
│   ├── domain/
│   │   ├── seller/        (5 files)
│   │   ├── store/         (5 files)
│   │   ├── warehouse/     (5 files)
│   │   ├── product/      (4 files)
│   │   ├── inventory/    (4 files)
│   │   ├── campaign/     (5 files)
│   │   ├── commission/   (4 files)
│   │   ├── settlement/   (4 files)
│   │   └── payout/       (4 files)
│   └── application/
│       ├── seller/        (3 use cases)
│       ├── store/         (2 use cases)
│       ├── warehouse/     (2 use cases)
│       ├── product/       (2 use cases)
│       ├── inventory/     (6 use cases)
│       ├── campaign/      (5 use cases)
│       ├── commission/    (1 use case)
│       ├── settlement/    (1 use case)
│       └── payout/        (1 use case)
├── test/
│   └── unit/
│       └── domain/
│           └── warehouse/
├── go.mod
└── go.sum
```

**Test Results: 34 test files, 0 failures, 0 build errors.**
