# Budget App – Backend Specification

**Last updated:** January 2025  
**Version:** 1.0

---

## 1. Purpose

A single-user (future multi-user) personal-budget application hosted on a Raspberry Pi. This document captures the complete backend architecture and implementation details for the Budget API.

---

## 2. Component Overview

| Component            | Language          | Scope                                                                                                            |
| -------------------- | ----------------- | ---------------------------------------------------------------------------------------------------------------- |
| **HTTP API**         | Go 1.22 + Gin     | CRUD for transactions, tags, recurring rules; reporting; admin scheduler endpoint                                |
| **Scheduler Runner** | Go (library func) | Invoked in-process and via systemd-timer. Materializes recurring rules, purges soft-deleted rows, nightly backup |

---

## 3. High-Level Architecture

```
                     HTTP/JSON                    ┌───────────────┐          ┌────────────┐
    ┌───────────┐ ────────────────▶ │ Go API │──────────▶│ SQLite DB │
    │ Client │ │ (Gin) │ └────────────┘
    └───────────┘ │ ↳ Scheduler │
    └────────┬──────┘ systemd-timer
    │ (hourly)
    ▼
    curl /admin/run-scheduler
```

### Directory Layout

```
cmd/budgetd/          # main.go – starts HTTP server & in-proc scheduler
deploy/              # docker-compose.yml, systemd unit/timer files
internal/
├── handler/         # HTTP handlers (DTO ↔ service)
├── service/         # Business rules & validation
├── repo/           # sqlc-generated query wrappers
├── scheduler/      # Run(ctx) function
└── docs/          # Swagger documentation
pkg/model/          # Shared structs (Transaction, Recurring…)
migrations/         # goose SQL files
```

---

## 4. Database Schema (SQLite)

### Core Tables

#### Users (Future Multi-User)

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    pw_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Settings (Global Key-Value Config)

```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
); -- e.g. ('payday_dom','28')
```

#### Tags (Flexible Categorization)

```sql
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);
```

#### Transactions (Immutable Ledger Rows)

```sql
CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    amount_pence INTEGER NOT NULL, -- £12.34 → 1234
    t_date DATE NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    source_recurring INTEGER REFERENCES recurring(id),
    deleted_at TIMESTAMP, -- NULL = live
    UNIQUE(source_recurring, t_date) -- idempotency guard
);
CREATE INDEX idx_txn_user_date_active ON transactions(user_id, t_date)
WHERE deleted_at IS NULL;
```

#### Transaction Tags (N-to-M Tags ↔ Transactions)

```sql
CREATE TABLE transaction_tags (
    transaction_id INTEGER NOT NULL REFERENCES transactions(id),
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    PRIMARY KEY (transaction_id, tag_id)
);
```

#### Recurring (Recurring Payment Rules)

```sql
CREATE TABLE recurring (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    amount_pence INTEGER NOT NULL,
    description TEXT,
    frequency TEXT NOT NULL, -- 'daily'|'weekly'|'monthly'|'yearly'
    interval_n INTEGER NOT NULL DEFAULT 1,
    first_due_date DATE NOT NULL,
    next_due_date DATE NOT NULL,
    end_date DATE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_recurring_next_due ON recurring(next_due_date)
WHERE active = TRUE;
```

#### Recurring Tags (N-to-M Tags ↔ Recurring Rules)

```sql
CREATE TABLE recurring_tags (
    recurring_id INTEGER NOT NULL REFERENCES recurring(id),
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    PRIMARY KEY (recurring_id, tag_id)
);
```

**Note:** Money stored as INT64 pence to avoid float rounding.

---

## 5. API Surface (/api/v1)

### Authentication

- **Header:** `X-API-Key: <secret>` vs env var `BUDGET_API_KEY`
- **Health endpoint:** `/health` (no auth required)

### Endpoints

#### Transactions

| Method & Path                                  | Purpose               | Request JSON                                               | Response        |
| ---------------------------------------------- | --------------------- | ---------------------------------------------------------- | --------------- |
| `POST /transactions`                           | Add manual txn        | `{amount:"-12.34", t_date:"2025-06-17", note?, tag_ids[]}` | `{id}`          |
| `GET /transactions`                            | List period           | `?from=YYYY-MM-DD&to=YYYY-MM-DD`                           | `[Txn]`         |
| `GET /transactions/:id`                        | Get single txn        | –                                                          | `Txn`           |
| `PATCH /transactions/:id`                      | Soft-delete or edit   | `{deleted:true}`                                           | `204`           |
| `GET /transactions/by-recurring/:recurring_id` | Get by recurring rule | –                                                          | `[Txn]`         |
| `GET /transactions/by-tag/:tag_id`             | Get by tag            | –                                                          | `[Txn]`         |
| `POST /transactions/purge`                     | Purge soft-deleted    | `{cutoff_date:"2025-01-01"}`                               | `{processed:N}` |

#### Tags

| Method & Path | Purpose    | Request JSON | Response |
| ------------- | ---------- | ------------ | -------- |
| `POST /tags`  | Create tag | `{name}`     | `{id}`   |
| `GET /tags`   | List tags  | –            | `[Tag]`  |

