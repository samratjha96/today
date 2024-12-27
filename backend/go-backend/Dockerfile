# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies only when needed
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache git gcc musl-dev

# Copy go mod files first to cache dependencies
COPY go.mod go.sum ./

ENV GOPROXY=direct
ENV CGO_ENABLED=1
ENV GOMODCACHE=/go/pkg/mod
ENV GOCACHE=/root/.cache/go-build

# Download dependencies with build cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source code
COPY . .

# Build the application with build cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache curl sqlite

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 3001

# Run the application
CMD ["./main"]
