package server

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/content/pb"
	"github.com/donaldnash/go-competitor/content/repository"
	"github.com/donaldnash/go-competitor/content/service"
)

// ContentServer implements the gRPC content service
type ContentServer struct {
	pb.UnimplementedContentServiceServer
	service service.ContentService
}

// NewContentServer creates a new ContentServer
func NewContentServer(service service.ContentService) *ContentServer {
	return &ContentServer{
		service: service,
	}
}

// GetContentFormats returns all content formats for a tenant
func (s *ContentServer) GetContentFormats(ctx context.Context, req *pb.GetContentFormatsRequest) (*pb.GetContentFormatsResponse, error) {
	formats, err := s.service.GetContentFormats(ctx, req.TenantId)
	if err != nil {
		return nil, err
	}

	pbFormats := make([]*pb.ContentFormat, len(formats))
	for i, format := range formats {
		pbFormats[i] = convertToPbFormat(format)
	}

	return &pb.GetContentFormatsResponse{
		Formats: pbFormats,
	}, nil
}

// GetContentFormat returns a specific content format
func (s *ContentServer) GetContentFormat(ctx context.Context, req *pb.GetContentFormatRequest) (*pb.GetContentFormatResponse, error) {
	format, err := s.service.GetContentFormat(ctx, req.TenantId, req.FormatId)
	if err != nil {
		return nil, err
	}

	return &pb.GetContentFormatResponse{
		Format: convertToPbFormat(*format),
	}, nil
}

