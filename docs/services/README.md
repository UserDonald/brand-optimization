# Microservice Documentation

The Strategic Brand Optimization Platform is composed of the following microservices:

## [GraphQL Gateway](./graphql.md)
API gateway service that provides a unified GraphQL interface for client applications.

## [Auth Service](./auth.md)
Authentication and user management service.

## [Competitor Service](./competitor.md)
Service for tracking and analyzing competitor data.

## [Engagement Service](./engagement.md)
Service for monitoring and analyzing engagement metrics.

## [Content Service](./content.md)
Service for managing content formats and scheduling.

## [Audience Service](./audience.md)
Service for audience segmentation and analysis.

## [Analytics Service](./analytics.md)
Service providing AI-powered recommendations and predictions.

## [Notification Service](./notification.md)
Service for sending alerts and notifications.

## [Scraper Service](./scraper.md)
Service for collecting data from social media platforms.

## Service Architecture Pattern

Each service follows a consistent architecture pattern:

```
service/
├── cmd/             # Command-line entrypoints
├── server/          # gRPC server implementation
├── client/          # gRPC client for other services to use
├── pb/              # Protocol Buffers definitions
├── service/         # Business logic
└── repository/      # Data access layer
```

This architecture promotes:
- Clear separation of concerns
- Independently deployable components
- Reusable clients for inter-service communication
- Contract-driven development with Protocol Buffers 