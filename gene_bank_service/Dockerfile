# Stage 1: Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy only go.mod and go.sum first to leverage Docker's caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application
COPY . ./

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger Docs
RUN swag init --output ./docs --generalInfo ./cmd/server/main.go

# Build the application
RUN go build -ldflags "-s -w" -o /main ./cmd/server/main.go

# Stage 2: Final Image (Slimmer)
FROM alpine:3.21.3

WORKDIR /app/

# Add a non-root user for security
RUN adduser -D appuser

# Copy the binary and Swagger docs from the build stage
COPY --from=builder /main /app/main
COPY --from=builder /app/docs /app/docs

# Set environment variables
ENV CONSUL_HOST=host.docker.internal \
    CONSUL_PORT=8500 \
    SERVICE_HOST=host.docker.internal \
    SERVICE_NAME=gene-bank-service \
    SERVICE_PORT=5050 \
    STATUS=passing \
    INTERVAL=15s \
    TIMEOUT=10s \
    DEREGISTER_CRITICAL_SERVICE_AFTER=5m \
    APP_PORT=5050 \
    RABBITMQ_URL=amqp://admin:naveen@2007@host.docker.internal:5672

# Set permissions and ownership
RUN chown appuser:appuser /app/main /app/docs && chmod +x /app/main

# Expose the application port
EXPOSE 5050

# Switch to a non-root user
USER appuser

# Command to run the application
CMD ["./main"]
