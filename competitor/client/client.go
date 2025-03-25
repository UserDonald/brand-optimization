package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/competitor/pb"
	"github.com/donaldnash/go-competitor/competitor/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CompetitorClient defines the interface for client communication with the competitor service
type CompetitorClient interface {
	GetCompetitors(ctx context.Context, tenantID string) ([]repository.Competitor, error)
	GetCompetitor(ctx context.Context, tenantID, competitorID string) (*repository.Competitor, error)
	GetCompetitorMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) ([]repository.CompetitorMetric, error)
	AddCompetitor(ctx context.Context, tenantID, name, platform string) (*repository.Competitor, error)
	UpdateCompetitor(ctx context.Context, tenantID, competitorID, name, platform string) (*repository.Competitor, error)
	DeleteCompetitor(ctx context.Context, tenantID, competitorID string) error
	UpdateCompetitorMetrics(ctx context.Context, tenantID, competitorID string, metrics []repository.CompetitorMetric) (int, error)
	CompareMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) (*ComparisonResult, error)
	Close() error
}

// ComparisonResult represents the result of comparing metrics
type ComparisonResult struct {
	CompetitorMetrics []repository.CompetitorMetric
	PersonalMetrics   []repository.CompetitorMetric
	Ratios            MetricRatios
}

// MetricRatios contains the calculated ratios between personal and competitor metrics
type MetricRatios struct {
	LikesRatio          float64
	SharesRatio         float64
	CommentsRatio       float64
	EngagementRateRatio float64
	WatchTimeRatio      float64
}

// GRPCCompetitorClient implements CompetitorClient using gRPC
type GRPCCompetitorClient struct {
	conn   *grpc.ClientConn
	client pb.CompetitorServiceClient
}

// NewGRPCCompetitorClient creates a new GRPCCompetitorClient
func NewGRPCCompetitorClient(serverAddr string) (*GRPCCompetitorClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to competitor service: %w", err)
	}

	client := pb.NewCompetitorServiceClient(conn)
	return &GRPCCompetitorClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *GRPCCompetitorClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetCompetitors retrieves all competitors for the current tenant
