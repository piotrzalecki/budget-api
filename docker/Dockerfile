# Build stage (ARMv7)
FROM --platform=linux/arm/v7 golang:1.23-alpine AS builder

# Toolchain for CGO + common deps
RUN apk add --no-cache build-base ca-certificates tzdata git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN mkdir -p /root/go/bin && \
    cp $(go env GOPATH)/bin/goose /root/go/bin/goose 2>/dev/null || \
    cp /go/bin/goose /root/go/bin/goose 2>/dev/null || \
    find /root -name "goose" -type f -executable -exec cp {} /root/go/bin/goose \; 2>/dev/null || \
    find /go -name "goose" -type f -executable -exec cp {} /root/go/bin/goose \; 2>/dev/null

# Copy source
COPY . .

# Compile for ARMv7 hard-float (Pi 32-bit)
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
    go build -trimpath -ldflags="-s -w" -o /out/budgetd ./cmd/budgetd

# Runtime stage (ARMv7)
FROM --platform=linux/arm/v7 alpine:3.20

RUN apk --no-cache add ca-certificates tzdata curl sqlite

# Non-root user
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /out/budgetd /app/budgetd
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /root/go/bin/goose /app/goose
COPY docker/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh && \
    mkdir -p /data && chown -R appuser:appgroup /data /app
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/app/entrypoint.sh"]