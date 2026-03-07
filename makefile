# ──────────────────────────────────────────────────────────────
# Budget App – consolidated Makefile (migrations + sqlc codegen)
# ──────────────────────────────────────────────────────────────

# 📂 Where your .sql migrations live
MIGR_DIR   := migrations

# 🔌 Default driver & connection string (SQLite in the repo root)
DB_DRIVER  ?= sqlite3
DB_STRING  ?= $(CURDIR)/dev.db        # override in CI/Prod

# ────────── Targets ───────────────────────────────────────────

.PHONY: migrate migrate-down migrate-dev new-migration generate db-status test install-timer docs

## Run unit tests
test:
	go test ./...

run:
	BUDGET_API_KEY=1234567890 CORS_ORIGINS="http://localhost:5173" \
	go run -ldflags="-X main.version=$(shell cat version)" ./cmd/budgetd

## Generate Swagger documentation
api-docs:
	~/go/bin/swag init -g cmd/budgetd/main.go -o internal/docs

## Apply all up migrations
migrate:
	goose -dir $(MIGR_DIR) $(DB_DRIVER) $(DB_STRING) up

reset:
	goose -dir $(MIGR_DIR) $(DB_DRIVER) $(DB_STRING) reset

## Apply schema migrations then seed with dev example data
migrate-dev: migrate
	goose -dir $(MIGR_DIR)/seeds $(DB_DRIVER) $(DB_STRING) up

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

## Install and enable systemd timer for budget scheduler
## Usage: sudo make install-timer
install-timer:
	@echo "Installing budget scheduler systemd units..."
	@cp deploy/budget-scheduler.service /etc/systemd/system/
	@cp deploy/budget-scheduler.timer /etc/systemd/system/
	@systemctl daemon-reload
	@systemctl enable budget-scheduler.timer
	@systemctl start budget-scheduler.timer
	@echo "✅ Budget scheduler timer installed and enabled"
	@echo "📊 Check status: systemctl status budget-scheduler.timer"
	@echo "📋 View logs: journalctl -u budget-scheduler.service"

## Generate architecture diagram from Mermaid
## Usage: make appdocs (requires mmdc to be installed)
appdocs:
	@echo "Generating architecture diagram..."
	@if command -v mmdc >/dev/null 2>&1; then \
		mmdc -i docs/architecture.mmd -o docs/img/architecture.svg; \
		echo "✅ Architecture diagram generated"; \
	else \
		echo "⚠️  mmdc not found. Install with: npm install -g @mermaid-js/mermaid-cli"; \
		echo "📝 Manual SVG available at: docs/img/architecture.svg"; \
	fi

# ────────── Helper shortcuts for Postgres (optional) ──────────
# Example: PG_DSN = "postgres://user:pw@host:5432/appdb?sslmode=disable"

migrate-pg:
	@$(MAKE) migrate DB_DRIVER=postgres DB_STRING="$(PG_DSN)"

db-status-pg:
	@$(MAKE) db-status DB_DRIVER=postgres DB_STRING="$(PG_DSN)"