func (c *GRPCCompetitorClient) GetCompetitors(ctx context.Context, tenantID string) ([]repository.Competitor, error) {
	resp, err := c.client.ListCompetitors(ctx, &pb.ListCompetitorsRequest{
		TenantId: tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get competitors: %w", err)
	}

	competitors := make([]repository.Competitor, len(resp.Competitors))
	for i, pbCompetitor := range resp.Competitors {
		competitors[i] = repository.Competitor{
			ID:        pbCompetitor.Id,
			TenantID:  pbCompetitor.TenantId,
			Name:      pbCompetitor.Name,
			Platform:  pbCompetitor.Platform,
			CreatedAt: pbCompetitor.CreatedAt.AsTime(),
		}
	}

	return competitors, nil
}

// GetCompetitor retrieves a specific competitor
func (c *GRPCCompetitorClient) GetCompetitor(ctx context.Context, tenantID, competitorID string) (*repository.Competitor, error) {
	resp, err := c.client.GetCompetitor(ctx, &pb.GetCompetitorRequest{
		TenantId:     tenantID,
		CompetitorId: competitorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor: %w", err)
	}

	return &repository.Competitor{
		ID:        resp.Id,
		TenantID:  resp.TenantId,
		Name:      resp.Name,
		Platform:  resp.Platform,
		CreatedAt: resp.CreatedAt.AsTime(),
	}, nil
}

// GetCompetitorMetrics retrieves metrics for a specific competitor
func (c *GRPCCompetitorClient) GetCompetitorMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) ([]repository.CompetitorMetric, error) {
	resp, err := c.client.GetCompetitorMetrics(ctx, &pb.GetCompetitorMetricsRequest{
		TenantId:     tenantID,
		CompetitorId: competitorID,
		StartDate:    timestamppb.New(startDate),
		EndDate:      timestamppb.New(endDate),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor metrics: %w", err)
	}

	metrics := make([]repository.CompetitorMetric, len(resp.Metrics))
	for i, pbMetric := range resp.Metrics {
		metrics[i] = repository.CompetitorMetric{
			ID:             pbMetric.Id,
			CompetitorID:   pbMetric.CompetitorId,
			PostID:         pbMetric.PostId,
			Likes:          int(pbMetric.Likes),
			Shares:         int(pbMetric.Shares),
			Comments:       int(pbMetric.Comments),
			CTR:            pbMetric.ClickThroughRate,
			AvgWatchTime:   pbMetric.AvgWatchTime,
			EngagementRate: pbMetric.EngagementRate,
			PostedAt:       pbMetric.PostedAt.AsTime(),
		}
	}

	return metrics, nil
}

// AddCompetitor adds a new competitor
func (c *GRPCCompetitorClient) AddCompetitor(ctx context.Context, tenantID, name, platform string) (*repository.Competitor, error) {
	resp, err := c.client.AddCompetitor(ctx, &pb.AddCompetitorRequest{
		TenantId: tenantID,
		Name:     name,
		Platform: platform,
		Metadata: make(map[string]string),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add competitor: %w", err)
	}

	return &repository.Competitor{
		ID:        resp.Id,
		TenantID:  resp.TenantId,
		Name:      resp.Name,
		Platform:  resp.Platform,
		CreatedAt: resp.CreatedAt.AsTime(),
	}, nil
}

// UpdateCompetitor updates an existing competitor
func (c *GRPCCompetitorClient) UpdateCompetitor(ctx context.Context, tenantID, competitorID, name, platform string) (*repository.Competitor, error) {
	resp, err := c.client.UpdateCompetitor(ctx, &pb.UpdateCompetitorRequest{
		TenantId:     tenantID,
		CompetitorId: competitorID,
		Name:         name,
		Platform:     platform,
		Metadata:     make(map[string]string),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update competitor: %w", err)
	}

	return &repository.Competitor{
		ID:        resp.Id,
		TenantID:  resp.TenantId,
		Name:      resp.Name,
		Platform:  resp.Platform,
		CreatedAt: resp.CreatedAt.AsTime(),
	}, nil
}

// DeleteCompetitor deletes a competitor
func (c *GRPCCompetitorClient) DeleteCompetitor(ctx context.Context, tenantID, competitorID string) error {
	_, err := c.client.DeleteCompetitor(ctx, &pb.DeleteCompetitorRequest{
		TenantId:     tenantID,
		CompetitorId: competitorID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete competitor: %w", err)
	}

	return nil
}

// UpdateCompetitorMetrics updates metrics for a specific competitor
func (c *GRPCCompetitorClient) UpdateCompetitorMetrics(ctx context.Context, tenantID, competitorID string, metrics []repository.CompetitorMetric) (int, error) {
	count := 0
	for _, metric := range metrics {
		_, err := c.client.TrackCompetitorPost(ctx, &pb.TrackCompetitorPostRequest{
			TenantId:         tenantID,
			CompetitorId:     competitorID,
			PostId:           metric.PostID,
			Likes:            int32(metric.Likes),
			Shares:           int32(metric.Shares),
			Comments:         int32(metric.Comments),
			ClickThroughRate: metric.CTR,
			AvgWatchTime:     metric.AvgWatchTime,
			EngagementRate:   metric.EngagementRate,
			PostedAt:         timestamppb.New(metric.PostedAt),
		})
		if err != nil {
			return count, fmt.Errorf("failed to update metric: %w", err)
		}
		count++
	}

	return count, nil
}

// CompareMetrics compares metrics between a competitor and personal brand
func (c *GRPCCompetitorClient) CompareMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) (*ComparisonResult, error) {
	resp, err := c.client.CompareMetrics(ctx, &pb.CompareMetricsRequest{
		TenantId:     tenantID,
		CompetitorId: competitorID,
		StartDate:    timestamppb.New(startDate),
		EndDate:      timestamppb.New(endDate),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to compare metrics: %w", err)
	}

	// Convert competitor metrics
	compMetrics := make([]repository.CompetitorMetric, len(resp.Competitor.Metrics))
	for i, pbMetric := range resp.Competitor.Metrics {
		compMetrics[i] = repository.CompetitorMetric{
			ID:             pbMetric.Id,
			CompetitorID:   pbMetric.CompetitorId,
			PostID:         pbMetric.PostId,
			Likes:          int(pbMetric.Likes),
			Shares:         int(pbMetric.Shares),
			Comments:       int(pbMetric.Comments),
			CTR:            pbMetric.ClickThroughRate,
			AvgWatchTime:   pbMetric.AvgWatchTime,
			EngagementRate: pbMetric.EngagementRate,
			PostedAt:       pbMetric.PostedAt.AsTime(),
		}
	}

	// Convert personal metrics
	persMetrics := make([]repository.CompetitorMetric, len(resp.Personal.Metrics))
	for i, pbMetric := range resp.Personal.Metrics {
		persMetrics[i] = repository.CompetitorMetric{
			ID:             pbMetric.Id,
			CompetitorID:   "", // Not applicable for personal metrics
			PostID:         pbMetric.PostId,
			Likes:          int(pbMetric.Likes),
			Shares:         int(pbMetric.Shares),
			Comments:       int(pbMetric.Comments),
			CTR:            pbMetric.ClickThroughRate,
			AvgWatchTime:   pbMetric.AvgWatchTime,
			EngagementRate: pbMetric.EngagementRate,
			PostedAt:       pbMetric.PostedAt.AsTime(),
		}
	}

	// Return the comparison result
	return &ComparisonResult{
		CompetitorMetrics: compMetrics,
		PersonalMetrics:   persMetrics,
		Ratios: MetricRatios{
			LikesRatio:          resp.Ratios.LikesRatio,
			SharesRatio:         resp.Ratios.SharesRatio,
			CommentsRatio:       resp.Ratios.CommentsRatio,
			EngagementRateRatio: resp.Ratios.EngagementRateRatio,
			WatchTimeRatio:      resp.Ratios.WatchTimeRatio,
		},
	}, nil
}
