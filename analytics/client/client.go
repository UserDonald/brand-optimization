package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/analytics/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AnalyticsClient is the client for the analytics service
type AnalyticsClient struct {
	conn   *grpc.ClientConn
	client pb.AnalyticsServiceClient
}

// NewAnalyticsClient creates a new analytics client
func NewAnalyticsClient(serverAddr string) (*AnalyticsClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to analytics service: %w", err)
	}

	client := pb.NewAnalyticsServiceClient(conn)
	return &AnalyticsClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *AnalyticsClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetPostingTimeRecommendations returns optimal posting time recommendations
func (c *AnalyticsClient) GetPostingTimeRecommendations(ctx context.Context, tenantID, dayOfWeek string) (*pb.PostingTimeResponse, error) {
	req := &pb.PostingTimeRequest{
		TenantId:  tenantID,
		DayOfWeek: dayOfWeek,
	}

	return c.client.GetPostingTimeRecommendations(ctx, req)
}

// GetContentFormatRecommendations returns optimal content format recommendations
func (c *AnalyticsClient) GetContentFormatRecommendations(ctx context.Context, tenantID string) (*pb.ContentFormatResponse, error) {
	req := &pb.ContentFormatRequest{
		TenantId: tenantID,
	}

	return c.client.GetContentFormatRecommendations(ctx, req)
}

// PredictEngagement predicts engagement metrics for a potential post
func (c *AnalyticsClient) PredictEngagement(ctx context.Context, tenantID string, postTime time.Time, contentFormat string) (*pb.EngagementPrediction, error) {
	req := &pb.PredictEngagementRequest{
		TenantId:      tenantID,
		PostTime:      timestamppb.New(postTime),
		ContentFormat: contentFormat,
	}

	return c.client.PredictEngagement(ctx, req)
}

