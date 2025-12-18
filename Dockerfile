# Step 1: Build the binary with the latest Go version
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the project files into the container
COPY . .

# Build the Go binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o playpi main.go

# Step 2: Create the final image
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 playpi && \
    adduser -D -u 1000 -G playpi playpi

# Copy the binary from the builder stage
COPY --from=builder /app/playpi .

# Change ownership to non-root user
RUN chown -R playpi:playpi /app

# Switch to non-root user
USER playpi

# Expose dashboard port and all service ports
EXPOSE 8000 8080 8081 8082 8084 8085 8086

# Default command to run the dashboard
CMD ["./playpi", "start", "dashboard"]
