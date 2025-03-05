# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o migration ./cmd/migration/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migration .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/docker-entrypoint.sh .

# Expose port
EXPOSE 8080

# Command to run
ENTRYPOINT ["sh", "/app/docker-entrypoint.sh"]
