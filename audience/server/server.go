package server

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/audience/pb"
	"github.com/donaldnash/go-competitor/audience/repository"
	"github.com/donaldnash/go-competitor/audience/service"
)

// AudienceServer implements the gRPC audience service
type AudienceServer struct {
	pb.UnimplementedAudienceServiceServer
	service service.AudienceService
}

// NewAudienceServer creates a new AudienceServer
func NewAudienceServer(service service.AudienceService) *AudienceServer {
	return &AudienceServer{
		service: service,
	}
}

// GetSegments returns all audience segments for a tenant
func (s *AudienceServer) GetSegments(ctx context.Context, req *pb.GetSegmentsRequest) (*pb.GetSegmentsResponse, error) {
	segments, err := s.service.GetSegments(ctx, req.TenantId)
	if err != nil {
		return nil, err
	}

	pbSegments := make([]*pb.AudienceSegment, len(segments))
	for i, segment := range segments {
		pbSegments[i] = convertToPbSegment(segment)
	}

	return &pb.GetSegmentsResponse{
		Segments: pbSegments,
	}, nil
}

// GetSegment returns a specific audience segment
func (s *AudienceServer) GetSegment(ctx context.Context, req *pb.GetSegmentRequest) (*pb.GetSegmentResponse, error) {
	segment, err := s.service.GetSegment(ctx, req.TenantId, req.SegmentId)
	if err != nil {
		return nil, err
	}

	return &pb.GetSegmentResponse{
		Segment: convertToPbSegment(*segment),
	}, nil
}

// CreateSegment creates a new audience segment
func (s *AudienceServer) CreateSegment(ctx context.Context, req *pb.CreateSegmentRequest) (*pb.CreateSegmentResponse, error) {
	segment, err := s.service.CreateSegment(ctx, req.TenantId, req.Name, req.Description, req.Type)
	if err != nil {
		return nil, err
	}

	return &pb.CreateSegmentResponse{
		Segment: convertToPbSegment(*segment),
	}, nil
}

// UpdateSegment updates an existing audience segment
func (s *AudienceServer) UpdateSegment(ctx context.Context, req *pb.UpdateSegmentRequest) (*pb.UpdateSegmentResponse, error) {
	segment, err := s.service.UpdateSegment(ctx, req.TenantId, req.SegmentId, req.Name, req.Description, req.Type)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSegmentResponse{
		Segment: convertToPbSegment(*segment),
	}, nil
}

// DeleteSegment deletes an audience segment
func (s *AudienceServer) DeleteSegment(ctx context.Context, req *pb.DeleteSegmentRequest) (*pb.DeleteSegmentResponse, error) {
	err := s.service.DeleteSegment(ctx, req.TenantId, req.SegmentId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteSegmentResponse{
		Success: true,
	}, nil
}

// GetSegmentMetrics returns metrics for a specific audience segment
func (s *AudienceServer) GetSegmentMetrics(ctx context.Context, req *pb.GetSegmentMetricsRequest) (*pb.GetSegmentMetricsResponse, error) {
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, err
	}

	metrics, err := s.service.GetSegmentMetrics(ctx, req.TenantId, req.SegmentId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	pbMetrics := make([]*pb.SegmentMetric, len(metrics))
	for i, metric := range metrics {
		pbMetrics[i] = convertToPbMetric(metric)
	}

	return &pb.GetSegmentMetricsResponse{
		Metrics: pbMetrics,
	}, nil
}

// UpdateSegmentMetrics updates metrics for a specific audience segment
func (s *AudienceServer) UpdateSegmentMetrics(ctx context.Context, req *pb.UpdateSegmentMetricsRequest) (*pb.UpdateSegmentMetricsResponse, error) {
	metrics := make([]repository.SegmentMetric, len(req.Metrics))
	for i, pbMetric := range req.Metrics {
		metrics[i] = convertFromPbMetric(pbMetric)
	}

	count, err := s.service.UpdateSegmentMetrics(ctx, req.TenantId, req.SegmentId, metrics)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSegmentMetricsResponse{
		UpdatedCount: int32(count),
	}, nil
}

// Helper functions to convert between domain and protobuf types

func convertToPbSegment(segment repository.AudienceSegment) *pb.AudienceSegment {
	return &pb.AudienceSegment{
		Id:          segment.ID,
		TenantId:    segment.TenantID,
		Name:        segment.Name,
		Description: segment.Description,
		Type:        segment.Type,
		CreatedAt:   segment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   segment.UpdatedAt.Format(time.RFC3339),
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