#### Recurring Rules

| Method & Path                   | Purpose             | Request JSON                                                                                                        | Response |
| ------------------------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------- | -------- |
| `POST /recurring`               | Create rule         | `{amount:"-50", description, frequency:"monthly", interval_n:1, first_due_date:"2025-07-01", end_date?, tag_ids[]}` | `{id}`   |
| `GET /recurring`                | List rules          | –                                                                                                                   | `[Rule]` |
| `GET /recurring/:id`            | Get single rule     | –                                                                                                                   | `Rule`   |
| `PATCH /recurring/:id`          | Pause/resume/update | `{active:false}` / field set                                                                                        | `204`    |
| `DELETE /recurring/:id`         | Delete rule         | –                                                                                                                   | `204`    |
| `GET /recurring/active`         | List active rules   | –                                                                                                                   | `[Rule]` |
| `PATCH /recurring/:id/toggle`   | Toggle active state | –                                                                                                                   | `204`    |
| `GET /recurring/due`            | Get due rules       | `?date=2025-01-01`                                                                                                  | `[Rule]` |
| `GET /recurring/by-tag/:tag_id` | Get by tag          | –                                                                                                                   | `[Rule]` |

#### Reports

| Method & Path                 | Purpose                | Request JSON  | Response                                   |
| ----------------------------- | ---------------------- | ------------- | ------------------------------------------ |
| `GET /reports/monthly`        | Aggregate spend/income | `?ym=YYYY-MM` | `{totalIn, totalOut, byTag[]}`             |
| `GET /reports/monthly/totals` | Monthly totals         | `?ym=YYYY-MM` | `{total_in, total_out, transaction_count}` |

#### Admin

| Method & Path               | Purpose                                  | Request JSON | Response        |
| --------------------------- | ---------------------------------------- | ------------ | --------------- |
| `POST /admin/run-scheduler` | Trigger materialization + purge + backup | –            | `{processed:N}` |

### Response Format

All responses use JSON envelope: `{data:…, error:null}` on success.

---

## 6. Scheduler Logic

```go
func RunScheduler(ctx context.Context, db *sql.DB, today time.Time) error {
    tx := BeginExclusive(db) // SQLite serializes
    rules := selectRulesDue(tx, today) // active && next_due_date ≤ today
    for _, r := range rules {
        insertTransaction(tx, r)
        advanceNextDue(tx, r)
    }
    purgeSoftDeleted(tx, today.AddDate(0,0,-30))
    maybeBackup(today)
    return tx.Commit()
}
```

**Exposed via:** `/admin/run-scheduler` and called hourly by systemd-timer.

---

## 7. Logging & Observability

- **Logger:** zap JSON → journald (SystemMaxUse=200M)
- **Scheduler emits:** `{"level":"info","msg":"scheduler","processed":7}`
- **Health checks:** Built-in endpoint at `/health`

---

## 8. Deployment & CI/CD

### Docker Deployment

- **Docker Compose** – single api service
- **Multi-arch support** – ARM64 optimized for Raspberry Pi
- **Data persistence** – SQLite in Docker volume
- **Health checks** – Built-in container health monitoring

### Systemd Timer

- **Independent scheduler** – Runs via systemd timer (daily at 2:05 AM)
- **Fault tolerance** – Continues working even if API process crashes
- **Configurable cadence** – Easy to adjust timing patterns

### Environment Variables

| Variable         | Example           | Purpose            |
| ---------------- | ----------------- | ------------------ |
| `BUDGET_API_KEY` | `8de7…`           | Header auth secret |
| `DB_PATH`        | `/data/budget.db` | SQLite location    |
| `TZ`             | `Europe/London`   | Local cron maths   |
| `PORT`           | `8080`            | Server port        |

---

## 9. Development Tools

| Tool              | Version | Purpose                          |
| ----------------- | ------- | -------------------------------- |
| **sqlc**          | v1.29.0 | Type-safe SQL code generation    |
| **goose**         | v3.24.3 | Database migrations              |
| **golangci-lint** | 2.1.6   | Code linting                     |
| **swag**          | Latest  | Swagger documentation generation |

---

## 10. Security Considerations

- **API Key Authentication** – Required for all endpoints except `/health`
- **Input Validation** – Comprehensive request validation with custom validators
- **SQL Injection Protection** – sqlc-generated queries prevent injection
- **Soft Deletes** – Data retention with purge capabilities
- **Idempotency** – Recurring transactions protected against duplicates

---

## 11. Future Enhancements (Non-Blocking)

- **PostgreSQL** – For multi-user support
- **JWT + Refresh Tokens** – Enhanced authentication
- **Budgets/Limits** – Spending limits and alerts
- **CSV Import** – Bulk transaction import
- **Multi-Currency** – International currency support
- **Observability Stack** – Prometheus + Loki integration
- **Automated Pi Deploy** – Watchtower or self-hosted GitHub runner

---

## 12. API Documentation

Interactive API documentation is available at `/docs` when running the application, powered by Swagger/OpenAPI 3.0.

---

_End of Backend Specification_
