# Build stage
FROM --platform=linux/arm64 golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application for ARM64
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o budgetd ./cmd/budgetd

# Runtime stage
FROM --platform=linux/arm64 alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl sqlite

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/budgetd .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Create data directory
RUN mkdir -p /data && chown -R appuser:appgroup /data /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1

# Default command
CMD ["./budgetd"] 