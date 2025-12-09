# Docker Setup for Budget API

This directory contains all Docker-related files for deploying the Budget API application.

## Architecture Support

This Docker setup is specifically configured for **Raspberry Pi (ARM64)** architecture, making it perfect for hosting on a Raspberry Pi device.

## Files

- `docker-compose.yaml` - Main Docker Compose configuration (ARM64 optimized)
- `Dockerfile` - Multi-stage build for ARM64 Raspberry Pi
- `.dockerignore` - Excludes unnecessary files from build context
- `env.example` - Environment variables template
- `deploy.sh` - Automated deployment script

## Quick Start

1. **Setup environment:**

   ```bash
   cp env.example .env
   # Edit .env and set your BUDGET_API_KEY
   ```

2. **Deploy:**

   ```bash
   ./deploy.sh
   ```

3. **Verify:**
   ```bash
   docker-compose ps
   curl http://localhost:8080/api/v1/health
   ```

## Manual Commands

```bash
# Build and start
docker-compose up -d --build

# View logs
docker-compose logs -f budget-api

# Stop services
docker-compose down

# Stop and remove volumes (⚠️ WARNING: This will delete all data)
docker-compose down -v
```

## Environment Variables

Required environment variables (set in `.env`):

- `BUDGET_API_KEY` - Secret key for API authentication
- `DB_PATH` - SQLite database path (default: `/data/budget.db`)
- `TZ` - Timezone for scheduler (default: `Europe/London`)
- `PORT` - Server port (default: `8080`)

## Data Persistence

The SQLite database is stored in the `budget_data` Docker volume, ensuring data persists across container restarts and updates.

## Health Checks

The service includes health checks that verify the API is responding. The health endpoint is available at `/api/v1/health` and doesn't require authentication.

## Raspberry Pi Optimization

- **ARM64 architecture**: Built specifically for Raspberry Pi 3/4/5
- **Alpine Linux**: Lightweight base image for better performance
- **SQLite**: Perfect for single-user budget tracking on Pi
- **Minimal dependencies**: Optimized for resource-constrained devices
