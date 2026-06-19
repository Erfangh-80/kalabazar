# Use Case: ReviewStore

## Purpose

After registration, the store is in `pending` status — it is not yet active and the seller cannot use it. The system administrator must review the store to allow it to operate. This prevents the registration of invalid or suspicious stores.

## Actor

System Administrator (Admin)

## Explanation

The administrator reviews the store registration request and makes a decision:

- **Approval:** `status` changes to `active` — the seller can register and sell inventory and products. The `store.activated` event is emitted.
- **Rejection:** `status` changes to `rejected` — the seller cannot use the store. The `store.rejected` event is emitted.

Until the administrator reviews it, the store remains in `pending` status and no virtual operations are possible on it.

## Example

The system administrator reviews and approves the "Electronics Shop" store. The `status` changes to `active` and the `store.activated` event is emitted.
