# Use Case: AddInventoryItem

## Purpose

The most important Use Case in the sales flow. A seller lists a specific product for sale in a specific store and warehouse — including price, inventory, sale status, and purchase terms.

## Actor

Seller

## Description

The seller registers the item. The system validates that the store is active, the warehouse is connected to the store, and the seller has category permissions. The item is registered with two independent statuses: `vendor_sale_status` (vendor-controlled) and `system_sale_status` (system-controlled). The item can only be sold if **both** are `active`. `final_price` is set to `base_price` (no discount). The `inventory.item_created` event is emitted.

## Example

The seller registers "X200 Wireless Headphones" with `base_price=500,000` and `instant_qty=50` in the store "ElectronicsShop". The item is saved with `id=1` and `final_price=500,000` and the `inventory.item_created` event is emitted.
