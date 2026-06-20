# Use Case: RecordReferencePrice

## Purpose

To analyze the market and make pricing decisions, the system stores competitor prices. These prices do not directly affect sales — they are used only for analytical reports and dashboards.

## Actor

Analyst / System Administrator (Analyst / Admin)

## Explanation

The user records the observed price of a product from an external source. A product can have multiple reference prices from different sources, all of which have a historical record. The reference price must be greater than zero. The `pricing.reference_price_recorded` event is emitted.

## Example

The analyst observes and records the price of the product "X200 Headphones" on Digikala for 550,000 Tomans. The system stores the record in `reference_prices` and emits the `pricing.reference_price_recorded` event.
