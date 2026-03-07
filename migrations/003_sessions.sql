-- +goose Up
-- +goose StatementBegin

ALTER TABLE users ADD COLUMN is_service BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE sessions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_token ON sessions(token);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE sessions;
ALTER TABLE users DROP COLUMN is_service;

-- +goose StatementEnd
