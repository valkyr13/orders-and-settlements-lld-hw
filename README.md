# Orders & Settlements — Local Setup

## Prerequisites

Install:

* Docker
* Docker Compose
* Node.js 18+
* npm

Go is **not required** to run the backend because it runs inside Docker.

---

## Project Structure

```text
orders-and-settlements-lld-hw/
├── backend/
│   ├── main.go
│   ├── routes/
│   ├── internal/
│   ├── migrations/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── docker-compose.yml
│
└── frontend/
    ├── src/
    ├── index.html
    ├── package.json
    └── vite.config.js
```

---

# 1. Start the Backend

Open a terminal:

```bash
cd orders-and-settlements-lld-hw/backend
```

Start the backend and PostgreSQL:

```bash
docker compose up --build
```

This starts:

```text
PostgreSQL
    ↓
Gin Backend
```

The backend should be available at:

```text
http://localhost:8080
```

### Health Check

```bash
curl http://localhost:8080/health
```

Expected:

```json
{
  "status": "ok"
}
```

Keep this terminal running.

---

# 2. Start the Frontend

Open a **second terminal**.

```bash
cd orders-and-settlements-lld-hw/frontend
```

Install dependencies:

```bash
npm install
```

Start Vite:

```bash
npm run dev
```

Vite will show something similar to:

```text
Local: http://localhost:5173/
```

Open:

```text
http://localhost:5173
```

---

# 3. Verify Authentication

The frontend communicates with:

```text
Frontend
http://localhost:5173
        ↓
Backend
http://localhost:8080
```

### Create an account

From the UI:

```text
Create Account
    ↓
Enter email + password
    ↓
Sign Up
```

### Login

Then:

```text
Login
    ↓
Enter credentials
    ↓
Login
```

The backend creates a `session_id` cookie.

---

# 4. Verify the Session

After logging in, the frontend calls:

```text
GET /auth/me
```

You can also test the backend directly:

```bash
curl -i \
  -b cookies.txt \
  http://localhost:8080/auth/me
```

A successful response should be:

```text
200 OK
```

---

# 5. Development Workflow

Whenever you change **backend code**:

```bash
docker compose up --build
```

Whenever you change **frontend code**, Vite automatically reloads the application.

Frontend:

```bash
npm run dev
```

Backend:

```bash
docker compose up --build
```

---

# 6. Stop Everything

### Backend

Press:

```text
Ctrl+C
```

or:

```bash
docker compose down
```

### Frontend

Press:

```text
Ctrl+C
```

---

# 7. Reset the Database

⚠️ This deletes all PostgreSQL data.

```bash
docker compose down -v
docker compose up --build
```

The database will be recreated and the initialization migrations will run again.

---

## Quick Start

Once the project is already set up, running the application requires two terminals:

### Terminal 1 — Backend

```bash
cd backend
docker compose up --build
```

### Terminal 2 — Frontend

```bash
cd frontend
npm run dev
```

Then open:

```text
http://localhost:5173
```



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
Why: Strong relational model and We don't need NoSQL's typical strengths like highly flexible schemas or massive distributed write scaling.
Alternative: NoSQL / in-memory
Trade-off: More schema management, but substantially stronger consistency guarantees.

### 006: Payment Concurrency

Decision: Protect the specific order row during payment validation and creation.
Chosen: Row-level protection inside a transaction.
Why: Prevents concurrent payments for the same order from bypassing the over-payment invariant.
Alternative: Atomic database constraint/update without explicit row locking.
Trade-off: Slightly more locking, but simpler correctness reasoning for this assignment.


### 007: Currency

Decision: Use a single fixed currency for the application. Assuming Dollar as currency
Why: Removes ambiguity while avoiding unnecessary multi-currency complexity.
Alternative: Support currency per order.
Trade-off: Less flexibility, significantly simpler financial model.

### 008: Money Representation

Decision: Use USD represented as integer cents.
Why: Exact arithmetic, simple comparisons, and fast operations without floating-point precision issues.
Alternative: Use a high-precision decimal/money library to represent monetary values.
Trade-off: Integer cents is simpler and faster for this assignment; a precision library would provide more flexibility but adds unnecessary complexity.


### 009: Authentication

Decision: Server-side sessions with secure HTTP-only cookies.
Why: Simplest authentication model for a browser-based take-home; easy revocation and straightforward ownership checks.
Alternative: JWT bearer tokens.
Trade-off: Requires server-side session state, but avoids token lifecycle and revocation complexity.


### 010: API Overview

| Method | Endpoint       | Description                    |
| ------ | -------------- | ------------------------------ |
| `POST` | `/auth/signup` | Register a user                |
| `POST` | `/auth/login`  | Login and create a session     |
| `POST` | `/auth/logout` | Invalidate the current session |

| Method   | Endpoint      | Description                                   |
| -------- | ------------- | --------------------------------------------- |
| `POST`   | `/orders`     | Create an order                               |
| `GET`    | `/orders`     | List user's orders; supports status filtering |
| `GET`    | `/orders/:id` | Get complete order details                    |
| `PUT`    | `/orders/:id` | Update an order before first payment          |
| `DELETE` | `/orders/:id` | Delete an order before first payment          |

| Method | Endpoint               | Description         |
| ------ | ---------------------- | ------------------- |
| `POST` | `/orders/:id/payments` | Record a payment    |
| `GET`  | `/orders/:id/payments` | Get payment history |

### 011: Status Derivation

| Condition                               | Status           |
| --------------------------------------- | ---------------- |
| Amount paid = 0 and due date not passed | `pending`        |
| Amount paid > 0 and amount paid < total | `partially_paid` |
| Amount paid = total                     | `paid`           |
| Due date passed and amount paid < total | `overdue`        |
| Due date passed and amount paid = total |  `paid`          |


### 012: Edge Cases & Business Rules

Here are the business rules as a numbered list:

1. An order must contain at least one line item.
2. Quantity must be a positive integer.
3. Unit price must be non-negative.
4. Order total must be greater than zero.
5. Order total is calculated by the server.
6. Monetary values are represented internally as integer cents.
7. Orders can be edited before the first payment.
8. Once a payment exists, customer, due date, and line items become immutable.
9. Orders with existing payments cannot be deleted.
10. A payment cannot exceed the remaining order balance.
11. Concurrent payments for the same order are serialized using row-level protection within a transaction.
12. A payment equal to the remaining balance changes the order to paid.
13. Past due dates are allowed.
14. A due date that is today is not considered overdue.
15. Customer is represented as a plain string; there is no separate customer entity.
16. Refunds are outside the core scope.

### 013: The current implementation is intentionally scoped for a take-home assignment. Before production, I would consider:

1. Add comprehensive audit logging for financial operations.
2. Add a proper refund/adjustment workflow.
3. Add stronger input validation and rate limiting.
4. Add CSRF protection where applicable.
5. Add structured logging, metrics, tracing, and alerting.
6. Add database backup and recovery procedures.
7. Add automated migration/deployment workflows.
8. Improve authorization and security hardening.
9. Add idempotency for payment requests.
10. Add more comprehensive concurrency and failure testing.
11. Consider pagination and query optimization for larger datasets.
12. Add multi-currency support only if the business requires it.
13. Add automated integration and end-to-end tests around financial invariants.