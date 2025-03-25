package server

import (
	"context"
	"fmt"

	"github.com/donaldnash/go-competitor/analytics/pb"
	"github.com/donaldnash/go-competitor/analytics/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AnalyticsServer is the gRPC server implementation for the analytics service
type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
	service service.AnalyticsService
}

// NewAnalyticsServer creates a new analytics gRPC server
func NewAnalyticsServer(svc service.AnalyticsService) (*AnalyticsServer, error) {
	if svc == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	return &AnalyticsServer{
		service: svc,
	}, nil
}

// GetPostingTimeRecommendations returns optimal posting time recommendations
func (s *AnalyticsServer) GetPostingTimeRecommendations(ctx context.Context, req *pb.PostingTimeRequest) (*pb.PostingTimeResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service layer
	recommendations, err := s.service.GetPostingTimeRecommendations(ctx, req.TenantId, req.DayOfWeek)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posting time recommendations: %v", err)
	}

	// Convert to protobuf response
	protoRecs := make([]*pb.PostingTimeRecommendation, 0, len(recommendations))
	for _, rec := range recommendations {
		protoRecs = append(protoRecs, &pb.PostingTimeRecommendation{
			Id:                      rec.ID,
			TenantId:                rec.TenantID,
			DayOfWeek:               rec.DayOfWeek,
			HourOfDay:               int32(rec.HourOfDay),
			PredictedEngagementRate: rec.PredictedEngagementRate,
			Confidence:              rec.Confidence,
			CreatedAt:               timestamppb.New(rec.CreatedAt),
		})
	}

	return &pb.PostingTimeResponse{
		Recommendations: protoRecs,
	}, nil
}

// GetContentFormatRecommendations returns optimal content format recommendations
func (s *AnalyticsServer) GetContentFormatRecommendations(ctx context.Context, req *pb.ContentFormatRequest) (*pb.ContentFormatResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service layer
	recommendations, err := s.service.GetContentFormatRecommendations(ctx, req.TenantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get content format recommendations: %v", err)
	}

	// Convert to protobuf response
	protoRecs := make([]*pb.ContentFormatRecommendation, 0, len(recommendations))
	for _, rec := range recommendations {
		protoRecs = append(protoRecs, &pb.ContentFormatRecommendation{
			Id:                      rec.ID,
			TenantId:                rec.TenantID,
			Format:                  rec.Format,
			TargetAudience:          rec.TargetAudience,
			PredictedEngagementRate: rec.PredictedEngagementRate,
			Confidence:              rec.Confidence,
			CreatedAt:               timestamppb.New(rec.CreatedAt),
		})
	}

	return &pb.ContentFormatResponse{
		Recommendations: protoRecs,
	}, nil
}

// PredictEngagement predicts engagement metrics for a potential post
func (s *AnalyticsServer) PredictEngagement(ctx context.Context, req *pb.PredictEngagementRequest) (*pb.EngagementPrediction, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.ContentFormat == "" {
		return nil, status.Error(codes.InvalidArgument, "content format is required")
	}

	if req.PostTime == nil {
		return nil, status.Error(codes.InvalidArgument, "post time is required")
	}

	// Convert timestamp to Go time
	postTime := req.PostTime.AsTime()

	// Call service layer
	prediction, err := s.service.PredictEngagement(ctx, req.TenantId, postTime, req.ContentFormat)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to predict engagement: %v", err)
	}

	// Convert to protobuf response
	return &pb.EngagementPrediction{
		Id:                prediction.ID,
		TenantId:          prediction.TenantID,
		PostTime:          timestamppb.New(prediction.PostTime),
		ContentFormat:     prediction.ContentFormat,
		PredictedLikes:    int32(prediction.PredictedLikes),
		PredictedShares:   int32(prediction.PredictedShares),
		PredictedComments: int32(prediction.PredictedComments),
		EngagementRate:    prediction.EngagementRate,
		Confidence:        prediction.Confidence,
		CreatedAt:         timestamppb.New(prediction.CreatedAt),
	}, nil
}

