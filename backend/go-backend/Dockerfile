# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
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

# Install curl for healthcheck and sqlite for database
RUN apk add --no-cache curl sqlite

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 3001

# Run the application
CMD ["./main"]
