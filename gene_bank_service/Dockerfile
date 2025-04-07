# Stage 1: Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy only go.mod and go.sum first to leverage Docker's caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy only the necessary files
COPY . ./

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger Docs
RUN swag init --output ./docs --generalInfo ./cmd/server/main.go

# Build the application
RUN go build -ldflags "-s -w" -o /main ./cmd/server/main.go

# Stage 2: Final Image (Slimmer)
FROM alpine:3.17.0

WORKDIR /app/

# Add a non-root user for security
RUN adduser -D appuser

# Copy the binary and other necessary files from the build stage
COPY --from=builder /main /app/main
COPY --from=builder /app/docs /app/docs

# Optionally, don't copy .env directly in production
COPY .env /app/.env

# Set permissions and ensure the correct ownership
RUN chmod 600 /app/.env && chown appuser:appuser /app/main /app/.env /app/docs && \
    chmod +x /app/main

EXPOSE 5050

# Switch to a non-root user for running the application
USER appuser

# Command to run the application
CMD ["./main"]
