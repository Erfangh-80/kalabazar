**Role:** Seller

**Description:**
The seller registers a specific item in a specific warehouse for sale. This is the first step in offering a product on Marketplace. Each item includes price, inventory, sales status, and purchase terms.

**Main flow:**

1. The seller selects the store, warehouse, and desired product
2. Completes the product information
3. The system registers the item

**Default values:**

- Seller status: `active`
- System status: `active`
- Campaign approval status: `pending`

**Result:**

- The `inventory.item_created` event is emitted
- The item is available for sale (if both statuses are active)

---

## Scenario 2: Inventory update

**Role:** Seller / System

**Description:**
The inventory of the item in the warehouse may change for various reasons. The seller can adjust the inventory manually, or the system can automatically reduce the inventory after each sale.

**Main flow:**

1. The current stock of the item is known (`instant_qty`)
2. The vendor or the system records the new quantity
3. The system updates the inventory
4. If the stock reaches zero, the item cannot be sold

**Result:**

- The `inventory.stock_updated` event is emitted
- The `instant_qty` value changes

---

## Scenario 3: Sales status control (two-state)

**Role:** Vendor / System Administrator

**Description:**
The item has two independent states, each controlled by a separate role:

### Vendor status (`vendor_sale_status`)

The vendor can set their item to one of the following states:

- `active` — The item is active for sale
- `inactive` — The vendor has stopped selling
- `scheduled` — The sale is scheduled
- `draft` — The item is not yet completed and can be sold None

**Flow:**

1. Seller changes status
2. System updates
3. Event `inventory.item_activated` or `inventory.item_deactivated` is emitted

### System status (`system_sale_status`)

System administrator can block or unblock the item:

- `active` — system allows selling
- `inactive` — system has prohibited selling (e.g. due to violation)

**Flow:**

1. System administrator changes status
2. System updates
3. Event `inventory.system_blocked` or `inventory.system_unblocked` is emitted

**Final decision:**
The item can only be sold if **both statuses** are `active`.

---

## Scenario 4: Sales Timing

**Role:** Seller

**Description:**
The seller wants to sell the item only during a specific time frame. For example, a seasonal product or a limited sale.

**Main flow:**

1. The seller sets the start (`start_at`) and end (`end_at`) dates for the sale
2. The system keeps the item active during the specified period
3. Outside the period, the item is automatically deactivated

**Options:**

- Specify only the start time (active until further notice)
- Specify both start and end
- No scheduling (always active)

**Result:**

- The `inventory.sale_scheduled` event is emitted

---

## Scenario 5: Change in product price

**Role:** Seller

**Description:**
The seller may change the base price of the item or the system may calculate and update the final price based on discounts.

**Main flow:**

1. Seller enters new `base_price`
2. Or system calculates `final_price` based on `base_price` and campaign discount
3. System updates price

**Result:**

- `inventory.price_updated` event is emitted
