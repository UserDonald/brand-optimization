# GraphQL Gateway Service

The GraphQL Gateway serves as the central entry point for client applications to communicate with the Strategic Brand Optimization Platform. It aggregates data from various microservices and presents a unified GraphQL API, simplifying client development.

## Features

- **Unified API**: Single endpoint for all client operations
- **Authentication Handling**: JWT-based authentication and tenant isolation
- **Request Delegation**: Forwards requests to appropriate microservices
- **Response Aggregation**: Combines data from multiple services into cohesive responses
- **CORS Support**: Cross-Origin Resource Sharing headers for browser access
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
graphql/
├── cmd/               # Entry point for the service
├── server/            # GraphQL server implementation
├── resolvers/         # GraphQL resolvers for queries and mutations
├── middleware/        # HTTP middleware (auth, logging, etc.)
├── models/            # Data models used by the service
└── schema.graphql     # GraphQL schema definition
```

## API Endpoints

- **GraphQL API**: `/query` - Main GraphQL endpoint
- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### GraphQL Schema

The GraphQL schema defines all types, queries, and mutations available through the API. It is located at `graphql/schema.graphql`. For full schema documentation, see the [GraphQL API Documentation](../api/graphql.md).

### Authentication Flow

1. Client authenticates using the `login` or `register` mutation
2. Auth service returns JWT tokens in response
3. Client includes token in `Authorization` header for subsequent requests
4. GraphQL gateway validates token with Auth service
5. User and tenant context is added to each request
6. Resolvers use this context to enforce tenant isolation

### Service Integration

The GraphQL gateway communicates with backend microservices using:

1. **gRPC**: For efficient binary communication with most services
2. **Context Propagation**: Tenant context is passed to all service calls
3. **Client Pooling**: Connection pooling for efficient resource usage

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GRAPHQL_PORT` | Port to listen on | `8080` |
| `AUTH_SERVICE_URL` | URL for Auth service | `localhost:9001` |
| `NOTIFICATION_SERVICE_URL` | URL for Notification service | `localhost:9002` |
| `COMPETITOR_SERVICE_URL` | URL for Competitor service | `localhost:9003` |
| `ENGAGEMENT_SERVICE_URL` | URL for Engagement service | `localhost:9004` |
| `CONTENT_SERVICE_URL` | URL for Content service | `localhost:9005` |
| `AUDIENCE_SERVICE_URL` | URL for Audience service | `localhost:9006` |
| `ANALYTICS_SERVICE_URL` | URL for Analytics service | `localhost:9007` |
| `SCRAPER_SERVICE_URL` | URL for Scraper service | `localhost:9008` |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run graphql/cmd/main.go
```

Or via Docker:

```bash
docker build -t graphql-service -f graphql/Dockerfile .
docker run -p 8080:8080 graphql-service
```

### Health Check

```bash
curl http://localhost:8080/health
```

Example response:
```json
{"status":"UP"}
```

### GraphQL Introspection

```bash
curl -X POST -H "Content-Type: application/json" --data '{"query": "{ __schema { types { name } } }"}' http://localhost:8080/query
```

## Error Handling

The service returns standard GraphQL errors with the following error codes:

- `AUTHENTICATION_ERROR`: Authentication failed or token expired
- `AUTHORIZATION_ERROR`: User doesn't have permission for the requested operation
- `VALIDATION_ERROR`: Input validation failed
- `NOT_FOUND`: Requested resource not found
- `INTERNAL_ERROR`: Server internal error

## Dependencies

- **graph-gophers/graphql-go**: GraphQL implementation for Go
- **Auth Service**: Required for token validation and user information
- **Service Clients**: Generated gRPC clients for each microservice 