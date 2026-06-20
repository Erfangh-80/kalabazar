# Use Case: ActivateScheduledCampaigns

## Purpose

Campaigns are created with `status = inactive` and manual activation is not scalable for large numbers. Scheduler automates this task.

## Actor

**System** (System / Scheduler) — Cron Job

## Description

The Scheduler periodically finds campaigns with `start_at >= NOW()` and `status = inactive`. For each, it changes `status` to `active` and emits the `promotion.campaign_activated` event. If `is_countdown = true`, a countdown timer is displayed in the UI.

## Example

The time reaches `2026-06-20T00:00:00Z`. Scheduler finds the `Nowruz Auction` campaign and changes `status` to `active`. The `promotion.campaign_activated` event is emitted.