// AnalyzeContentPerformance returns performance analysis for different content types
func (s *AnalyticsServer) AnalyzeContentPerformance(ctx context.Context, req *pb.ContentPerformanceRequest) (*pb.ContentPerformanceResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.StartDate == nil || req.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start date and end date are required")
	}

	// Convert timestamps to Go time
	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	// Call service layer
	performances, err := s.service.AnalyzeContentPerformance(ctx, req.TenantId, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze content performance: %v", err)
	}

	// Convert to protobuf response
	protoPerfs := make([]*pb.ContentPerformance, 0, len(performances))
	for _, perf := range performances {
		protoPerfs = append(protoPerfs, &pb.ContentPerformance{
			Format:            perf.Format,
			TotalPosts:        int32(perf.TotalPosts),
			AvgEngagementRate: perf.AvgEngagementRate,
			AvgLikes:          perf.AvgLikes,
			AvgShares:         perf.AvgShares,
			AvgComments:       perf.AvgComments,
			PerformanceScore:  perf.PerformanceScore,
			PerformanceTrend:  perf.PerformanceTrend,
		})
	}

	return &pb.ContentPerformanceResponse{
		Performances: protoPerfs,
	}, nil
}

// CreateRecommendation creates a new recommendation
func (s *AnalyticsServer) CreateRecommendation(ctx context.Context, req *pb.CreateRecommendationRequest) (*pb.Recommendation, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "recommendation type is required")
	}

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	// Call service layer
	rec, err := s.service.CreateRecommendation(ctx, req.TenantId, req.Type, req.Title, req.Description, req.ExpectedImprovement)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create recommendation: %v", err)
	}

	// Convert to protobuf response
	return &pb.Recommendation{
		Id:                  rec.ID,
		TenantId:            rec.TenantID,
		Type:                rec.Type,
		Title:               rec.Title,
		Description:         rec.Description,
		ExpectedImprovement: rec.ExpectedImprovement,
		Status:              rec.Status,
		CreatedAt:           timestamppb.New(rec.CreatedAt),
		UpdatedAt:           timestamppb.New(rec.UpdatedAt),
	}, nil
}

// GetRecommendations returns recommendations filtered by status
func (s *AnalyticsServer) GetRecommendations(ctx context.Context, req *pb.GetRecommendationsRequest) (*pb.RecommendationsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service layer
	recommendations, err := s.service.GetRecommendations(ctx, req.TenantId, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get recommendations: %v", err)
	}

	// Convert to protobuf response
	protoRecs := make([]*pb.Recommendation, 0, len(recommendations))
	for _, rec := range recommendations {
		protoRecs = append(protoRecs, &pb.Recommendation{
			Id:                  rec.ID,
			TenantId:            rec.TenantID,
			Type:                rec.Type,
			Title:               rec.Title,
			Description:         rec.Description,
			ExpectedImprovement: rec.ExpectedImprovement,
			Status:              rec.Status,
			CreatedAt:           timestamppb.New(rec.CreatedAt),
			UpdatedAt:           timestamppb.New(rec.UpdatedAt),
		})
	}

	return &pb.RecommendationsResponse{
		Recommendations: protoRecs,
	}, nil
}

// UpdateRecommendationStatus updates the status of a recommendation
func (s *AnalyticsServer) UpdateRecommendationStatus(ctx context.Context, req *pb.UpdateRecommendationStatusRequest) (*pb.UpdateRecommendationStatusResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.RecommendationId == "" {
		return nil, status.Error(codes.InvalidArgument, "recommendation ID is required")
	}

	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	// Call service layer
	err := s.service.UpdateRecommendationStatus(ctx, req.TenantId, req.RecommendationId, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update recommendation status: %v", err)
	}

	return &pb.UpdateRecommendationStatusResponse{
		Success: true,
	}, nil
}
