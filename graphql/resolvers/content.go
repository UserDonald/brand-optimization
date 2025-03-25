package resolvers

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/content/client"
	"github.com/donaldnash/go-competitor/graphql/models"
)

// ContentResolver handles all content-related GraphQL queries and mutations
type ContentResolver struct {
	client client.ContentClient
}

// NewContentResolver creates a new ContentResolver
func NewContentResolver(client client.ContentClient) *ContentResolver {
	return &ContentResolver{
		client: client,
	}
}

// GetContentFormats retrieves all content formats for the current tenant
func (r *ContentResolver) GetContentFormats(ctx context.Context, tenantID string) ([]*models.ContentFormat, error) {
	formats, err := r.client.GetContentFormats(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var result []*models.ContentFormat
	for _, f := range formats {
		result = append(result, &models.ContentFormat{
			ID:          f.ID,
			Name:        f.Name,
			Description: f.Description,
		})
	}
	return result, nil
}

// GetContentFormat retrieves a specific content format
func (r *ContentResolver) GetContentFormat(ctx context.Context, tenantID, formatID string) (*models.ContentFormat, error) {
	format, err := r.client.GetContentFormat(ctx, tenantID, formatID)
	if err != nil {
		return nil, err
	}

	return &models.ContentFormat{
		ID:          format.ID,
		Name:        format.Name,
		Description: format.Description,
	}, nil
}

// GetFormatPerformance retrieves performance metrics for a specific content format
func (r *ContentResolver) GetFormatPerformance(ctx context.Context, tenantID, formatID string, dateRange *models.DateRange) ([]*models.FormatPerformance, error) {
	startDate, err := time.Parse(time.RFC3339, dateRange.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(time.RFC3339, dateRange.EndDate)
	if err != nil {
		return nil, err
	}

	performance, err := r.client.GetFormatPerformance(ctx, tenantID, formatID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var result []*models.FormatPerformance
	for _, p := range performance {
		result = append(result, &models.FormatPerformance{
			FormatID:       p.FormatID,
			EngagementRate: p.EngagementRate,
			ReachRate:      p.ReachRate,
			ConversionRate: p.ConversionRate,
		})
	}
	return result, nil
}

// CreateContentFormat creates a new content format
func (r *ContentResolver) CreateContentFormat(ctx context.Context, input *models.CreateContentFormatInput) (*models.ContentFormat, error) {
	format, err := r.client.CreateContentFormat(ctx, input.TenantID, input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	return &models.ContentFormat{
		ID:          format.ID,
		Name:        format.Name,
		Description: format.Description,
	}, nil
}

// UpdateContentFormat updates an existing content format
func (r *ContentResolver) UpdateContentFormat(ctx context.Context, tenantID, formatID string, input *models.UpdateContentFormatInput) (*models.ContentFormat, error) {
	format, err := r.client.UpdateContentFormat(ctx, tenantID, formatID, input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	return &models.ContentFormat{
		ID:          format.ID,
		Name:        format.Name,
		Description: format.Description,
	}, nil
}

// DeleteContentFormat deletes a content format
func (r *ContentResolver) DeleteContentFormat(ctx context.Context, tenantID, formatID string) (bool, error) {
	err := r.client.DeleteContentFormat(ctx, tenantID, formatID)
	return err == nil, err
}

// GetScheduledPosts retrieves all scheduled posts for the current tenant
func (r *ContentResolver) GetScheduledPosts(ctx context.Context, tenantID string) ([]*models.ScheduledPost, error) {
	posts, err := r.client.GetScheduledPosts(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var result []*models.ScheduledPost
	for _, p := range posts {
		result = append(result, &models.ScheduledPost{
			ID:            p.ID,
			Content:       p.Content,
			ScheduledTime: p.ScheduledTime.Format(time.RFC3339),
			Platform:      p.Platform,
			Format:        p.Format,
			Status:        p.Status,
		})
	}
	return result, nil
}

// GetScheduledPost retrieves a specific scheduled post
func (r *ContentResolver) GetScheduledPost(ctx context.Context, tenantID, postID string) (*models.ScheduledPost, error) {
	post, err := r.client.GetScheduledPost(ctx, tenantID, postID)
	if err != nil {
		return nil, err
	}

	return &models.ScheduledPost{
		ID:            post.ID,
		Content:       post.Content,
		ScheduledTime: post.ScheduledTime.Format(time.RFC3339),
		Platform:      post.Platform,
		Format:        post.Format,
		Status:        post.Status,
	}, nil
}

// SchedulePost schedules a new post
func (r *ContentResolver) SchedulePost(ctx context.Context, tenantID string, input *models.SchedulePostInput) (*models.ScheduledPost, error) {
	scheduledTime, err := time.Parse(time.RFC3339, input.ScheduledTime)
	if err != nil {
		return nil, err
	}

	post, err := r.client.SchedulePost(ctx, tenantID, input.Content, input.Platform, input.Format, scheduledTime)
	if err != nil {
		return nil, err
	}

	return &models.ScheduledPost{
		ID:            post.ID,
		Content:       post.Content,
		ScheduledTime: post.ScheduledTime.Format(time.RFC3339),
		Platform:      post.Platform,
		Format:        post.Format,
		Status:        post.Status,
	}, nil
}

// UpdateScheduledPost updates an existing scheduled post
func (r *ContentResolver) UpdateScheduledPost(ctx context.Context, tenantID, postID string, input *models.UpdateScheduledPostInput) (*models.ScheduledPost, error) {
	var scheduledTime time.Time
	if input.ScheduledTime != "" {
		var err error
		scheduledTime, err = time.Parse(time.RFC3339, input.ScheduledTime)
		if err != nil {
			return nil, err
		}
	}

	post, err := r.client.UpdateScheduledPost(ctx, tenantID, postID, input.Content, input.Platform, input.Format, input.Status, scheduledTime)
	if err != nil {
		return nil, err
	}

	return &models.ScheduledPost{
		ID:            post.ID,
		Content:       post.Content,
		ScheduledTime: post.ScheduledTime.Format(time.RFC3339),
		Platform:      post.Platform,
		Format:        post.Format,
		Status:        post.Status,
	}, nil
}

// DeleteScheduledPost deletes a scheduled post
func (r *ContentResolver) DeleteScheduledPost(ctx context.Context, tenantID, postID string) (bool, error) {
	err := r.client.DeleteScheduledPost(ctx, tenantID, postID)
	return err == nil, err
}
