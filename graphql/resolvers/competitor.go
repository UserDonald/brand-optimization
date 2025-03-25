package resolvers

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/competitor/client"
	"github.com/donaldnash/go-competitor/graphql/models"
)

// CompetitorResolver handles all competitor-related GraphQL queries and mutations
type CompetitorResolver struct {
	client client.CompetitorClient
}

// NewCompetitorResolver creates a new CompetitorResolver
func NewCompetitorResolver(client client.CompetitorClient) *CompetitorResolver {
	return &CompetitorResolver{
		client: client,
	}
}

// GetCompetitors retrieves all competitors for the current tenant
func (r *CompetitorResolver) GetCompetitors(ctx context.Context, tenantID string) ([]*models.Competitor, error) {
	competitors, err := r.client.GetCompetitors(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var result []*models.Competitor
	for _, c := range competitors {
		result = append(result, &models.Competitor{
			ID:        c.ID,
			TenantID:  c.TenantID,
			Name:      c.Name,
			Platform:  c.Platform,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

// GetCompetitor retrieves a specific competitor
func (r *CompetitorResolver) GetCompetitor(ctx context.Context, tenantID, competitorID string) (*models.Competitor, error) {
	competitor, err := r.client.GetCompetitor(ctx, tenantID, competitorID)
	if err != nil {
		return nil, err
	}

	return &models.Competitor{
		ID:        competitor.ID,
		TenantID:  competitor.TenantID,
		Name:      competitor.Name,
		Platform:  competitor.Platform,
		CreatedAt: competitor.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetCompetitorMetrics retrieves metrics for a specific competitor
func (r *CompetitorResolver) GetCompetitorMetrics(ctx context.Context, tenantID, competitorID string, dateRange *models.DateRange) ([]*models.CompetitorMetric, error) {
	startDate, err := time.Parse(time.RFC3339, dateRange.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(time.RFC3339, dateRange.EndDate)
	if err != nil {
		return nil, err
	}

	metrics, err := r.client.GetCompetitorMetrics(ctx, tenantID, competitorID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var result []*models.CompetitorMetric
	for _, m := range metrics {
		result = append(result, &models.CompetitorMetric{
			ID:               m.ID,
			CompetitorID:     m.CompetitorID,
			PostID:           m.PostID,
			Likes:            m.Likes,
			Shares:           m.Shares,
			Comments:         m.Comments,
			ClickThroughRate: m.CTR,
			AvgWatchTime:     m.AvgWatchTime,
			EngagementRate:   m.EngagementRate,
			PostedAt:         m.PostedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

// AddCompetitor adds a new competitor
func (r *CompetitorResolver) AddCompetitor(ctx context.Context, input *models.AddCompetitorInput) (*models.Competitor, error) {
	competitor, err := r.client.AddCompetitor(ctx, input.TenantID, input.Name, input.Platform)
	if err != nil {
		return nil, err
	}

	return &models.Competitor{
		ID:        competitor.ID,
		TenantID:  competitor.TenantID,
		Name:      competitor.Name,
		Platform:  competitor.Platform,
		CreatedAt: competitor.CreatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateCompetitor updates an existing competitor
func (r *CompetitorResolver) UpdateCompetitor(ctx context.Context, tenantID, competitorID string, input *models.UpdateCompetitorInput) (*models.Competitor, error) {
	competitor, err := r.client.UpdateCompetitor(ctx, tenantID, competitorID, input.Name, input.Platform)
	if err != nil {
		return nil, err
	}

	return &models.Competitor{
		ID:        competitor.ID,
		TenantID:  competitor.TenantID,
		Name:      competitor.Name,
		Platform:  competitor.Platform,
		CreatedAt: competitor.CreatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteCompetitor deletes a competitor
func (r *CompetitorResolver) DeleteCompetitor(ctx context.Context, tenantID, competitorID string) (bool, error) {
	err := r.client.DeleteCompetitor(ctx, tenantID, competitorID)
	return err == nil, err
}

// CompareMetrics compares metrics between a competitor and the client's own brand
func (r *CompetitorResolver) CompareMetrics(ctx context.Context, tenantID, competitorID string, dateRange *models.DateRange) (*models.ComparisonResult, error) {
	// This would call the competitor service and potentially other services to get the comparison data
	// For simplicity, we'll just return a placeholder
	return &models.ComparisonResult{
		Competitor: &models.CompetitorComparison{
			Metrics: []*models.CompetitorMetric{},
			Aggregates: &models.MetricAggregates{
				TotalLikes:        0,
				TotalShares:       0,
				TotalComments:     0,
				AvgEngagementRate: 0,
				AvgWatchTime:      0,
			},
		},
		Personal: &models.PersonalComparison{
			Metrics: []*models.PersonalMetric{},
			Aggregates: &models.MetricAggregates{
				TotalLikes:        0,
				TotalShares:       0,
				TotalComments:     0,
				AvgEngagementRate: 0,
				AvgWatchTime:      0,
			},
		},
		Ratios: &models.ComparisonRatios{
			LikesRatio:          0,
			SharesRatio:         0,
			CommentsRatio:       0,
			EngagementRateRatio: 0,
			WatchTimeRatio:      0,
		},
	}, nil
}