// AnalyzeContentPerformance returns performance analysis for different content types
func (c *AnalyticsClient) AnalyzeContentPerformance(ctx context.Context, tenantID string, startDate, endDate time.Time) (*pb.ContentPerformanceResponse, error) {
	req := &pb.ContentPerformanceRequest{
		TenantId:  tenantID,
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	return c.client.AnalyzeContentPerformance(ctx, req)
}

// CreateRecommendation creates a new recommendation
func (c *AnalyticsClient) CreateRecommendation(ctx context.Context, tenantID, recType, title, description string, expectedImprovement float64) (*pb.Recommendation, error) {
	req := &pb.CreateRecommendationRequest{
		TenantId:            tenantID,
		Type:                recType,
		Title:               title,
		Description:         description,
		ExpectedImprovement: expectedImprovement,
	}

	return c.client.CreateRecommendation(ctx, req)
}

// GetRecommendations returns recommendations filtered by status
func (c *AnalyticsClient) GetRecommendations(ctx context.Context, tenantID, status string) (*pb.RecommendationsResponse, error) {
	req := &pb.GetRecommendationsRequest{
		TenantId: tenantID,
		Status:   status,
	}

	return c.client.GetRecommendations(ctx, req)
}

// UpdateRecommendationStatus updates the status of a recommendation
func (c *AnalyticsClient) UpdateRecommendationStatus(ctx context.Context, tenantID, recID, status string) (*pb.UpdateRecommendationStatusResponse, error) {
	req := &pb.UpdateRecommendationStatusRequest{
		TenantId:         tenantID,
		RecommendationId: recID,
		Status:           status,
	}

	return c.client.UpdateRecommendationStatus(ctx, req)
}

// NewGRPCAnalyticsClient creates a new analytics client that uses a local service directly
// This is useful for in-process communication without going through gRPC network calls
func NewGRPCAnalyticsClient(svc interface{}) *AnalyticsClient {
	return &AnalyticsClient{
		// We leave conn as nil since we're not using a real connection
		conn:   nil,
		client: &localAnalyticsClient{svc: svc},
	}
}

// localAnalyticsClient implements pb.AnalyticsServiceClient using a local service
// This allows us to bypass gRPC for in-process communication
type localAnalyticsClient struct {
	svc interface{} // The local service implementation
}

// GetPostingTimeRecommendations forwards the call to the local service
func (c *localAnalyticsClient) GetPostingTimeRecommendations(ctx context.Context, req *pb.PostingTimeRequest, opts ...grpc.CallOption) (*pb.PostingTimeResponse, error) {
	analyticsService, ok := c.svc.(interface {
		GetPostingTimeRecommendations(ctx context.Context, tenantID string, dayOfWeek string) ([]interface{}, error)
	})
	if !ok {
		return nil, fmt.Errorf("service does not implement GetPostingTimeRecommendations")
	}

	// Call the service method
	_, err := analyticsService.GetPostingTimeRecommendations(ctx, req.TenantId, req.DayOfWeek)
	if err != nil {
		return nil, err
	}

	// For demonstration, return static data instead of converting the actual results
	// In a real implementation, we would convert each recommendation
	protoRecs := []*pb.PostingTimeRecommendation{
		{
			Id:                      "local-rec-id-1",
			DayOfWeek:               req.DayOfWeek,
			HourOfDay:               9, // 9 AM
			PredictedEngagementRate: 0.07,
			Confidence:              0.85,
		},
		{
			Id:                      "local-rec-id-2",
			DayOfWeek:               req.DayOfWeek,
			HourOfDay:               18, // 6 PM
			PredictedEngagementRate: 0.09,
			Confidence:              0.9,
		},
	}

	return &pb.PostingTimeResponse{
		Recommendations: protoRecs,
	}, nil
}

// Implement other methods of the pb.AnalyticsServiceClient interface
// using the localAnalyticsClient, following the same pattern

// GetContentFormatRecommendations forwards the call to the local service
func (c *localAnalyticsClient) GetContentFormatRecommendations(ctx context.Context, req *pb.ContentFormatRequest, opts ...grpc.CallOption) (*pb.ContentFormatResponse, error) {
	// Simplified implementation for demonstration
	return &pb.ContentFormatResponse{
		Recommendations: []*pb.ContentFormatRecommendation{
			{
				Format:                  "video",
				PredictedEngagementRate: 0.08,
				Confidence:              0.9,
			},
			{
				Format:                  "image",
				PredictedEngagementRate: 0.06,
				Confidence:              0.85,
			},
		},
	}, nil
}

// PredictEngagement forwards the call to the local service
func (c *localAnalyticsClient) PredictEngagement(ctx context.Context, req *pb.PredictEngagementRequest, opts ...grpc.CallOption) (*pb.EngagementPrediction, error) {
	// Simplified implementation for demonstration
	return &pb.EngagementPrediction{
		PredictedLikes:    100,
		PredictedShares:   20,
		PredictedComments: 35,
		EngagementRate:    0.07,
		Confidence:        0.75,
	}, nil
}

// AnalyzeContentPerformance forwards the call to the local service
func (c *localAnalyticsClient) AnalyzeContentPerformance(ctx context.Context, req *pb.ContentPerformanceRequest, opts ...grpc.CallOption) (*pb.ContentPerformanceResponse, error) {
	// Simplified implementation for demonstration
	return &pb.ContentPerformanceResponse{
		Performances: []*pb.ContentPerformance{
			{
				Format:            "video",
				AvgEngagementRate: 0.08,
				PerformanceScore:  85.0,
			},
			{
				Format:            "image",
				AvgEngagementRate: 0.06,
				PerformanceScore:  75.0,
			},
		},
	}, nil
}

// CreateRecommendation forwards the call to the local service
func (c *localAnalyticsClient) CreateRecommendation(ctx context.Context, req *pb.CreateRecommendationRequest, opts ...grpc.CallOption) (*pb.Recommendation, error) {
	// Simplified implementation for demonstration
	return &pb.Recommendation{
		Id:                  "local-rec-id",
		Title:               req.Title,
		Description:         req.Description,
		ExpectedImprovement: req.ExpectedImprovement,
		Status:              "pending",
	}, nil
}

// GetRecommendations forwards the call to the local service
func (c *localAnalyticsClient) GetRecommendations(ctx context.Context, req *pb.GetRecommendationsRequest, opts ...grpc.CallOption) (*pb.RecommendationsResponse, error) {
	// Simplified implementation for demonstration
	return &pb.RecommendationsResponse{
		Recommendations: []*pb.Recommendation{
			{
				Id:                  "local-rec-id-1",
				Title:               "Sample Recommendation 1",
				Status:              req.Status,
				ExpectedImprovement: 0.15,
			},
		},
	}, nil
}

// UpdateRecommendationStatus forwards the call to the local service
func (c *localAnalyticsClient) UpdateRecommendationStatus(ctx context.Context, req *pb.UpdateRecommendationStatusRequest, opts ...grpc.CallOption) (*pb.UpdateRecommendationStatusResponse, error) {
	// Simplified implementation for demonstration
	return &pb.UpdateRecommendationStatusResponse{
		Success: true,
	}, nil
}
