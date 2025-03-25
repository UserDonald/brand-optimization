package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/common/db"
	"github.com/google/uuid"
)

// Common errors
var (
	ErrTenantIDRequired = errors.New("tenant ID is required")
	ErrFormatNotFound   = errors.New("content format not found")
	ErrPostNotFound     = errors.New("scheduled post not found")
)

// ContentRepository defines the interface for content data access
type ContentRepository interface {
	// Content format management
	GetContentFormats(ctx context.Context, tenantID string) ([]ContentFormat, error)
	GetContentFormat(ctx context.Context, tenantID, formatID string) (*ContentFormat, error)
	CreateContentFormat(ctx context.Context, format *ContentFormat) (*ContentFormat, error)
	UpdateContentFormat(ctx context.Context, format *ContentFormat) (*ContentFormat, error)
	DeleteContentFormat(ctx context.Context, tenantID, formatID string) error

	// Content format performance
	GetFormatPerformance(ctx context.Context, tenantID, formatID string, startDate, endDate time.Time) ([]FormatPerformance, error)
	UpdateFormatPerformance(ctx context.Context, tenantID, formatID string, performance []FormatPerformance) (int, error)

	// Scheduled posts
	GetScheduledPosts(ctx context.Context, tenantID string) ([]ScheduledPost, error)
	GetScheduledPost(ctx context.Context, tenantID, postID string) (*ScheduledPost, error)
	CreateScheduledPost(ctx context.Context, post *ScheduledPost) (*ScheduledPost, error)
	UpdateScheduledPost(ctx context.Context, post *ScheduledPost) (*ScheduledPost, error)
	DeleteScheduledPost(ctx context.Context, tenantID, postID string) error
	GetPostsDue(ctx context.Context, tenantID string, before time.Time) ([]ScheduledPost, error)
}

// ContentFormat represents a content format entity
type ContentFormat struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FormatPerformance represents performance metrics for a content format
type FormatPerformance struct {
	ID              string    `json:"id"`
	FormatID        string    `json:"format_id"`
	EngagementRate  float64   `json:"engagement_rate"`
	ReachRate       float64   `json:"reach_rate"`
	ConversionRate  float64   `json:"conversion_rate"`
	AudienceType    string    `json:"audience_type"`
	MeasurementDate time.Time `json:"measurement_date"`
}

// ScheduledPost represents a scheduled post entity
type ScheduledPost struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	Content       string    `json:"content"`
	ScheduledTime time.Time `json:"scheduled_time"`
	Platform      string    `json:"platform"`
	Format        string    `json:"format"`
	Status        string    `json:"status"` // Pending, Published, Failed, Cancelled
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SupabaseContentRepository implements ContentRepository using Supabase
type SupabaseContentRepository struct {
	client *db.SupabaseClient
}

