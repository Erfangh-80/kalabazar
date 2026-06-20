# Event Handler: OnStockUpdatedHandler → CommissionCalculator

## Purpose

After each successful sale, the system should calculate its share. This Event Handler automatically and asynchronously calculates the commission based on the `sale_commission` rule so as not to block the sale.

## Nature

Event Handler — consumes the `inventory.stock_updated` event.

## Description

This Handler is triggered when `inventory.stock_updated` is emitted with `reason = "sale"`. First, the commission rule for `inventory_id` is found, then the conditions are checked (`sale_model`, `min_price`, `max_price`, `min_qty`). If the conditions are met: `commission = sale_amount × rate%`. The `commission.calculated` event is emitted.

## Example

Sale of "X200 headphones" for 400,000 Tomans. 10% commission rule: `commission = 400,000 × 10% = 40,000`. The `commission.calculated` event is emitted with `commission_amount=40,000`.
