# ──────────────────────────────────────────────────────────────
# Budget App – consolidated Makefile (migrations + sqlc codegen)
# ──────────────────────────────────────────────────────────────

# 📂 Where your .sql migrations live
MIGR_DIR   := migrations

# 🔌 Default driver & connection string (SQLite in the repo root)
DB_DRIVER  ?= sqlite3
DB_STRING  ?= $(CURDIR)/dev.db        # override in CI/Prod

# ────────── Targets ───────────────────────────────────────────

.PHONY: migrate migrate-down new-migration generate db-status test

## Run unit tests
test:
	go test ./...

run:
	BUDGET_API_KEY=1234567890 go run ./cmd/budgetd

## Apply all up migrations
migrate:
	goose -dir $(MIGR_DIR) $(DB_DRIVER) $(DB_STRING) up

reset:
	goose -dir $(MIGR_DIR) $(DB_DRIVER) $(DB_STRING) reset

## Roll back the last migration
migrate-down:
	goose -dir $(MIGR_DIR) $(DB_DRIVER) $(DB_STRING) down

## Show current migration status
db-status:
	goose -dir $(MIGR_DIR) $(DB_DRIVER) $(DB_STRING) status

## Create a new timestamped migration file
## Usage: make new-migration name=<snake_case>
new-migration:
ifndef name
	$(error Specify name via 'make new-migration name=add_users')
endif
	goose -dir $(MIGR_DIR) create $(name) sql

## Regenerate sqlc-produced Go code
generate:
	sqlc generate

generate-pg:
	sqlc generate file sqlc-prd.yaml

# ────────── Helper shortcuts for Postgres (optional) ──────────
# Example: PG_DSN = "postgres://user:pw@host:5432/appdb?sslmode=disable"

migrate-pg:
	@$(MAKE) migrate DB_DRIVER=postgres DB_STRING="$(PG_DSN)"

db-status-pg:
	@$(MAKE) db-status DB_DRIVER=postgres DB_STRING="$(PG_DSN)"