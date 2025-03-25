package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/content/pb"
	"github.com/donaldnash/go-competitor/content/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ContentClient defines the interface for client communication with the content service
type ContentClient interface {
	// Content format management
	GetContentFormats(ctx context.Context, tenantID string) ([]repository.ContentFormat, error)
	GetContentFormat(ctx context.Context, tenantID, formatID string) (*repository.ContentFormat, error)
	CreateContentFormat(ctx context.Context, tenantID, name, description string) (*repository.ContentFormat, error)
	UpdateContentFormat(ctx context.Context, tenantID, formatID, name, description string) (*repository.ContentFormat, error)
	DeleteContentFormat(ctx context.Context, tenantID, formatID string) error

	// Content format performance
	GetFormatPerformance(ctx context.Context, tenantID, formatID string, startDate, endDate time.Time) ([]repository.FormatPerformance, error)
	UpdateFormatPerformance(ctx context.Context, tenantID, formatID string, performance []repository.FormatPerformance) (int, error)

	// Scheduled posts
	GetScheduledPosts(ctx context.Context, tenantID string) ([]repository.ScheduledPost, error)
	GetScheduledPost(ctx context.Context, tenantID, postID string) (*repository.ScheduledPost, error)
	SchedulePost(ctx context.Context, tenantID, content, platform, format string, scheduledTime time.Time) (*repository.ScheduledPost, error)
	UpdateScheduledPost(ctx context.Context, tenantID, postID, content, platform, format, status string, scheduledTime time.Time) (*repository.ScheduledPost, error)
	DeleteScheduledPost(ctx context.Context, tenantID, postID string) error
	GetPostsDue(ctx context.Context, tenantID string, before time.Time) ([]repository.ScheduledPost, error)

	// Close closes the client connection
	Close() error
}

// GRPCContentClient implements ContentClient using gRPC
type GRPCContentClient struct {
	conn   *grpc.ClientConn
	client pb.ContentServiceClient
}

// NewGRPCContentClient creates a new GRPCContentClient
func NewGRPCContentClient(serverAddr string) (*GRPCContentClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to content service: %w", err)
	}

	client := pb.NewContentServiceClient(conn)
	return &GRPCContentClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *GRPCContentClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetContentFormats retrieves all content formats for the current tenant
