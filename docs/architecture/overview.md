# Architecture Overview

The Strategic Brand Optimization platform follows a microservices architecture pattern with clear separation of concerns and domain-driven design principles. This document outlines the high-level architecture, communication patterns, and design decisions.

## Architectural Principles

1. **Domain-Driven Design**: Services are organized around business capabilities and domains
2. **Microservice Independence**: Each service can be developed, deployed, and scaled independently
3. **Data Sovereignty**: Each service owns its data and exposes it through well-defined APIs
4. **API Gateway Pattern**: GraphQL serves as the centralized API gateway for client applications
5. **Multi-tenancy**: All services and data are tenant-aware through Supabase RLS
6. **Event-Driven**: Certain workflows use event-driven patterns for asynchronous operations
7. **Health Monitoring**: All services expose health check endpoints for monitoring and reliability

## System Components

### Core Services

![Architecture Diagram](../diagrams/architecture.png)

1. **GraphQL Gateway (Port 8080)**
   - Serves as the entry point for all client applications
   - Handles authentication and tenant context
   - Delegates requests to appropriate microservices
   - Aggregates responses into cohesive GraphQL responses
   - Exposes a health check endpoint at `/health`

2. **Auth Service (Port 9001)**
   - Manages user authentication via JWT tokens
   - Handles tenant association and permissions
   - Provides middleware for token validation
   - Manages organization/tenant data
   - Exposes a health check endpoint at `/health`

3. **Notification Service (Port 9002)**
   - Manages alert thresholds and preferences
   - Sends notifications via various channels
   - Generates scheduled reports
   - Handles event subscriptions
   - Exposes a health check endpoint at `/health`

4. **Competitor Service (Port 9003)**
   - Stores competitor profiles and relationships
   - Tracks competitor metrics over time
   - Provides comparison capabilities
   - Manages competitor categories and tags
   - Exposes a health check endpoint at `/health`

5. **Engagement Service (Port 9004)**
   - Tracks all engagement metrics (likes, shares, comments, etc.)
   - Calculates derived metrics (engagement rates, ratios)
   - Provides time-series analysis capabilities
   - Supports filtering by various dimensions
   - Exposes a health check endpoint at `/health`

6. **Content Service (Port 9005)**
   - Manages content formats and categorization
   - Handles content scheduling and posting
   - Tracks content performance metrics
   - Provides A/B testing capabilities
   - Exposes a health check endpoint at `/health`

7. **Audience Service (Port 9006)**
   - Defines audience segments and their behaviors
   - Tracks audience growth and engagement patterns
   - Maps content preferences to segments
   - Provides demographic insights
   - Exposes a health check endpoint at `/health`

8. **Analytics Service (Port 9007)**
   - Generates AI-powered recommendations
   - Predicts optimal posting times and formats
   - Identifies trends and patterns
   - Provides "what-if" scenario analysis
   - Exposes a health check endpoint at `/health`

9. **Scraper Service (Port 9008)**
   - Integrates with social media APIs
   - Adheres to rate limits and quotas
   - Normalizes data from different platforms
   - Schedules and executes data collection jobs
   - Exposes a health check endpoint at `/health`

### Service Architecture Pattern

Each service follows a consistent architecture pattern:

```
service/
├── cmd/               # Command-line entrypoints
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer
```

This pattern promotes:
- Clear separation of concerns
- Independently deployable components
- Reusable clients for inter-service communication
- Contract-driven development with Protocol Buffers

### Infrastructure Components

1. **Supabase (PostgreSQL + RLS)**
   - Primary data store with row-level security
   - Tenant isolation through RLS policies
   - Real-time capabilities through subscriptions
   - Edge functions for serverless processing

2. **Redis**
   - Used for caching frequently accessed data
   - Manages distributed locks for scheduled tasks
   - Provides pub/sub capabilities for event broadcasting
   - Stores session data and real-time analytics

3. **Docker & Kubernetes/Swarm**
   - Containerization of all services
   - Orchestration for scaling and management
   - Blue/green deployment strategies
   - Auto-scaling based on load

## Communication Patterns

### Synchronous Communication

1. **Service-to-Service**: gRPC for efficient binary communication
2. **Client-to-Gateway**: GraphQL over HTTP/HTTPS
3. **Gateway-to-Service**: gRPC 

### Asynchronous Communication

1. **Event Broadcasting**: Redis pub/sub for lightweight events
2. **Workflow Orchestration**: Temporal for complex workflows
3. **Background Jobs**: Worker pools for CPU-intensive tasks

## Tenant Isolation Strategy

Tenant isolation is implemented at multiple levels:

1. **Authentication Layer**:
   - JWT tokens contain tenant ID
   - Auth middleware validates tenant context

2. **Data Access Layer**:
   - Supabase RLS policies enforce tenant boundaries
   - All database queries include tenant filter

3. **Service Layer**:
   - Services validate tenant context in requests
   - Cross-tenant operations are explicitly prohibited

4. **API Gateway**:
   - GraphQL resolvers include tenant context
   - Field-level authorization for sensitive data

## Scalability Considerations

1. **Horizontal Scaling**:
   - Stateless services can scale horizontally
   - Database connection pooling for efficiency

2. **Caching Strategy**:
   - Multi-level caching (application, Redis, database)
   - Cache invalidation based on write operations

3. **Database Optimization**:
   - Indexing strategy for tenant-aware queries
   - Partitioning for large tenants if needed

4. **Background Processing**:
   - Offload intensive operations to worker pools
   - Rate limiting for scraping operations

## Security Measures

1. **Authentication**:
   - JWT-based authentication
   - Short-lived tokens with refresh capability

2. **Authorization**:
   - Role-based access control within tenants
   - Row-level security for data access

3. **Transport Security**:
   - TLS for all service communication
   - API gateway with rate limiting and DDoS protection

4. **Data Protection**:
   - Encryption for sensitive data
   - Audit logging for key operations

## Health Monitoring and Observability

1. **Health Checks**:
   - Each service exposes a `/health` endpoint
   - Returns status information in a standardized format
   - Used by container orchestration for service status

2. **Metrics Collection**:
   - Service-level metrics (requests, latency, errors)
   - Business metrics (active tenants, scraping jobs, etc.)

3. **Distributed Tracing**:
   - OpenTelemetry for end-to-end request tracing
   - Correlation IDs across service boundaries

4. **Logging**:
   - Structured logging with tenant context
   - Centralized log aggregation

5. **Alerting**:
   - SLO-based alerting for critical services
   - Business-level alerts for tenant-specific issues

## Deployment Pipeline

1. **CI/CD Workflow**:
   - Automated testing for each service
   - Integration tests for cross-service functionality
   - Container building and versioning

2. **Environment Promotion**:
   - Development → Staging → Production
   - Feature flags for controlled rollouts

3. **Database Migrations**:
   - Version-controlled migration scripts
   - Zero-downtime migration strategy

4. **Rollback Procedures**:
   - Automated rollback for failed deployments
   - Database migration rollback capabilities

## Future Considerations

1. **Multi-region Deployment**:
   - Geographic distribution for lower latency
   - Data residency considerations

2. **Enhanced ML Capabilities**:
   - Integration with specialized ML services
   - Real-time recommendation engines

3. **Additional Social Platforms**:
   - Expand scraper support to emerging platforms
   - Normalize metrics across more sources

4. **Advanced Visualization**:
   - Interactive dashboards with drill-down capabilities
   - Custom report generation 