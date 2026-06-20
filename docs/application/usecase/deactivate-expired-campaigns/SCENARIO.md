# Use Case: DeactivateExpiredCampaigns

## Purpose

Campaigns should be automatically deactivated after `end_at` or the discount will continue incorrectly. This Use Case deactivates expired campaigns and resets `final_price` to `base_price`.

## Actor

**System** (System / Scheduler) — Cron Job

## Description

The Scheduler periodically finds campaigns with `end_at >= NOW()` and `status = active`. For each one, it changes `status` to `inactive` and calls `FinalPriceCalculator` with `action = "reset"` to reset `final_price` to `base_price`. The `promotion.campaign_deactivated` and `inventory.price_updated` events are emitted. If `expire_sale_with_promotion = true`, `vendor_sale_status` will also change to `inactive`.

## Example

The time reaches `2026-07-05T23:59:59Z`. Scheduler deactivates the "Nowruz Auction" campaign. `final_price` of "X200 Headphones" will be changed from 400,000 to 500,000.
