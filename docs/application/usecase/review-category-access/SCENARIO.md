# Use Case: ReviewCategoryAccess

## Purpose

Category permission requests submitted by the seller must be reviewed by the administrator. This Use Case provides the administrator's decision-making mechanism — approval allows the seller to operate, and rejection prohibits the seller with a reason.

## Actor

Admin

## Explanation

The administrator reviews the request and makes a decision:

- **Approve:** `status = approved` ← The seller can register in this product category
- **Reject:** `status = rejected` ← The administrator can write a `support_note` (e.g. "Complete your documents")

## Example

The administrator approves the request of the store "ElectronicsShop" for the category "Digital Goods". The `status` changes to `approved` and the seller can register in this product category.
