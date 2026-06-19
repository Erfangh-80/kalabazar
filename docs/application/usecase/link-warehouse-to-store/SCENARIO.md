# Use Case: LinkWarehouseToStore

## Purpose

Warehouse and store are independent entities — a warehouse can be connected to multiple stores, and a store can have multiple warehouses. This Use Case manages a Many-to-Many relationship.

## Actor

Seller

## Description

A seller connects a warehouse to a store with a specified `relation_type`:

- **`primary`:** The store's default warehouse for shipping goods
- **`secondary`:** A backup or supplementary warehouse

The connection is recorded in the `store_warehouse_links` table, and each `(store_id, warehouse_id)` pair can only be recorded once. The `warehouse.linked_to_store` event is emitted.

## Example

A seller connects the warehouse "Tehran Central Warehouse" to the store "Electronics Shop" with `relation_type = primary`. The system registers the connection and emits the `warehouse.linked_to_store` event.
