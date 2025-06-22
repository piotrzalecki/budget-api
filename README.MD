# budget-api

## Authentication

Every request must include header `X-API-Key: <value of $BUDGET_API_KEY>`.
The server panics on startup if `BUDGET_API_KEY` is unset.

**Note:** The `/health` endpoint does not require authentication.

## Dependencies

# Development Tools

| Tool Name     | Version |
| ------------- | ------- |
| sqlc          | v1.29.0 |
| goose         | v3.24.3 |
| golangci-lint | 2.1.6   |
