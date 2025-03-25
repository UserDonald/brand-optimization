package service

import (
	"context"
	"errors"
	"time"

	"github.com/donaldnash/go-competitor/content/repository"
)

// ContentService defines the interface for content service operations
type ContentService interface {
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
}

// contentService implements the ContentService interface
type contentService struct {
	repo repository.ContentRepository
}

// NewContentService creates a new ContentService instance
func NewContentService(repo repository.ContentRepository) (ContentService, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}

	return &contentService{
		repo: repo,
	}, nil
}

// GetContentFormats retrieves all content formats for the current tenant
func (s *contentService) GetContentFormats(ctx context.Context, tenantID string) ([]repository.ContentFormat, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.GetContentFormats(ctx, tenantID)
}

// GetContentFormat retrieves a specific content format
func (s *contentService) GetContentFormat(ctx context.Context, tenantID, formatID string) (*repository.ContentFormat, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if formatID == "" {
		return nil, errors.New("format ID is required")
	}

	return s.repo.GetContentFormat(ctx, tenantID, formatID)
}

// CreateContentFormat creates a new content format
func (s *contentService) CreateContentFormat(ctx context.Context, tenantID, name, description string) (*repository.ContentFormat, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	format := &repository.ContentFormat{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
	}

	return s.repo.CreateContentFormat(ctx, format)
}

// UpdateContentFormat updates an existing content format
func (s *contentService) UpdateContentFormat(ctx context.Context, tenantID, formatID, name, description string) (*repository.ContentFormat, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if formatID == "" {
		return nil, errors.New("format ID is required")
	}

	// Get existing format
	existingFormat, err := s.repo.GetContentFormat(ctx, tenantID, formatID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != "" {
		existingFormat.Name = name
	}

	if description != "" {
		existingFormat.Description = description
	}

	return s.repo.UpdateContentFormat(ctx, existingFormat)
}

// DeleteContentFormat deletes a content format
func (s *contentService) DeleteContentFormat(ctx context.Context, tenantID, formatID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if formatID == "" {
		return errors.New("format ID is required")
	}

	return s.repo.DeleteContentFormat(ctx, tenantID, formatID)
}

// GetFormatPerformance retrieves performance metrics for a specific content format
func (s *contentService) GetFormatPerformance(ctx context.Context, tenantID, formatID string, startDate, endDate time.Time) ([]repository.FormatPerformance, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if formatID == "" {
		return nil, errors.New("format ID is required")
	}

	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	return s.repo.GetFormatPerformance(ctx, tenantID, formatID, startDate, endDate)
}

// UpdateFormatPerformance updates performance metrics for a specific content format
func (s *contentService) UpdateFormatPerformance(ctx context.Context, tenantID, formatID string, performance []repository.FormatPerformance) (int, error) {
	if tenantID == "" {
		return 0, errors.New("tenant ID is required")
	}

	if formatID == "" {
		return 0, errors.New("format ID is required")
	}

	if len(performance) == 0 {
		return 0, errors.New("no performance data provided")
	}

	return s.repo.UpdateFormatPerformance(ctx, tenantID, formatID, performance)
}

// GetScheduledPosts retrieves all scheduled posts for the current tenant
func (s *contentService) GetScheduledPosts(ctx context.Context, tenantID string) ([]repository.ScheduledPost, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.GetScheduledPosts(ctx, tenantID)
}

// GetScheduledPost retrieves a specific scheduled post
func (s *contentService) GetScheduledPost(ctx context.Context, tenantID, postID string) (*repository.ScheduledPost, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if postID == "" {
		return nil, errors.New("post ID is required")
	}

	return s.repo.GetScheduledPost(ctx, tenantID, postID)
}

// SchedulePost schedules a new post
func (s *contentService) SchedulePost(ctx context.Context, tenantID, content, platform, format string, scheduledTime time.Time) (*repository.ScheduledPost, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if content == "" {
		return nil, errors.New("content is required")
	}

	if platform == "" {
		return nil, errors.New("platform is required")
	}

	if format == "" {
		return nil, errors.New("format is required")
	}

	if scheduledTime.IsZero() || scheduledTime.Before(time.Now()) {
		return nil, errors.New("scheduled time must be in the future")
	}

	post := &repository.ScheduledPost{
		TenantID:      tenantID,
		Content:       content,
		Platform:      platform,
		Format:        format,
		ScheduledTime: scheduledTime,
		Status:        "Pending",
	}

	return s.repo.CreateScheduledPost(ctx, post)
}

// UpdateScheduledPost updates an existing scheduled post
func (s *contentService) UpdateScheduledPost(ctx context.Context, tenantID, postID, content, platform, format, status string, scheduledTime time.Time) (*repository.ScheduledPost, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if postID == "" {
		return nil, errors.New("post ID is required")
	}

	// Get existing post
	existingPost, err := s.repo.GetScheduledPost(ctx, tenantID, postID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if content != "" {
		existingPost.Content = content
	}

	if platform != "" {
		existingPost.Platform = platform
	}

	if format != "" {
		existingPost.Format = format
	}

	if status != "" {
		validStatuses := map[string]bool{
			"Pending":   true,
			"Published": true,
			"Failed":    true,
			"Cancelled": true,
		}

		if !validStatuses[status] {
			return nil, errors.New("invalid status: must be Pending, Published, Failed, or Cancelled")
		}

		existingPost.Status = status
	}

	if !scheduledTime.IsZero() {
		if scheduledTime.Before(time.Now()) && existingPost.Status == "Pending" {
			return nil, errors.New("scheduled time must be in the future for pending posts")
		}
		existingPost.ScheduledTime = scheduledTime
	}

	return s.repo.UpdateScheduledPost(ctx, existingPost)
}

// DeleteScheduledPost deletes a scheduled post
func (s *contentService) DeleteScheduledPost(ctx context.Context, tenantID, postID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if postID == "" {
		return errors.New("post ID is required")
	}

	return s.repo.DeleteScheduledPost(ctx, tenantID, postID)
}

// GetPostsDue retrieves all scheduled posts that are due for publishing
func (s *contentService) GetPostsDue(ctx context.Context, tenantID string, before time.Time) ([]repository.ScheduledPost, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if before.IsZero() {
		before = time.Now()
	}

	return s.repo.GetPostsDue(ctx, tenantID, before)
}
