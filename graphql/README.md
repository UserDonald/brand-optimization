# GraphQL API Gateway

This service acts as the central API gateway for the Go Competitor platform. It provides a unified GraphQL API that integrates all microservices into a single, consistent interface.

## Purpose

The GraphQL API Gateway:

- Presents a single endpoint for frontend clients
- Handles authentication and authorization
- Routes requests to appropriate microservices
- Aggregates data from multiple services into cohesive responses
- Enforces tenant isolation for multi-tenancy

## Architecture

The service follows a clean architecture pattern:

- `cmd/` - Entry point and application setup
- `server/` - Core server implementation
- `middleware/` - HTTP middleware (auth, logging, etc.)
- `resolvers/` - GraphQL resolvers for each domain
- `models/` - GraphQL response types
- `schema.graphql` - GraphQL schema definition

## Development

### Running Locally

```bash
cd cmd
go run main.go
```

The server will start on port 8080 by default. You can override this with the `GRAPHQL_PORT` environment variable.

### GraphQL Endpoint

The main GraphQL endpoint is available at:

```
http://localhost:8080/query
```

You can use tools like GraphiQL, Insomnia, or Postman to interact with the API.

### Health Check

A simple health check endpoint is available at:

```
http://localhost:8080/health
```

## Authentication

The service uses JWT authentication. Include a Bearer token in the Authorization header:

```
Authorization: Bearer <your-token>
```

Tokens can be obtained through the auth service's login endpoint.

## Service Dependencies

This gateway connects to the following microservices:

- Auth Service - User authentication and management
- Competitor Service - Competitor tracking and analysis
- Content Service - Content management and scheduling
- Audience Service - Audience segmentation
- Analytics Service - Performance analytics and recommendations
- Engagement Service - Engagement tracking and metrics
- Notification Service - User notifications
- Scraper Service - Data collection from external sources

## Schema Development

The GraphQL schema is defined in `schema.graphql`. As the platform evolves, new types and operations will be added to this schema.

The current implementation uses a simplified schema for initial development, with comments indicating the full schema that will be implemented incrementally as the services mature. 