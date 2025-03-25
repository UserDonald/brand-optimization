package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/engagement/pb"
	"github.com/donaldnash/go-competitor/engagement/repository"
	"github.com/donaldnash/go-competitor/engagement/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EngagementClient defines the interface for client communication with the engagement service
type EngagementClient interface {
	GetPersonalMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]repository.PersonalMetric, error)
	AddPersonalMetric(ctx context.Context, tenantID string, metric *repository.PersonalMetric) (*repository.PersonalMetric, error)
	UpdatePersonalMetric(ctx context.Context, tenantID string, metric *repository.PersonalMetric) (*repository.PersonalMetric, error)
	DeletePersonalMetric(ctx context.Context, tenantID, metricID string) error
	CompareMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) (*repository.ComparisonResult, error)
	GetEngagementTrends(ctx context.Context, tenantID, period string, startDate, endDate time.Time) ([]repository.EngagementTrend, error)
	GetEngagementInsights(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]service.EngagementInsight, error)
	Close() error
}

// GrpcEngagementClient implements EngagementClient using gRPC
type GrpcEngagementClient struct {
	conn   *grpc.ClientConn
	client pb.EngagementServiceClient
}

// NewEngagementClient creates a new gRPC client for the engagement service
func NewEngagementClient(serverAddr string) (*GrpcEngagementClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to engagement service: %w", err)
	}

	client := pb.NewEngagementServiceClient(conn)
	return &GrpcEngagementClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *GrpcEngagementClient) Close() error {
	return c.conn.Close()
}

