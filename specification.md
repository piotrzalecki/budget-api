Budget App – Backend & Scheduler Specification

Last updated: 17 June 2025

⸻

1. Purpose

A single-user (future multi-user) personal-budget application hosted on a Raspberry Pi. This document captures all agreed backend decisions so implementation can begin without ambiguity.

⸻

2. Component Outline

# Component Language Scope

A HTTP API Go 1.22 + Gin CRUD for transactions, tags, recurring rules; reporting; admin scheduler endpoint
B Scheduler Runner Go (library func) invoked in-process and via systemd-timer Materialise recurring rules, purge soft-deleted rows, nightly backup

⸻

3.  High-Level Architecture

                     HTTP/JSON                    ┌───────────────┐          ┌────────────┐

    ┌───────────┐ ────────────────▶ │ Go API │──────────▶│ SQLite DB │
    │ Client │ │ (Gin) │ └────────────┘
    └───────────┘ │ ↳ Scheduler │
    └────────┬──────┘ systemd-timer
    │ (hourly)
    ▼
    curl /admin/run-scheduler

Directory layout

cmd/budgetd # main.go – starts HTTP server & in-proc scheduler
deploy/ # docker-compose.yml, systemd unit/timer files
internal/
handler/ # HTTP handlers (DTO ↔ service)
service/ # Business rules & validation
repo/ # sqlc-generated query wrappers
scheduler/ # Run(ctx) function
pkg/model/ # Shared structs (Transaction, Recurring…)
migrations/ # goose SQL files

⸻

4. Database Schema (SQLite)

-- users (future multi-user)
CREATE TABLE users (
id INTEGER PRIMARY KEY AUTOINCREMENT,
email TEXT UNIQUE NOT NULL,
pw_hash TEXT NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- global key-value config
CREATE TABLE settings (
key TEXT PRIMARY KEY,
value TEXT NOT NULL
); -- e.g. ('payday_dom','28')

-- flexible categorisation
CREATE TABLE tags (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT UNIQUE NOT NULL
);

-- immutable ledger rows
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

-- N-to-M tags ↔ transactions
CREATE TABLE transaction_tags (
transaction_id INTEGER NOT NULL REFERENCES transactions(id),
tag_id INTEGER NOT NULL REFERENCES tags(id),
PRIMARY KEY (transaction_id, tag_id)
);

-- recurring payment rules (generators)
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

-- N-to-M tags ↔ recurring rules
CREATE TABLE recurring_tags (
recurring_id INTEGER NOT NULL REFERENCES recurring(id),
tag_id INTEGER NOT NULL REFERENCES tags(id),
PRIMARY KEY (recurring_id, tag_id)
);

Money stored as INT64 pence to avoid float rounding.

⸻

5. API Surface (/api/v1)

Authentication
• Header X-API-Key: <secret> vs env var BUDGET_API_KEY.

Endpoints

Method & Path Purpose Request JSON Response
POST /transactions Add manual txn {amount:"-12.34", t_date:"2025-06-17", note?, tag_ids[]} {id}
GET /transactions List period ?from=YYYY-MM-DD&to=YYYY-MM-DD [Txn]
PATCH /transactions/{id} Soft-delete or edit note/tags {deleted:true} 204
POST /tags Create tag {name} {id}
GET /tags List tags – [Tag]
POST /recurring Create rule {amount:"-50", description, frequency:"monthly", interval_n:1, first_due_date:"2025-07-01", end_date?, tag_ids[]} {id}
PATCH /recurring/{id} Pause / resume / update {active:false} / field set 204
GET /recurring List rules – [Rule]
GET /reports/monthly Aggregate spend / income ?ym=YYYY-MM {totalIn, totalOut, byTag[]}
POST /admin/run-scheduler Trigger materialisation + purge + backup – (API key protected) {processed:N}

All responses JSON with envelope {data:…, error:null} on success.

⸻

6. Scheduler Logic

func RunScheduler(ctx context.Context, db \*sql.DB, today time.Time) error {
tx := BeginExclusive(db) // SQLite serialises
rules := selectRulesDue(tx, today) // active && next*due_date ≤ today
for *, r := range rules {
insertTransaction(tx, r)
advanceNextDue(tx, r)
}
purgeSoftDeleted(tx, today.AddDate(0,0,-30))
maybeBackup(today)
return tx.Commit()
}

Exposed via /admin/run-scheduler and called hourly by systemd-timer.

⸻

7. Logging & Observability
   • zap JSON → journald (SystemMaxUse=200M).
   • Scheduler emits: {"level":"info","msg":"scheduler","processed":7}.

⸻

8. Deployment & CI/CD
   • Docker Compose – single api service.
   • GitHub Actions – Buildx multi-arch image → ghcr.io.
   • Manual deploy (deploy.sh) – SSH, pull, restart, prune.
   • Migrations – goose up in docker-entrypoint.sh.

⸻

9. Environment Variables

Var Example Purpose
BUDGET_API_KEY 8de7… Header auth secret
DB_PATH /data/budget.db SQLite location
TZ Europe/London Local cron maths

⸻

10. Future Enhancements (non-blocking)
    • PostgreSQL, JWT + refresh tokens, budgets/limits, CSV import, multi-currency.
    • Observability stack (Prometheus + Loki).
    • Automated Pi deploy (watchtower or self-hosted GitHub runner).

⸻

End of specification.
