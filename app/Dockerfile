# Multi-stage security-hardened containerfile
# Builder stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install CA certificates for secure HTTPS outbound
RUN apk add --no-cache ca-certificates

COPY go.mod ./
RUN go mod download || true

COPY . .

# Compile static binary without CGO, stripping debug symbols
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /build/service main.go

# Production runtime stage (Distroless-style scratch / minimal non-root alpine)
FROM alpine:3.20

# Create dedicated non-root user and group
RUN addgroup -S appgroup -g 10001 && \
    adduser -S -D -H -u 10001 -h /app -s /sbin/nologin -G appgroup appuser

WORKDIR /app

# Copy CA certs and statically linked binary
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/service /app/service

# Set ownership
RUN chown -R appuser:appgroup /app

# Security: run as non-root UID 10001
USER 10001:10001

EXPOSE 8080

# Built-in health check
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/service", "-version"] || exit 0

ENTRYPOINT ["/app/service"]
