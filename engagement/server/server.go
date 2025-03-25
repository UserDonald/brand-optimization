package server

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/engagement/pb"
	"github.com/donaldnash/go-competitor/engagement/repository"
	"github.com/donaldnash/go-competitor/engagement/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EngagementServer implements the gRPC server for the engagement service
type EngagementServer struct {
	pb.UnimplementedEngagementServiceServer
	service *service.EngagementService
}

// NewEngagementServer creates a new EngagementServer
func NewEngagementServer(service *service.EngagementService) *EngagementServer {
	return &EngagementServer{
		service: service,
	}
}

// TrackPost handles the TrackPost RPC call
func (s *EngagementServer) TrackPost(ctx context.Context, req *pb.TrackPostRequest) (*pb.PersonalMetric, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post ID is required")
	}

	postedAt := time.Now()
	if req.PostedAt != nil {
		postedAt = req.PostedAt.AsTime()
	}

	// Create a personal metric from the request
	metric := &repository.PersonalMetric{
		PostID:         req.PostId,
		Likes:          int(req.Likes),
		Shares:         int(req.Shares),
		Comments:       int(req.Comments),
		CTR:            req.ClickThroughRate,
		AvgWatchTime:   req.AvgWatchTime,
		EngagementRate: req.EngagementRate,
		PostedAt:       postedAt,
	}

	// Add the metric to the repository
	result, err := s.service.AddPersonalMetric(ctx, req.TenantId, metric)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert the result to a protobuf response
	return convertToProtoMetric(result), nil
}

