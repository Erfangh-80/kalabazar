# End-to-End Test Scenario

> Seller `Saeed` registers the store `ElectronicsShop` and sells the product `X200 Wireless Headphones` at a 20% discount.

---

## Step 1: Register a store

**Input:**

```json
{
  "store_name": "Electronics Shop",
  "user_id": 42,
  "contact_phone": "09121234567"
}
```

**Output:**

- Record in `stores` with `id=1`, `status=active`
- Event `store.created` with `store_id=1`

---

## Step 1.5: Category Permission

**Input:**

```json
{
  "store_id": 1,
  "category_id": 7
}
```

**Output:**

- Record in `store_allowed_categories` with `status=pending`
- After admin approval: `status=approved`
- Event `store.category_allowed` with `store_id=1`, `category_id=7`

---

## Step 2: Create warehouse and connect

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

**Connect to store:**

```json
{
  "store_id": 1,
  "warehouse_id": 1,
  "relation_type": "primary"
}
```

**Output:**

- Record in `store_warehouse_links`
- Event `warehouse.linked_to_store` with `relation_type=primary`

---

## Step 3: Register goods

**Input:**

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
  "attributes": { "color": "black", "warranty_months": 18 }
}
```

**Output:**

- Record in `inventory` with `id=1`
- `final_price = 500000.00` (equal to base_price)
- `vendor_sale_status = active`, `system_sale_status = active`
- `promotion_status = pending`, `promotion_id = null`
- Event `inventory.item_created` with `inventory_id=1`

---

## Step 4: Record reference price

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

## Step 5: Create a campaign and connect to the product

**Campaign input:**

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

**Campaign output:**

- Record in `promotions` with `id=1`, `status=inactive`
- Event `promotion.campaign_created` with `promotion_id=1`

**Link campaign to product:**

```json
{
  "inventory_id": 1,
  "promotion_id": 1
}
```

**Connection output:**

- `inventory.promotion_id = 1`
- `inventory.promotion_linked` event
- `promotion.campaign_linked_to_product` event

---

## Step 6: Campaign approval

**Input:**

```json
{
  "inventory_id": 1,
  "promotion_status": "approved"
}
```

**Output:**

- `inventory.promotion_status = approved`
- `inventory.promotion_status_changed` event: `pending → approved`
- `promotion.campaign_approved` event

---

## Step 7: Calculate final price

**Calculate:**

```
final_price = 500,000 - (500,000 × 20 / 100) = 400,000
```

**Output:**

- `inventory.final_price = 400000.00`
- Event `inventory.price_updated`: `500000 → 400000`

---

## Step 8: Automatically activate the campaign

**Condition:** `NOW() >= start_at (2026-06-20T00:00:00Z)`

**Output:**

- `promotions.status = active`
- Event `promotion.campaign_activated`

---

## Step 9: Sales

**Sales input:**

```json
{
  "inventory_id": 1,
  "qty": 1
}
```

**Validation:**

- `vendor_sale_status = active` ✅
- `system_sale_status = active` ✅
- `instant_qty (50) >= 1` ✅
- `min_order_qty (1) <= 1 <= max_order_qty (null)` ✅

**Output:**

- `inventory.instant_qty = 49` (decrease by 1 unit)
- `inventory.stock_updated` event: `50 → 49`

---

## Step 10: Calculate commission

**Predefined commission rule:**

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

**Calculation:**

```
400,000 >= 100,000 ✅
400,000 <= 1,000,000 ✅
1 >= 1 ✅
Commission = 400,000 × 10% = 40,000
```

**Output:**

- Event `commission.calculated`: `commission_amount=40000`

---

## Step 11: End of campaign

**Condition:** `NOW() > end_at (2026-07-05T23:59:59Z)`

**Output:**

- `promotions.status = inactive`

- `inventory.final_price = 500000.00` (return to base_price)
- Event `promotion.campaign_deactivated`
- Event `inventory.price_updated`: `400000 → 500000`

---