// CreateContentFormat creates a new content format
func (s *ContentServer) CreateContentFormat(ctx context.Context, req *pb.CreateContentFormatRequest) (*pb.CreateContentFormatResponse, error) {
	format, err := s.service.CreateContentFormat(ctx, req.TenantId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &pb.CreateContentFormatResponse{
		Format: convertToPbFormat(*format),
	}, nil
}

// UpdateContentFormat updates an existing content format
func (s *ContentServer) UpdateContentFormat(ctx context.Context, req *pb.UpdateContentFormatRequest) (*pb.UpdateContentFormatResponse, error) {
	format, err := s.service.UpdateContentFormat(ctx, req.TenantId, req.FormatId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateContentFormatResponse{
		Format: convertToPbFormat(*format),
	}, nil
}

// DeleteContentFormat deletes a content format
func (s *ContentServer) DeleteContentFormat(ctx context.Context, req *pb.DeleteContentFormatRequest) (*pb.DeleteContentFormatResponse, error) {
	err := s.service.DeleteContentFormat(ctx, req.TenantId, req.FormatId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteContentFormatResponse{
		Success: true,
	}, nil
}

// GetFormatPerformance returns performance metrics for a specific content format
func (s *ContentServer) GetFormatPerformance(ctx context.Context, req *pb.GetFormatPerformanceRequest) (*pb.GetFormatPerformanceResponse, error) {
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, err
	}

	performance, err := s.service.GetFormatPerformance(ctx, req.TenantId, req.FormatId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	pbPerformance := make([]*pb.FormatPerformance, len(performance))
	for i, p := range performance {
		pbPerformance[i] = convertToPbPerformance(p)
	}

	return &pb.GetFormatPerformanceResponse{
		Performance: pbPerformance,
	}, nil
}

// UpdateFormatPerformance updates performance metrics for a specific content format
func (s *ContentServer) UpdateFormatPerformance(ctx context.Context, req *pb.UpdateFormatPerformanceRequest) (*pb.UpdateFormatPerformanceResponse, error) {
	performance := make([]repository.FormatPerformance, len(req.Performance))
	for i, p := range req.Performance {
		performance[i] = convertFromPbPerformance(p)
	}

	count, err := s.service.UpdateFormatPerformance(ctx, req.TenantId, req.FormatId, performance)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateFormatPerformanceResponse{
		UpdatedCount: int32(count),
	}, nil
}

// GetScheduledPosts returns all scheduled posts for a tenant
func (s *ContentServer) GetScheduledPosts(ctx context.Context, req *pb.GetScheduledPostsRequest) (*pb.GetScheduledPostsResponse, error) {
	posts, err := s.service.GetScheduledPosts(ctx, req.TenantId)
	if err != nil {
		return nil, err
	}

	pbPosts := make([]*pb.ScheduledPost, len(posts))
	for i, post := range posts {
		pbPosts[i] = convertToPbPost(post)
	}

	return &pb.GetScheduledPostsResponse{
		Posts: pbPosts,
	}, nil
}

// GetScheduledPost returns a specific scheduled post
func (s *ContentServer) GetScheduledPost(ctx context.Context, req *pb.GetScheduledPostRequest) (*pb.GetScheduledPostResponse, error) {
	post, err := s.service.GetScheduledPost(ctx, req.TenantId, req.PostId)
	if err != nil {
		return nil, err
	}

	return &pb.GetScheduledPostResponse{
		Post: convertToPbPost(*post),
	}, nil
}

// SchedulePost schedules a new post
func (s *ContentServer) SchedulePost(ctx context.Context, req *pb.SchedulePostRequest) (*pb.SchedulePostResponse, error) {
	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
	if err != nil {
		return nil, err
	}

	post, err := s.service.SchedulePost(ctx, req.TenantId, req.Content, req.Platform, req.Format, scheduledTime)
	if err != nil {
		return nil, err
	}

	return &pb.SchedulePostResponse{
		Post: convertToPbPost(*post),
	}, nil
}

// UpdateScheduledPost updates an existing scheduled post
func (s *ContentServer) UpdateScheduledPost(ctx context.Context, req *pb.UpdateScheduledPostRequest) (*pb.UpdateScheduledPostResponse, error) {
	var scheduledTime time.Time
	if req.ScheduledTime != "" {
		var err error
		scheduledTime, err = time.Parse(time.RFC3339, req.ScheduledTime)
		if err != nil {
			return nil, err
		}
	}

	post, err := s.service.UpdateScheduledPost(ctx, req.TenantId, req.PostId, req.Content, req.Platform, req.Format, req.Status, scheduledTime)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateScheduledPostResponse{
		Post: convertToPbPost(*post),
	}, nil
}

// DeleteScheduledPost deletes a scheduled post
func (s *ContentServer) DeleteScheduledPost(ctx context.Context, req *pb.DeleteScheduledPostRequest) (*pb.DeleteScheduledPostResponse, error) {
	err := s.service.DeleteScheduledPost(ctx, req.TenantId, req.PostId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteScheduledPostResponse{
		Success: true,
	}, nil
}

// GetPostsDue returns all scheduled posts that are due for publishing
func (s *ContentServer) GetPostsDue(ctx context.Context, req *pb.GetPostsDueRequest) (*pb.GetPostsDueResponse, error) {
	var before time.Time
	if req.Before != "" {
		var err error
		before, err = time.Parse(time.RFC3339, req.Before)
		if err != nil {
			return nil, err
		}
	} else {
		before = time.Now()
	}

	posts, err := s.service.GetPostsDue(ctx, req.TenantId, before)
	if err != nil {
		return nil, err
	}

	pbPosts := make([]*pb.ScheduledPost, len(posts))
	for i, post := range posts {
		pbPosts[i] = convertToPbPost(post)
	}

	return &pb.GetPostsDueResponse{
		Posts: pbPosts,
	}, nil
}

// Helper functions to convert between domain and protobuf types

func convertToPbFormat(format repository.ContentFormat) *pb.ContentFormat {
	return &pb.ContentFormat{
		Id:          format.ID,
		TenantId:    format.TenantID,
		Name:        format.Name,
		Description: format.Description,
		CreatedAt:   format.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   format.UpdatedAt.Format(time.RFC3339),
	}
}

func convertToPbPerformance(performance repository.FormatPerformance) *pb.FormatPerformance {
	return &pb.FormatPerformance{
		Id:              performance.ID,
		FormatId:        performance.FormatID,
		EngagementRate:  performance.EngagementRate,
		ReachRate:       performance.ReachRate,
		ConversionRate:  performance.ConversionRate,
		AudienceType:    performance.AudienceType,
		MeasurementDate: performance.MeasurementDate.Format(time.RFC3339),
	}
}

func convertFromPbPerformance(pbPerformance *pb.FormatPerformance) repository.FormatPerformance {
	measurementDate, _ := time.Parse(time.RFC3339, pbPerformance.MeasurementDate)

	return repository.FormatPerformance{
		ID:              pbPerformance.Id,
		FormatID:        pbPerformance.FormatId,
		EngagementRate:  pbPerformance.EngagementRate,
		ReachRate:       pbPerformance.ReachRate,
		ConversionRate:  pbPerformance.ConversionRate,
		AudienceType:    pbPerformance.AudienceType,
		MeasurementDate: measurementDate,
	}
}

func convertToPbPost(post repository.ScheduledPost) *pb.ScheduledPost {
	return &pb.ScheduledPost{
		Id:            post.ID,
		TenantId:      post.TenantID,
		Content:       post.Content,
		ScheduledTime: post.ScheduledTime.Format(time.RFC3339),
		Platform:      post.Platform,
		Format:        post.Format,
		Status:        post.Status,
		CreatedAt:     post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     post.UpdatedAt.Format(time.RFC3339),
	}
}
