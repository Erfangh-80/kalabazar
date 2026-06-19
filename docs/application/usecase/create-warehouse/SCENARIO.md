# Use Case: CreateWarehouse

## Purpose

A seller needs a warehouse to store and manage their goods. Each warehouse is a physical location for storing goods. This Use Case allows the seller to register one or more warehouses with (public/private) access.

## Actor

Seller

## Description

A seller registers a new warehouse. The warehouse can be `is_public = true` (customers can buy from it) or `is_public = false` (internal management only). After creation, the warehouse is not yet connected to any store — the connection will be done in the next Use Case. The `warehouse.created` event is emitted.

## Example

A seller registers a warehouse "Tehran Central Warehouse" with `is_public = true`. The warehouse is created with `id=1` and the `warehouse.created` event is emitted.
