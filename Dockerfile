# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code including database/ subdirectory
COPY *.go ./
COPY database/ ./database/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /devops-interview

# Runtime stage - use minimal Alpine image
FROM alpine:3.19

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /devops-interview /app/devops-interview

# Set environment variables with defaults matching .env.example
ENV HTTP_ADDR=:8080 \
    DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable \
    REDIS_ADDR=localhost:6379

# Use non-root user
USER appuser

# Document the port the application listens on
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/devops-interview"]