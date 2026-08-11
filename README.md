# Architecture Decision Log

## ADR-001: Order Mutability After Payment

**Decision:** Orders become immutable after the first payment.

**Why:** Once money has been recorded against an order, changing the facts that determined the amount owed makes the historical payment data ambiguous or corruptible.

---

## ADR-002: Order Deletion

**Decision:** Orders may be deleted only before the first payment.

**Why:** Once payments exist, deleting the order would destroy the context for the payment history.

---