// GetPersonalMetrics retrieves personal metrics for a date range
func (c *GrpcEngagementClient) GetPersonalMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]repository.PersonalMetric, error) {
	resp, err := c.client.GetPersonalMetrics(ctx, &pb.GetPersonalMetricsRequest{
		TenantId:  tenantID,
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get personal metrics: %w", err)
	}

	metrics := make([]repository.PersonalMetric, len(resp.Metrics))
	for i, pbMetric := range resp.Metrics {
		metrics[i] = repository.PersonalMetric{
			ID:             pbMetric.Id,
			TenantID:       pbMetric.TenantId,
			PostID:         pbMetric.PostId,
			Likes:          int(pbMetric.Likes),
			Shares:         int(pbMetric.Shares),
			Comments:       int(pbMetric.Comments),
			CTR:            pbMetric.ClickThroughRate,
			AvgWatchTime:   pbMetric.AvgWatchTime,
			EngagementRate: pbMetric.EngagementRate,
			PostedAt:       pbMetric.PostedAt.AsTime(),
			CreatedAt:      pbMetric.CreatedAt.AsTime(),
		}
	}

	return metrics, nil
}

// AddPersonalMetric adds a new personal metric
func (c *GrpcEngagementClient) AddPersonalMetric(ctx context.Context, tenantID string, metric *repository.PersonalMetric) (*repository.PersonalMetric, error) {
	resp, err := c.client.TrackPost(ctx, &pb.TrackPostRequest{
		TenantId:         tenantID,
		PostId:           metric.PostID,
		Platform:         "",
		ContentType:      "",
		ContentLength:    "",
		PostedAt:         timestamppb.New(metric.PostedAt),
		Metadata:         make(map[string]string),
		Likes:            int32(metric.Likes),
		Shares:           int32(metric.Shares),
		Comments:         int32(metric.Comments),
		ClickThroughRate: metric.CTR,
		AvgWatchTime:     metric.AvgWatchTime,
		EngagementRate:   metric.EngagementRate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add personal metric: %w", err)
	}

	return &repository.PersonalMetric{
		ID:             resp.Id,
		TenantID:       resp.TenantId,
		PostID:         resp.PostId,
		Likes:          int(resp.Likes),
		Shares:         int(resp.Shares),
		Comments:       int(resp.Comments),
		CTR:            resp.ClickThroughRate,
		AvgWatchTime:   resp.AvgWatchTime,
		EngagementRate: resp.EngagementRate,
		PostedAt:       resp.PostedAt.AsTime(),
		CreatedAt:      resp.CreatedAt.AsTime(),
	}, nil
}

// UpdatePersonalMetric updates an existing personal metric
func (c *GrpcEngagementClient) UpdatePersonalMetric(ctx context.Context, tenantID string, metric *repository.PersonalMetric) (*repository.PersonalMetric, error) {
	resp, err := c.client.UpdatePostMetrics(ctx, &pb.UpdatePostMetricsRequest{
		TenantId:         tenantID,
		PostId:           metric.PostID,
		Likes:            int32(metric.Likes),
		Shares:           int32(metric.Shares),
		Comments:         int32(metric.Comments),
		ClickThroughRate: metric.CTR,
		AvgWatchTime:     metric.AvgWatchTime,
		EngagementRate:   metric.EngagementRate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update personal metric: %w", err)
	}

	return &repository.PersonalMetric{
		ID:             resp.Id,
		TenantID:       resp.TenantId,
		PostID:         resp.PostId,
		Likes:          int(resp.Likes),
		Shares:         int(resp.Shares),
		Comments:       int(resp.Comments),
		CTR:            resp.ClickThroughRate,
		AvgWatchTime:   resp.AvgWatchTime,
		EngagementRate: resp.EngagementRate,
		PostedAt:       resp.PostedAt.AsTime(),
		CreatedAt:      resp.CreatedAt.AsTime(),
	}, nil
}

// DeletePersonalMetric deletes a personal metric
func (c *GrpcEngagementClient) DeletePersonalMetric(ctx context.Context, tenantID, metricID string) error {
	// For simplicity, we'll assume that metricID is actually a postID since our
	// DeletePostMetrics API operates on post IDs, not metric IDs
	_, err := c.client.DeletePostMetrics(ctx, &pb.DeletePostMetricsRequest{
		TenantId: tenantID,
		PostId:   metricID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete personal metric: %w", err)
	}

	return nil
}

// CompareMetrics compares personal metrics with competitor metrics
func (c *GrpcEngagementClient) CompareMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) (*repository.ComparisonResult, error) {
	// In this implementation, we'll use the service directly since we don't have
	// a direct gRPC endpoint for this functionality in the protobuf definition

	// Create a temporary service and repository
	repo, err := repository.NewSupabaseEngagementRepository(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	svc := service.NewEngagementService(repo)
	return svc.CompareMetrics(ctx, tenantID, competitorID, startDate, endDate)
}

// GetEngagementTrends retrieves engagement trends over time
func (c *GrpcEngagementClient) GetEngagementTrends(ctx context.Context, tenantID, period string, startDate, endDate time.Time) ([]repository.EngagementTrend, error) {
	resp, err := c.client.GetEngagementTrends(ctx, &pb.GetEngagementTrendsRequest{
		TenantId:   tenantID,
		StartDate:  timestamppb.New(startDate),
		EndDate:    timestamppb.New(endDate),
		MetricType: "engagement_rate",
		Interval:   period,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement trends: %w", err)
	}

	// Convert from proto response to repository model
	trends := make([]repository.EngagementTrend, len(resp.Points))
	for i, point := range resp.Points {
		trends[i] = repository.EngagementTrend{
			Date:            point.Date.AsTime(),
			EngagementRate:  point.Value,
			Likes:           int(point.PostCount), // Using PostCount as a proxy for likes
			Shares:          0,                    // Not provided in the response
			Comments:        0,                    // Not provided in the response
			ComparisonValue: 0,                    // Not provided in the response
		}
	}

	return trends, nil
}

// GetEngagementInsights generates insights from engagement data
func (c *GrpcEngagementClient) GetEngagementInsights(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]service.EngagementInsight, error) {
	// In this implementation, we'll use the service directly since we don't have
	// a direct gRPC endpoint for this functionality in the protobuf definition

	// Create a temporary service and repository
	repo, err := repository.NewSupabaseEngagementRepository(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	svc := service.NewEngagementService(repo)
	return svc.GetEngagementInsights(ctx, tenantID, startDate, endDate)
}
