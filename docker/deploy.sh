#!/bin/bash

# Budget API Deployment Script
# Usage: ./deploy.sh [environment]

set -e

ENVIRONMENT=${1:-production}
COMPOSE_FILE="docker-compose.yaml"

echo "🚀 Deploying Budget API to $ENVIRONMENT..."

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed"
    exit 1
fi

# Check if .env file exists (look in parent directory first, then current)
if [ -f ../.env ]; then
    ENV_FILE="../.env"
elif [ -f .env ]; then
    ENV_FILE=".env"
else
    echo "⚠️  .env file not found. Please copy env.example to .env and configure it."
    echo "   cp env.example .env"
    exit 1
fi

# Export environment variables from .env file
if [ -f "$ENV_FILE" ]; then
    export $(cat "$ENV_FILE" | grep -v '^#' | xargs)
fi

# Pull latest images (if using remote registry)
echo "📥 Pulling latest images..."
docker-compose -f $COMPOSE_FILE pull || echo "⚠️  No remote images to pull"

# Build and start services
echo "🔨 Building and starting services..."
docker-compose -f $COMPOSE_FILE up -d --build

# Wait for health check
echo "⏳ Waiting for service to be healthy..."
timeout=60
counter=0
while [ $counter -lt $timeout ]; do
    if docker-compose -f $COMPOSE_FILE ps | grep -q "healthy"; then
        echo "✅ Service is healthy!"
        break
    fi
    sleep 2
    counter=$((counter + 2))
done

if [ $counter -eq $timeout ]; then
    echo "⚠️  Service health check timeout. Check logs with: docker-compose logs"
fi

# Clean up old images
echo "🧹 Cleaning up old images..."
docker image prune -f

echo "🎉 Deployment completed!"
echo "📊 Service status:"
docker-compose -f $COMPOSE_FILE ps

echo "📝 Logs: docker-compose logs -f budget-api"
echo "🛑 Stop: docker-compose down" 