# Strategic Brand Optimization Platform

A Go-based microservice architecture for tracking, analyzing, and optimizing brand performance across social media platforms compared to competitors. This system enables organizations to make data-driven decisions using AI-powered insights and interactive dashboards.

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Services](#services)
5. [Data Model](#data-model)
6. [Tenant Isolation](#tenant-isolation)
7. [GraphQL API](#graphql-api)
8. [Setup and Installation](#setup-and-installation)
9. [Development](#development)
10. [Deployment](#deployment)
11. [Documentation](#documentation)
12. [Contributing](#contributing)

## Overview

The Strategic Brand Optimization Platform focuses on data preparation and competitor tracking in a B2B scenario. Each business client has:
- Their own brand data (posting frequency, engagement metrics, audience segments, etc.)
- Competitor data relevant to their market or niche

By adopting a Go microservice design with GraphQL gateway and Supabase for data storage, this platform:
- Supports tenant isolation through Row Level Security (RLS)
- Collects and normalizes social media insights via scraping tasks
- Provides real-time analytics and AI-powered recommendations
- Exposes structured data through a GraphQL API for visualizations and side-by-side comparisons

## Features

- **Competitor Tracking**: Monitor competitors' social media performance with metrics like engagement rates, posting frequency, and content formats
- **Side-by-Side Comparison**: Compare your brand's performance against competitors with ratio calculations and trend analysis
- **Audience Segmentation**: Understand different audience segments (Passive Viewers, Reactors, Conversationalists, Content Creators) and how they engage
- **Content Optimization**: Analyze which content formats and posting times generate the highest engagement
- **Predictive Analytics**: AI-powered recommendations for optimal posting times and content formats
- **Scheduled Posting**: Plan and schedule content based on data-driven insights
- **Multi-tenant Architecture**: Secure isolation of data between different clients using Supabase RLS
- **Health Monitoring**: Each service exposes a health check endpoint for monitoring and reliability
- **Interactive Dashboard**: User-friendly interface to visualize and act on insights

## Architecture

The platform follows a microservice architecture with the following components:

```
go-competitor/
├── auth/              # Authentication & user management (port 9001)
├── notification/      # Alerts and scheduled notifications (port 9002)
├── competitor/        # Competitor data tracking & analysis (port 9003)
├── engagement/        # Metrics & engagement analytics (port 9004)
├── content/           # Content scheduling & optimization (port 9005)
├── audience/          # Audience segmentation & insights (port 9006)
├── analytics/         # AI-powered predictive analytics (port 9007)
├── scraper/           # Data collection from social platforms (port 9008)
├── graphql/           # GraphQL gateway & schema definitions (port 8080)
├── common/            # Shared utilities and middleware
└── docs/              # Documentation
```

Each service follows a consistent structure:

```
service/
├── cmd/               # Command-line entrypoints
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer
```

Each service is containerized and can be deployed independently, communicating through gRPC APIs.

## Services

1. **Auth Service (Port 9001)**: Manages authentication, authorization, and tenant management with JWT tokens
2. **Notification Service (Port 9002)**: Sends alerts and generates scheduled reports through various channels
3. **Competitor Service (Port 9003)**: Tracks and analyzes competitor data across social platforms
4. **Engagement Service (Port 9004)**: Monitors and compares engagement metrics over time
5. **Content Service (Port 9005)**: Analyzes content formats and optimizes scheduling
6. **Audience Service (Port 9006)**: Segments audiences and tracks their behaviors
7. **Analytics Service (Port 9007)**: Provides AI-powered recommendations and predictions
8. **Scraper Service (Port 9008)**: Collects data from social media platforms
9. **GraphQL Gateway (Port 8080)**: Unified API entry point for frontend applications

## Data Model

The system uses Supabase (PostgreSQL) for structured data storage with the following primary tables:

- **Organizations**: Represents tenants/clients using the platform
- **Competitors**: Stores information about competitors being tracked
- **Competitor Metrics**: Records engagement metrics for competitor posts
- **Personal Metrics**: Stores metrics for the client's own social media posts
- **Audience Segments**: Defines different audience segments and their characteristics
- **Content Formats**: Tracks performance of different content formats
- **Scheduled Posts**: Manages content scheduling and posting status

For detailed schema information, see [Supabase Schema Documentation](docs/db/supabase_schema.sql).

## Tenant Isolation

Data security and tenant isolation are achieved through Supabase's Row Level Security (RLS) policies:

1. **Authentication**: Users are authenticated through Supabase Auth
2. **Tenant Association**: Each user is associated with a specific tenant (organization)
3. **RLS Policies**: Database tables have RLS policies that filter data based on the authenticated user's tenant ID
4. **Service Requests**: All service-to-service requests include tenant context

Example RLS policy:
```sql
create policy "Users can only access their tenant's competitors"
on public.competitors
for all
using (auth.uid() = tenant_id);
```

## GraphQL API

The platform exposes a unified GraphQL API for clients to interact with the system. Key query types include:

- Competitor data and metrics
- Personal brand performance
- Side-by-side comparisons
- Audience segment analysis
- Content format performance
- Posting time recommendations

Example query:
```graphql
query CompareMetrics($competitorId: ID!, $dateRange: DateRangeInput!) {
  compareMetrics(competitorId: $competitorId, dateRange: $dateRange) {
    competitor {
      metrics {
        likes
        shares
        comments
        engagementRate
      }
      aggregates {
        totalLikes
        avgEngagementRate
      }
    }
    personal {
      metrics {
        likes
        shares
        comments
        engagementRate
      }
      aggregates {
        totalLikes
        avgEngagementRate
      }
    }
    ratios {
      likesRatio
      sharesRatio
      engagementRateRatio
    }
  }
}
```

For the complete schema, see [GraphQL Schema](graphql/schema.graphql) and [API Documentation](docs/api/api.md).

## Setup and Installation

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Supabase account

### Environment Variables

Create a `.env` file with the following variables:

```
SUPABASE_URL=your_supabase_url
SUPABASE_ANON_KEY=your_supabase_anon_key
SUPABASE_SERVICE_ROLE=your_supabase_service_role
JWT_SECRET=your_jwt_secret
```

### Local Development

1. Clone the repository:
   ```
   git clone https://github.com/your-username/go-competitor.git
   cd go-competitor
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Start the services with Docker Compose:
   ```
   docker-compose up
   ```

4. Verify the health of services:
   ```
   curl http://localhost:9001/health  # Auth service
   curl http://localhost:8080/health  # GraphQL gateway
   ```

5. The GraphQL API will be available at `http://localhost:8080/query`

For a more detailed guide, see the [Development Guide](docs/development/development.md).

## Development

### Adding a New Service

1. Create a new directory for your service
2. Implement the service logic
3. Add the service to `docker-compose.yaml`
4. Update GraphQL schema if necessary

### Testing

Run tests with:
```
go test ./...
```

## Deployment

The platform can be deployed to various cloud environments:

### Kubernetes

Kubernetes deployment manifests are available in the `deployment/kubernetes` directory.

### Docker Swarm

For simpler deployments, Docker Swarm can be used:
```
docker stack deploy -c docker-compose.yaml brand-optimization
```

For detailed deployment instructions, see our [Deployment Guide](docs/deployment/deployment.md).

## Documentation

Comprehensive documentation is available in the `docs` directory:

- [Architecture Overview](docs/architecture/overview.md)
- [API Documentation](docs/api/api.md)
- [Service Documentation](docs/services/README.md)
- [Deployment Guide](docs/deployment/deployment.md)
- [Development Guide](docs/development/development.md)

## Contributing

We welcome contributions via pull requests and feature proposals. Key areas include:
- Additional microservices for specific analytics needs
- Integration with more social media platforms
- Advanced ML/AI features for better predictions
- Improved data visualization tools or front-end integration

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for more details.
