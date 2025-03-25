FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o audience-service ./audience/cmd/server

# Use a minimal alpine image
FROM alpine:3.18

# Add ca-certificates for secure connections
RUN apk --no-cache add ca-certificates

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/audience-service .

# Use non-root user
USER appuser

# Set environment variable for port
ENV PORT=9006

# Expose the port
EXPOSE 9006

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Set the entrypoint
ENTRYPOINT ["./audience-service"] 