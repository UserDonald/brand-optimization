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
	ErrSegmentNotFound  = errors.New("audience segment not found")
)

// AudienceRepository defines the interface for audience data access
type AudienceRepository interface {
	// Segment management
	GetSegments(ctx context.Context, tenantID string) ([]AudienceSegment, error)
	GetSegment(ctx context.Context, tenantID, segmentID string) (*AudienceSegment, error)
	CreateSegment(ctx context.Context, segment *AudienceSegment) (*AudienceSegment, error)
	UpdateSegment(ctx context.Context, segment *AudienceSegment) (*AudienceSegment, error)
	DeleteSegment(ctx context.Context, tenantID, segmentID string) error

	// Segment metrics
	GetSegmentMetrics(ctx context.Context, tenantID, segmentID string, startDate, endDate time.Time) ([]SegmentMetric, error)
	UpdateSegmentMetrics(ctx context.Context, tenantID, segmentID string, metrics []SegmentMetric) (int, error)
}

// AudienceSegment represents an audience segment entity
type AudienceSegment struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // Passive, Reactor, Conversationalist, Content Creator
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SegmentMetric represents engagement metrics for an audience segment
type SegmentMetric struct {
	ID                string    `json:"id"`
	SegmentID         string    `json:"segment_id"`
	Size              int       `json:"size"`
	EngagementRate    float64   `json:"engagement_rate"`
	ContentPreference string    `json:"content_preference"`
	ResponseTime      float64   `json:"response_time"`
	ConversionRate    float64   `json:"conversion_rate"`
	TopicalInterest   string    `json:"topical_interest"`
	DeviceType        string    `json:"device_type"`
	EngagementFreq    string    `json:"engagement_frequency"`
	SentimentTendency string    `json:"sentiment_tendency"`
	MeasurementDate   time.Time `json:"measurement_date"`
}

// SupabaseAudienceRepository implements AudienceRepository using Supabase
type SupabaseAudienceRepository struct {
	client *db.SupabaseClient
}

// NewSupabaseAudienceRepository creates a new SupabaseAudienceRepository
func NewSupabaseAudienceRepository(tenantID string) (*SupabaseAudienceRepository, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	client, err := db.NewSupabaseClient(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	return &SupabaseAudienceRepository{
		client: client,
	}, nil
}

// GetSegments retrieves all audience segments for the current tenant
func (r *SupabaseAudienceRepository) GetSegments(ctx context.Context, tenantID string) ([]AudienceSegment, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	var segments []AudienceSegment
	err := r.client.Query("audience_segments").Select("*").Execute(&segments)
	if err != nil {
		return nil, fmt.Errorf("failed to get audience segments: %w", err)
	}
	return segments, nil
}

// GetSegment retrieves a specific audience segment
func (r *SupabaseAudienceRepository) GetSegment(ctx context.Context, tenantID, segmentID string) (*AudienceSegment, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	if segmentID == "" {
		return nil, errors.New("segment ID is required")
	}

	var segments []AudienceSegment
	err := r.client.Query("audience_segments").
		Select("*").
		Where("id", "eq", segmentID).
		Execute(&segments)

	if err != nil {
		return nil, fmt.Errorf("failed to get audience segment: %w", err)
	}

	if len(segments) == 0 {
		return nil, ErrSegmentNotFound
	}

	return &segments[0], nil
}

// CreateSegment creates a new audience segment
func (r *SupabaseAudienceRepository) CreateSegment(ctx context.Context, segment *AudienceSegment) (*AudienceSegment, error) {
	if segment == nil {
		return nil, errors.New("segment cannot be nil")
	}

	if segment.ID == "" {
		segment.ID = uuid.New().String()
	}

	segment.TenantID = r.client.TenantID
	if segment.CreatedAt.IsZero() {
		now := time.Now()
		segment.CreatedAt = now
		segment.UpdatedAt = now
	}

	err := r.client.Insert(ctx, "audience_segments", segment)
	if err != nil {
		return nil, fmt.Errorf("failed to create audience segment: %w", err)
	}

	return segment, nil
}

// UpdateSegment updates an existing audience segment
func (r *SupabaseAudienceRepository) UpdateSegment(ctx context.Context, segment *AudienceSegment) (*AudienceSegment, error) {
	if segment == nil {
		return nil, errors.New("segment cannot be nil")
	}

	if segment.ID == "" {
		return nil, errors.New("segment ID is required")
	}

	// Make sure the segment belongs to the tenant
	existing, err := r.GetSegment(ctx, r.client.TenantID, segment.ID)
	if err != nil {
		return nil, err
	}

	// Keep original tenant ID and created date
	segment.TenantID = existing.TenantID
	segment.CreatedAt = existing.CreatedAt
	segment.UpdatedAt = time.Now()

	err = r.client.Update(ctx, "audience_segments", "id", segment.ID, segment)
	if err != nil {
		return nil, fmt.Errorf("failed to update audience segment: %w", err)
	}

	return segment, nil
}

// DeleteSegment deletes an audience segment
func (r *SupabaseAudienceRepository) DeleteSegment(ctx context.Context, tenantID, segmentID string) error {
	// Verify the segment exists and belongs to the tenant
	_, err := r.GetSegment(ctx, tenantID, segmentID)
	if err != nil {
		return err
	}

	err = r.client.Delete(ctx, "audience_segments", "id", segmentID)
	if err != nil {
		return fmt.Errorf("failed to delete audience segment: %w", err)
	}

	return nil
}

// GetSegmentMetrics retrieves metrics for a specific audience segment within a date range
func (r *SupabaseAudienceRepository) GetSegmentMetrics(ctx context.Context, tenantID, segmentID string, startDate, endDate time.Time) ([]SegmentMetric, error) {
	if tenantID == "" {
		return nil, ErrTenantIDRequired
	}

	if segmentID == "" {
		return nil, errors.New("segment ID is required")
	}

	// First verify the segment exists and belongs to the tenant
	_, err := r.GetSegment(ctx, tenantID, segmentID)
	if err != nil {
		return nil, err
	}

	var metrics []SegmentMetric
	err = r.client.Query("segment_metrics").
		Select("*").
		Where("segment_id", "eq", segmentID).
		Where("measurement_date", "gte", startDate.Format(time.RFC3339)).
		Where("measurement_date", "lte", endDate.Format(time.RFC3339)).
		Order("measurement_date", false).
		Execute(&metrics)

	if err != nil {
		return nil, fmt.Errorf("failed to get segment metrics: %w", err)
	}

	return metrics, nil
}

// UpdateSegmentMetrics updates metrics for a specific audience segment
func (r *SupabaseAudienceRepository) UpdateSegmentMetrics(ctx context.Context, tenantID, segmentID string, metrics []SegmentMetric) (int, error) {
	if tenantID == "" {
		return 0, ErrTenantIDRequired
	}

	if segmentID == "" {
		return 0, errors.New("segment ID is required")
	}

	if len(metrics) == 0 {
		return 0, errors.New("no metrics provided")
	}

	// Verify the segment exists and belongs to the tenant
	_, err := r.GetSegment(ctx, tenantID, segmentID)
	if err != nil {
		return 0, err
	}

	// Insert each metric
	for i := range metrics {
		// Set IDs and ensure segment ID is set
		if metrics[i].ID == "" {
			metrics[i].ID = uuid.New().String()
		}
		metrics[i].SegmentID = segmentID

		err = r.client.Insert(ctx, "segment_metrics", metrics[i])
		if err != nil {
			return i, fmt.Errorf("failed to update metric at index %d: %w", i, err)
		}
	}

	return len(metrics), nil
}