// NewSupabaseContentRepository creates a new SupabaseContentRepository
func NewSupabaseContentRepository(tenantID string) (*SupabaseContentRepository, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	client, err := db.NewSupabaseClient(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	return &SupabaseContentRepository{
		client: client,
	}, nil
}

// GetContentFormats retrieves all content formats for the current tenant
func (r *SupabaseContentRepository) GetContentFormats(ctx context.Context, tenantID string) ([]ContentFormat, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	var formats []ContentFormat
	err := r.client.Query("content_formats").Select("*").Execute(&formats)
	if err != nil {
		return nil, fmt.Errorf("failed to get content formats: %w", err)
	}
	return formats, nil
}

// GetContentFormat retrieves a specific content format
func (r *SupabaseContentRepository) GetContentFormat(ctx context.Context, tenantID, formatID string) (*ContentFormat, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	if formatID == "" {
		return nil, errors.New("format ID is required")
	}

	var formats []ContentFormat
	err := r.client.Query("content_formats").
		Select("*").
		Where("id", "eq", formatID).
		Execute(&formats)

	if err != nil {
		return nil, fmt.Errorf("failed to get content format: %w", err)
	}

	if len(formats) == 0 {
		return nil, ErrFormatNotFound
	}

	return &formats[0], nil
}

// CreateContentFormat creates a new content format
func (r *SupabaseContentRepository) CreateContentFormat(ctx context.Context, format *ContentFormat) (*ContentFormat, error) {
	if format == nil {
		return nil, errors.New("format cannot be nil")
	}

	if format.ID == "" {
		format.ID = uuid.New().String()
	}

	format.TenantID = r.client.TenantID
	if format.CreatedAt.IsZero() {
		now := time.Now()
		format.CreatedAt = now
		format.UpdatedAt = now
	}

	err := r.client.Insert(ctx, "content_formats", format)
	if err != nil {
		return nil, fmt.Errorf("failed to create content format: %w", err)
	}

	return format, nil
}

// UpdateContentFormat updates an existing content format
func (r *SupabaseContentRepository) UpdateContentFormat(ctx context.Context, format *ContentFormat) (*ContentFormat, error) {
	if format == nil {
		return nil, errors.New("format cannot be nil")
	}

	if format.ID == "" {
		return nil, errors.New("format ID is required")
	}

	// Make sure the format belongs to the tenant
	existing, err := r.GetContentFormat(ctx, r.client.TenantID, format.ID)
	if err != nil {
		return nil, err
	}

	// Keep original tenant ID and created date
	format.TenantID = existing.TenantID
	format.CreatedAt = existing.CreatedAt
	format.UpdatedAt = time.Now()

	err = r.client.Update(ctx, "content_formats", "id", format.ID, format)
	if err != nil {
		return nil, fmt.Errorf("failed to update content format: %w", err)
	}

	return format, nil
}

// DeleteContentFormat deletes a content format
func (r *SupabaseContentRepository) DeleteContentFormat(ctx context.Context, tenantID, formatID string) error {
	if tenantID == "" {
		return ErrTenantIDRequired
	}

	if formatID == "" {
		return errors.New("format ID is required")
	}

	// Verify the format exists and belongs to the tenant
	_, err := r.GetContentFormat(ctx, tenantID, formatID)
	if err != nil {
		return err
	}

	err = r.client.Delete(ctx, "content_formats", "id", formatID)
	if err != nil {
		return fmt.Errorf("failed to delete content format: %w", err)
	}

	return nil
}

// GetFormatPerformance retrieves performance metrics for a specific content format within a date range
func (r *SupabaseContentRepository) GetFormatPerformance(ctx context.Context, tenantID, formatID string, startDate, endDate time.Time) ([]FormatPerformance, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	if formatID == "" {
		return nil, errors.New("format ID is required")
	}

	// First verify the format exists and belongs to the tenant
	_, err := r.GetContentFormat(ctx, tenantID, formatID)
	if err != nil {
		return nil, err
	}

	var performance []FormatPerformance
	err = r.client.Query("format_performance").
		Select("*").
		Where("format_id", "eq", formatID).
		Where("measurement_date", "gte", startDate.Format(time.RFC3339)).
		Where("measurement_date", "lte", endDate.Format(time.RFC3339)).
		Order("measurement_date", false).
		Execute(&performance)

	if err != nil {
		return nil, fmt.Errorf("failed to get format performance: %w", err)
	}

	return performance, nil
}

// UpdateFormatPerformance updates performance metrics for a specific content format
func (r *SupabaseContentRepository) UpdateFormatPerformance(ctx context.Context, tenantID, formatID string, performance []FormatPerformance) (int, error) {
	if tenantID == "" {
		return 0, ErrTenantIDRequired
	}

	if formatID == "" {
		return 0, errors.New("format ID is required")
	}

	if len(performance) == 0 {
		return 0, errors.New("no performance data provided")
	}

	// Verify the format exists and belongs to the tenant
	_, err := r.GetContentFormat(ctx, tenantID, formatID)
	if err != nil {
		return 0, err
	}

	// Insert each performance record
	for i := range performance {
		// Set IDs and ensure format ID is set
		if performance[i].ID == "" {
			performance[i].ID = uuid.New().String()
		}
		performance[i].FormatID = formatID

		err = r.client.Insert(ctx, "format_performance", performance[i])
		if err != nil {
			return i, fmt.Errorf("failed to update performance at index %d: %w", i, err)
		}
	}

	return len(performance), nil
}

// GetScheduledPosts retrieves all scheduled posts for the current tenant
func (r *SupabaseContentRepository) GetScheduledPosts(ctx context.Context, tenantID string) ([]ScheduledPost, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	var posts []ScheduledPost
	err := r.client.Query("scheduled_posts").Select("*").Execute(&posts)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled posts: %w", err)
	}
	return posts, nil
}

