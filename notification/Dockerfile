FROM golang:1.21 as builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./notification/cmd/server

# Use minimal alpine image
FROM alpine:latest

# Add ca-certificates for HTTPS and wget for health check
RUN apk --no-cache add ca-certificates wget

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /notification-service .

# Set environment variables
ENV PORT=9002
ENV ENV=production

# Expose port
EXPOSE 9002

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run the service
CMD ["./notification-service"] 