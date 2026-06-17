This document describes a complete business flow for an e-commerce system.
It is written as a **test scenario specification (TDD input)**.

Each step represents a business behavior that must be covered by tests before implementation.

---

# 🏬 Step 1 — Register a Store

### Scenario

A seller named Saeed registers a new store:

- Store name: "Electronics Shop"
- Initial status: active

### Expected Behavior

- Store is successfully created
- Store status is set to active
- A `store.created` event is emitted

---

# 🏢 Step 2 — Create Warehouse and Link to Store

### Scenario

Saeed creates a warehouse:

- Warehouse name: "Tehran Central Warehouse"
- Access type: public

Then he links this warehouse to his store as the primary warehouse.

### Expected Behavior

- Warehouse is successfully created
- A `warehouse.created` event is emitted
- Warehouse is linked to the store
- A `warehouse.linked_to_store` event is emitted

---

# 📦 Step 3 — Register Inventory Item

### Scenario

Saeed registers a product:

- Name: "X200 Wireless Headphones"
- Base price: 500,000 Tomans
- Instant quantity: 50 units
- Sale model: retail
- Vendor sale status: active

### Expected Behavior

- Product is successfully created in inventory
- Stock is initialized to 50
- Product is available for sale
- `inventory.item_created` event is emitted

---

# 💰 Step 4 — Register Reference Market Price

### Scenario

A market analyst records competitor pricing:

- Product: "X200 Wireless Headphones"
- Competitor price: 550,000 Tomans (Digikala)

### Expected Behavior

- Reference price is stored for analysis
- `pricing.reference_price_recorded` event is emitted

---

# 🎯 Step 5 — Create Discount Campaign

### Scenario

Saeed creates a campaign:

- Name: "Nowruz Auction"
- Discount: 20%
- Requires approval: true

The campaign is linked to product "X200 Wireless Headphones"

### Expected Behavior

- Campaign is created with status: inactive
- Campaign is linked to product
- `promotion.campaign_created` event is emitted
- `inventory.promotion_linked` event is emitted
- `promotion.campaign_linked_to_product` event is emitted

---

# ✅ Step 6 — Campaign Approval

### Scenario

An administrator approves the campaign.

### Expected Behavior

- Campaign status changes to approved
- `promotion.campaign_approved` event is emitted
- `inventory.promotion_status_changed` event is emitted

---

# 🔄 Step 7 — Final Price Calculation

### Scenario

System calculates final product price:

- Base price: 500,000
- Discount: 20%

### Expected Behavior

- Final price is calculated: 400,000 Tomans
- Product price is updated
- `inventory.price_updated` event is emitted

---

# ⏰ Step 8 — Campaign Activation (Time-Based)

### Scenario

At campaign start time, system activates the campaign.

### Expected Behavior

- Campaign becomes active
- Discount is applied to product
- `promotion.campaign_activated` event is emitted

---

# 🛒 Step 9 — Successful Purchase

### Scenario

A customer purchases:

- Product: "X200 Wireless Headphones"
- Price paid: 400,000 Tomans

### Expected Behavior

- Stock is reduced (50 → 49)
- Purchase is completed successfully
- `inventory.stock_updated` event is emitted

---

# 💸 Step 10 — Commission Calculation

### Scenario

System applies commission rule:

- Rate: 10%
- Applicable range: 100,000 – 1,000,000 Tomans
- Sale price: 400,000 Tomans

### Expected Behavior

- Commission is calculated: 40,000 Tomans
- `commission.calculated` event is emitted

---

# 🔚 Step 11 — Campaign End

### Scenario

At campaign end time:

- Campaign is automatically deactivated
- Product pricing returns to base price

If `expire_sale_with_promotion = true`:

- Product is removed from sale

### Expected Behavior

- Campaign is deactivated
- `promotion.campaign_deactivated` event is emitted
- Product price resets to base price
- Optional removal from sale is applied based on rule
