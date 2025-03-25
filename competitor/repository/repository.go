package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/common/db"
	"github.com/google/uuid"
)

// CompetitorRepository defines the interface for competitor data access
type CompetitorRepository interface {
	GetCompetitors(ctx context.Context, tenantID string) ([]Competitor, error)
	GetCompetitor(ctx context.Context, tenantID, competitorID string) (*Competitor, error)
	AddCompetitor(ctx context.Context, competitor *Competitor) (*Competitor, error)
	UpdateCompetitor(ctx context.Context, competitor *Competitor) (*Competitor, error)
	DeleteCompetitor(ctx context.Context, tenantID, competitorID string) error
	GetCompetitorMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) ([]CompetitorMetric, error)
	UpdateCompetitorMetrics(ctx context.Context, tenantID, competitorID string, metrics []CompetitorMetric) (int, error)
}

// Competitor represents a competitor entity
type Competitor struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"created_at"`
}

// CompetitorMetric represents engagement metrics for a competitor
type CompetitorMetric struct {
	ID             string    `json:"id"`
	CompetitorID   string    `json:"competitor_id"`
	PostID         string    `json:"post_id"`
	Likes          int       `json:"likes"`
	Shares         int       `json:"shares"`
	Comments       int       `json:"comments"`
	CTR            float64   `json:"ctr"`
	AvgWatchTime   float64   `json:"avg_watch_time"`
	EngagementRate float64   `json:"engagement_rate"`
	PostedAt       time.Time `json:"posted_at"`
}

// SupabaseCompetitorRepository implements CompetitorRepository using Supabase
type SupabaseCompetitorRepository struct {
	client *db.SupabaseClient
}

// NewSupabaseCompetitorRepository creates a new SupabaseCompetitorRepository
func NewSupabaseCompetitorRepository(tenantID string) (*SupabaseCompetitorRepository, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	client, err := db.NewSupabaseClient(tenantID)
	if err != nil {
		return nil, err
	}

	return &SupabaseCompetitorRepository{
		client: client,
	}, nil
}

// GetCompetitors retrieves all competitors for the current tenant
func (r *SupabaseCompetitorRepository) GetCompetitors(ctx context.Context, tenantID string) ([]Competitor, error) {
	var competitors []Competitor
	err := r.client.Query("competitors").Select("*").Execute(&competitors)
	if err != nil {
		return nil, fmt.Errorf("failed to get competitors: %w", err)
	}
	return competitors, nil
}

// GetCompetitor retrieves a specific competitor
func (r *SupabaseCompetitorRepository) GetCompetitor(ctx context.Context, tenantID, competitorID string) (*Competitor, error) {
	var competitors []Competitor
	err := r.client.Query("competitors").
		Select("*").
		Where("id", "eq", competitorID).
		Execute(&competitors)

	if err != nil {
		return nil, fmt.Errorf("failed to get competitor: %w", err)
	}

	if len(competitors) == 0 {
		return nil, errors.New("competitor not found")
	}

	return &competitors[0], nil
}

// AddCompetitor adds a new competitor
func (r *SupabaseCompetitorRepository) AddCompetitor(ctx context.Context, competitor *Competitor) (*Competitor, error) {
	if competitor.ID == "" {
		competitor.ID = uuid.New().String()
	}

	competitor.TenantID = r.client.TenantID
	if competitor.CreatedAt.IsZero() {
		competitor.CreatedAt = time.Now()
	}

	err := r.client.Insert(ctx, "competitors", competitor)
	if err != nil {
		return nil, fmt.Errorf("failed to add competitor: %w", err)
	}

	return competitor, nil
}

// UpdateCompetitor updates an existing competitor
func (r *SupabaseCompetitorRepository) UpdateCompetitor(ctx context.Context, competitor *Competitor) (*Competitor, error) {
	// Make sure the competitor belongs to the tenant
	existing, err := r.GetCompetitor(ctx, r.client.TenantID, competitor.ID)
	if err != nil {
		return nil, err
	}

	// Keep original tenant ID and created date
	competitor.TenantID = existing.TenantID
	competitor.CreatedAt = existing.CreatedAt

	err = r.client.Update(ctx, "competitors", "id", competitor.ID, competitor)
	if err != nil {
		return nil, fmt.Errorf("failed to update competitor: %w", err)
	}

	return competitor, nil
}

// DeleteCompetitor deletes a competitor
func (r *SupabaseCompetitorRepository) DeleteCompetitor(ctx context.Context, tenantID, competitorID string) error {
	// Verify the competitor exists and belongs to the tenant
	_, err := r.GetCompetitor(ctx, tenantID, competitorID)
	if err != nil {
		return err
	}

	err = r.client.Delete(ctx, "competitors", "id", competitorID)
	if err != nil {
		return fmt.Errorf("failed to delete competitor: %w", err)
	}

	return nil
}

// GetCompetitorMetrics retrieves metrics for a specific competitor within a date range
func (r *SupabaseCompetitorRepository) GetCompetitorMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) ([]CompetitorMetric, error) {
	// First verify the competitor exists and belongs to the tenant
	_, err := r.GetCompetitor(ctx, tenantID, competitorID)
	if err != nil {
		return nil, err
	}

	var metrics []CompetitorMetric
	err = r.client.Query("competitor_metrics").
		Select("*").
		Where("competitor_id", "eq", competitorID).
		Where("posted_at", "gte", startDate.Format(time.RFC3339)).
		Where("posted_at", "lte", endDate.Format(time.RFC3339)).
		Order("posted_at", false).
		Execute(&metrics)

	if err != nil {
		return nil, fmt.Errorf("failed to get competitor metrics: %w", err)
	}

	return metrics, nil
}

// UpdateCompetitorMetrics updates metrics for a specific competitor
func (r *SupabaseCompetitorRepository) UpdateCompetitorMetrics(ctx context.Context, tenantID, competitorID string, metrics []CompetitorMetric) (int, error) {
	// Verify the competitor exists and belongs to the tenant
	_, err := r.GetCompetitor(ctx, tenantID, competitorID)
	if err != nil {
		return 0, err
	}

	// Insert each metric
	for i := range metrics {
		// Set IDs and ensure competitor ID is set
		if metrics[i].ID == "" {
			metrics[i].ID = uuid.New().String()
		}
		metrics[i].CompetitorID = competitorID

		err = r.client.Insert(ctx, "competitor_metrics", metrics[i])
		if err != nil {
			return i, fmt.Errorf("failed to update metric at index %d: %w", i, err)
		}
	}

	return len(metrics), nil
}
