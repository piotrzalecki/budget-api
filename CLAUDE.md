# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Single-user personal budget management REST API built in Go, designed for Raspberry Pi deployment. Runs on SQLite with SQLC for type-safe query generation.

## Common Commands

```bash
# Run all tests
make test                    # go test ./...

# Run a single test or package
go test ./internal/handler/... -run TestTransactionHandler
go test ./internal/scheduler/... -run TestRunScheduler

# Start the server locally (uses BUDGET_API_KEY=1234567890)
make run

# Database migrations
make migrate                 # Apply pending migrations
make migrate-down            # Rollback last migration
make db-status               # Show migration status
make new-migration name=<x>  # Create a new migration file

# Code generation (run after editing internal/repo/query.sql)
make generate                # Regenerate SQLC code for SQLite
make generate-pg             # Regenerate SQLC code for PostgreSQL

# Linting
golangci-lint run            # Uses .golangci.yml config

# API documentation
make api-docs                # Regenerate Swagger docs from handler annotations
```

## Architecture

**Layered architecture with the repository pattern:**

```
Gin Router (cmd/budgetd/routes.go)
    → Middleware (auth, validation) (internal/handler/middleware.go)
    → Handlers (internal/handler/*.go)
    → Repository interface (internal/repo/interface.go)
    → SQLC-generated queries (internal/repo/query.sql.go)
    → SQLite via database/sql
```

**Key layers:**
- `cmd/budgetd/` — entry point, server setup, route registration, Swagger config
- `internal/handler/` — HTTP handlers, one file per domain (transactions, tags, recurring, reports, scheduler); `handler.go` wires dependencies
- `internal/repo/` — repository interface + implementation wrapping SQLC; `repo.go` adds `WithTx()` for ACID transaction support; edit `query.sql` then run `make generate`
- `internal/scheduler/` — materializes recurring rules into transactions; runs in-process (hourly) and via systemd timer (2:05 AM daily)
- `pkg/model/` — shared DTOs (`dto.go`) and helpers (`utils.go`); amounts always stored and passed as integer pence

## Database Conventions

- **Amounts are stored as `INT64` pence** — use `pkg/model/utils.go` helpers (`PenceToString`, `StringToPence`) for conversion; never use floats for money
- **Soft-delete pattern**: `deleted_at` timestamp, `NULL` means active; scheduler purges records older than 30 days
- **Idempotency**: `(source_recurring, t_date)` unique constraint prevents duplicate materialization of recurring rules
- Schema lives in `migrations/001_init_schema.sql` (goose format with `-- +goose Up` / `-- +goose Down` markers)

## Required Environment Variable

`BUDGET_API_KEY` must be set — the server panics on startup if missing. Use `docker/env.example` as a template.

## Testing Patterns

- Handler tests mock the `Repository` interface (defined in `internal/repo/interface.go`)
- Integration tests in `internal/scheduler/scheduler_integration_test.go` create an actual in-memory SQLite DB
- Test files live alongside the source they test (`*_test.go` in the same package)

## Code Generation

`internal/repo/query.sql.go` and `internal/repo/models.go` are **generated** — never edit them directly. Make changes in `internal/repo/query.sql`, then run `make generate`.

Swagger docs in `internal/docs/` are also generated — run `make api-docs` after changing handler annotations.

## Deployment

- Docker: multi-stage ARMv7 build targeting Raspberry Pi (`docker/Dockerfile`)
- CI/CD: GitHub Actions builds and pushes to AWS ECR on `main` branch; version read from the `version` file
- Systemd timer provides a backup scheduler trigger (`deploy/budget-scheduler.service` + `.timer`)
