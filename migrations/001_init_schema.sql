-- +goose Up
-- +goose StatementBegin

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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS recurring_tags;
DROP TABLE IF EXISTS recurring;
DROP TABLE IF EXISTS transaction_tags;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
