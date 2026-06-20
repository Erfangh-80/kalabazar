# Use Case: DefineCommissionRule

## Purpose

Marketplace charges commission for each sale. This Use Case allows the admin to define commission rules for each item. Without this rule, the system doesn't know how to calculate its share.

## Actor

System Administrator (Admin)

## Description

The admin defines a commission rule for an `inventory_id`: percentage (`rate_percent`), price range (`min_price`, `max_price`), and minimum quantity (`min_qty`). Each item has only one commission rule (`UNIQUE`). The rule conditions are checked at the time of sale — if `final_price` is in the range and the quantity is greater than `min_qty`, the commission is calculated. The `commission.rule_defined` event is emitted.

## Example

Admin defines a 10% commission rule for "X200 Headphones" with `min_price=100,000` and `max_price=1,000,000`. The `commission.rule_defined` event is emitted.
