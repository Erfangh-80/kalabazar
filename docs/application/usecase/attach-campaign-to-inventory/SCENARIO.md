# Use Case: AttachCampaignToInventory

## Purpose

A discount campaign alone is not enough — it must be specified which products it will apply to. This Use Case attaches the campaign to a specific product.

## Actor

Seller

## Description

The seller attaches the campaign to a product. Attachment mechanism: `inventory.promotion_id` is initialized. Each product can have only one `promotion_id` — if the product already has an active campaign, the new campaign will replace it. Two events are emitted from two separate Bounded Contexts: `inventory.promotion_linked` and `promotion.campaign_linked_to_product`.

## Example

The campaign "New Year's Auction" (`promotion_id=1`) is attached to the product "X200 Headphones" (`inventory_id=1`). `inventory.promotion_id = 1` and two events are emitted.
