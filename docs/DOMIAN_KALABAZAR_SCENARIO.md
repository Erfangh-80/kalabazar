````md
# Vendor Service - End To End Scenario

> فروشنده "SportLine" وارد مارکت‌پلیس می‌شود، فروشگاه ایجاد می‌کند، محصول ثبت می‌کند، کمپین اجرا می‌کند و پس از فروش موفق تسویه دریافت می‌کند.

---

# 1. Seller Registration

Input:

```json
{
  "user_id": 101,
  "store_name": "SportLine",
  "phone": "09120001111"
}
```
````

Output:

- Store Created
- Store Status = Pending

Events:

```text
store.created
```

---

# 2. KYC Verification

Input:

```json
{
  "seller_id": 1,
  "kyc_status": "approved"
}
```

Output:

```text
seller.status = VERIFIED
```

Events:

```text
seller.verified
```

---

# 3. Store Approval

Input:

```json
{
  "store_id": 1,
  "decision": "approved"
}
```

Output:

```text
store.status = ACTIVE
```

Events:

```text
store.activated
```

---

# 4. Category Permission

Input:

```json
{
  "store_id": 1,
  "category_id": 12
}
```

Output:

```text
store_allowed_category.created
status = APPROVED
```

Events:

```text
store.category_allowed
```

---

# 5. Warehouse Creation

Input:

```json
{
  "name": "Tehran Central Warehouse",
  "capacity": 10000
}
```

Output:

```text
warehouse.created
```

Events:

```text
warehouse.created
```

---

# 6. Warehouse Linking

Input:

```json
{
  "store_id": 1,
  "warehouse_id": 1,
  "type": "primary"
}
```

Output:

```text
warehouse linked to store
```

Events:

```text
warehouse.linked_to_store
```

---

# 7. Product Creation

Input:

```json
{
  "store_id": 1,
  "title": "Runner 3000",
  "category_id": 12,
  "brand": "SportLine"
}
```

Output:

```text
product.status = PENDING_REVIEW
```

Events:

```text
product.created
```

---

# 8. Product Approval

Input:

```json
{
  "product_id": 10,
  "decision": "approved"
}
```

Output:

```text
product.status = ACTIVE
```

Events:

```text
product.approved
```

---

# 9. Inventory Creation

Input:

```json
{
  "product_id": 10,
  "warehouse_id": 1,
  "base_price": 1200000,
  "stock": 50
}
```

Output:

```text
available_stock = 50
final_price = 1200000
```

Events:

```text
inventory.created
inventory.stock_in
```

---

# 10. Reference Price Recording

Input:

```json
{
  "product_id": 10,
  "price": 1300000,
  "source": "Marketplace"
}
```

Output:

```text
reference_price.created
```

Events:

```text
pricing.reference_price_recorded
```

---

# 11. Campaign Creation

Input:

```json
{
  "title": "Launch Discount",
  "discount_type": "percentage",
  "value": 15,
  "start_at": "...",
  "end_at": "..."
}
```

Output:

```text
campaign.status = INACTIVE
```

Events:

```text
campaign.created
```

---

# 12. Campaign Linking

Input:

```json
{
  "campaign_id": 1,
  "inventory_id": 1
}
```

Output:

```text
inventory linked to campaign
```

Events:

```text
campaign.linked_to_inventory
```

---

# 13. Campaign Approval

Input:

```json
{
  "campaign_id": 1,
  "decision": "approved"
}
```

Output:

```text
campaign.approval_status = APPROVED
```

Events:

```text
campaign.approved
```

---

# 14. Campaign Activation

Condition:

```text
NOW >= start_at
```

Output:

```text
campaign.status = ACTIVE
```

Events:

```text
campaign.activated
```

---

# 15. Price Calculation

Calculation:

```text
base_price = 1,200,000

discount = 15%

final_price = 1,020,000
```

Output:

```text
inventory.final_price updated
```

Events:

```text
inventory.price_updated
```

---

# 16. Receive Order Paid Event

Source:

```text
Order Service
```

Incoming Event:

```text
order.paid
```

Action:

```text
reserve stock
```

Output:

```text
available_stock = 48
reserved_stock = 2
```

Events:

```text
inventory.reserved
```

---

# 17. Receive Order Delivered Event

Source:

```text
Order Service
```

Incoming Event:

```text
order.delivered
```

Action:

```text
finalize sale
decrease reserved stock
perform stock out
```

Output:

```text
stock_out = 2
```

Events:

```text
inventory.stock_out
```

---

# 18. Commission Calculation

Rule:

```text
commission_rate = 10%
```

Calculation:

```text
sales = 2,040,000

commission = 204,000
```

Events:

```text
commission.calculated
```

---

# 19. Settlement Creation

Output:

```text
gross_sales = 2,040,000

commission = 204,000

net_amount = 1,836,000
```

Events:

```text
settlement.created
```

---

# 20. Payout Execution

Output:

```text
seller_wallet += 1,836,000
```

Events:

```text
payout.executed
```

---

# 21. Seller Ranking

Factors:

- Sales Volume
- SLA
- Return Rate
- Customer Satisfaction

Output:

```text
score = 4.7

rank = A
```

Events:

```text
seller.rank.updated
```

---

# 22. Campaign End

Condition:

```text
NOW > end_at
```

Output:

```text
campaign.status = INACTIVE

inventory.final_price = base_price
```

Events:

```text
campaign.ended
inventory.price_updated
```

---

# Event Flow

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

````

نکته مهم معماری:

- `order.created`
- `payment.completed`
- `shipment.created`
- `shipment.delivered`

نباید داخل Vendor Domain باشند.

Vendor فقط Consumer این Eventهاست و معمولاً این دو Event برایش کافی هستند:

```text
order.paid
order.delivered
````
