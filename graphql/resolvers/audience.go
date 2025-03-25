package resolvers

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/audience/client"
	"github.com/donaldnash/go-competitor/graphql/models"
)

// AudienceResolver handles all audience-related GraphQL queries and mutations
type AudienceResolver struct {
	client client.AudienceClient
}

// NewAudienceResolver creates a new AudienceResolver
func NewAudienceResolver(client client.AudienceClient) *AudienceResolver {
	return &AudienceResolver{
		client: client,
	}
}

// GetAudienceSegments retrieves all audience segments for the current tenant
func (r *AudienceResolver) GetAudienceSegments(ctx context.Context, tenantID string) ([]*models.AudienceSegment, error) {
	segments, err := r.client.GetSegments(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var result []*models.AudienceSegment
	for _, s := range segments {
		result = append(result, &models.AudienceSegment{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
		})
	}
	return result, nil
}

// GetAudienceSegment retrieves a specific audience segment
func (r *AudienceResolver) GetAudienceSegment(ctx context.Context, tenantID, segmentID string) (*models.AudienceSegment, error) {
	segment, err := r.client.GetSegment(ctx, tenantID, segmentID)
	if err != nil {
		return nil, err
	}

	return &models.AudienceSegment{
		ID:          segment.ID,
		Name:        segment.Name,
		Description: segment.Description,
	}, nil
}

// GetSegmentMetrics retrieves metrics for a specific audience segment
func (r *AudienceResolver) GetSegmentMetrics(ctx context.Context, tenantID, segmentID string, dateRange *models.DateRange) ([]*models.SegmentMetric, error) {
	startDate, err := time.Parse(time.RFC3339, dateRange.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(time.RFC3339, dateRange.EndDate)
	if err != nil {
		return nil, err
	}

	metrics, err := r.client.GetSegmentMetrics(ctx, tenantID, segmentID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var result []*models.SegmentMetric
	for _, m := range metrics {
		result = append(result, &models.SegmentMetric{
			SegmentID:         m.SegmentID,
			Size:              m.Size,
			EngagementRate:    m.EngagementRate,
			ContentPreference: m.ContentPreference,
		})
	}
	return result, nil
}

// CreateAudienceSegment creates a new audience segment
func (r *AudienceResolver) CreateAudienceSegment(ctx context.Context, input *models.CreateAudienceSegmentInput) (*models.AudienceSegment, error) {
	segment, err := r.client.CreateSegment(ctx, input.TenantID, input.Name, input.Description, input.Type)
	if err != nil {
		return nil, err
	}

	return &models.AudienceSegment{
		ID:          segment.ID,
		Name:        segment.Name,
		Description: segment.Description,
	}, nil
}

// UpdateAudienceSegment updates an existing audience segment
func (r *AudienceResolver) UpdateAudienceSegment(ctx context.Context, tenantID, segmentID string, input *models.UpdateAudienceSegmentInput) (*models.AudienceSegment, error) {
	segment, err := r.client.UpdateSegment(ctx, tenantID, segmentID, input.Name, input.Description, input.Type)
	if err != nil {
		return nil, err
	}

	return &models.AudienceSegment{
		ID:          segment.ID,
		Name:        segment.Name,
		Description: segment.Description,
	}, nil
}

// DeleteAudienceSegment deletes an audience segment
func (r *AudienceResolver) DeleteAudienceSegment(ctx context.Context, tenantID, segmentID string) (bool, error) {
	err := r.client.DeleteSegment(ctx, tenantID, segmentID)
	return err == nil, err
}
