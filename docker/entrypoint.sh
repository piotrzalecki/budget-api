#!/bin/sh
set -e

echo "Starting budget API initialization..."
# Check if database file exists, create directory if needed
if [ ! -f "$DB_PATH" ]; then
    echo "Database file not found at $DB_PATH, creating directory structure..."
    mkdir -p "$(dirname "$DB_PATH")"
fi

# Run database migrations
echo "Running database migrations..."
/app/goose -dir /app/migrations sqlite3 "$DB_PATH" up

# Start the application
echo "Starting budget API server..."
exec /app/budgetd
