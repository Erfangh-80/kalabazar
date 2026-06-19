# Use Case: RequestCategoryAccess

## Purpose

After registering a store, a seller cannot register in every product category — some categories require special permissions. This Use Case manages the permission request process so that the system administrator can monitor the variety of products for each seller.

## Actor

Seller

## Explanation

A seller requests permission for a specific `category_id`. The system creates a record in `store_allowed_categories` with `status = pending`. Each `(store_id, category_id)` pair is unique — a store cannot have two separate requests for the same category. Until the administrator approves, the seller cannot register in this product category.

## Example

The store "ElectronicsShop" requests permission for the category "Digital Goods" (`category_id=7`). The system creates a record with `status=pending` and waits for the administrator to approve.
