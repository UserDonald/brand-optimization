# Scraper Service

The Scraper Service is responsible for collecting data from social media platforms to track competitor activity and market trends. It integrates with various social media APIs, respects rate limits, and normalizes data for use by other services in the platform.

## Features

- **Social Media Integration**: Connect to multiple social media platforms (Instagram, Twitter, Facebook, LinkedIn, TikTok)
- **Scheduled Scraping**: Automatically collect data on predefined schedules using cron expressions
- **Rate Limit Management**: Respect API rate limits and quotas
- **Data Normalization**: Convert platform-specific data to a standardized format
- **Incremental Collection**: Only collect new or updated data
- **Error Recovery**: Gracefully handle API failures and retry operations
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
scraper/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Scraper service exposes a gRPC API defined in `scraper/pb/scraper.proto`.

Key operations:
- **Job Management**
  - `CreateScraperJob`: Create a new scraper job with scheduling
  - `GetScraperJob`: Retrieve information about a specific job
  - `ListScraperJobs`: List all scraper jobs with optional filters
  - `CancelScraperJob`: Cancel a scheduled or running job
  - `DeleteScraperJob`: Remove a scraper job from the system

- **Platform Operations**
  - `ListSupportedPlatforms`: List all platforms supported by the scraper
  - `GetPlatformStatus`: Check API status for a platform

- **Data Retrieval**
  - `GetScrapedData`: Retrieve data collected by a scraper job

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Scraper Job Entity

```go
type ScraperJob struct {
    ID        string            `json:"id"`
    TenantID  string            `json:"tenant_id"`
    Platform  string            `json:"platform"`  // "instagram", "twitter", "facebook", etc.
    TargetID  string            `json:"target_id"` // Platform-specific ID
    JobType   JobType           `json:"job_type"`  // Profile, Posts, Engagement, etc.
    Status    JobStatus         `json:"status"`    // Pending, Scheduled, Running, etc.
    Schedule  ScraperSchedule   `json:"schedule"`
    LastError string            `json:"last_error"`
    RunCount  int               `json:"run_count"`
    LastRunAt time.Time         `json:"last_run_at"`
    NextRunAt time.Time         `json:"next_run_at"`
    Metadata  map[string]string `json:"metadata"`
    CreatedAt time.Time         `json:"created_at"`
    UpdatedAt time.Time         `json:"updated_at"`
}

type ScraperSchedule struct {
    CronExpression string            `json:"cron_expression"`
    Frequency      ScheduleFrequency `json:"frequency"` // Once, Hourly, Daily, Weekly
    StartDate      time.Time         `json:"start_date"`
    EndDate        time.Time         `json:"end_date"`
}
```

### Scraped Data Item

```go
type ScrapedDataItem struct {
    ID                string            `json:"id"`
    JobID             string            `json:"job_id"`
    TenantID          string            `json:"tenant_id"`
    Platform          string            `json:"platform"`
    TargetID          string            `json:"target_id"`
    PostID            string            `json:"post_id"`
    DataType          DataType          `json:"data_type"` // Profile, Post, Story, etc.
    PostedAt          time.Time         `json:"posted_at"`
    Likes             int               `json:"likes"`
    Shares            int               `json:"shares"`
    Comments          int               `json:"comments"`
    CTR               float64           `json:"ctr"`
    AvgWatchTime      float64           `json:"avg_watch_time"`
    EngagementRate    float64           `json:"engagement_rate"`
    ContentType       string            `json:"content_type"`
    ContentURL        string            `json:"content_url"`
    ContentAttributes map[string]string `json:"content_attributes"`
    ScrapedAt         time.Time         `json:"scraped_at"`
    CreatedAt         time.Time         `json:"created_at"`
}
```

### Platform Status

```go
type PlatformStatus struct {
    Platform      string           `json:"platform"`
    Available     bool             `json:"available"`
    StatusMessage string           `json:"status_message"`
    RateLimits    PlatformRateLimits `json:"rate_limits"`
    LastChecked   time.Time        `json:"last_checked"`
}

type PlatformRateLimits struct {
    RequestsPerMinute int       `json:"requests_per_minute"`
    RequestsPerHour   int       `json:"requests_per_hour"`
    RequestsPerDay    int       `json:"requests_per_day"`
    AvailableRequests int       `json:"available_requests"`
    ResetAt           time.Time `json:"reset_at"`
}
```

