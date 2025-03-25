package service

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/competitor/repository"
)

// CompetitorService provides business logic for competitor operations
type CompetitorService struct {
	repo repository.CompetitorRepository
}

// NewCompetitorService creates a new CompetitorService
func NewCompetitorService(repo repository.CompetitorRepository) *CompetitorService {
	return &CompetitorService{
		repo: repo,
	}
}

// GetCompetitors retrieves all competitors for the current tenant
func (s *CompetitorService) GetCompetitors(ctx context.Context, tenantID string) ([]repository.Competitor, error) {
	return s.repo.GetCompetitors(ctx, tenantID)
}

// GetCompetitor retrieves details about a specific competitor
func (s *CompetitorService) GetCompetitor(ctx context.Context, tenantID, competitorID string) (*repository.Competitor, error) {
	return s.repo.GetCompetitor(ctx, tenantID, competitorID)
}

// GetCompetitorMetrics retrieves metrics for a specific competitor
func (s *CompetitorService) GetCompetitorMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) ([]repository.CompetitorMetric, error) {
	return s.repo.GetCompetitorMetrics(ctx, tenantID, competitorID, startDate, endDate)
}

// AddCompetitor adds a new competitor for the current tenant
func (s *CompetitorService) AddCompetitor(ctx context.Context, tenantID, name, platform string) (*repository.Competitor, error) {
	competitor := &repository.Competitor{
		TenantID: tenantID,
		Name:     name,
		Platform: platform,
	}
	return s.repo.AddCompetitor(ctx, competitor)
}

// UpdateCompetitor updates details about an existing competitor
func (s *CompetitorService) UpdateCompetitor(ctx context.Context, tenantID, competitorID, name, platform string) (*repository.Competitor, error) {
	// First retrieve the competitor to ensure it exists and belongs to the tenant
	existingCompetitor, err := s.repo.GetCompetitor(ctx, tenantID, competitorID)
	if err != nil {
		return nil, err
	}

	// Update the fields
	if name != "" {
		existingCompetitor.Name = name
	}
	if platform != "" {
		existingCompetitor.Platform = platform
	}

	return s.repo.UpdateCompetitor(ctx, existingCompetitor)
}

// DeleteCompetitor deletes a competitor from tracking
func (s *CompetitorService) DeleteCompetitor(ctx context.Context, tenantID, competitorID string) error {
	return s.repo.DeleteCompetitor(ctx, tenantID, competitorID)
}

// UpdateCompetitorMetrics updates metrics for a specific competitor
func (s *CompetitorService) UpdateCompetitorMetrics(ctx context.Context, tenantID, competitorID string, metrics []repository.CompetitorMetric) (int, error) {
	// First check if the competitor exists and belongs to the tenant
	_, err := s.repo.GetCompetitor(ctx, tenantID, competitorID)
	if err != nil {
		return 0, err
	}

	// Update the metrics
	return s.repo.UpdateCompetitorMetrics(ctx, tenantID, competitorID, metrics)
}
