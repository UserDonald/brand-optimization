package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/audience/pb"
	"github.com/donaldnash/go-competitor/audience/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AudienceClient defines the interface for client communication with the audience service
type AudienceClient interface {
	// Segment management
	GetSegments(ctx context.Context, tenantID string) ([]repository.AudienceSegment, error)
	GetSegment(ctx context.Context, tenantID, segmentID string) (*repository.AudienceSegment, error)
	CreateSegment(ctx context.Context, tenantID, name, description, segmentType string) (*repository.AudienceSegment, error)
	UpdateSegment(ctx context.Context, tenantID, segmentID, name, description, segmentType string) (*repository.AudienceSegment, error)
	DeleteSegment(ctx context.Context, tenantID, segmentID string) error

	// Segment metrics
	GetSegmentMetrics(ctx context.Context, tenantID, segmentID string, startDate, endDate time.Time) ([]repository.SegmentMetric, error)
	UpdateSegmentMetrics(ctx context.Context, tenantID, segmentID string, metrics []repository.SegmentMetric) (int, error)

	// Close the client connection
	Close() error
}

// GRPCAudienceClient implements AudienceClient using gRPC
type GRPCAudienceClient struct {
	conn   *grpc.ClientConn
	client pb.AudienceServiceClient
}

// NewGRPCAudienceClient creates a new GRPCAudienceClient
func NewGRPCAudienceClient(serviceURL string) (*GRPCAudienceClient, error) {
	conn, err := grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to audience service: %w", err)
	}

	client := pb.NewAudienceServiceClient(conn)
	return &GRPCAudienceClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *GRPCAudienceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetSegments retrieves all audience segments for the current tenant
func (c *GRPCAudienceClient) GetSegments(ctx context.Context, tenantID string) ([]repository.AudienceSegment, error) {
	resp, err := c.client.GetSegments(ctx, &pb.GetSegmentsRequest{
		TenantId: tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get segments: %w", err)
	}

	segments := make([]repository.AudienceSegment, len(resp.Segments))
	for i, s := range resp.Segments {
		segments[i] = convertFromPbSegment(s)
	}

	return segments, nil
}

// GetSegment retrieves a specific audience segment
func (c *GRPCAudienceClient) GetSegment(ctx context.Context, tenantID, segmentID string) (*repository.AudienceSegment, error) {
	resp, err := c.client.GetSegment(ctx, &pb.GetSegmentRequest{
		TenantId:  tenantID,
		SegmentId: segmentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get segment: %w", err)
	}

	segment := convertFromPbSegment(resp.Segment)
	return &segment, nil
}

// CreateSegment creates a new audience segment
func (c *GRPCAudienceClient) CreateSegment(ctx context.Context, tenantID, name, description, segmentType string) (*repository.AudienceSegment, error) {
	resp, err := c.client.CreateSegment(ctx, &pb.CreateSegmentRequest{
		TenantId:    tenantID,
		Name:        name,
		Description: description,
		Type:        segmentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create segment: %w", err)
	}

	segment := convertFromPbSegment(resp.Segment)
	return &segment, nil
}

// UpdateSegment updates an existing audience segment
func (c *GRPCAudienceClient) UpdateSegment(ctx context.Context, tenantID, segmentID, name, description, segmentType string) (*repository.AudienceSegment, error) {
	resp, err := c.client.UpdateSegment(ctx, &pb.UpdateSegmentRequest{
		TenantId:    tenantID,
		SegmentId:   segmentID,
		Name:        name,
		Description: description,
		Type:        segmentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update segment: %w", err)
	}

	segment := convertFromPbSegment(resp.Segment)
	return &segment, nil
}

// DeleteSegment deletes an audience segment
func (c *GRPCAudienceClient) DeleteSegment(ctx context.Context, tenantID, segmentID string) error {
	_, err := c.client.DeleteSegment(ctx, &pb.DeleteSegmentRequest{
		TenantId:  tenantID,
		SegmentId: segmentID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete segment: %w", err)
	}

	return nil
}

// GetSegmentMetrics retrieves metrics for a specific audience segment
func (c *GRPCAudienceClient) GetSegmentMetrics(ctx context.Context, tenantID, segmentID string, startDate, endDate time.Time) ([]repository.SegmentMetric, error) {
	resp, err := c.client.GetSegmentMetrics(ctx, &pb.GetSegmentMetricsRequest{
		TenantId:  tenantID,
		SegmentId: segmentID,
		StartDate: startDate.Format(time.RFC3339),
		EndDate:   endDate.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get segment metrics: %w", err)
	}

	metrics := make([]repository.SegmentMetric, len(resp.Metrics))
	for i, m := range resp.Metrics {
		metrics[i] = convertFromPbMetric(m)
	}

	return metrics, nil
}

// UpdateSegmentMetrics updates metrics for a specific audience segment
func (c *GRPCAudienceClient) UpdateSegmentMetrics(ctx context.Context, tenantID, segmentID string, metrics []repository.SegmentMetric) (int, error) {
	pbMetrics := make([]*pb.SegmentMetric, len(metrics))
	for i, m := range metrics {
		pbMetrics[i] = convertToPbMetric(m)
	}

	resp, err := c.client.UpdateSegmentMetrics(ctx, &pb.UpdateSegmentMetricsRequest{
		TenantId:  tenantID,
		SegmentId: segmentID,
		Metrics:   pbMetrics,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to update segment metrics: %w", err)
	}

	return int(resp.UpdatedCount), nil
}

// Helper functions to convert between domain and protobuf types

func convertFromPbSegment(pbSegment *pb.AudienceSegment) repository.AudienceSegment {
	createdAt, _ := time.Parse(time.RFC3339, pbSegment.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbSegment.UpdatedAt)

	return repository.AudienceSegment{
		ID:          pbSegment.Id,
		TenantID:    pbSegment.TenantId,
		Name:        pbSegment.Name,
		Description: pbSegment.Description,
		Type:        pbSegment.Type,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func convertFromPbMetric(pbMetric *pb.SegmentMetric) repository.SegmentMetric {
	measurementDate, _ := time.Parse(time.RFC3339, pbMetric.MeasurementDate)

	return repository.SegmentMetric{
		ID:                pbMetric.Id,
		SegmentID:         pbMetric.SegmentId,
		Size:              int(pbMetric.Size),
		EngagementRate:    pbMetric.EngagementRate,
		ContentPreference: pbMetric.ContentPreference,
		ResponseTime:      pbMetric.ResponseTime,
		ConversionRate:    pbMetric.ConversionRate,
		TopicalInterest:   pbMetric.TopicalInterest,
		DeviceType:        pbMetric.DeviceType,
		EngagementFreq:    pbMetric.EngagementFreq,
		SentimentTendency: pbMetric.SentimentTendency,
		MeasurementDate:   measurementDate,
	}
}

func convertToPbMetric(metric repository.SegmentMetric) *pb.SegmentMetric {
	return &pb.SegmentMetric{
		Id:                metric.ID,
		SegmentId:         metric.SegmentID,
		Size:              int32(metric.Size),
		EngagementRate:    metric.EngagementRate,
		ContentPreference: metric.ContentPreference,
		ResponseTime:      metric.ResponseTime,
		ConversionRate:    metric.ConversionRate,
		TopicalInterest:   metric.TopicalInterest,
		DeviceType:        metric.DeviceType,
		EngagementFreq:    metric.EngagementFreq,
		SentimentTendency: metric.SentimentTendency,
		MeasurementDate:   metric.MeasurementDate.Format(time.RFC3339),
	}
}
