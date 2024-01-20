# Use the official Go image for building
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy Go modules manifests and build the application
COPY go.* ./
RUN go mod download

# Copy the local files to the container's working directory
COPY . .

# Build the application
RUN go build -o webhook-service

# Create a minimal image
FROM alpine:3.14

# Set the working directory
WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/webhook-service /app/webhook-service

# Copy the .env file
COPY config.env /app/config.env

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["/app/webhook-service"]
