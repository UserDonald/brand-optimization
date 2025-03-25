FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /engagement-service ./engagement/cmd

# Use a minimal alpine image for the final stage
FROM alpine:3.18

# Add ca-certificates for HTTPS and wget for health check
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /engagement-service .

# Set environment variable for port
ENV PORT=9004

# Expose the service port
EXPOSE 9004

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run the service
CMD ["./engagement-service"] 