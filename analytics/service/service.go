package service

import (
	"context"
	"errors"
	"time"

	"github.com/donaldnash/go-competitor/analytics/repository"
	"github.com/google/uuid"
)

// AnalyticsService defines the interface for analytics service operations
type AnalyticsService interface {
	// Predictive analytics
	GetPostingTimeRecommendations(ctx context.Context, tenantID string, dayOfWeek string) ([]repository.PostingTimeRecommendation, error)
	GetContentFormatRecommendations(ctx context.Context, tenantID string) ([]repository.ContentFormatRecommendation, error)

	// Performance predictions
	PredictEngagement(ctx context.Context, tenantID string, postTime time.Time, contentFormat string) (*repository.EngagementPrediction, error)

	// Content analysis
	AnalyzeContentPerformance(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]repository.ContentPerformance, error)

	// Recommendation management
	CreateRecommendation(ctx context.Context, tenantID, recType, title, description string, expectedImprovement float64) (*repository.Recommendation, error)
	GetRecommendations(ctx context.Context, tenantID string, status string) ([]repository.Recommendation, error)
	UpdateRecommendationStatus(ctx context.Context, tenantID, recID, status string) error
}

// analyticsService implements the AnalyticsService interface
type analyticsService struct {
	repo repository.AnalyticsRepository
}

// NewAnalyticsService creates a new AnalyticsService instance
func NewAnalyticsService(repo repository.AnalyticsRepository) (AnalyticsService, error) {
	if repo == nil {
		return nil, errors.New("repository cannot be nil")
	}

	return &analyticsService{
		repo: repo,
	}, nil
}

// GetPostingTimeRecommendations returns recommended posting times with predicted engagement
func (s *analyticsService) GetPostingTimeRecommendations(ctx context.Context, tenantID string, dayOfWeek string) ([]repository.PostingTimeRecommendation, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Call repository to get recommendations
	return s.repo.GetPostingTimeRecommendations(ctx, tenantID, dayOfWeek)
}

// GetContentFormatRecommendations returns recommended content formats with predicted engagement
func (s *analyticsService) GetContentFormatRecommendations(ctx context.Context, tenantID string) ([]repository.ContentFormatRecommendation, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Call repository to get recommendations
	return s.repo.GetContentFormatRecommendations(ctx, tenantID)
}

// PredictEngagement predicts engagement metrics for a potential post
func (s *analyticsService) PredictEngagement(ctx context.Context, tenantID string, postTime time.Time, contentFormat string) (*repository.EngagementPrediction, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if contentFormat == "" {
		return nil, errors.New("content format is required")
	}

	// Validate postTime is in the future
	if postTime.Before(time.Now()) {
		return nil, errors.New("post time must be in the future")
	}

	// Call repository to get prediction
	return s.repo.PredictEngagement(ctx, tenantID, postTime, contentFormat)
}

// AnalyzeContentPerformance returns performance analysis for different content types
func (s *analyticsService) AnalyzeContentPerformance(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]repository.ContentPerformance, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	// Call repository to get performance analysis
	return s.repo.AnalyzeContentPerformance(ctx, tenantID, startDate, endDate)
}

// CreateRecommendation creates a new recommendation
func (s *analyticsService) CreateRecommendation(ctx context.Context, tenantID, recType, title, description string, expectedImprovement float64) (*repository.Recommendation, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	if recType == "" {
		return nil, errors.New("recommendation type is required")
	}

	if title == "" {
		return nil, errors.New("title is required")
	}

	// Create recommendation object
	now := time.Now()
	rec := &repository.Recommendation{
		ID:                  uuid.New().String(),
		TenantID:            tenantID,
		Type:                recType,
		Title:               title,
		Description:         description,
		ExpectedImprovement: expectedImprovement,
		Status:              "pending",
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Save to repository
	return s.repo.SaveRecommendation(ctx, rec)
}

// GetRecommendations returns recommendations filtered by status
func (s *analyticsService) GetRecommendations(ctx context.Context, tenantID string, status string) ([]repository.Recommendation, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Call repository to get recommendations
	return s.repo.GetRecommendations(ctx, tenantID, status)
}

// UpdateRecommendationStatus updates the status of a recommendation
func (s *analyticsService) UpdateRecommendationStatus(ctx context.Context, tenantID, recID, status string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if recID == "" {
		return errors.New("recommendation ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"applied":   true,
		"dismissed": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status: must be 'pending', 'applied', or 'dismissed'")
	}

	// Call repository to update status
	return s.repo.UpdateRecommendationStatus(ctx, tenantID, recID, status)
}
