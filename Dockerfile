# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Enable Go modules and static build
ENV CGO_ENABLED=0

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary from the real entrypoint (main.go in root)
RUN go build -trimpath -ldflags="-s -w" -o /app/purchase-cart-service ./main.go

# ---------- Runtime stage ----------
FROM alpine:3.19
WORKDIR /app

# CA certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

# Copy binary and configuration
COPY --from=builder /app/purchase-cart-service /app/purchase-cart-service
COPY config.json /app/config.json

# Run as non-root user
RUN addgroup -S app && adduser -S -G app app && chown -R app:app /app
USER app:app

# Service exposed port
EXPOSE 8080

# Useful environment variables
ENV GIN_MODE=release
ENV CONFIG_PATH=/app/config.json

# Simple healthcheck on the health endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/v1/health >/dev/null 2>&1 || exit 1

# Start the service
ENTRYPOINT ["/app/purchase-cart-service"]
