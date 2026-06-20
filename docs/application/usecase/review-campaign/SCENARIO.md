# Use Case: ReviewCampaign

## Purpose

If the campaign is created with `requires_approval = true`, it must be reviewed by the system administrator before the discount is applied. This Use Case provides an approval/rejection mechanism.

## Actor

System Administrator (Admin)

## Description

The administrator changes `inventory.promotion_status`:

- **Approve:** `promotion_status = approved` ← `FinalPriceCalculator` is activated and `final_price` is calculated with the discount
- **Reject:** `promotion_status = rejected` ← The item remains without the discount

The `inventory.promotion_status_changed` and `promotion.campaign_approved` or `promotion.campaign_rejected` events are emitted.

## Example

The administrator approves the `Nowruz Auction` campaign for `X200 Headphones`. `promotion_status` changes to `approved` and `final_price` decreases from 500,000 to 400,000.
