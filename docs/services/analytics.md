# Analytics Service

The Analytics Service provides AI-powered recommendations and predictions for optimizing social media content and strategy. It analyzes data from multiple services to generate actionable insights that help improve engagement and reach.

## Features

- **Posting Time Recommendations**: Identify optimal times to post for maximum engagement
- **Content Format Recommendations**: Suggest content formats that perform best
- **Engagement Prediction**: Forecast expected engagement for planned content
- **Trend Analysis**: Identify emerging trends in audience behavior
- **Competitor Benchmarking**: Compare performance against competitors
- **What-If Scenario Analysis**: Model potential outcomes of strategy changes
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
analytics/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Analytics service exposes a gRPC API defined in `analytics/pb/analytics.proto`.

Key operations:
- `GetRecommendedPostingTimes`: Get optimal posting times for maximum engagement
- `GetRecommendedContentFormats`: Get content format recommendations
- `PredictPostEngagement`: Predict engagement for a planned post
- `GetContentPerformanceAnalysis`: Analyze content performance over time
- `GetCompetitorBenchmark`: Compare performance with competitors
- `GetAudienceInsights`: Get advanced audience insights

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Posting Time Recommendation

```go
type PostingTimeRecommendation struct {
    DayOfWeek              string    `json:"day_of_week"`
    TimeOfDay              string    `json:"time_of_day"`
    PredictedEngagementRate float64   `json:"predicted_engagement_rate"`
    Confidence             float64   `json:"confidence"` // 0-1 scale
    SampleSize             int       `json:"sample_size"`
    Platform               string    `json:"platform"`
}
```

### Content Format Recommendation

```go
type ContentFormatRecommendation struct {
    Format                 string    `json:"format"`
    PredictedEngagementRate float64   `json:"predicted_engagement_rate"`
    TargetAudience         string    `json:"target_audience"` // Segment ID or name
    Confidence             float64   `json:"confidence"` // 0-1 scale
    SupportingEvidence     []string  `json:"supporting_evidence"` // Post IDs that performed well
    Platform               string    `json:"platform"`
}
```

### Engagement Prediction

```go
type EngagementPrediction struct {
    PredictedLikes         int       `json:"predicted_likes"`
    PredictedShares        int       `json:"predicted_shares"`
    PredictedComments      int       `json:"predicted_comments"`
    PredictedEngagementRate float64   `json:"predicted_engagement_rate"`
    PredictedReach         int       `json:"predicted_reach"`
    Confidence             float64   `json:"confidence"` // 0-1 scale
    FactorsInfluencing     []string  `json:"factors_influencing"` // Factors affecting prediction
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9007` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |
| `AUTH_SERVICE_URL` | URL for Auth service | `localhost:9001` |
| `COMPETITOR_SERVICE_URL` | URL for Competitor service | `localhost:9003` |
| `ENGAGEMENT_SERVICE_URL` | URL for Engagement service | `localhost:9004` |
| `CONTENT_SERVICE_URL` | URL for Content service | `localhost:9005` |
| `AUDIENCE_SERVICE_URL` | URL for Audience service | `localhost:9006` |
| `MODEL_REFRESH_INTERVAL` | How often to retrain models (in hours) | `24` |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run analytics/cmd/main.go
```

Or via Docker:

```bash
docker build -t analytics-service -f analytics/Dockerfile .
docker run -p 9007:9007 analytics-service
```

### Health Check

```bash
curl http://localhost:9007/health
```

Example response:
```json
{"status":"UP"}
```

### Accessing via gRPC Client

To use the Analytics service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/analytics/client"

func main() {
    // Create analytics client
    analyticsClient, err := client.NewAnalyticsClient("localhost:9007")
    if err != nil {
        log.Fatalf("Failed to create analytics client: %v", err)
    }
    
    // Get recommended posting times
    recommendations, err := analyticsClient.GetRecommendedPostingTimes(ctx, &pb.RecommendedPostingTimesRequest{
        TenantId:  tenantID,
        Platform:  "instagram",
        DayOfWeek: "all",
    })
    if err != nil {
        log.Fatalf("Failed to get recommendations: %v", err)
    }
    
    // Display recommendations
    for _, rec := range recommendations.Times {
        log.Printf("%s at %s: %.2f%% engagement (confidence: %.2f)",
            rec.DayOfWeek, rec.TimeOfDay, rec.PredictedEngagementRate*100, rec.Confidence)
    }
}
```

## Analysis Algorithms

The Analytics service uses several machine learning algorithms:

1. **Time Series Analysis**: For identifying patterns in posting times
2. **Regression Models**: For predicting engagement metrics
3. **Collaborative Filtering**: For content format recommendations
4. **Sentiment Analysis**: For audience response analysis
5. **Bayesian Networks**: For what-if scenario modeling

Models are retrained periodically based on the `MODEL_REFRESH_INTERVAL` setting to incorporate new data.

## Data Storage

The Analytics service uses Supabase (PostgreSQL) for data storage with these tables:

1. `recommendation_models`: Stores trained model metadata
2. `historical_recommendations`: Stores historical recommendations and their accuracy
3. `prediction_records`: Tracks predictions for comparison with actual outcomes
4. `analysis_jobs`: Manages background analysis tasks

Row Level Security (RLS) policies ensure that tenants can only access their own data.

## Error Handling

Common error scenarios:

- **Insufficient Data**: Returned when there's not enough data for reliable analysis
- **Model Training Error**: Returned when a model fails to train properly
- **Prediction Error**: Returned when a prediction cannot be made
- **Authorization Error**: Returned when a user doesn't have permission for an operation
- **Database Error**: Returned when there's an issue with the database operation

## Dependencies

- **Auth Service**: For token validation and tenant context
- **Competitor Service**: For competitor data analysis
- **Engagement Service**: For historical engagement data
- **Content Service**: For content performance correlation
- **Audience Service**: For audience insights
- **gRPC**: For service API
- **PostgreSQL**: Via Supabase for data storage 