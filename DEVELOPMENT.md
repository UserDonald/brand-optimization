# Development Guide

This guide explains how to set up your development environment and start working on the Strategic Brand Optimization Platform.

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Supabase account
- PostgreSQL (if developing locally without Docker)
- Protocol Buffers Compiler (protoc) and related Go plugins (for generating gRPC code)

## Environment Setup

1. Clone the repository:
   ```
   git clone https://github.com/your-username/go-competitor.git
   cd go-competitor
   ```

2. Create a `.env` file in the root directory with the following variables:
   ```
   SUPABASE_URL=your_supabase_url
   SUPABASE_ANON_KEY=your_supabase_anon_key
   SUPABASE_SERVICE_ROLE=your_supabase_service_role
   ```

3. Install Go dependencies:
   ```
   go mod tidy
   ```

4. Generate Protocol Buffer code:
   ```
   ./scripts/generate_proto.sh
   ```

## Database Setup

1. Create a new project in Supabase

2. Initialize the database using the SQL script:
   - Navigate to the SQL Editor in Supabase
   - Copy the contents of `docs/db/supabase_schema.sql`
   - Run the script to create tables and set up Row Level Security (RLS)

## Development Workflow

### Running Services Locally

You can run individual services locally for development:

```bash
# Run the competitor service
go run ./competitor/cmd/main.go

# Run the engagement service
go run ./engagement/cmd/main.go

# Run the GraphQL server
go run ./cmd/server/main.go
```

### Running with Docker Compose

For a full development environment, use Docker Compose:

```bash
docker-compose up
```

This will start all services and the GraphQL gateway.

### Making Changes to Protocol Buffers

If you modify any `.proto` files, regenerate the code:

```bash
./scripts/generate_proto.sh
```

### Adding a New Service

1. Create a new directory structure:
   ```
   new-service/
   ├── cmd/           # Application entry point
   ├── pb/            # Protocol Buffers definitions
   ├── client/        # Client for other services to use
   ├── repository/    # Data access layer
   ├── server/        # gRPC server implementation
   ├── service/       # Business logic
   └── Dockerfile     # Container definition
   ```

2. Create a `.proto` file in the `pb` directory
3. Generate code using the script
4. Implement the service
5. Add the service to `docker-compose.yaml`
6. Update the GraphQL schema if necessary

## Testing

Run tests for all packages:

```bash
go test ./...
```

Or for a specific package:

```bash
go test ./competitor/...
```

## Common Development Tasks

### Creating a New Repository

1. Create a repository interface in the service's `repository` directory
2. Implement the interface using Supabase
3. Make sure to handle tenant isolation correctly

### Creating a New Service Method

1. Add the method to the `.proto` file
2. Regenerate the code
3. Implement the method in the service
4. Add the gRPC handler in the server
5. Update the client as needed

### Adding a GraphQL Query or Mutation

1. Update the GraphQL schema in `graphql/schema.graphql`
2. Implement the resolver function
3. Register the resolver in the main GraphQL server

## Troubleshooting

### gRPC Connectivity Issues

If services can't connect to each other, check:
- Network configurations in Docker Compose
- Service addresses and ports
- Firewall settings

### Database Issues

For Supabase connectivity problems:
- Verify environment variables are set correctly
- Check the Supabase console for RLS policy issues
- Review database logs for query errors

### Code Generation Problems

If `protoc` fails:
- Ensure all dependencies are installed
- Check for syntax errors in `.proto` files
- Make sure `$PATH` includes the Go binary directory

## CI/CD Integration

The project uses GitHub Actions for CI/CD. When making changes, ensure:
- All tests pass
- Linting rules are followed
- Docker images build successfully

## Project Status

### Implemented Features
- Complete microservice architecture with gRPC communication between all services
- Authentication service with JWT support
- GraphQL API gateway with schema and resolvers
- Competitor tracking functionality (basic implementation)
- Content analysis capabilities (basic implementation)
- Audience segmentation features (basic implementation)
- Analytics service for data processing
- Notification system for alerts and updates
- Engagement tracking functionality
- Web scraper service for data collection
- Tenant isolation via Supabase
- Common database utilities for all services
- Health check endpoints for all services

### Current Development Focus
- All microservices have basic implementations with repositories, services, and gRPC servers
- Each service includes health check endpoints
- GraphQL API has resolvers for connecting to backend services
- Auth service provides development fallbacks for testing without Supabase
- Database repositories use Supabase for persistence

### Next Steps
1. Enhance competitor tracking with more advanced analytics capabilities
2. Improve content analysis with AI-based sentiment analysis and trend detection
3. Expand audience segmentation with targeting algorithms and custom segments
4. Build a more comprehensive analytics and recommendation engine
5. Develop more sophisticated scraping capabilities with scheduling and targeting
6. Implement real-time notification delivery through multiple channels
7. Add more comprehensive GraphQL resolvers for complex queries
8. Develop frontend application to consume the GraphQL API
9. Add comprehensive testing across all services
10. Implement monitoring and observability tools

### Testing Without Supabase
During development, you can test the services without a Supabase backend. The auth service provides fallback implementations that return dummy data, making it possible to develop and test service integrations before setting up the database. All services will show SUPABASE credential errors when run without proper credentials, which is expected and can be ignored during initial development. 