package service

import (
	"context"
	"errors"
	"time"

	"github.com/donaldnash/go-competitor/audience/repository"
)

// AudienceService defines the interface for audience service operations
type AudienceService interface {
	// Segment management
	GetSegments(ctx context.Context, tenantID string) ([]repository.AudienceSegment, error)
	GetSegment(ctx context.Context, tenantID, segmentID string) (*repository.AudienceSegment, error)
	CreateSegment(ctx context.Context, tenantID, name, description, segmentType string) (*repository.AudienceSegment, error)
	UpdateSegment(ctx context.Context, tenantID, segmentID, name, description, segmentType string) (*repository.AudienceSegment, error)
	DeleteSegment(ctx context.Context, tenantID, segmentID string) error

	// Segment metrics
	GetSegmentMetrics(ctx context.Context, tenantID, segmentID string, startDate, endDate time.Time) ([]repository.SegmentMetric, error)
	UpdateSegmentMetrics(ctx context.Context, tenantID, segmentID string, metrics []repository.SegmentMetric) (int, error)
}

// audienceService implements the AudienceService interface
type audienceService struct {
	repo repository.AudienceRepository
}

// NewAudienceService creates a new AudienceService instance
func NewAudienceService(repo repository.AudienceRepository) (AudienceService, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}

	return &audienceService{
		repo: repo,
	}, nil
}

// GetSegments retrieves all audience segments for the current tenant
func (s *audienceService) GetSegments(ctx context.Context, tenantID string) ([]repository.AudienceSegment, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.GetSegments(ctx, tenantID)
}

// GetSegment retrieves details about a specific audience segment
func (s *audienceService) GetSegment(ctx context.Context, tenantID, segmentID string) (*repository.AudienceSegment, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if segmentID == "" {
		return nil, errors.New("segment ID is required")
	}

	return s.repo.GetSegment(ctx, tenantID, segmentID)
}

// CreateSegment creates a new audience segment
func (s *audienceService) CreateSegment(ctx context.Context, tenantID, name, description, segmentType string) (*repository.AudienceSegment, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	// Validate segment type
	if segmentType == "" {
		segmentType = "Passive" // Default segment type
	}

	validTypes := map[string]bool{
		"Passive":           true,
		"Reactor":           true,
		"Conversationalist": true,
		"Content Creator":   true,
	}

	if !validTypes[segmentType] {
		return nil, errors.New("invalid segment type: must be Passive, Reactor, Conversationalist, or Content Creator")
	}

	now := time.Now()
	segment := &repository.AudienceSegment{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		Type:        segmentType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.CreateSegment(ctx, segment)
}

// UpdateSegment updates an existing audience segment
func (s *audienceService) UpdateSegment(ctx context.Context, tenantID, segmentID, name, description, segmentType string) (*repository.AudienceSegment, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if segmentID == "" {
		return nil, errors.New("segment ID is required")
	}

	// Get existing segment
	existingSegment, err := s.repo.GetSegment(ctx, tenantID, segmentID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != "" {
		existingSegment.Name = name
	}

	if description != "" {
		existingSegment.Description = description
	}

	if segmentType != "" {
		// Validate segment type
		validTypes := map[string]bool{
			"Passive":           true,
			"Reactor":           true,
			"Conversationalist": true,
			"Content Creator":   true,
		}

		if !validTypes[segmentType] {
			return nil, errors.New("invalid segment type: must be Passive, Reactor, Conversationalist, or Content Creator")
		}

		existingSegment.Type = segmentType
	}

	existingSegment.UpdatedAt = time.Now()

	return s.repo.UpdateSegment(ctx, existingSegment)
}

// DeleteSegment deletes an audience segment
func (s *audienceService) DeleteSegment(ctx context.Context, tenantID, segmentID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if segmentID == "" {
		return errors.New("segment ID is required")
	}

	return s.repo.DeleteSegment(ctx, tenantID, segmentID)
}

// GetSegmentMetrics retrieves metrics for a specific audience segment
func (s *audienceService) GetSegmentMetrics(ctx context.Context, tenantID, segmentID string, startDate, endDate time.Time) ([]repository.SegmentMetric, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if segmentID == "" {
		return nil, errors.New("segment ID is required")
	}

	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	return s.repo.GetSegmentMetrics(ctx, tenantID, segmentID, startDate, endDate)
}

// UpdateSegmentMetrics updates metrics for a specific audience segment
func (s *audienceService) UpdateSegmentMetrics(ctx context.Context, tenantID, segmentID string, metrics []repository.SegmentMetric) (int, error) {
	if tenantID == "" {
		return 0, errors.New("tenant ID is required")
	}

	if segmentID == "" {
		return 0, errors.New("segment ID is required")
	}

	if len(metrics) == 0 {
		return 0, errors.New("no metrics provided")
	}

	return s.repo.UpdateSegmentMetrics(ctx, tenantID, segmentID, metrics)
}
