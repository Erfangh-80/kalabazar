# End-to-End Test Scenario

> Seller `Saeed` registers the store `ElectronicsShop`, gets approval from the marketplace, receives permission to sell in a category, creates inventory, joins a promotional campaign, and sells the product `X200 Wireless Headphones` with a 20% discount.

---

## Step 1: Register a Store

**Input:**

```json
{
  "store_name": "Electronics Shop",
  "user_id": 42,
  "contact_phone": "09121234567"
}
```

**Output:**

- Record in `stores` with `id=1`, `status=pending_review`
- Event `store.created` with `store_id=1`
- Store is not operational yet

---

## Step 1.1: Store Approval

**Input:**

```json
{
  "store_id": 1,
  "decision": "approve"
}
```

**Output:**

- `stores.status = active`
- Event `store.activated` with `store_id=1`
- Store becomes operational

---

## Step 1.2: Category Permission Request

**Input:**

```json
{
  "store_id": 1,
  "category_id": 7
}
```

**Output:**

- Record in `store_allowed_categories`
- `status = pending`
- Event `store.category_requested`

---

## Step 1.3: Category Approval

**Input:**

```json
{
  "store_id": 1,
  "category_id": 7,
  "decision": "approve"
}
```

**Output:**

- `store_allowed_categories.status = approved`
- Event `store.category_allowed`
- Store is allowed to sell products in category `7`

---

## Step 2: Create Warehouse and Connect

**Input:**

```json
{
  "warehouse_name": "Tehran Central Warehouse",
  "is_public": true
}
```

**Output:**

- Record in `warehouses` with `id=1`
- Event `warehouse.created` with `warehouse_id=1`

### Connect to Store

**Input:**

```json
{
  "store_id": 1,
  "warehouse_id": 1,
  "relation_type": "primary"
}
```

**Output:**

- Record in `store_warehouse_links`
- Event `warehouse.linked_to_store`
- `relation_type = primary`

---

## Step 3: Register Goods

### Validation

- Store status = `active` ✅
- Category permission status = `approved` ✅
- Warehouse linked to store ✅

### Input

```json
{
  "store_id": 1,
  "warehouse_id": 1,
  "product_id": 99,
  "sale_model": "retail",
  "base_price": 500000.0,
  "instant_qty": 50,
  "condition": "new",
  "min_order_qty": 1,
  "attributes": {
    "color": "black",
    "warranty_months": 18
  }
}
```

### Output

- Record in `inventory` with `id=1`
- `final_price = 500000.00`
- `vendor_sale_status = active`
- `system_sale_status = active`
- `promotion_status = pending`
- `promotion_id = null`
- Event `inventory.item_created` with `inventory_id=1`

---

## Step 4: Record Reference Price

**Input:**

```json
{
  "product_id": 99,
  "price": 550000.0,
  "source": "DigiKala"
}
```

**Output:**

- Record in `reference_prices`
- Event `pricing.reference_price_recorded`

---

## Step 5: Create Campaign

### Campaign Input

```json
{
  "title": "Nowruz Auction",
  "start_at": "2026-06-20T00:00:00Z",
  "end_at": "2026-07-05T23:59:59Z",
  "requires_approval": true,
  "is_countdown": true,
  "expire_sale_with_promotion": false
}
```

### Campaign Output

- Record in `promotions`
- `id = 1`
- `status = inactive`
- Event `promotion.campaign_created`

---

## Step 5.1: Connect Campaign to Product

**Input:**

```json
{
  "inventory_id": 1,
  "promotion_id": 1
}
```

**Output:**

- `inventory.promotion_id = 1`
- Event `inventory.promotion_linked`
- Event `promotion.campaign_linked_to_product`

---

## Step 6: Campaign Approval

**Input:**

```json
{
  "inventory_id": 1,
  "promotion_status": "approved"
}
```

**Output:**

- `inventory.promotion_status = approved`
- Event `inventory.promotion_status_changed`
- Event `promotion.campaign_approved`

---

## Step 7: Calculate Final Price

### Calculation

```text
final_price = 500,000 - (500,000 × 20 / 100)
final_price = 400,000
```

### Output

- `inventory.final_price = 400000.00`
- Event `inventory.price_updated`
- Previous price: `500000`
- New price: `400000`

---

## Step 8: Automatically Activate Campaign

### Condition

```text
NOW() >= start_at
```

### Output

- `promotions.status = active`
- Event `promotion.campaign_activated`

---

## Step 9: Sales

### Sales Input

```json
{
  "inventory_id": 1,
  "qty": 1
}
```

### Validation

- Store status = `active` ✅
- Category permission status = `approved` ✅
- Vendor sale status = `active` ✅
- System sale status = `active` ✅
- Inventory quantity available ✅
- `instant_qty (50) >= 1` ✅
- `min_order_qty (1) <= 1 <= max_order_qty (null)` ✅

### Output

- `inventory.instant_qty = 49`
- Event `inventory.stock_updated`
- Previous quantity: `50`
- New quantity: `49`

---

## Step 10: Calculate Commission

### Commission Rule

```json
{
  "inventory_id": 1,
  "category_commission_rule_id": 5,
  "sale_model": "retail",
  "rate_percent": 10.0,
  "min_price": 100000.0,
  "max_price": 1000000.0,
  "min_qty": 1
}
```

### Validation

```text
400,000 >= 100,000 ✅
400,000 <= 1,000,000 ✅
1 >= 1 ✅
```

### Calculation

```text
Commission = 400,000 × 10%
Commission = 40,000
```

### Output

- Event `commission.calculated`
- `commission_amount = 40000`

---

## Step 11: End of Campaign

### C