func (c *GRPCContentClient) GetContentFormats(ctx context.Context, tenantID string) ([]repository.ContentFormat, error) {
	resp, err := c.client.GetContentFormats(ctx, &pb.GetContentFormatsRequest{
		TenantId: tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get content formats: %w", err)
	}

	formats := make([]repository.ContentFormat, len(resp.Formats))
	for i, f := range resp.Formats {
		formats[i] = convertFromPbFormat(f)
	}

	return formats, nil
}

// GetContentFormat retrieves a specific content format
func (c *GRPCContentClient) GetContentFormat(ctx context.Context, tenantID, formatID string) (*repository.ContentFormat, error) {
	resp, err := c.client.GetContentFormat(ctx, &pb.GetContentFormatRequest{
		TenantId: tenantID,
		FormatId: formatID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get content format: %w", err)
	}

	format := convertFromPbFormat(resp.Format)
	return &format, nil
}

// CreateContentFormat creates a new content format
func (c *GRPCContentClient) CreateContentFormat(ctx context.Context, tenantID, name, description string) (*repository.ContentFormat, error) {
	resp, err := c.client.CreateContentFormat(ctx, &pb.CreateContentFormatRequest{
		TenantId:    tenantID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create content format: %w", err)
	}

	format := convertFromPbFormat(resp.Format)
	return &format, nil
}

// UpdateContentFormat updates an existing content format
func (c *GRPCContentClient) UpdateContentFormat(ctx context.Context, tenantID, formatID, name, description string) (*repository.ContentFormat, error) {
	resp, err := c.client.UpdateContentFormat(ctx, &pb.UpdateContentFormatRequest{
		TenantId:    tenantID,
		FormatId:    formatID,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update content format: %w", err)
	}

	format := convertFromPbFormat(resp.Format)
	return &format, nil
}

// DeleteContentFormat deletes a content format
func (c *GRPCContentClient) DeleteContentFormat(ctx context.Context, tenantID, formatID string) error {
	resp, err := c.client.DeleteContentFormat(ctx, &pb.DeleteContentFormatRequest{
		TenantId: tenantID,
		FormatId: formatID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete content format: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("delete operation was not successful")
	}

	return nil
}

// GetFormatPerformance retrieves performance metrics for a specific content format
func (c *GRPCContentClient) GetFormatPerformance(ctx context.Context, tenantID, formatID string, startDate, endDate time.Time) ([]repository.FormatPerformance, error) {
	resp, err := c.client.GetFormatPerformance(ctx, &pb.GetFormatPerformanceRequest{
		TenantId:  tenantID,
		FormatId:  formatID,
		StartDate: startDate.Format(time.RFC3339),
		EndDate:   endDate.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get format performance: %w", err)
	}

	performance := make([]repository.FormatPerformance, len(resp.Performance))
	for i, p := range resp.Performance {
		performance[i] = convertFromPbPerformance(p)
	}

	return performance, nil
}

// UpdateFormatPerformance updates performance metrics for a specific content format
func (c *GRPCContentClient) UpdateFormatPerformance(ctx context.Context, tenantID, formatID string, performance []repository.FormatPerformance) (int, error) {
	pbPerformance := make([]*pb.FormatPerformance, len(performance))
	for i, p := range performance {
		pbPerformance[i] = convertToPbPerformance(p)
	}

	resp, err := c.client.UpdateFormatPerformance(ctx, &pb.UpdateFormatPerformanceRequest{
		TenantId:    tenantID,
		FormatId:    formatID,
		Performance: pbPerformance,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to update format performance: %w", err)
	}

	return int(resp.UpdatedCount), nil
}

// GetScheduledPosts retrieves all scheduled posts for the current tenant
func (c *GRPCContentClient) GetScheduledPosts(ctx context.Context, tenantID string) ([]repository.ScheduledPost, error) {
	resp, err := c.client.GetScheduledPosts(ctx, &pb.GetScheduledPostsRequest{
		TenantId: tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled posts: %w", err)
	}

	posts := make([]repository.ScheduledPost, len(resp.Posts))
	for i, p := range resp.Posts {
		posts[i] = convertFromPbPost(p)
	}

	return posts, nil
}

// GetScheduledPost retrieves a specific scheduled post
func (c *GRPCContentClient) GetScheduledPost(ctx context.Context, tenantID, postID string) (*repository.ScheduledPost, error) {
	resp, err := c.client.GetScheduledPost(ctx, &pb.GetScheduledPostRequest{
		TenantId: tenantID,
		PostId:   postID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled post: %w", err)
	}

	post := convertFromPbPost(resp.Post)
	return &post, nil
}

// SchedulePost schedules a new post
func (c *GRPCContentClient) SchedulePost(ctx context.Context, tenantID, content, platform, format string, scheduledTime time.Time) (*repository.ScheduledPost, error) {
	resp, err := c.client.SchedulePost(ctx, &pb.SchedulePostRequest{
		TenantId:      tenantID,
		Content:       content,
		Platform:      platform,
		Format:        format,
		ScheduledTime: scheduledTime.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to schedule post: %w", err)
	}

	post := convertFromPbPost(resp.Post)
	return &post, nil
}

// UpdateScheduledPost updates an existing scheduled post
func (c *GRPCContentClient) UpdateScheduledPost(ctx context.Context, tenantID, postID, content, platform, format, status string, scheduledTime time.Time) (*repository.ScheduledPost, error) {
	scheduledTimeStr := ""
	if !scheduledTime.IsZero() {
		scheduledTimeStr = scheduledTime.Format(time.RFC3339)
	}

	resp, err := c.client.UpdateScheduledPost(ctx, &pb.UpdateScheduledPostRequest{
		TenantId:      tenantID,
		PostId:        postID,
		Content:       content,
		Platform:      platform,
		Format:        format,
		Status:        status,
		ScheduledTime: scheduledTimeStr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update scheduled post: %w", err)
	}

	post := convertFromPbPost(resp.Post)
	return &post, nil
}

// DeleteScheduledPost deletes a scheduled post
func (c *GRPCContentClient) DeleteScheduledPost(ctx context.Context, tenantID, postID string) error {
	resp, err := c.client.DeleteScheduledPost(ctx, &pb.DeleteScheduledPostRequest{
		TenantId: tenantID,
		PostId:   postID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete scheduled post: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("delete operation was not successful")
	}

	return nil
}

// GetPostsDue retrieves all scheduled posts that are due for publishing
func (c *GRPCContentClient) GetPostsDue(ctx context.Context, tenantID string, before time.Time) ([]repository.ScheduledPost, error) {
	resp, err := c.client.GetPostsDue(ctx, &pb.GetPostsDueRequest{
		TenantId: tenantID,
		Before:   before.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get posts due: %w", err)
	}

	posts := make([]repository.ScheduledPost, len(resp.Posts))
	for i, p := range resp.Posts {
		posts[i] = convertFromPbPost(p)
	}

	return posts, nil
}

// Helper functions to convert between domain and protobuf types

func convertFromPbFormat(pbFormat *pb.ContentFormat) repository.ContentFormat {
	createdAt, _ := time.Parse(time.RFC3339, pbFormat.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbFormat.UpdatedAt)

	return repository.ContentFormat{
		ID:          pbFormat.Id,
		TenantID:    pbFormat.TenantId,
		Name:        pbFormat.Name,
		Description: pbFormat.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
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

func convertFromPbPost(pbPost *pb.ScheduledPost) repository.ScheduledPost {
	scheduledTime, _ := time.Parse(time.RFC3339, pbPost.ScheduledTime)
	createdAt, _ := time.Parse(time.RFC3339, pbPost.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbPost.UpdatedAt)

	return repository.ScheduledPost{
		ID:            pbPost.Id,
		TenantID:      pbPost.TenantId,
		Content:       pbPost.Content,
		ScheduledTime: scheduledTime,
		Platform:      pbPost.Platform,
		Format:        pbPost.Format,
		Status:        pbPost.Status,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
