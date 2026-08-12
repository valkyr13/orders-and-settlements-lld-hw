# Architecture Decision Log

## 001: Order Mutability After Payment

**Decision:** Orders become immutable after the first payment.

**Why:** Once money has been recorded against an order, changing the facts that determined the amount owed makes the historical payment data ambiguous or corruptible.

---

## 002: Order Deletion

**Decision:** Orders may be deleted only before the first payment.

**Why:** Once payments exist, deleting the order would destroy the context for the payment history.

---
## 003: Decomposition of Problem Statement`

| Responsibility     | Owns                                          | Depends on                      |
| ------------------ | --------------------------------------------- | ------------------------------- |
| **Authentication** | User identity/credentials/session             | Credentials                     |
| **Authorization**  | Resource ownership rules                      | Authenticated user + resource   |
| **Order**          | Order definition, line items, due date, total | Owner                           |
| **Payment**        | Payment records + payment validation          | Order                           |
| **Status**         | Status derivation rules                       | Order + payments + current date |
| **Dashboard**      | Presentation/query of order information       | Order + payment/status data     |

### 004: Application Architecture

Chosen: Modular monolith
Why: Appropriate scale and time constraint; keeps tightly coupled order/payment transactions simple.
Alternative: Microservices
Trade-off: Sacrifices independent scaling for simplicity and correctness.

### 005: Persistence

Chosen: PostgreSQL / relational database
Why: Strong relational model and transactions fit order/payment correctness requirements.
Alternative: NoSQL / in-memory
Trade-off: More schema management, but substantially stronger consistency guarantees.

