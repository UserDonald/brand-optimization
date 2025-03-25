# Audience Service

The Audience Service is responsible for segmenting and analyzing audience data, tracking audience growth, and providing insights into audience behaviors and preferences. It enables organizations to better understand their audience and tailor content strategies accordingly.

## Features

- **Audience Segmentation**: Define and analyze different audience segments
- **Demographic Analysis**: Track demographic information about audiences
- **Audience Growth Tracking**: Monitor audience size and growth over time
- **Content Preferences**: Analyze which content resonates with different segments
- **Cross-platform Insights**: Unify audience data across social platforms
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
audience/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Audience service exposes a gRPC API defined in `audience/pb/audience.proto`.

Key operations:
- `CreateAudienceSegment`: Define a new audience segment
- `GetAudienceSegment`: Retrieve a specific audience segment
- `ListAudienceSegments`: List all audience segments for a tenant
- `UpdateAudienceSegment`: Update audience segment information
- `DeleteAudienceSegment`: Remove an audience segment
- `AddAudienceMetrics`: Record audience metrics for a time period
- `GetAudienceMetrics`: Retrieve audience metrics by segment
- `GetContentAffinity`: Find content types that resonate with a segment
- `GetAudienceGrowth`: Track audience growth over time

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Audience Segment Entity

```go
type AudienceSegment struct {
    ID          string    `json:"id"`
    TenantID    string    `json:"tenant_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Type        string    `json:"type"` // "demographic", "behavioral", "custom"
    Criteria    map[string]interface{} `json:"criteria"`
    Platform    string    `json:"platform"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Audience Metric Entity

```go
type AudienceMetric struct {
    ID              string    `json:"id"`
    TenantID        string    `json:"tenant_id"`
    SegmentID       string    `json:"segment_id"`
    Size            int       `json:"size"`
    Growth          float64   `json:"growth"` // Percentage change from previous period
    EngagementRate  float64   `json:"engagement_rate"`
    ReachRate       float64   `json:"reach_rate"`
    ConversionRate  float64   `json:"conversion_rate"`
    ContentPreference []string `json:"content_preference"` // Top content types
    Platform        string    `json:"platform"`
    PeriodStart     time.Time `json:"period_start"`
    PeriodEnd       time.Time `json:"period_end"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### Content Affinity Result

```go
type ContentAffinityResult struct {
    SegmentID       string  `json:"segment_id"`
    SegmentName     string  `json:"segment_name"`
    ContentFormat   string  `json:"content_format"`
    AffinityScore   float64 `json:"affinity_score"` // 0-1 scale
    EngagementRate  float64 `json:"engagement_rate"`
    SampleSize      int     `json:"sample_size"`
    Confidence      float64 `json:"confidence"` // 0-1 scale
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9006` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |
| `AUTH_SERVICE_URL` | URL for Auth service | `localhost:9001` |
| `CONTENT_SERVICE_URL` | URL for Content service | `localhost:9005` |
| `ENGAGEMENT_SERVICE_URL` | URL for Engagement service | `localhost:9004` |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run audience/cmd/main.go
```

Or via Docker:

```bash
docker build -t audience-service -f audience/Dockerfile .
docker run -p 9006:9006 audience-service
```

### Health Check

```bash
curl http://localhost:9006/health
```

Example response:
```json
{"status":"UP"}
```

### Accessing via gRPC Client

To use the Audience service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/audience/client"

func main() {
    // Create audience client
    audienceClient, err := client.NewAudienceClient("localhost:9006")
    if err != nil {
        log.Fatalf("Failed to create audience client: %v", err)
    }
    
    // Create a new audience segment
    segment, err := audienceClient.CreateAudienceSegment(ctx, &pb.CreateAudienceSegmentRequest{
        TenantId:    tenantID,
        Name:        "Active Engagers",
        Description: "Users who engage frequently with content",
        Type:        "behavioral",
        Criteria: map[string]string{
            "min_engagements": "5",
            "timeframe":      "30_days",
        },
        Platform:    "instagram",
    })
    if err != nil {
        log.Fatalf("Failed to create audience segment: %v", err)
    }
    
    log.Printf("Created segment with ID: %s", segment.Id)
}
```

## Segmentation Algorithms

The Audience service uses several algorithms to create and analyze audience segments:

1. **RFM Analysis** (Recency, Frequency, Monetary) - For behavioral segmentation
2. **Clustering** - For identifying natural segments in audience data
3. **Engagement Patterns** - For categorizing users by their engagement style
4. **Demographic Grouping** - For segmenting by demographic attributes

Each algorithm can be customized with specific parameters through the segment criteria field.

## Data Storage

The Audience service uses Supabase (PostgreSQL) for data storage with these tables:

1. `audience_segments`: Stores segment definitions
2. `audience_metrics`: Stores audience metrics over time
3. `content_affinity`: Stores relationships between segments and content
4. `segment_membership`: Maps individual audience members to segments

Row Level Security (RLS) policies ensure that tenants can only access their own data.

## Error Handling

Common error scenarios:

- **Not Found**: Returned when a requested segment does not exist
- **Validation Error**: Returned when input data does not meet validation requirements
- **Insufficient Data**: Returned when there's not enough data for reliable analysis
- **Authorization Error**: Returned when a user doesn't have permission for an operation
- **Database Error**: Returned when there's an issue with the database operation

## Dependencies

- **Auth Service**: For token validation and tenant context
- **Content Service**: For content correlation
- **Engagement Service**: For engagement data correlation
- **gRPC**: For service API
- **PostgreSQL**: Via Supabase for data storage 