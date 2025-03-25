# Development Guide

This guide provides instructions for setting up a development environment and contributing to the Strategic Brand Optimization Platform.

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Project Structure](#project-structure)
3. [Adding a New Service](#adding-a-new-service)
4. [Modifying an Existing Service](#modifying-an-existing-service)
5. [Protocol Buffers and gRPC](#protocol-buffers-and-grpc)
6. [Testing](#testing)
7. [Local Development Workflow](#local-development-workflow)
8. [Coding Standards](#coding-standards)
9. [Documentation Standards](#documentation-standards)
10. [Common Development Tasks](#common-development-tasks)
11. [Troubleshooting](#troubleshooting)

## Development Environment Setup

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Protocol Buffers compiler (`protoc`)
- Supabase CLI (optional, for local Supabase)
- Git

### Installation Steps

1. **Install Go**:
   - Download from [go.dev](https://go.dev/dl/)
   - Verify installation: `go version`

2. **Install Docker and Docker Compose**:
   - Follow instructions at [docker.com](https://docs.docker.com/get-docker/)
   - Verify installation: `docker --version` and `docker-compose --version`

3. **Install Protocol Buffers compiler**:
   - Follow instructions at [protocolbuffers/protobuf](https://github.com/protocolbuffers/protobuf/releases)
   - Verify installation: `protoc --version`

4. **Install Go gRPC plugins**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
   ```

5. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/go-competitor.git
   cd go-competitor
   ```

6. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your development settings
   ```

## Project Structure

The project follows a microservice architecture with these main components:

```
go-competitor/
├── auth/              # Authentication & user management
├── competitor/        # Competitor data tracking & analysis
├── engagement/        # Metrics & engagement analytics
├── content/           # Content scheduling & optimization
├── audience/          # Audience segmentation & insights
├── analytics/         # AI-powered predictive analytics
├── notification/      # Alerts and scheduled notifications
├── scraper/           # Data collection from social platforms
├── graphql/           # GraphQL gateway & schema definitions
├── common/            # Shared utilities and middleware
├── docs/              # Documentation
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
├── docker-compose.yaml  # Docker Compose configuration
└── README.md          # Project readme
```

Each service follows a consistent structure:

```
service/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer
```

## Adding a New Service

To add a new service to the platform:

1. **Create the service directory structure**:
   ```bash
   mkdir -p newservice/{cmd,server,client,pb,service,repository}
   ```

2. **Define the Protocol Buffers**:
   Create `newservice/pb/newservice.proto`:
   ```protobuf
   syntax = "proto3";
   
   package newservice;
   option go_package = "github.com/donaldnash/go-competitor/newservice/pb";
   
   service NewService {
     // Define your RPC methods here
     rpc Example(ExampleRequest) returns (ExampleResponse);
   }
   
   message ExampleRequest {
     string tenant_id = 1;
     // Add more fields as needed
   }
   
   message ExampleResponse {
     string message = 1;
     // Add more fields as needed
   }
   ```

3. **Generate gRPC code**:
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       newservice/pb/newservice.proto
   ```

4. **Implement the server**:
   Create `newservice/server/server.go` with the gRPC server implementation.

5. **Implement the client**:
   Create `newservice/client/client.go` with the gRPC client implementation.

6. **Implement the business logic**:
   Create service implementations in `newservice/service/`.

7. **Implement the repository layer**:
   Create data access functions in `newservice/repository/`.

8. **Create the main entry point**:
   Create `newservice/cmd/main.go`.

9. **Add Dockerfile**:
   Create `newservice/Dockerfile`.

10. **Update docker-compose.yaml**:
    Add the new service to the Docker Compose configuration.

11. **Document the service**:
    Create documentation in `docs/services/newservice.md`.

## Modifying an Existing Service

When modifying an existing service:

1. **Understand the current implementation**:
   - Read the service documentation
   - Review the Protocol Buffers definition
   - Understand the business logic and dependencies

2. **Make changes to Protocol Buffers (if needed)**:
   - Update the `.proto` file
   - Regenerate the gRPC code
   - Update server and client implementations

3. **Modify the service logic**:
   - Update the business logic in the `service/` directory
   - Update the repository layer if data access changes

4. **Test your changes**:
   - Write unit tests for new functionality
   - Test the service in isolation
   - Test integration with other services

5. **Update documentation**:
   - Update the service documentation to reflect changes
   - Update the GraphQL schema if the change affects the API

## Protocol Buffers and gRPC

The platform uses Protocol Buffers and gRPC for service-to-service communication.

### Updating Protocol Buffers

1. **Edit the .proto file**:
   Make changes to the service definition in the `.proto` file.

2. **Generate code**:
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       service/pb/service.proto
   ```

3. **Implement changes**:
   Update server and client implementations to match the new interface.

### gRPC Best Practices

- Keep messages focused and concise
- Use appropriate field types
- Consider backward compatibility when making changes
- Use consistent naming conventions
- Document all fields and methods

## Testing

The platform uses Go's testing framework for various levels of testing.

### Unit Testing

Write unit tests for business logic:

```go
func TestExample(t *testing.T) {
    // Test code here
}
```

Run unit tests:
```bash
go test ./...
```

### Integration Testing

Integration tests can be found in the `tests/` directory.

Run integration tests:
```bash
go test ./tests/...
```

### End-to-End Testing

End-to-end tests can be run against a full deployment:

```bash
docker-compose up -d
go test ./e2e/...
```

## Local Development Workflow

1. **Start the services**:
   ```bash
   docker-compose up -d
   ```

2. **Make changes to code**:
   Edit the code for the service you're working on.

3. **Rebuild the service**:
   ```bash
   docker-compose build service-name
   docker-compose up -d service-name
   ```

4. **Test your changes**:
   Use the GraphQL API or direct gRPC calls to test functionality.

5. **View logs**:
   ```bash
   docker-compose logs -f service-name
   ```

## Coding Standards

### Go Code Style

Follow the standard Go code style:
- Use `gofmt` or `go fmt` to format your code
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Run `golint` and `go vet` to catch common issues

### Code Organization

- Keep functions focused and concise
- Use meaningful names for functions, variables, and types
- Group related functionality in packages
- Use interfaces to define clean boundaries

### Error Handling

- Return errors rather than using panic
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Document error return conditions

## Documentation Standards

When writing documentation:

1. **Service Documentation**:
   - Document all service capabilities
   - Provide examples of how to use the service
   - Document API endpoints and parameters
   - Explain environment variables and configuration

2. **Code Documentation**:
   - Document all exported functions, types, and constants
   - Explain non-obvious implementations
   - Document known limitations or edge cases

3. **README files**:
   - Each service should have a README.md with basic usage information
   - Document any special setup requirements

## Common Development Tasks

### Accessing the Database

The platform uses Supabase for data storage. To access the database:

1. **Connect through Supabase UI**:
   - Log in to your Supabase account
   - Go to SQL Editor to run queries

2. **Connect through `psql`**:
   ```bash
   psql -h db.your-instance.supabase.co -U postgres -d postgres
   ```

### Debugging a Service

1. **Enable debug logs**:
   Set the environment variable `LOG_LEVEL=debug` for the service.

2. **View logs**:
   ```bash
   docker-compose logs -f service-name
   ```

3. **Debug with Delve**:
   For more advanced debugging, you can use Delve:
   ```bash
   # Run service outside of container
   dlv debug ./service/cmd/main.go
   ```

### Generating Mock Clients

For testing, you can generate mock clients:

```bash
go install github.com/golang/mock/mockgen@v1.6.0

mockgen -source=service/pb/service_grpc.pb.go -destination=service/mocks/mock_client.go -package=mocks
```

## Troubleshooting

### Common Issues

#### Protocol Buffer Generation Errors

**Issue**: `protoc-gen-go: program not found or is not executable`

**Solution**:
- Ensure `protoc-gen-go` is installed: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- Add `$GOPATH/bin` to your PATH

#### Docker Compose Errors

**Issue**: Service fails to start

**Solution**:
- Check logs: `docker-compose logs service-name`
- Ensure environment variables are set correctly
- Verify port availability

#### Build Errors

**Issue**: `go: cannot find module providing package`

**Solution**:
- Run `go mod tidy` to update dependencies
- Check for missing imports or packages

#### Service Communication Errors

**Issue**: Services can't communicate with each other

**Solution**:
- Verify service names and ports in environment variables
- Check network connectivity between containers
- Ensure gRPC clients are configured correctly 