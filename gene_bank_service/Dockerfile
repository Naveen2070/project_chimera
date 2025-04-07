# Use the official Golang image as the base image for building the app
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./ 

# Download dependencies (go mod tidy)
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Install swag (Swagger code generation tool)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Run swag init (this will generate swagger docs inside the app)
RUN swag init --output ./docs --generalInfo ./cmd/server/main.go

# Build the Go app
RUN go build -ldflags "-s -w" -o /main ./cmd/server/main.go

# Start a new stage to create a smaller image with no build tools
FROM alpine:3.17.0

# Set the Current Working Directory inside the container
WORKDIR /app/

# Add a non-root user (avoid running as root for security)
RUN adduser -D appuser

# Copy the pre-built binary from the builder stage to /app/
COPY --from=builder /main /app/main

# Copy the generated Swagger docs into the container
COPY --from=builder /app/docs /app/docs

# Copy the .env file to /app/
COPY .env /app/.env

# Ensure proper permissions for the .env file and the binary
RUN chmod 600 /app/.env && chown appuser:appuser /app/main /app/.env /app/docs && \
    chmod +x /app/main

# Expose the necessary port(s)
EXPOSE 5050

# Switch to the non-root user
USER appuser

# Command to run the executable
CMD ["./main"]
