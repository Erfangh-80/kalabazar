# Use Case: CreateCampaign

## Purpose

A seller needs to create discount campaigns to attract more customers and increase sales. This Use Case allows defining a campaign with a title, time frame, and specific settings.

## Actor

Seller / Admin

## Description

User creates a new discount campaign. The campaign is created with `status = inactive` — activation is done later by the Scheduler. Important settings:

- **`requires_approval`:** If `true`, admin approval is required before applying the discount
- **`is_countdown`:** Display a countdown in the UI
- **`expire_sale_with_promotion`:** If `true`, the product will be removed from sale after the campaign ends

The `promotion.campaign_created` event is emitted.

## Example

A seller creates a "Nowruz Auction" campaign with 20% discount for the period June 20 to July 5. The campaign is created with `id=1` and `status=inactive` and the `promotion.campaign_created` event is emitted.
