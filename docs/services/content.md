# Content Service

The Content Service is responsible for managing content formats, categorizing content, and scheduling posts. It enables organizations to optimize their content strategy based on historical performance and analytics recommendations.

## Features

- **Content Format Management**: Define and track different content formats (video, carousel, etc.)
- **Content Categorization**: Categorize content for better analysis
- **Content Scheduling**: Plan and schedule posts for optimal engagement
- **A/B Testing Support**: Track variants of content to determine effectiveness
- **Performance Tracking**: Associate content with engagement metrics
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
content/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Content service exposes a gRPC API defined in `content/pb/content.proto`.

Key operations:
- `CreateContentFormat`: Define a new content format
- `GetContentFormat`: Retrieve a specific content format
- `ListContentFormats`: List all content formats for a tenant
- `UpdateContentFormat`: Update content format information
- `DeleteContentFormat`: Remove a content format
- `SchedulePost`: Schedule a post for publishing
- `GetScheduledPost`: Retrieve a specific scheduled post
- `ListScheduledPosts`: List all scheduled posts with filters
- `UpdateScheduledPost`: Update a scheduled post
- `CancelScheduledPost`: Cancel a scheduled post

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Content Format Entity

```go
type ContentFormat struct {
    ID          string    `json:"id"`
    TenantID    string    `json:"tenant_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Platform    string    `json:"platform"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Scheduled Post Entity

```go
type ScheduledPost struct {
    ID           string    `json:"id"`
    TenantID     string    `json:"tenant_id"`
    Content      string    `json:"content"`
    MediaURLs    []string  `json:"media_urls"`
    Platform     string    `json:"platform"`
    FormatID     string    `json:"format_id"`
    Category     string    `json:"category"`
    Tags         []string  `json:"tags"`
    ScheduledAt  time.Time `json:"scheduled_at"`
    PublishedAt  *time.Time `json:"published_at"`
    Status       string    `json:"status"` // "scheduled", "published", "cancelled", "failed"
    ErrorMessage string    `json:"error_message"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Content Performance Result

```go
type ContentPerformanceResult struct {
    FormatID        string  `json:"format_id"`
    FormatName      string  `json:"format_name"`
    PostCount       int     `json:"post_count"`
    AvgEngagement   float64 `json:"avg_engagement"`
    AvgReach        float64 `json:"avg_reach"`
    AvgConversion   float64 `json:"avg_conversion"`
    TopPerformers   []string `json:"top_performers"` // IDs of top performing posts
    RecommendedTime string  `json:"recommended_time"` // Day and time for best performance
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9005` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |
| `AUTH_SERVICE_URL` | URL for Auth service | `localhost:9001` |
| `ENGAGEMENT_SERVICE_URL` | URL for Engagement service | `localhost:9004` |
| `ANALYTICS_SERVICE_URL` | URL for Analytics service | `localhost:9007` |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run content/cmd/main.go
```

Or via Docker:

```bash
docker build -t content-service -f content/Dockerfile .
docker run -p 9005:9005 content-service
```

### Health Check

```bash
curl http://localhost:9005/health
```

Example response:
```json
{"status":"UP"}
```

### Accessing via gRPC Client

To use the Content service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/content/client"

func main() {
    // Create content client
    contentClient, err := client.NewContentClient("localhost:9005")
    if err != nil {
        log.Fatalf("Failed to create content client: %v", err)
    }
    
    // Schedule a post
    post, err := contentClient.SchedulePost(ctx, &pb.SchedulePostRequest{
        TenantId:    tenantID,
        Content:     "Check out our new product launch!",
        MediaUrls:   []string{"https://example.com/image.jpg"},
        Platform:    "instagram",
        FormatId:    "carousel-format-id",
        Category:    "product-launch",
        Tags:        []string{"newproduct", "launch"},
        ScheduledAt: timestamppb.New(scheduledTime),
    })
    if err != nil {
        log.Fatalf("Failed to schedule post: %v", err)
    }
    
    log.Printf("Post scheduled with ID: %s for %s", post.Id, scheduledTime.Format(time.RFC3339))
}
```

## Scheduling Mechanism

The Content service uses a background worker to manage the posting schedule:

1. Posts are stored in the database with status "scheduled"
2. A background worker periodically checks for posts that are due to be published
3. When a post is due, the worker changes its status to "processing"
4. The worker integrates with social media APIs to publish the post
5. Once published, the status is updated to "published" with the actual publish time

In case of failure, the post status is set to "failed" with an error message, and notifications are sent via the Notification service.

## Data Storage

The Content service uses Supabase (PostgreSQL) for data storage with these tables:

1. `content_formats`: Stores format definitions
2. `scheduled_posts`: Stores scheduled and published posts
3. `content_categories`: Stores content categorization taxonomy
4. `content_performance`: Stores pre-calculated performance metrics by format

Row Level Security (RLS) policies ensure that tenants can only access their own data.

## Error Handling

Common error scenarios:

- **Not Found**: Returned when a requested format or post does not exist
- **Validation Error**: Returned when input data does not meet validation requirements
- **Scheduling Error**: Returned when there's an issue with scheduling a post
- **Authorization Error**: Returned when a user doesn't have permission for an operation
- **Database Error**: Returned when there's an issue with the database operation

## Dependencies

- **Auth Service**: For token validation and tenant context
- **Engagement Service**: For performance data correlation
- **Analytics Service**: For recommended posting times
- **gRPC**: For service API
- **PostgreSQL**: Via Supabase for data storage 