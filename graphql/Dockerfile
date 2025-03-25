FROM golang:1.21-alpine AS build

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o graphql-service ./graphql/cmd/main.go

# Final stage
FROM alpine:latest

# Install CA certificates and wget for health check
RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

# Copy the binary from the build stage
COPY --from=build /app/graphql-service .
COPY --from=build /app/graphql/schema.graphql ./graphql/schema.graphql

# Expose the port
EXPOSE 8080

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run
CMD ["./graphql-service"] 