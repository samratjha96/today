# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

ENV GOPROXY=direct

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install curl for healthcheck
RUN apk add --no-cache curl

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 3001

# Run the application
CMD ["./main"]
