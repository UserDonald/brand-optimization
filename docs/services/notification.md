# Notification Service

The Notification Service manages notifications, alerts, and scheduled reports for the Strategic Brand Optimization Platform. It ensures that users receive timely updates about important events, scheduled posts, and performance insights.

## Features

- **Alert Management**: Define and trigger alerts based on various conditions
- **Scheduled Reports**: Generate and send recurring performance reports
- **Event Notifications**: Notify users about system events
- **Notification Preferences**: Allow users to customize notification settings
- **Multi-channel Delivery**: Send notifications via various channels (email, in-app, etc.)
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
notification/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Notification service exposes a gRPC API defined in `notification/pb/notification.proto`.

Key operations:
- `SendNotification`: Send a notification to a user
- `CreateAlertRule`: Create a new alert rule
- `GetAlertRule`: Retrieve a specific alert rule
- `ListAlertRules`: List all alert rules for a tenant
- `UpdateAlertRule`: Update an alert rule
- `DeleteAlertRule`: Remove an alert rule
- `CreateScheduledReport`: Create a new scheduled report
- `GetScheduledReport`: Retrieve a specific scheduled report
- `UpdateScheduledReport`: Update a scheduled report
- `DeleteScheduledReport`: Remove a scheduled report
- `GetUserPreferences`: Get notification preferences for a user
- `UpdateUserPreferences`: Update notification preferences

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Notification Entity

```go
type Notification struct {
    ID          string    `json:"id"`
    TenantID    string    `json:"tenant_id"`
    UserID      string    `json:"user_id"`
    Title       string    `json:"title"`
    Message     string    `json:"message"`
    Type        string    `json:"type"` // "alert", "report", "system", "scheduled_post"
    Severity    string    `json:"severity"` // "info", "warning", "critical"
    ResourceID  string    `json:"resource_id"` // Related resource (e.g., post ID)
    ResourceType string   `json:"resource_type"` // Type of related resource
    Channels    []string  `json:"channels"` // "email", "in_app", "sms"
    Status      string    `json:"status"` // "pending", "sent", "failed"
    SentAt      *time.Time `json:"sent_at"`
    ReadAt      *time.Time `json:"read_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Alert Rule Entity

```go
type AlertRule struct {
    ID          string    `json:"id"`
    TenantID    string    `json:"tenant_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Type        string    `json:"type"` // "engagement_drop", "competitor_activity", etc.
    Conditions  map[string]interface{} `json:"conditions"`
    Severity    string    `json:"severity"` // "info", "warning", "critical"
    Channels    []string  `json:"channels"` // "email", "in_app", "sms"
    UserIDs     []string  `json:"user_ids"` // Users to notify
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Scheduled Report Entity

```go
type ScheduledReport struct {
    ID          string    `json:"id"`
    TenantID    string    `json:"tenant_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    ReportType  string    `json:"report_type"` // "engagement", "competitor", "audience"
    Parameters  map[string]interface{} `json:"parameters"`
    Schedule    string    `json:"schedule"` // CRON expression
    Format      string    `json:"format"` // "pdf", "csv", "json"
    Channels    []string  `json:"channels"` // "email", "in_app"
    UserIDs     []string  `json:"user_ids"` // Users to receive the report
    IsActive    bool      `json:"is_active"`
    LastRunAt   *time.Time `json:"last_run_at"`
    NextRunAt   *time.Time `json:"next_run_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9002` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |
| `AUTH_SERVICE_URL` | URL for Auth service | `localhost:9001` |
| `SMTP_HOST` | SMTP server for email notifications | - |
| `SMTP_PORT` | SMTP port | `587` |
| `SMTP_USERNAME` | SMTP username | - |
| `SMTP_PASSWORD` | SMTP password | - |
| `SMTP_FROM` | Email sender address | - |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run notification/cmd/main.go
```

Or via Docker:

```bash
docker build -t notification-service -f notification/Dockerfile .
docker run -p 9002:9002 notification-service
```

### Health Check

```bash
curl http://localhost:9002/health
```

Example response:
```json
{"status":"UP"}
```

### Accessing via gRPC Client

To use the Notification service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/notification/client"

func main() {
    // Create notification client
    notificationClient, err := client.NewNotificationClient("localhost:9002")
    if err != nil {
        log.Fatalf("Failed to create notification client: %v", err)
    }
    
    // Send a notification
    _, err = notificationClient.SendNotification(ctx, &pb.SendNotificationRequest{
        TenantId:     tenantID,
        UserId:       userID,
        Title:        "Engagement Spike Detected",
        Message:      "Your recent post has 200% higher engagement than usual.",
        Type:         "alert",
        Severity:     "info",
        ResourceId:   postID,
        ResourceType: "post",
        Channels:     []string{"email", "in_app"},
    })
    if err != nil {
        log.Fatalf("Failed to send notification: %v", err)
    }
    
    log.Println("Notification sent successfully")
}
```

## Notification Channels

The service supports these notification channels:

1. **In-App**: Real-time notifications within the application
2. **Email**: Email notifications using SMTP
3. **SMS**: Text message notifications (requires additional configuration)
4. **Webhook**: HTTP webhook callbacks for integration with external systems

New channels can be added by implementing the NotificationChannel interface.

## Scheduling System

The Notification service uses a background worker to manage scheduled reports:

1. Reports are defined with a CRON schedule expression
2. A background worker checks for reports due to be generated
3. When a report is due, the worker generates the report using data from other services
4. The report is delivered through the specified channels
5. The next run time is calculated and updated

## Data Storage

The Notification service uses Supabase (PostgreSQL) for data storage with these tables:

1. `notifications`: Stores all notifications
2. `alert_rules`: Stores alert definitions
3. `scheduled_reports`: Stores report definitions
4. `notification_preferences`: Stores user preferences
5. `notification_logs`: Stores delivery logs

Row Level Security (RLS) policies ensure that tenants can only access their own data.

## Error Handling

Common error scenarios:

- **Not Found**: Returned when a requested entity does not exist
- **Validation Error**: Returned when input data does not meet validation requirements
- **Delivery Error**: Returned when a notification can't be delivered via a channel
- **Authorization Error**: Returned when a user doesn't have permission for an operation
- **Database Error**: Returned when there's an issue with the database operation

## Dependencies

- **Auth Service**: For token validation and user information
- **gRPC**: For service API
- **PostgreSQL**: Via Supabase for data storage
- **SMTP Server**: For email notifications 