# Competitor Service

The Competitor Service is responsible for tracking and analyzing competitor data across social media platforms. It enables organizations to monitor competitors' social media activities and compare performance metrics.

## Features

- **Competitor Management**: Add, update, and delete competitor profiles
- **Competitor Metrics Tracking**: Store and retrieve metrics for competitor posts
- **Categorization**: Tag and categorize competitors for better organization
- **Multi-tenant Isolation**: Ensures data privacy between different tenants
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
competitor/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Competitor service exposes a gRPC API defined in `competitor/pb/competitor.proto`.

Key operations:
- `CreateCompetitor`: Add a new competitor to track
- `GetCompetitor`: Retrieve a specific competitor's information
- `ListCompetitors`: List all competitors for a tenant
- `UpdateCompetitor`: Update competitor information
- `DeleteCompetitor`: Remove a competitor from tracking
- `AddCompetitorMetrics`: Add performance metrics for a competitor's post
- `GetCompetitorMetrics`: Retrieve metrics for a specific competitor

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Competitor Entity

```go
type Competitor struct {
    ID        string    `json:"id"`
    TenantID  string    `json:"tenant_id"`
    Name      string    `json:"name"`
    Platform  string    `json:"platform"`
    URL       string    `json:"url"`
    Category  string    `json:"category"`
    Tags      []string  `json:"tags"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Competitor Metrics

```go
type CompetitorMetric struct {
    ID              string    `json:"id"`
    CompetitorID    string    `json:"competitor_id"`
    PostID          string    `json:"post_id"`
    Likes           int       `json:"likes"`
    Shares          int       `json:"shares"`
    Comments        int       `json:"comments"`
    ClickThroughRate float64   `json:"click_through_rate"`
    AvgWatchTime    float64   `json:"avg_watch_time"`
    EngagementRate  float64   `json:"engagement_rate"`
    PostedAt        time.Time `json:"posted_at"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9003` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |
| `AUTH_SERVICE_URL` | URL for Auth service | `localhost:9001` |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run competitor/cmd/main.go
```

Or via Docker:

```bash
docker build -t competitor-service -f competitor/Dockerfile .
docker run -p 9003:9003 competitor-service
```

### Health Check

```bash
curl http://localhost:9003/health
```

Example response:
```json
{"status":"UP"}
```

### Accessing via gRPC Client

To use the Competitor service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/competitor/client"

func main() {
    // Create competitor client
    competitorClient, err := client.NewCompetitorClient("localhost:9003")
    if err != nil {
        log.Fatalf("Failed to create competitor client: %v", err)
    }
    
    // List competitors
    competitors, err := competitorClient.ListCompetitors(ctx, tenantID)
    if err != nil {
        log.Fatalf("Failed to list competitors: %v", err)
    }
    
    // Display competitors
    for _, competitor := range competitors {
        log.Printf("Competitor: %s, Platform: %s", competitor.Name, competitor.Platform)
    }
}
```

## Data Storage

The Competitor service uses Supabase (PostgreSQL) for data storage with these tables:

1. `competitors`: Stores competitor profiles
2. `competitor_metrics`: Stores metrics for each competitor post
3. `competitor_tags`: Stores tags associated with competitors

Row Level Security (RLS) policies ensure that tenants can only access their own data.

## Error Handling

Common error scenarios:

- **Not Found**: Returned when a requested competitor does not exist
- **Validation Error**: Returned when input data does not meet validation requirements
- **Authorization Error**: Returned when a user doesn't have permission for an operation
- **Database Error**: Returned when there's an issue with the database operation

## Dependencies

- **Auth Service**: For token validation and tenant context
- **gRPC**: For service API
- **PostgreSQL**: Via Supabase for data storage
- **Scraper Service**: Data source for competitor metrics (optional integration) 