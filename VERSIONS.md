# Changelog

## 0.1.1

- Default admin user seeded on startup via `ADMIN_EMAIL` + `ADMIN_PASSWORD` env vars (skipped if any regular users already exist)

## 0.1.0

- Added session-based authentication: `POST /api/v1/auth/login` and `POST /api/v1/auth/logout`
- `/api/v1/*` routes now require `Authorization: Bearer <token>` instead of `X-API-Key`
- `/admin/*` keeps `X-API-Key` authentication (used by systemd timer)
- Added user management endpoints: `GET/POST /api/v1/users`, `GET/PATCH/DELETE /api/v1/users/:id`
- Sessions stored in new `sessions` DB table (migration 003); supports logout/revocation
- Added `is_service` flag to `users` table for permanent service accounts
- Service user seeded on startup via `SERVICE_USER_EMAIL` + `SERVICE_USER_TOKEN` env vars (permanent session, idempotent)
- Passwords hashed with bcrypt

## 0.0.12

- Added `PATCH /tags/:id` and `DELETE /tags/:id` endpoints
- Fixed `POST /tags` to return `201 Created` with full tag object (`id` + `name`)

## 0.0.11

- Version is now injected at build time via Go ldflags (`-X main.version=`) instead of being hardcoded
- `version` file is the single source of truth: read by CI/CD (GHA), Docker build (`ARG APP_VERSION`), and `make run`
- `GET /health` response now reports the actual build version instead of the hardcoded `"1.0.0"`
- `make run` injects the version from the `version` file during local development

## 0.0.10

- Initial release with transaction tracking, recurring payments, tags, and monthly reports
- SQLite storage with SQLC-generated queries
- Scheduler for materialising recurring rules into transactions
- Docker multi-stage ARMv7 build targeting Raspberry Pi
- CI/CD via GitHub Actions with AWS ECR push