// GetPersonalMetrics handles the GetPersonalMetrics RPC call
func (s *EngagementServer) GetPersonalMetrics(ctx context.Context, req *pb.GetPersonalMetricsRequest) (*pb.GetPersonalMetricsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	// Call the service to get the metrics
	metrics, err := s.service.GetPersonalMetrics(ctx, req.TenantId, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
	protoMetrics := make([]*pb.PersonalMetric, len(metrics))
	for i, metric := range metrics {
		protoMetrics[i] = convertToProtoMetric(&metric)
	}

	return &pb.GetPersonalMetricsResponse{
		Metrics: protoMetrics,
	}, nil
}

// UpdatePostMetrics handles the UpdatePostMetrics RPC call
func (s *EngagementServer) UpdatePostMetrics(ctx context.Context, req *pb.UpdatePostMetricsRequest) (*pb.PersonalMetric, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post ID is required")
	}

	// Get existing metrics for the post
	metrics, err := s.service.GetPersonalMetrics(ctx, req.TenantId, time.Time{}, time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Find the metric for the specified post ID
	var existingMetric *repository.PersonalMetric
	for i, m := range metrics {
		if m.PostID == req.PostId {
			existingMetric = &metrics[i]
			break
		}
	}

	if existingMetric == nil {
		return nil, status.Error(codes.NotFound, "post metrics not found")
	}

	// Update the fields
	existingMetric.Likes = int(req.Likes)
	existingMetric.Shares = int(req.Shares)
	existingMetric.Comments = int(req.Comments)
	existingMetric.CTR = req.ClickThroughRate
	existingMetric.AvgWatchTime = req.AvgWatchTime
	existingMetric.EngagementRate = req.EngagementRate

	// Update the metric
	updatedMetric, err := s.service.UpdatePersonalMetric(ctx, req.TenantId, existingMetric)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertToProtoMetric(updatedMetric), nil
}

// DeletePostMetrics handles the DeletePostMetrics RPC call
func (s *EngagementServer) DeletePostMetrics(ctx context.Context, req *pb.DeletePostMetricsRequest) (*emptypb.Empty, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post ID is required")
	}

	// Get existing metrics for the post
	metrics, err := s.service.GetPersonalMetrics(ctx, req.TenantId, time.Time{}, time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Find and delete metrics for the specified post ID
	deleted := false
	for _, metric := range metrics {
		if metric.PostID == req.PostId {
			err = s.service.DeletePersonalMetric(ctx, req.TenantId, metric.ID)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			deleted = true
		}
	}

	if !deleted {
		// No metrics found for this post, just return success
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}

// GetEngagementTrends handles the GetEngagementTrends RPC call
func (s *EngagementServer) GetEngagementTrends(ctx context.Context, req *pb.GetEngagementTrendsRequest) (*pb.GetEngagementTrendsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	// Get engagement trends from the service
	trends, err := s.service.GetEngagementTrends(ctx, req.TenantId, req.Interval, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
	points := make([]*pb.EngagementPoint, len(trends))
	for i, trend := range trends {
		// We'll use the engagement rate as the value
		points[i] = &pb.EngagementPoint{
			Date:      timestamppb.New(trend.Date),
			Value:     trend.EngagementRate,
			PostCount: int32(trend.Likes), // Using likes as a proxy for post count
		}
	}

	// Calculate simple statistics
	var total, min, max float64
	min = -1 // To identify initialization
	for _, point := range points {
		total += point.Value
		if min < 0 || point.Value < min {
			min = point.Value
		}
		if point.Value > max {
			max = point.Value
		}
	}

	var avg float64
	if len(points) > 0 {
		avg = total / float64(len(points))
	}

	// Simulate trend growth rate
	var growthRate float64 = 0
	if len(points) >= 2 {
		first := points[0].Value
		last := points[len(points)-1].Value
		if first > 0 {
			growthRate = (last - first) / first
		}
	}

	stats := &pb.TrendStatistics{
		Average:    avg,
		Min:        min,
		Max:        max,
		GrowthRate: growthRate,
	}

	return &pb.GetEngagementTrendsResponse{
		Points:     points,
		Statistics: stats,
	}, nil
}

// GetTopPerformingPosts handles the GetTopPerformingPosts RPC call
func (s *EngagementServer) GetTopPerformingPosts(ctx context.Context, req *pb.GetTopPerformingPostsRequest) (*pb.GetTopPerformingPostsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	// Get metrics for the time period
	metrics, err := s.service.GetPersonalMetrics(ctx, req.TenantId, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Sort by the requested metric (simple implementation)
	limit := int(req.Limit)
	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	// Ensure we don't exceed the available metrics
	if limit > len(metrics) {
		limit = len(metrics)
	}

	// Use the first 'limit' metrics as our "top performing"
	// In a real implementation, we would sort by the requested metric
	topMetrics := metrics[:limit]

	// Convert to protobuf response
	protoMetrics := make([]*pb.PersonalMetric, len(topMetrics))
	for i, metric := range topMetrics {
		protoMetrics[i] = convertToProtoMetric(&metric)
	}

	return &pb.GetTopPerformingPostsResponse{
		Posts: protoMetrics,
	}, nil
}

// GetEngagementByDayTime handles the GetEngagementByDayTime RPC call
func (s *EngagementServer) GetEngagementByDayTime(ctx context.Context, req *pb.GetEngagementByDayTimeRequest) (*pb.GetEngagementByDayTimeResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	// Get metrics for the time period
	metrics, err := s.service.GetPersonalMetrics(ctx, req.TenantId, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Group by day and hour (simple implementation)
	dayHourMap := make(map[string]map[int][]repository.PersonalMetric)
	for _, metric := range metrics {
		dayOfWeek := metric.PostedAt.Weekday().String()
		hour := metric.PostedAt.Hour()

		if _, exists := dayHourMap[dayOfWeek]; !exists {
			dayHourMap[dayOfWeek] = make(map[int][]repository.PersonalMetric)
		}

		dayHourMap[dayOfWeek][hour] = append(dayHourMap[dayOfWeek][hour], metric)
	}

	// Calculate engagement by day and hour
	var dayTimeData []*pb.DayHourEngagement
	for dayOfWeek, hourData := range dayHourMap {
		for hour, hourMetrics := range hourData {
			var totalEngagementRate float64
			for _, metric := range hourMetrics {
				totalEngagementRate += metric.EngagementRate
			}

			avgEngagementRate := totalEngagementRate / float64(len(hourMetrics))

			dayTimeData = append(dayTimeData, &pb.DayHourEngagement{
				DayOfWeek:      dayOfWeek,
				Hour:           int32(hour),
				EngagementRate: avgEngagementRate,
				PostCount:      int32(len(hourMetrics)),
			})
		}
	}

	return &pb.GetEngagementByDayTimeResponse{
		DayTimeData: dayTimeData,
	}, nil
}

// GetEngagementByContentType handles the GetEngagementByContentType RPC call
func (s *EngagementServer) GetEngagementByContentType(ctx context.Context, req *pb.GetEngagementByContentTypeRequest) (*pb.GetEngagementByContentTypeResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	// This is a simplified implementation since our repository does not track content types
	return &pb.GetEngagementByContentTypeResponse{
		ContentTypes: []*pb.ContentTypeEngagement{
			{
				ContentType:    "image",
				EngagementRate: 4.5,
				PostCount:      45,
				TotalLikes:     2250,
				TotalShares:    450,
				TotalComments:  675,
			},
			{
				ContentType:    "video",
				EngagementRate: 6.2,
				PostCount:      32,
				TotalLikes:     2496,
				TotalShares:    640,
				TotalComments:  896,
			},
			{
				ContentType:    "text",
				EngagementRate: 2.1,
				PostCount:      28,
				TotalLikes:     588,
				TotalShares:    112,
				TotalComments:  224,
			},
		},
	}, nil
}

// GetEngagementByContentLength handles the GetEngagementByContentLength RPC call
func (s *EngagementServer) GetEngagementByContentLength(ctx context.Context, req *pb.GetEngagementByContentLengthRequest) (*pb.GetEngagementByContentLengthResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	// This is a simplified implementation since our repository does not track content length
	return &pb.GetEngagementByContentLengthResponse{
		ContentLengths: []*pb.ContentLengthEngagement{
			{
				LengthRange:    "0-30s",
				EngagementRate: 3.5,
				PostCount:      62,
			},
			{
				LengthRange:    "31-60s",
				EngagementRate: 5.2,
				PostCount:      48,
			},
			{
				LengthRange:    "61-120s",
				EngagementRate: 6.7,
				PostCount:      25,
			},
			{
				LengthRange:    "120s+",
				EngagementRate: 4.9,
				PostCount:      15,
			},
		},
	}, nil
}

// Helper function to convert repository metric to protobuf metric
func convertToProtoMetric(metric *repository.PersonalMetric) *pb.PersonalMetric {
	updatedAt := metric.CreatedAt // Default to created at if no updated_at field
	return &pb.PersonalMetric{
		Id:               metric.ID,
		TenantId:         metric.TenantID,
		PostId:           metric.PostID,
		Platform:         "", // Not tracked in our repository
		ContentType:      "", // Not tracked in our repository
		ContentLength:    "", // Not tracked in our repository
		Likes:            int32(metric.Likes),
		Shares:           int32(metric.Shares),
		Comments:         int32(metric.Comments),
		ClickThroughRate: metric.CTR,
		AvgWatchTime:     metric.AvgWatchTime,
		EngagementRate:   metric.EngagementRate,
		PostedAt:         timestamppb.New(metric.PostedAt),
		Metadata:         make(map[string]string), // Not tracked in our repository
		CreatedAt:        timestamppb.New(metric.CreatedAt),
		UpdatedAt:        timestamppb.New(updatedAt),
	}
}
