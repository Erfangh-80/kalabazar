# Use Case: RegisterStore

## Purpose

Every seller needs a store to operate in the Marketplace. This Use Case is the seller's entry point into the system — without a store, no further operations (warehouse, goods, sales) are possible.

## Actor

Seller

## Explanation

The seller requests to register a store. The system registers the store with the status `pending` and connects it to the seller's account. In this case, the store is not yet active — the seller cannot register or sell goods. The `store.created` event is emitted.

After this step, the administrator must review and approve the store (see the next Use Case).

## Example

Seller `Said` registers the store "Electronics Shop". The system creates the store with `id=1` and `status=pending` and the `store.created` event is emitted.