### Supported Job Types

The scraper service supports the following job types:

- `JOB_TYPE_PROFILE`: Scrape profile information
- `JOB_TYPE_POSTS`: Scrape posts
- `JOB_TYPE_ENGAGEMENT`: Scrape engagement metrics
- `JOB_TYPE_COMMENTS`: Scrape comments
- `JOB_TYPE_FOLLOWERS`: Scrape followers information

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9008` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run scraper/cmd/main.go
```

Or via Docker:

```bash
docker build -t scraper-service -f scraper/Dockerfile .
docker run -p 9008:9008 scraper-service
```

### Health Check

```bash
curl http://localhost:9008/health
```

Example response:
```json
{"status":"UP"}
```

### Using the gRPC Client

To use the Scraper service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/scraper/client"

func main() {
    // Create scraper client
    scraperClient, err := client.NewGRPCScraperClient("localhost:9008")
    if err != nil {
        log.Fatalf("Failed to create scraper client: %v", err)
    }
    defer scraperClient.Close()
    
    // Create a new scraper job
    job, err := scraperClient.CreateScraperJob(
        ctx,
        "tenant-123",
        "instagram",
        "competitor_username",
        repository.JobTypePosts,
        repository.ScraperSchedule{
            CronExpression: "0 */6 * * *", // Every 6 hours
            Frequency:      repository.FrequencyDaily,
        },
        map[string]string{
            "include_comments": "true",
            "post_limit": "50",
        },
    )
    if err != nil {
        log.Fatalf("Failed to create scraper job: %v", err)
    }
    
    log.Printf("Created scraper job with ID: %s, next run at: %v", 
        job.ID, job.NextRunAt)
}
```

## Supported Platforms

The Scraper service supports data collection from these platforms:

1. **Instagram**: Public posts, engagement metrics
2. **Twitter (X)**: Tweets, replies, retweets, likes
3. **Facebook**: Public page posts, engagement metrics
4. **LinkedIn**: Company posts, engagement metrics
5. **TikTok**: Videos, engagement metrics

Support for additional platforms can be added by implementing new provider integrations.

## Scraping Workflow

The Scraper service follows this workflow:

1. **Schedule**: Jobs are created with a frequency defined by a cron expression
2. **Execute**: The scheduler triggers jobs based on their next run time
3. **Collect**: Platform-specific APIs are used to collect data
4. **Process**: Data is normalized to a standard format
5. **Store**: Collected data is stored in the repository
6. **Update**: Job status and statistics are updated

## Rate Limiting Strategy

Platform-specific rate limits are defined for each supported social media:

| Platform  | Requests/Min | Requests/Hour | Requests/Day |
|-----------|--------------|---------------|--------------|
| Instagram | 30           | 500           | 5000         |
| Twitter   | 50           | 1500          | 10000        |
| Facebook  | 20           | 200           | 2000         |
| LinkedIn  | 10           | 100           | 1000         |
| TikTok    | 15           | 150           | 1500         |

Rate limits are tracked and enforced for each tenant and platform.

## Data Storage

The Scraper service uses Supabase (PostgreSQL) for data storage with these tables:

1. `scraper_jobs`: Stores job definitions and schedules
2. `scraped_data`: Stores the data collected by scraper jobs

Row Level Security (RLS) policies ensure that tenants can only access their own data.

## Error Handling

Common error scenarios:

- **Validation Error**: Returned when input parameters are invalid
- **Not Found Error**: Returned when a requested resource doesn't exist
- **Rate Limit Error**: Returned when API rate limits are hit
- **Authentication Error**: Returned when API credentials are invalid
- **Platform Error**: Returned when a social platform API returns an error

## Dependencies

- **Supabase**: For data storage
- **cron/v3**: For scheduling jobs
- **gRPC**: For service API
- **Platform APIs**: For data collection 