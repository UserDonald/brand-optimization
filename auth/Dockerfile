FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./auth/cmd

# Use a minimal alpine image for the final stage
FROM alpine:3.18

# Add wget for healthcheck
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /auth-service .

# Set environment variable for port
ENV PORT=9001

# Expose the service port
EXPOSE 9001

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run the service
CMD ["./auth-service"] 