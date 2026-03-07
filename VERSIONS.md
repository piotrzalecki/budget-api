# Changelog

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
