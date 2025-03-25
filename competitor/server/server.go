package server

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/competitor/pb"
	"github.com/donaldnash/go-competitor/competitor/repository"
	"github.com/donaldnash/go-competitor/competitor/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CompetitorServer implements the gRPC server for competitor service
type CompetitorServer struct {
	pb.UnimplementedCompetitorServiceServer
	service *service.CompetitorService
}

// NewCompetitorServer creates a new CompetitorServer
func NewCompetitorServer(service *service.CompetitorService) *CompetitorServer {
	return &CompetitorServer{
		service: service,
	}
}

// AddCompetitor handles the AddCompetitor RPC call
func (s *CompetitorServer) AddCompetitor(ctx context.Context, req *pb.AddCompetitorRequest) (*pb.Competitor, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Platform == "" {
		return nil, status.Error(codes.InvalidArgument, "platform is required")
	}

	competitor, err := s.service.AddCompetitor(ctx, req.TenantId, req.Name, req.Platform)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Competitor{
		Id:        competitor.ID,
		TenantId:  competitor.TenantID,
		Name:      competitor.Name,
		Platform:  competitor.Platform,
		CreatedAt: timestamppb.New(competitor.CreatedAt),
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
}

// GetCompetitor handles the GetCompetitor RPC call
func (s *CompetitorServer) GetCompetitor(ctx context.Context, req *pb.GetCompetitorRequest) (*pb.Competitor, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.CompetitorId == "" {
		return nil, status.Error(codes.InvalidArgument, "competitor ID is required")
	}

	competitor, err := s.service.GetCompetitor(ctx, req.TenantId, req.CompetitorId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Competitor{
		Id:        competitor.ID,
		TenantId:  competitor.TenantID,
		Name:      competitor.Name,
		Platform:  competitor.Platform,
		CreatedAt: timestamppb.New(competitor.CreatedAt),
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
}

// ListCompetitors handles the ListCompetitors RPC call
func (s *CompetitorServer) ListCompetitors(ctx context.Context, req *pb.ListCompetitorsRequest) (*pb.ListCompetitorsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	competitors, err := s.service.GetCompetitors(ctx, req.TenantId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbCompetitors := make([]*pb.Competitor, len(competitors))
	for i, competitor := range competitors {
		pbCompetitors[i] = &pb.Competitor{
			Id:        competitor.ID,
			TenantId:  competitor.TenantID,
			Name:      competitor.Name,
			Platform:  competitor.Platform,
			CreatedAt: timestamppb.New(competitor.CreatedAt),
			UpdatedAt: timestamppb.New(time.Now()),
		}
	}

	return &pb.ListCompetitorsResponse{
		Competitors: pbCompetitors,
	}, nil
}

// UpdateCompetitor handles the UpdateCompetitor RPC call
func (s *CompetitorServer) UpdateCompetitor(ctx context.Context, req *pb.UpdateCompetitorRequest) (*pb.Competitor, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.CompetitorId == "" {
		return nil, status.Error(codes.InvalidArgument, "competitor ID is required")
	}

	competitor, err := s.service.UpdateCompetitor(ctx, req.TenantId, req.CompetitorId, req.Name, req.Platform)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Competitor{
		Id:        competitor.ID,
		TenantId:  competitor.TenantID,
		Name:      competitor.Name,
		Platform:  competitor.Platform,
		CreatedAt: timestamppb.New(competitor.CreatedAt),
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
}

// DeleteCompetitor handles the DeleteCompetitor RPC call
func (s *CompetitorServer) DeleteCompetitor(ctx context.Context, req *pb.DeleteCompetitorRequest) (*emptypb.Empty, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.CompetitorId == "" {
		return nil, status.Error(codes.InvalidArgument, "competitor ID is required")
	}

	err := s.service.DeleteCompetitor(ctx, req.TenantId, req.CompetitorId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// GetCompetitorMetrics handles the GetCompetitorMetrics RPC call
func (s *CompetitorServer) GetCompetitorMetrics(ctx context.Context, req *pb.GetCompetitorMetricsRequest) (*pb.GetCompetitorMetricsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.CompetitorId == "" {
		return nil, status.Error(codes.InvalidArgument, "competitor ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	metrics, err := s.service.GetCompetitorMetrics(ctx, req.TenantId, req.CompetitorId, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbMetrics := make([]*pb.CompetitorMetric, len(metrics))
	for i, metric := range metrics {
		pbMetrics[i] = &pb.CompetitorMetric{
			Id:               metric.ID,
			CompetitorId:     metric.CompetitorID,
			TenantId:         req.TenantId,
			PostId:           metric.PostID,
			Likes:            int32(metric.Likes),
			Shares:           int32(metric.Shares),
			Comments:         int32(metric.Comments),
			ClickThroughRate: metric.CTR,
			AvgWatchTime:     metric.AvgWatchTime,
			EngagementRate:   metric.EngagementRate,
			PostedAt:         timestamppb.New(metric.PostedAt),
			CreatedAt:        timestamppb.New(time.Now()),
		}
	}

	return &pb.GetCompetitorMetricsResponse{
		Metrics: pbMetrics,
	}, nil
}

// TrackCompetitorPost handles the TrackCompetitorPost RPC call
func (s *CompetitorServer) TrackCompetitorPost(ctx context.Context, req *pb.TrackCompetitorPostRequest) (*pb.CompetitorMetric, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.CompetitorId == "" {
		return nil, status.Error(codes.InvalidArgument, "competitor ID is required")
	}

	if req.PostId == "" {
		return nil, status.Error(codes.InvalidArgument, "post ID is required")
	}

	if req.PostedAt == nil {
		return nil, status.Error(codes.InvalidArgument, "posted at date is required")
	}

	// Create the metric
	metric := repository.CompetitorMetric{
		CompetitorID:   req.CompetitorId,
		PostID:         req.PostId,
		Likes:          int(req.Likes),
		Shares:         int(req.Shares),
		Comments:       int(req.Comments),
		CTR:            req.ClickThroughRate,
		AvgWatchTime:   req.AvgWatchTime,
		EngagementRate: req.EngagementRate,
		PostedAt:       req.PostedAt.AsTime(),
	}

	// Save the metric
	metrics := []repository.CompetitorMetric{metric}
	_, err := s.service.UpdateCompetitorMetrics(ctx, req.TenantId, req.CompetitorId, metrics)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return the saved metric
	return &pb.CompetitorMetric{
		Id:               metric.ID,
		CompetitorId:     metric.CompetitorID,
		TenantId:         req.TenantId,
		PostId:           metric.PostID,
		Likes:            int32(metric.Likes),
		Shares:           int32(metric.Shares),
		Comments:         int32(metric.Comments),
		ClickThroughRate: metric.CTR,
		AvgWatchTime:     metric.AvgWatchTime,
		EngagementRate:   metric.EngagementRate,
		PostedAt:         timestamppb.New(metric.PostedAt),
		CreatedAt:        timestamppb.New(time.Now()),
	}, nil
}

// CompareMetrics handles the CompareMetrics RPC call
func (s *CompetitorServer) CompareMetrics(ctx context.Context, req *pb.CompareMetricsRequest) (*pb.CompareMetricsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.CompetitorId == "" {
		return nil, status.Error(codes.InvalidArgument, "competitor ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end dates are required")
	}

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	// Get competitor metrics
	competitorMetrics, err := s.service.GetCompetitorMetrics(ctx, req.TenantId, req.CompetitorId, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get competitor metrics: "+err.Error())
	}

	// For this example, we'll simulate the personal metrics as just the competitor metrics
	// In a real implementation, you would fetch the personal metrics from another service
	var personalMetrics []repository.CompetitorMetric = competitorMetrics

	// Convert competitor metrics to protobuf format
	pbCompetitorMetrics := make([]*pb.CompetitorMetric, len(competitorMetrics))
	for i, metric := range competitorMetrics {
		pbCompetitorMetrics[i] = &pb.CompetitorMetric{
			Id:               metric.ID,
			CompetitorId:     metric.CompetitorID,
			TenantId:         req.TenantId,
			PostId:           metric.PostID,
			Likes:            int32(metric.Likes),
			Shares:           int32(metric.Shares),
			Comments:         int32(metric.Comments),
			ClickThroughRate: metric.CTR,
			AvgWatchTime:     metric.AvgWatchTime,
			EngagementRate:   metric.EngagementRate,
			PostedAt:         timestamppb.New(metric.PostedAt),
			CreatedAt:        timestamppb.New(time.Now()),
		}
	}

	// Convert personal metrics to protobuf format
	pbPersonalMetrics := make([]*pb.PersonalMetric, len(personalMetrics))
	for i, metric := range personalMetrics {
		pbPersonalMetrics[i] = &pb.PersonalMetric{
			Id:               metric.ID,
			TenantId:         req.TenantId,
			PostId:           metric.PostID,
			Likes:            int32(metric.Likes),
			Shares:           int32(metric.Shares),
			Comments:         int32(metric.Comments),
			ClickThroughRate: metric.CTR,
			AvgWatchTime:     metric.AvgWatchTime,
			EngagementRate:   metric.EngagementRate,
			PostedAt:         timestamppb.New(metric.PostedAt),
			CreatedAt:        timestamppb.New(time.Now()),
		}
	}

	// Calculate aggregates
	var totalCompetitorLikes, totalCompetitorShares, totalCompetitorComments int32
	var totalPersonalLikes, totalPersonalShares, totalPersonalComments int32
	var totalCompetitorEngagementRate, totalCompetitorWatchTime float64
	var totalPersonalEngagementRate, totalPersonalWatchTime float64

	for _, metric := range pbCompetitorMetrics {
		totalCompetitorLikes += metric.Likes
		totalCompetitorShares += metric.Shares
		totalCompetitorComments += metric.Comments
		totalCompetitorEngagementRate += metric.EngagementRate
		totalCompetitorWatchTime += metric.AvgWatchTime
	}

	for _, metric := range pbPersonalMetrics {
		totalPersonalLikes += metric.Likes
		totalPersonalShares += metric.Shares
		totalPersonalComments += metric.Comments
		totalPersonalEngagementRate += metric.EngagementRate
		totalPersonalWatchTime += metric.AvgWatchTime
	}

	// Calculate averages
	competitorCount := float64(len(pbCompetitorMetrics))
	personalCount := float64(len(pbPersonalMetrics))

	var avgCompetitorEngagementRate, avgCompetitorWatchTime float64
	var avgPersonalEngagementRate, avgPersonalWatchTime float64

	if competitorCount > 0 {
		avgCompetitorEngagementRate = totalCompetitorEngagementRate / competitorCount
		avgCompetitorWatchTime = totalCompetitorWatchTime / competitorCount
	}

	if personalCount > 0 {
		avgPersonalEngagementRate = totalPersonalEngagementRate / personalCount
		avgPersonalWatchTime = totalPersonalWatchTime / personalCount
	}

	// Calculate ratios
	var likesRatio, sharesRatio, commentsRatio, engagementRateRatio, watchTimeRatio float64

	if totalCompetitorLikes > 0 {
		likesRatio = float64(totalPersonalLikes) / float64(totalCompetitorLikes)
	}
	if totalCompetitorShares > 0 {
		sharesRatio = float64(totalPersonalShares) / float64(totalCompetitorShares)
	}
	if totalCompetitorComments > 0 {
		commentsRatio = float64(totalPersonalComments) / float64(totalCompetitorComments)
	}
	if avgCompetitorEngagementRate > 0 {
		engagementRateRatio = avgPersonalEngagementRate / avgCompetitorEngagementRate
	}
	if avgCompetitorWatchTime > 0 {
		watchTimeRatio = avgPersonalWatchTime / avgCompetitorWatchTime
	}

	// Build response
	return &pb.CompareMetricsResponse{
		Competitor: &pb.CompetitorComparison{
			Metrics: pbCompetitorMetrics,
			Aggregates: &pb.MetricAggregates{
				TotalLikes:        totalCompetitorLikes,
				TotalShares:       totalCompetitorShares,
				TotalComments:     totalCompetitorComments,
				AvgEngagementRate: avgCompetitorEngagementRate,
				AvgWatchTime:      avgCompetitorWatchTime,
			},
		},
		Personal: &pb.PersonalComparison{
			Metrics: pbPersonalMetrics,
			Aggregates: &pb.MetricAggregates{
				TotalLikes:        totalPersonalLikes,
				TotalShares:       totalPersonalShares,
				TotalComments:     totalPersonalComments,
				AvgEngagementRate: avgPersonalEngagementRate,
				AvgWatchTime:      avgPersonalWatchTime,
			},
		},
		Ratios: &pb.ComparisonRatios{
			LikesRatio:          likesRatio,
			SharesRatio:         sharesRatio,
			CommentsRatio:       commentsRatio,
			EngagementRateRatio: engagementRateRatio,
			WatchTimeRatio:      watchTimeRatio,
		},
	}, nil
}
