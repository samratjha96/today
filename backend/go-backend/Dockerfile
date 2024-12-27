# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies including gcc and musl-dev for CGO
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

ENV GOPROXY=direct
ENV CGO_ENABLED=1

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required runtime dependencies
RUN apk add --no-cache sqlite

# Create data directory for SQLite database
RUN mkdir -p /app/data && chown -R nobody:nobody /app/data

# Copy binary from builder
COPY --from=builder /app/main .

# Use non-root user
USER nobody

# Expose port
EXPOSE 3001

# Run the application
CMD ["./main"]
