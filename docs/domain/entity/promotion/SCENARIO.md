## Scenario 1: Create a discount campaign

**Role:** Seller / System Administrator

**Description:**
A seller or system administrator wants to create a discount campaign to attract more customers. The campaign includes a title, a time period, and its own settings.

**Main flow:**

1. User requests to create a campaign
2. Completes the campaign information
3. The system records the campaign

---

## Scenario 2: Link the campaign to the product

**Role:** Seller

**Description:**
After creating the campaign, the seller must specify which products this campaign will apply to. A campaign can be linked to multiple products, but each product can only have one campaign at a time.

**Main flow:**

1. Seller selects the desired campaign and product
2. System sets `promotion_id` to the product in the `inventory` table
3. If the product already has another campaign, it must be replaced

---

## Scenario 3: Approve or reject the campaign

**Role:** System Administrator

**Description:**
If the campaign is created with `requires_approval = true`, it must be reviewed and approved by the system administrator. Until approval, the campaign will not be applied to the product.

**Main flow:**

1. The manager checks the campaign
2. If approved, `inventory.promotion_status` changes to `approved`
3. If rejected, `inventory.promotion_status` changes to `rejected`

---

## Scenario 4: Campaign activation and deactivation based on time

**Role:** System (scheduler)

**Description:**
Campaigns are automatically activated and deactivated based on their time frame.

**Main flow (activation):**

1. The current time has reached `start_at`
2. The system activates the campaign
3. The discount is applied to the connected products
4. The `promotion.campaign_activated` event is emitted
