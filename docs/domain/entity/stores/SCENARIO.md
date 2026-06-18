**Description:**
A new seller wants to start working on Marketplace. He must first register his store in order to list his goods for sale. The store acts as the seller's main entity in the system and all subsequent activities (warehouse, inventory, sales) are connected to it.
Store creation is not immediately active. Every newly registered store must go through an admin review process before it becomes operational.

🧩 Store Lifecycle States
active - inactive - rejected
pending_review: Newly created store waiting for admin approval
active: Approved and operational store
inactive: Temporarily disabled store (can be reactivated)
rejected: Store rejected by admin and not allowed to operate

Scenario 1: Store Registration

Role: Seller

Description:

A seller submits a request to register a store. The store is created but cannot operate until it is approved by an administrator

**Main flow:**

Seller requests store registration
System creates store with status = pending_review
Store is linked to user_id
Admin is notified for review
Store remains non-operational until approval

---

## Scenario 2: Update store information

**Role:** Seller

**Description:**
The seller may need Has to change his store information after the initial registration. For example, register a new contact number, change the store name, or add new media.

**Main flow:**

1. Seller requests to update store information
2. System validates new fields
3. System updates store information

**Result:**

- `store.updated` event is emitted
- New information is visible immediately

---

## Scenario 3: Activating and deactivating the store

**Role:** Seller / System Administrator

**Description:**
The seller may want to temporarily stop selling without deleting the store. The system administrator may also deactivate the store in case of violation.

**Main flow (activation):**

1. Store is in `inactive` state
2. Seller requests activation
3. System changes status to `active`
4. `store.activated` event is emitted

**Main flow (deactivation):**

1. Store is in `active` state
2. Seller or admin requests deactivation
3. System changes status to `inactive`
4. `store.deactivated` event is emitted

**Impact:**

- In `inactive` state, no products can be sold from this store
- Store inventory and information are preserved
- Seller can reactivate later
