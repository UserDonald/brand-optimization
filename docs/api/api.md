# Strategic Brand Optimization Platform API Guide

This comprehensive guide documents the API for the Strategic Brand Optimization platform. The platform uses GraphQL as its primary API technology, providing a flexible and powerful interface for client applications.

## Table of Contents

1. [API Overview](#api-overview)
2. [Getting Started](#getting-started)
3. [Authentication](#authentication)
4. [GraphQL Endpoint](#graphql-endpoint)
5. [Health Check Endpoint](#health-check-endpoint)
6. [Using GraphQL](#using-graphql)
7. [Core Queries](#core-queries)
8. [Core Mutations](#core-mutations)
9. [API Concepts](#api-concepts)
10. [Service Integration](#service-integration)
11. [Error Handling](#error-handling)
12. [Rate Limiting](#rate-limiting)
13. [Best Practices](#best-practices)
14. [Examples](#examples)

## API Overview

The Strategic Brand Optimization Platform exposes a unified GraphQL API that integrates data from multiple backend microservices:

- Auth Service (Port 9001)
- Notification Service (Port 9002)
- Competitor Service (Port 9003)
- Engagement Service (Port 9004)
- Content Service (Port 9005)
- Audience Service (Port 9006)
- Analytics Service (Port 9007)
- Scraper Service (Port 9008)

The GraphQL gateway (Port 8080) aggregates these services into a coherent API, handling authentication, request routing, and response aggregation.

## Getting Started

To use the API, you'll need:

1. An account on the Strategic Brand Optimization platform
2. API credentials (JWT tokens)
3. A GraphQL client (Apollo Client, Relay, or simple HTTP requests)

Basic request flow:
1. Authenticate to receive a JWT token
2. Include the token in the Authorization header
3. Send GraphQL queries/mutations to the endpoint
4. Process the JSON responses

## Authentication

All API requests (except for the login mutation) require authentication using JWT tokens:

```
Authorization: Bearer <your_jwt_token>
```

### Obtaining a Token

Use the `login` mutation to authenticate and receive a token:

```graphql
mutation Login($email: String!, $password: String!) {
  login(email: $email, password: $password) {
    accessToken
    refreshToken
    tokenType
    expiresIn
    user {
      id
      email
      firstName
      lastName
    }
  }
}
```

Example variables:
```json
{
  "email": "user@example.com",
  "password": "your-secure-password"
}
```

### Token Refresh

When an access token expires, use the `refreshToken` mutation:

```graphql
mutation RefreshToken($refreshToken: String!) {
  refreshToken(refreshToken: $refreshToken) {
    accessToken
    refreshToken
    tokenType
    expiresIn
  }
}
```

## GraphQL Endpoint

The GraphQL API is available at:

**Production:**
```
https://api.example.com/query
```

**Local Development:**
```
http://localhost:8080/query
```

All GraphQL operations should be sent as POST requests to this endpoint with:
- `Content-Type: application/json`
- Request body containing `query`, `variables`, and optional `operationName`

Example request:
```json
{
  "query": "query GetCompetitors { competitors { id name platform } }",
  "variables": {},
  "operationName": "GetCompetitors"
}
```

## Health Check Endpoint

Each service, including the GraphQL gateway, exposes a health check endpoint:

```
GET http://localhost:8080/health
```

Response:
```json
{
  "status": "UP"
}
```

This endpoint can be used to verify service availability.

## Using GraphQL

GraphQL allows clients to request exactly the data they need. The entire API schema can be explored using tools like GraphiQL or GraphQL Playground.

### Introspection Query

To explore available types and operations:

```graphql
query IntrospectionQuery {
  __schema {
    types {
      name
      description
    }
    queryType {
      name
      fields {
        name
        description
      }
    }
    mutationType {
      name
      fields {
        name
        description
      }
    }
  }
}
```

## Core Queries

### Competitor Data

#### Get All Competitors

Retrieves all competitors for the current tenant.

```graphql
query GetCompetitors {
  competitors {
    id
    name
    platform
  }
}
```

#### Get Competitor Details

Retrieves details about a specific competitor.

```graphql
query GetCompetitor($id: ID!) {
  competitor(id: $id) {
    id
    name
    platform
  }
}
```

#### Get Competitor Metrics

Retrieves metrics for a specific competitor within a date range.

```graphql
query GetCompetitorMetrics($competitorId: ID!, $dateRange: DateRangeInput!) {
  competitorMetrics(competitorID: $competitorId, dateRange: $dateRange) {
    competitorID
    postID
    likes
    shares
    comments
    clickThroughRate
    avgWatchTime
    engagementRate
    postedAt
  }
}
```

### Personal Brand Metrics

#### Get Personal Metrics

Retrieves metrics for the tenant's own brand within a date range.

```graphql
query GetPersonalMetrics($dateRange: DateRangeInput!) {
  personalMetrics(dateRange: $dateRange) {
    postID
    likes
    shares
    comments
    clickThroughRate
    avgWatchTime
    engagementRate
    postedAt
  }
}
```

### Comparison Data

#### Compare Metrics

Compares metrics between a competitor and the tenant's brand within a date range.

```graphql
query CompareMetrics($competitorId: ID!, $dateRange: DateRangeInput!, $lockedVars: LockedVariablesInput) {
  compareMetrics(competitorID: $competitorId, dateRange: $dateRange, lockedVars: $lockedVars) {
    competitor {
      metrics {
        likes
        shares
        comments
        engagementRate
      }
      aggregates {
        totalLikes
        totalShares
        totalComments
        avgEngagementRate
        avgWatchTime
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
        totalShares
        totalComments
        avgEngagementRate
        avgWatchTime
      }
    }
    ratios {
      likesRatio
      sharesRatio
      commentsRatio
      engagementRateRatio
      watchTimeRatio
    }
  }
}
```

### Audience Insights

#### Get Audience Segments

Retrieves all audience segments for the current tenant.

```graphql
query GetAudienceSegments {
  audienceSegments {
    id
    name
    description
  }
}
```

#### Get Segment Metrics

Retrieves metrics for a specific audience segment within a date range.

```graphql
query GetSegmentMetrics($segmentId: ID!, $dateRange: DateRangeInput!) {
  audienceSegmentMetrics(segmentID: $segmentId, dateRange: $dateRange) {
    segmentID
    size
    engagementRate
    contentPreference
  }
}
```

### Content Analytics

#### Get Content Formats

Retrieves all content formats tracked by the system.

```graphql
query GetContentFormats {
  contentFormats {
    id
    name
    description
  }
}
```

#### Get Format Performance

Retrieves performance metrics for a specific content format within a date range.

```graphql
query GetFormatPerformance($formatId: ID!, $dateRange: DateRangeInput!) {
  contentFormatPerformance(formatID: $formatId, dateRange: $dateRange) {
    formatID
    engagementRate
    reachRate
    conversionRate
  }
}
```

### Recommendations

#### Get Recommended Posting Times

Retrieves AI-recommended posting times, optionally filtered by day of week.

```graphql
query GetRecommendedPostingTimes($dayOfWeek: String) {
  recommendedPostingTimes(dayOfWeek: $dayOfWeek) {
    dayOfWeek
    timeOfDay
    predictedEngagementRate
    confidence
  }
}
```

#### Get Recommended Content Formats

Retrieves AI-recommended content formats based on historical performance.

```graphql
query GetRecommendedContentFormats {
  recommendedContentFormats {
    format
    predictedEngagementRate
    targetAudience
    confidence
  }
}
```

## Core Mutations

### Competitor Management

#### Add Competitor

Adds a new competitor to track.

```graphql
mutation AddCompetitor($input: AddCompetitorInput!) {
  addCompetitor(input: $input) {
    id
    name
    platform
  }
}
```

Input format:
```graphql
input AddCompetitorInput {
  name: String!
  platform: String!
}
```

#### Update Competitor

Updates information about an existing competitor.

```graphql
mutation UpdateCompetitor($id: ID!, $input: UpdateCompetitorInput!) {
  updateCompetitor(id: $id, input: $input) {
    id
    name
    platform
  }
}
```

Input format:
```graphql
input UpdateCompetitorInput {
  name: String
  platform: String
}
```

#### Delete Competitor

Deletes a competitor from tracking.

```graphql
mutation DeleteCompetitor($id: ID!) {
  deleteCompetitor(id: $id)
}
```

### Personal Brand Management

#### Update Personal Data

Updates metrics for a personal brand post.

```graphql
mutation UpdatePersonalData($input: UpdatePersonalDataInput!) {
  updatePersonalData(input: $input)
}
```

Input format:
```graphql
input UpdatePersonalDataInput {
  postID: String!
  metrics: MetricsInput!
}

input MetricsInput {
  likes: Int
  shares: Int
  comments: Int
  clickThroughRate: Float
  avgWatchTime: Float
}
```

### Content Management

#### Schedule Post

Schedules a post for publishing.

```graphql
mutation SchedulePost($input: SchedulePostInput!) {
  schedulePost(input: $input) {
    id
    content
    scheduledTime
    platform
    format
    status
  }
}
```

Input format:
```graphql
input SchedulePostInput {
  content: String!
  scheduledTime: String!
  platform: String!
  format: String!
}
```

#### Cancel Scheduled Post

Cancels a previously scheduled post.

```graphql
mutation CancelScheduledPost($id: ID!) {
  cancelScheduledPost(id: $id)
}
```

## API Concepts

### Key Input Types

#### DateRangeInput

Used for specifying time periods for queries.

```graphql
input DateRangeInput {
  startDate: String!  # ISO 8601 format (YYYY-MM-DD)
  endDate: String!    # ISO 8601 format (YYYY-MM-DD)
}
```

#### LockedVariablesInput

Used for filtering comparison data by specific dimensions.

```graphql
input LockedVariablesInput {
  dayOfWeek: String       # e.g. "Monday", "Tuesday", etc.
  timeGap: String         # e.g. "morning", "afternoon", "evening"
  contentCategory: String # e.g. "promotional", "educational", "entertainment"
  contentFormat: String   # e.g. "video", "image", "carousel"
}
```

### Tenant Context

All API requests are executed in the context of the authenticated user's tenant. This isolation is enforced throughout the system:

1. JWT tokens contain the tenant ID
2. GraphQL resolvers include tenant context in service calls
3. Each backend service enforces tenant boundaries
4. Supabase RLS policies filter data by tenant

## Service Integration

The GraphQL gateway delegates operations to specific backend services:

- **Auth Service (9001)**: Authentication and user management
- **Notification Service (9002)**: Alerts and report generation
- **Competitor Service (9003)**: Competitor tracking and analysis
- **Engagement Service (9004)**: Engagement metrics processing
- **Content Service (9005)**: Content scheduling and optimization
- **Audience Service (9006)**: Audience segmentation and insights
- **Analytics Service (9007)**: AI-powered recommendations
- **Scraper Service (9008)**: Data collection from platforms

## Error Handling

Errors are returned in the standard GraphQL format:

```json
{
  "errors": [
    {
      "message": "Error message",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": ["fieldName"],
      "extensions": {
        "code": "ERROR_CODE"
      }
    }
  ]
}
```

### Common Error Codes

- `AUTHENTICATION_ERROR`: Authentication failed or token expired
- `AUTHORIZATION_ERROR`: User doesn't have permission for the requested operation
- `VALIDATION_ERROR`: Input validation failed
- `NOT_FOUND`: Requested resource not found
- `INTERNAL_ERROR`: Server internal error
- `RATE_LIMIT_EXCEEDED`: API rate limit reached
- `SERVICE_UNAVAILABLE`: Backend service is unavailable

### Error Handling Strategy

1. Check for `errors` in the GraphQL response
2. Examine the error `code` in extensions to determine the error type
3. Handle authentication errors by refreshing the token
4. Retry with exponential backoff for transient errors

## Rate Limiting

The API is rate limited based on your subscription tier:

| Tier | Requests per Minute |
|------|---------------------|
| Standard | 100 |
| Professional | 500 |
| Enterprise | 2000 |

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Total requests allowed per minute
- `X-RateLimit-Remaining`: Requests remaining in the current window
- `X-RateLimit-Reset`: Time (in seconds) until the rate limit resets

### Rate Limit Strategy

1. Monitor the rate limit headers in responses
2. Implement client-side throttling when approaching limits
3. Prioritize critical operations when rate limited
4. Consider upgrading your subscription tier for higher limits

## Best Practices

### Efficient Querying

1. **Request only needed fields**: GraphQL allows precise field selection
   ```graphql
   # Instead of
   query { competitor(id: "123") { id name platform url category tags createdAt updatedAt } }
   
   # Request only what you need
   query { competitor(id: "123") { id name platform } }
   ```

2. **Use pagination for large result sets**:
   ```graphql
   query {
     competitors(first: 10, after: "cursor") {
       edges {
         node {
           id
           name
         }
       }
       pageInfo {
         hasNextPage
         endCursor
       }
     }
   }
   ```

3. **Batch related queries** in a single request:
   ```graphql
   query {
     competitors {
       id
       name
     }
     audienceSegments {
       id
       name
     }
   }
   ```

### Error Handling

1. Always check for errors in responses
2. Implement token refresh for authentication errors
3. Use exponential backoff for retrying failed requests

### Security

1. Store tokens securely (HttpOnly cookies or secure storage)
2. Implement token refresh before expiry
3. Never expose tokens in client-side code or URLs

## Examples

### Complete Client Example (JavaScript/Apollo)

```javascript
import { ApolloClient, InMemoryCache, HttpLink, ApolloLink } from '@apollo/client';
import { onError } from '@apollo/client/link/error';

// Auth utils
const getAccessToken = () => localStorage.getItem('accessToken');
const setTokens = (accessToken, refreshToken) => {
  localStorage.setItem('accessToken', accessToken);
  localStorage.setItem('refreshToken', refreshToken);
};

// Create the auth link
const authLink = new ApolloLink((operation, forward) => {
  const token = getAccessToken();
  operation.setContext({
    headers: {
      authorization: token ? `Bearer ${token}` : ''
    }
  });
  return forward(operation);
});

// Error handling with token refresh
const errorLink = onError(({ graphQLErrors, networkError, operation, forward }) => {
  if (graphQLErrors) {
    for (let err of graphQLErrors) {
      if (err.extensions.code === 'AUTHENTICATION_ERROR') {
        // Perform token refresh
        return refreshToken().then(newToken => {
          setTokens(newToken.accessToken, newToken.refreshToken);
          
          // Retry with new token
          const oldHeaders = operation.getContext().headers;
          operation.setContext({
            headers: {
              ...oldHeaders,
              authorization: `Bearer ${newToken.accessToken}`
            }
          });
          return forward(operation);
        });
      }
    }
  }
});

// HTTP link
const httpLink = new HttpLink({ uri: 'http://localhost:8080/query' });

// Create Apollo Client
const client = new ApolloClient({
  link: ApolloLink.from([errorLink, authLink, httpLink]),
  cache: new InMemoryCache()
});

// Login example
async function login(email, password) {
  const result = await client.mutate({
    mutation: gql`
      mutation Login($email: String!, $password: String!) {
        login(email: $email, password: $password) {
          accessToken
          refreshToken
          expiresIn
        }
      }
    `,
    variables: { email, password }
  });
  
  const { accessToken, refreshToken } = result.data.login;
  setTokens(accessToken, refreshToken);
  return result.data.login;
}

// Example query
async function getCompetitors() {
  return client.query({
    query: gql`
      query {
        competitors {
          id
          name
          platform
        }
      }
    `
  });
}
```

### cURL Examples

#### Login and Get JWT Token

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  --data '{"query": "mutation Login($email: String!, $password: String!) { login(email: $email, password: $password) { accessToken refreshToken tokenType expiresIn } }", "variables": {"email": "user@example.com", "password": "your-password"}}'
```

#### Get Competitors List

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  --data '{"query": "query { competitors { id name platform } }"}'
```

#### Health Check

```bash
curl http://localhost:8080/health
``` 