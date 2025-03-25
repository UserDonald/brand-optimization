# Go Client Integration Guide

This guide provides a detailed walkthrough for integrating with the Strategic Brand Optimization Platform using Go. It covers authentication, data retrieval, and common use cases for Go applications.

## Table of Contents

1. [Setup](#setup)
2. [Authentication](#authentication)
3. [Working with Services](#working-with-services)
4. [GraphQL Integration](#graphql-integration)
5. [Common Use Cases](#common-use-cases)

## Setup

### Dependencies

First, add the required dependencies to your Go project:

```go
import (
    "context"
    "fmt"
    "time"
    
    "github.com/machinebox/graphql"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)
```

Install the required packages:

```bash
go get github.com/machinebox/graphql
go get google.golang.org/grpc
```

### Project Structure

We recommend organizing your Go client integration with the following structure:

```
myapp/
├── auth/            # Authentication client and utilities
├── clients/         # Service clients
│   ├── competitor/  # Competitor service client
│   ├── analytics/   # Analytics service client
│   └── ...          # Other service clients
├── graphql/         # GraphQL client integration
└── main.go          # Application entry point
```

## Authentication

### Creating the Auth Client

```go
// auth/client.go
package auth

import (
    "context"
    "fmt"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// Client represents the authentication client
type Client struct {
    conn          *grpc.ClientConn
    serviceClient AuthServiceClient
    accessToken   string
    refreshToken  string
    expiresAt     time.Time
}

// NewClient creates a new auth client
func NewClient(serverAddr string) (*Client, error) {
    conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to auth service: %w", err)
    }
    
    client := NewAuthServiceClient(conn)
    
    return &Client{
        conn:          conn,
        serviceClient: client,
    }, nil
}

// Close closes the connection
func (c *Client) Close() error {
    if c.conn != nil {
        return c.conn.Close()
    }
    return nil
}

// Login authenticates a user and stores the tokens
func (c *Client) Login(ctx context.Context, email, password string) (*User, error) {
    resp, err := c.serviceClient.Login(ctx, &LoginRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        return nil, fmt.Errorf("login failed: %w", err)
    }
    
    c.accessToken = resp.AccessToken
    c.refreshToken = resp.RefreshToken
    c.expiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
    
    user := &User{
        ID:        resp.User.Id,
        Email:     resp.User.Email,
        FirstName: resp.User.FirstName,
        LastName:  resp.User.LastName,
        TenantID:  resp.User.TenantId,
        Role:      resp.User.Role,
    }
    
    return user, nil
}

// GetAccessToken returns the current access token, refreshing if needed
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
    // If token is still valid, return it
    if time.Now().Before(c.expiresAt) {
        return c.accessToken, nil
    }
    
    // Otherwise, refresh the token
    resp, err := c.serviceClient.RefreshToken(ctx, &RefreshTokenRequest{
        RefreshToken: c.refreshToken,
    })
    if err != nil {
        return "", fmt.Errorf("token refresh failed: %w", err)
    }
    
    c.accessToken = resp.AccessToken
    c.refreshToken = resp.RefreshToken
    c.expiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
    
    return c.accessToken, nil
}
```

### Integration with Auth Module

```go
// main.go
package main

import (
    "context"
    "fmt"
    "log"
    
    "myapp/auth"
)

func main() {
    // Create auth client
    authClient, err := auth.NewClient("localhost:9001")
    if err != nil {
        log.Fatalf("Failed to create auth client: %v", err)
    }
    defer authClient.Close()
    
    // Login
    ctx := context.Background()
    user, err := authClient.Login(ctx, "user@example.com", "password123")
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    // Get token for other service calls
    token, err := authClient.GetAccessToken(ctx)
    if err != nil {
        log.Fatalf("Failed to get access token: %v", err)
    }
    
    // Use the token for other service calls...
}
```

## Working with Services

### Creating a Service Client (Competitor Example)

```go
// clients/competitor/client.go
package competitor

import (
    "context"
    "fmt"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/metadata"
)

// Client represents the competitor service client
type Client struct {
    conn          *grpc.ClientConn
    serviceClient CompetitorServiceClient
    getToken      func(context.Context) (string, error)
}

// NewClient creates a new competitor client
func NewClient(serverAddr string, tokenProvider func(context.Context) (string, error)) (*Client, error) {
    conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to competitor service: %w", err)
    }
    
    client := NewCompetitorServiceClient(conn)
    
    return &Client{
        conn:          conn,
        serviceClient: client,
        getToken:      tokenProvider,
    }, nil
}

// Close closes the connection
func (c *Client) Close() error {
    if c.conn != nil {
        return c.conn.Close()
    }
    return nil
}

// authenticatedContext adds the auth token to the context
func (c *Client) authenticatedContext(ctx context.Context) (context.Context, error) {
    token, err := c.getToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get auth token: %w", err)
    }
    
    return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), nil
}

// GetCompetitors retrieves all competitors for the current tenant
func (c *Client) GetCompetitors(ctx context.Context, tenantID string) ([]Competitor, error) {
    ctx, err := c.authenticatedContext(ctx)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.serviceClient.ListCompetitors(ctx, &ListCompetitorsRequest{
        TenantId: tenantID,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get competitors: %w", err)
    }
    
    competitors := make([]Competitor, len(resp.Competitors))
    for i, pbCompetitor := range resp.Competitors {
        competitors[i] = Competitor{
            ID:       pbCompetitor.Id,
            Name:     pbCompetitor.Name,
            Platform: pbCompetitor.Platform,
        }
    }
    
    return competitors, nil
}

// AddCompetitor adds a new competitor
func (c *Client) AddCompetitor(ctx context.Context, tenantID, name, platform string) (*Competitor, error) {
    ctx, err := c.authenticatedContext(ctx)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.serviceClient.AddCompetitor(ctx, &AddCompetitorRequest{
        TenantId: tenantID,
        Name:     name,
        Platform: platform,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to add competitor: %w", err)
    }
    
    return &Competitor{
        ID:       resp.Id,
        Name:     resp.Name,
        Platform: resp.Platform,
    }, nil
}
```

### Using Service Clients Together

```go
// main.go
package main

import (
    "context"
    "fmt"
    "log"
    
    "myapp/auth"
    "myapp/clients/competitor"
)

func main() {
    // Create auth client
    authClient, err := auth.NewClient("localhost:9001")
    if err != nil {
        log.Fatalf("Failed to create auth client: %v", err)
    }
    defer authClient.Close()
    
    // Login
    ctx := context.Background()
    user, err := authClient.Login(ctx, "user@example.com", "password123")
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    // Create competitor client with token provider function
    competitorClient, err := competitor.NewClient("localhost:9003", authClient.GetAccessToken)
    if err != nil {
        log.Fatalf("Failed to create competitor client: %v", err)
    }
    defer competitorClient.Close()
    
    // Get competitors
    competitors, err := competitorClient.GetCompetitors(ctx, user.TenantID)
    if err != nil {
        log.Fatalf("Failed to get competitors: %v", err)
    }
    
    fmt.Println("Competitors:")
    for _, comp := range competitors {
        fmt.Printf("- %s (%s)\n", comp.Name, comp.Platform)
    }
    
    // Add a new competitor
    newCompetitor, err := competitorClient.AddCompetitor(ctx, user.TenantID, "New Competitor", "instagram")
    if err != nil {
        log.Fatalf("Failed to add competitor: %v", err)
    }
    
    fmt.Printf("Added competitor: %s (ID: %s)\n", newCompetitor.Name, newCompetitor.ID)
}
```

## GraphQL Integration

For more complex operations, you might want to use the GraphQL API directly. Here's how to set up a GraphQL client in Go:

```go
// graphql/client.go
package graphql

import (
    "context"
    "fmt"
    
    "github.com/machinebox/graphql"
)

// Client represents a GraphQL client
type Client struct {
    client    *graphql.Client
    getToken  func(context.Context) (string, error)
    serverURL string
}

// NewClient creates a new GraphQL client
func NewClient(serverURL string, tokenProvider func(context.Context) (string, error)) *Client {
    client := graphql.NewClient(serverURL)
    
    return &Client{
        client:    client,
        getToken:  tokenProvider,
        serverURL: serverURL,
    }
}

// Execute executes a GraphQL query or mutation
func (c *Client) Execute(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
    // Create request
    req := graphql.NewRequest(query)
    
    // Add variables
    for key, value := range variables {
        req.Var(key, value)
    }
    
    // Add auth header
    token, err := c.getToken(ctx)
    if err != nil {
        return fmt.Errorf("failed to get auth token: %w", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)
    
    // Run the query
    if err := c.client.Run(ctx, req, response); err != nil {
        return fmt.Errorf("graphql query failed: %w", err)
    }
    
    return nil
}
```

### Example GraphQL Query

```go
// main.go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "myapp/auth"
    "myapp/graphql"
)

// Response types
type CompareMetricsResponse struct {
    CompareMetrics struct {
        Competitor struct {
            Metrics []struct {
                Likes          int     `json:"likes"`
                Shares         int     `json:"shares"`
                Comments       int     `json:"comments"`
                EngagementRate float64 `json:"engagementRate"`
            } `json:"metrics"`
            Aggregates struct {
                TotalLikes       int     `json:"totalLikes"`
                AvgEngagementRate float64 `json:"avgEngagementRate"`
            } `json:"aggregates"`
        } `json:"competitor"`
        Personal struct {
            Metrics []struct {
                Likes          int     `json:"likes"`
                Shares         int     `json:"shares"`
                Comments       int     `json:"comments"`
                EngagementRate float64 `json:"engagementRate"`
            } `json:"metrics"`
            Aggregates struct {
                TotalLikes       int     `json:"totalLikes"`
                AvgEngagementRate float64 `json:"avgEngagementRate"`
            } `json:"aggregates"`
        } `json:"personal"`
        Ratios struct {
            LikesRatio          float64 `json:"likesRatio"`
            SharesRatio         float64 `json:"sharesRatio"`
            EngagementRateRatio float64 `json:"engagementRateRatio"`
        } `json:"ratios"`
    } `json:"compareMetrics"`
}

func main() {
    // Create auth client
    authClient, err := auth.NewClient("localhost:9001")
    if err != nil {
        log.Fatalf("Failed to create auth client: %v", err)
    }
    defer authClient.Close()
    
    // Login
    ctx := context.Background()
    user, err := authClient.Login(ctx, "user@example.com", "password123")
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    // Create GraphQL client
    gqlClient := graphql.NewClient("http://localhost:8080/query", authClient.GetAccessToken)
    
    // Define the query
    query := `
        query CompareMetrics($competitorId: ID!, $dateRange: DateRangeInput!) {
            compareMetrics(competitorID: $competitorId, dateRange: $dateRange) {
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
    `
    
    // Set up variables
    variables := map[string]interface{}{
        "competitorId": "competitor-123",
        "dateRange": map[string]interface{}{
            "startDate": time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
            "endDate":   time.Now().Format("2006-01-02"),
        },
    }
    
    // Execute the query
    var response CompareMetricsResponse
    if err := gqlClient.Execute(ctx, query, variables, &response); err != nil {
        log.Fatalf("Failed to execute GraphQL query: %v", err)
    }
    
    // Process the results
    result := response.CompareMetrics
    fmt.Printf("Comparison Results:\n")
    fmt.Printf("Your Total Likes: %d\n", result.Personal.Aggregates.TotalLikes)
    fmt.Printf("Competitor Total Likes: %d\n", result.Competitor.Aggregates.TotalLikes)
    fmt.Printf("Likes Ratio: %.2f%%\n", result.Ratios.LikesRatio*100)
    fmt.Printf("Engagement Rate Ratio: %.2f%%\n", result.Ratios.EngagementRateRatio*100)
}
```

## Common Use Cases

### 1. Fetching Engagement Trends

```go
// clients/engagement/client.go
package engagement

import (
    "context"
    "fmt"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "google.golang.org/protobuf/types/known/timestamppb"
)

// Client for the engagement service
type Client struct {
    conn          *grpc.ClientConn
    serviceClient EngagementServiceClient
    getToken      func(context.Context) (string, error)
}

// NewClient creates a new engagement client
func NewClient(serverAddr string, tokenProvider func(context.Context) (string, error)) (*Client, error) {
    conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to engagement service: %w", err)
    }
    
    client := NewEngagementServiceClient(conn)
    
    return &Client{
        conn:          conn,
        serviceClient: client,
        getToken:      tokenProvider,
    }, nil
}

// authenticatedContext adds the auth token to the context
func (c *Client) authenticatedContext(ctx context.Context) (context.Context, error) {
    token, err := c.getToken(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get auth token: %w", err)
    }
    
    return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), nil
}

// GetEngagementTrends gets engagement trends for a specified period
func (c *Client) GetEngagementTrends(ctx context.Context, tenantID, period string, startDate, endDate time.Time) ([]EngagementTrend, error) {
    ctx, err := c.authenticatedContext(ctx)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.serviceClient.GetEngagementTrends(ctx, &GetEngagementTrendsRequest{
        TenantId:  tenantID,
        Period:    period,
        StartDate: timestamppb.New(startDate),
        EndDate:   timestamppb.New(endDate),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get engagement trends: %w", err)
    }
    
    trends := make([]EngagementTrend, len(resp.Trends))
    for i, trend := range resp.Trends {
        trends[i] = EngagementTrend{
            Date:           trend.Date.AsTime(),
            Likes:          trend.Likes,
            Shares:         trend.Shares,
            Comments:       trend.Comments,
            EngagementRate: trend.EngagementRate,
        }
    }
    
    return trends, nil
}
```

### 2. Scheduling Content Posts

```go
// Example of scheduling a content post
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "myapp/auth"
    "myapp/clients/content"
)

func main() {
    // Initialize auth client
    authClient, err := auth.NewClient("localhost:9001")
    if err != nil {
        log.Fatalf("Failed to create auth client: %v", err)
    }
    defer authClient.Close()
    
    // Login
    ctx := context.Background()
    user, err := authClient.Login(ctx, "user@example.com", "password123")
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    // Create content client
    contentClient, err := content.NewClient("localhost:9005", authClient.GetAccessToken)
    if err != nil {
        log.Fatalf("Failed to create content client: %v", err)
    }
    defer contentClient.Close()
    
    // Schedule a post for tomorrow
    scheduledTime := time.Now().AddDate(0, 0, 1).Round(time.Hour)
    post, err := contentClient.SchedulePost(
        ctx,
        user.TenantID,
        "Check out our new product launch! #innovation",
        "instagram",
        "image",
        scheduledTime,
    )
    if err != nil {
        log.Fatalf("Failed to schedule post: %v", err)
    }
    
    fmt.Printf("Post scheduled for %s with ID: %s\n", post.ScheduledTime.Format(time.RFC3339), post.ID)
    
    // List all scheduled posts
    posts, err := contentClient.GetScheduledPosts(ctx, user.TenantID)
    if err != nil {
        log.Fatalf("Failed to get scheduled posts: %v", err)
    }
    
    fmt.Println("Scheduled posts:")
    for _, p := range posts {
        fmt.Printf("- %s: %s (%s)\n", p.ScheduledTime.Format("2006-01-02 15:04"), p.Content, p.Status)
    }
}
```

### 3. Getting AI-Powered Recommendations

```go
// Example of getting recommendations from the analytics service
package main

import (
    "context"
    "fmt"
    "log"
    
    "myapp/auth"
    "myapp/clients/analytics"
)

func main() {
    // Initialize auth client
    authClient, err := auth.NewClient("localhost:9001")
    if err != nil {
        log.Fatalf("Failed to create auth client: %v", err)
    }
    defer authClient.Close()
    
    // Login
    ctx := context.Background()
    user, err := authClient.Login(ctx, "user@example.com", "password123")
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    
    // Create analytics client
    analyticsClient, err := analytics.NewClient("localhost:9007", authClient.GetAccessToken)
    if err != nil {
        log.Fatalf("Failed to create analytics client: %v", err)
    }
    defer analyticsClient.Close()
    
    // Get posting time recommendations
    timeResp, err := analyticsClient.GetPostingTimeRecommendations(ctx, user.TenantID, "monday")
    if err != nil {
        log.Fatalf("Failed to get posting time recommendations: %v", err)
    }
    
    fmt.Println("Recommended posting times (Monday):")
    for _, rec := range timeResp.Recommendations {
        fmt.Printf("- %s: %.2f%% engagement rate\n", rec.Time, rec.ExpectedEngagementRate*100)
    }
    
    // Get content format recommendations
    formatResp, err := analyticsClient.GetContentFormatRecommendations(ctx, user.TenantID)
    if err != nil {
        log.Fatalf("Failed to get content format recommendations: %v", err)
    }
    
    fmt.Println("\nRecommended content formats:")
    for _, rec := range formatResp.Recommendations {
        fmt.Printf("- %s: %.2f%% engagement rate\n", rec.Format, rec.ExpectedEngagementRate*100)
    }
}
```

For more examples and detailed API documentation, refer to the [API Documentation](../api.md) and the available client implementations in the source code. 