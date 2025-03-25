# Engagement Service

The Engagement Service is responsible for tracking and analyzing engagement metrics for both the client's own social media posts and competitor posts. It provides insights into audience interaction and helps identify patterns in engagement performance.

## Features

- **Engagement Metrics Tracking**: Store and retrieve engagement metrics (likes, shares, comments, etc.)
- **Comparative Analysis**: Compare engagement metrics over time and between competitors
- **Derived Metrics Calculation**: Calculate derived metrics such as engagement rates
- **Time-Series Analysis**: Track metric trends over different time periods
- **Multi-tenant Isolation**: Ensures data privacy between different tenants
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
engagement/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Engagement service exposes a gRPC API defined in `engagement/pb/engagement.proto`.

Key operations:
- `AddEngagementMetric`: Record a new engagement metric for a post
- `GetEngagementMetric`: Retrieve a specific engagement metric
- `ListEngagementMetrics`: List engagement metrics by various filters
- `GetComparisonMetrics`: Compare metrics between tenant's posts and competitor posts
- `CalculateEngagementRate`: Calculate the engagement rate for a post
- `GetEngagementTrends`: Retrieve trend analysis for engagement metrics

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Engagement Metric Entity

```go
type EngagementMetric struct {
    ID              string    `json:"id"`
    TenantID        string    `json:"tenant_id"`
    PostID          string    `json:"post_id"`
    Source          string    `json:"source"` // "own" or "competitor"
    SourceID        string    `json:"source_id"` // CompetitorID if source is "competitor"
    Likes           int       `json:"likes"`
    Shares          int       `json:"shares"`
    Comments        int       `json:"comments"`
    Clicks          int       `json:"clicks"`
    Impressions     int       `json:"impressions"`
    Reach           int       `json:"reach"`
    ClickThroughRate float64   `json:"click_through_rate"`
    AvgWatchTime    float64   `json:"avg_watch_time"`
    EngagementRate  float64   `json:"engagement_rate"`
    PostedAt        time.Time `json:"posted_at"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### Comparison Result

```go
type ComparisonResult struct {
    Own struct {
        Metrics    []EngagementMetric `json:"metrics"`
        Aggregates MetricAggregates   `json:"aggregates"`
    } `json:"own"`
    Competitor struct {
        Metrics    []EngagementMetric `json:"metrics"`
        Aggregates MetricAggregates   `json:"aggregates"`
    } `json:"competitor"`
    Ratios struct {
        LikesRatio          float64 `json:"likes_ratio"`
        SharesRatio         float64 `json:"shares_ratio"`
        CommentsRatio       float64 `json:"comments_ratio"`
        EngagementRateRatio float64 `json:"engagement_rate_ratio"`
    } `json:"ratios"`
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9004` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |
| `AUTH_SERVICE_URL` | URL for Auth service | `localhost:9001` |
| `COMPETITOR_SERVICE_URL` | URL for Competitor service | `localhost:9003` |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run engagement/cmd/main.go
```

Or via Docker:

```bash
docker build -t engagement-service -f engagement/Dockerfile .
docker run -p 9004:9004 engagement-service
```

### Health Check

```bash
curl http://localhost:9004/health
```

Example response:
```json
{"status":"UP"}
```

### Accessing via gRPC Client

To use the Engagement service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/engagement/client"

func main() {
    // Create engagement client
    engagementClient, err := client.NewEngagementClient("localhost:9004")
    if err != nil {
        log.Fatalf("Failed to create engagement client: %v", err)
    }
    
    // Get comparison metrics
    result, err := engagementClient.GetComparisonMetrics(ctx, tenantID, competitorID, startDate, endDate)
    if err != nil {
        log.Fatalf("Failed to get comparison metrics: %v", err)
    }
    
    // Process results
    log.Printf("Own engagement rate: %.2f, Competitor engagement rate: %.2f, Ratio: %.2f",
        result.Own.Aggregates.AvgEngagementRate,
        result.Competitor.Aggregates.AvgEngagementRate,
        result.Ratios.EngagementRateRatio)
}
```

## Calculation Methods

### Engagement Rate Calculation

The service uses different formulas for calculating engagement rates based on the platform:

- **Standard Engagement Rate**:
  `(Likes + Comments + Shares) / Impressions * 100`

- **Extended Engagement Rate**:
  `(Likes + Comments + Shares + Clicks + Saves) / Impressions * 100`

- **Watch Time Engagement**:
  `(Total Watch Time / (Video Duration * Views)) * 100`

## Data Storage

The Engagement service uses Supabase (PostgreSQL) for data storage with these tables:

1. `engagement_metrics`: Stores all engagement metrics
2. `engagement_trends`: Stores pre-calculated trend data

Row Level Security (RLS) policies ensure that tenants can only access their own data.

## Error Handling

Common error scenarios:

- **Not Found**: Returned when a requested metric does not exist
- **Validation Error**: Returned when input data does not meet validation requirements
- **Authorization Error**: Returned when a user doesn't have permission for an operation
- **Database Error**: Returned when there's an issue with the database operation

## Dependencies

- **Auth Service**: For token validation and tenant context
- **Competitor Service**: For competitor data correlation
- **gRPC**: For service API
- **PostgreSQL**: Via Supabase for data storage 