// GetScheduledPost retrieves a specific scheduled post
func (r *SupabaseContentRepository) GetScheduledPost(ctx context.Context, tenantID, postID string) (*ScheduledPost, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	if postID == "" {
		return nil, errors.New("post ID is required")
	}

	var posts []ScheduledPost
	err := r.client.Query("scheduled_posts").
		Select("*").
		Where("id", "eq", postID).
		Execute(&posts)

	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled post: %w", err)
	}

	if len(posts) == 0 {
		return nil, ErrPostNotFound
	}

	return &posts[0], nil
}

// CreateScheduledPost creates a new scheduled post
func (r *SupabaseContentRepository) CreateScheduledPost(ctx context.Context, post *ScheduledPost) (*ScheduledPost, error) {
	if post == nil {
		return nil, errors.New("post cannot be nil")
	}

	if post.ID == "" {
		post.ID = uuid.New().String()
	}

	post.TenantID = r.client.TenantID

	// Set default status if not provided
	if post.Status == "" {
		post.Status = "Pending"
	}

	if post.CreatedAt.IsZero() {
		now := time.Now()
		post.CreatedAt = now
		post.UpdatedAt = now
	}

	err := r.client.Insert(ctx, "scheduled_posts", post)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduled post: %w", err)
	}

	return post, nil
}

// UpdateScheduledPost updates an existing scheduled post
func (r *SupabaseContentRepository) UpdateScheduledPost(ctx context.Context, post *ScheduledPost) (*ScheduledPost, error) {
	if post == nil {
		return nil, errors.New("post cannot be nil")
	}

	if post.ID == "" {
		return nil, errors.New("post ID is required")
	}

	// Make sure the post belongs to the tenant
	existing, err := r.GetScheduledPost(ctx, r.client.TenantID, post.ID)
	if err != nil {
		return nil, err
	}

	// Keep original tenant ID and created date
	post.TenantID = existing.TenantID
	post.CreatedAt = existing.CreatedAt
	post.UpdatedAt = time.Now()

	err = r.client.Update(ctx, "scheduled_posts", "id", post.ID, post)
	if err != nil {
		return nil, fmt.Errorf("failed to update scheduled post: %w", err)
	}

	return post, nil
}

// DeleteScheduledPost deletes a scheduled post
func (r *SupabaseContentRepository) DeleteScheduledPost(ctx context.Context, tenantID, postID string) error {
	if tenantID == "" {
		return ErrTenantIDRequired
	}

	if postID == "" {
		return errors.New("post ID is required")
	}

	// Verify the post exists and belongs to the tenant
	_, err := r.GetScheduledPost(ctx, tenantID, postID)
	if err != nil {
		return err
	}

	err = r.client.Delete(ctx, "scheduled_posts", "id", postID)
	if err != nil {
		return fmt.Errorf("failed to delete scheduled post: %w", err)
	}

	return nil
}

// GetPostsDue retrieves all scheduled posts that are due for publishing before the specified time
func (r *SupabaseContentRepository) GetPostsDue(ctx context.Context, tenantID string, before time.Time) ([]ScheduledPost, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	var posts []ScheduledPost
	err := r.client.Query("scheduled_posts").
		Select("*").
		Where("status", "eq", "Pending").
		Where("scheduled_time", "lte", before.Format(time.RFC3339)).
		Execute(&posts)

	if err != nil {
		return nil, fmt.Errorf("failed to get posts due: %w", err)
	}

	return posts, nil
}
