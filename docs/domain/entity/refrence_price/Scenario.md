## Scenario 1: Registering the base price of the product

**Role:** Salesperson

**Description:**
The base price is determined when the product is registered in the `inventory` table. This price determines subsequent calculations such as final price discounts.

---

## Scenario 1: Registering the reference price

**Role:** System Administrator / Analyst

**Description:**
To analyze the market and make decisions about pricing, the system allows you to register reference prices from various sources. These prices are purely analytical and can be used directly in sales.

**Main flow:**

1. User selects the product of interest
2. Records the price observed in the market and its source
3. The system stores the reference price

**Features:**

- A product can be a reference price reference from different sources
- The reference price must be greater than zero
- Reference prices have a historical record (`created_at`)

---

## Scenario 3: Calculating the final price

**Role:** System

**Description:**
The price that the customer finally sees (`final_price`) is not necessarily equal to the base price. The final price must be calculated by combining different systems.

**Factors on the final price:**

1. **Base price** (`base_price`) — The source of the calculation
2. **Campaign discount** — If the product is associated with a campaign, the discount percentage is applied.

**How ​​to calculate:**

```
Base_price = Base_price - (Base_price * Discount percentage)
```
