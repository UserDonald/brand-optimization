package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/common/db"
	"github.com/google/uuid"
)

// EngagementRepository defines the interface for engagement data access
type EngagementRepository interface {
	// Personal metrics
	GetPersonalMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]PersonalMetric, error)
	AddPersonalMetric(ctx context.Context, metric *PersonalMetric) (*PersonalMetric, error)
	UpdatePersonalMetric(ctx context.Context, metric *PersonalMetric) (*PersonalMetric, error)
	DeletePersonalMetric(ctx context.Context, tenantID, metricID string) error

	// Comparison analytics
	CompareMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) (*ComparisonResult, error)

	// Engagement trends
	GetEngagementTrends(ctx context.Context, tenantID string, period string, startDate, endDate time.Time) ([]EngagementTrend, error)
}

// PersonalMetric represents engagement metrics for the client's own social media
type PersonalMetric struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	PostID         string    `json:"post_id"`
	Likes          int       `json:"likes"`
	Shares         int       `json:"shares"`
	Comments       int       `json:"comments"`
	CTR            float64   `json:"click_through_rate"`
	AvgWatchTime   float64   `json:"avg_watch_time"`
	EngagementRate float64   `json:"engagement_rate"`
	PostedAt       time.Time `json:"posted_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// MetricAggregate represents aggregated metrics for a time period
type MetricAggregate struct {
	TotalLikes        int     `json:"total_likes"`
	TotalShares       int     `json:"total_shares"`
	TotalComments     int     `json:"total_comments"`
	AvgEngagementRate float64 `json:"avg_engagement_rate"`
	AvgWatchTime      float64 `json:"avg_watch_time"`
}

// ComparisonRatio represents the ratio between competitor and personal metrics
type ComparisonRatio struct {
	LikesRatio          float64 `json:"likes_ratio"`
	SharesRatio         float64 `json:"shares_ratio"`
	CommentsRatio       float64 `json:"comments_ratio"`
	EngagementRateRatio float64 `json:"engagement_rate_ratio"`
	WatchTimeRatio      float64 `json:"watch_time_ratio"`
}

// ComparisonResult represents the result of comparing competitor and personal metrics
type ComparisonResult struct {
	CompetitorMetrics MetricAggregate `json:"competitor_metrics"`
	PersonalMetrics   MetricAggregate `json:"personal_metrics"`
	Ratios            ComparisonRatio `json:"ratios"`
}

// EngagementTrend represents a trend in engagement metrics over time
type EngagementTrend struct {
	Date            time.Time `json:"date"`
	EngagementRate  float64   `json:"engagement_rate"`
	Likes           int       `json:"likes"`
	Shares          int       `json:"shares"`
	Comments        int       `json:"comments"`
	ComparisonValue float64   `json:"comparison_value"`
}

// SupabaseEngagementRepository implements EngagementRepository using Supabase
type SupabaseEngagementRepository struct {
	client *db.SupabaseClient
}

// NewSupabaseEngagementRepository creates a new SupabaseEngagementRepository
func NewSupabaseEngagementRepository(tenantID string) (*SupabaseEngagementRepository, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	client, err := db.NewSupabaseClient(tenantID)
	if err != nil {
		return nil, err
	}

	return &SupabaseEngagementRepository{
		client: client,
	}, nil
}

// GetPersonalMetrics retrieves personal metrics for a date range
func (r *SupabaseEngagementRepository) GetPersonalMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]PersonalMetric, error) {
	var metrics []PersonalMetric
	err := r.client.Query("personal_metrics").
		Select("*").
		Where("tenant_id", "eq", tenantID).
		Where("posted_at", "gte", startDate.Format(time.RFC3339)).
		Where("posted_at", "lte", endDate.Format(time.RFC3339)).
		Order("posted_at", false).
		Execute(&metrics)

	if err != nil {
		return nil, fmt.Errorf("failed to get personal metrics: %w", err)
	}

	return metrics, nil
}

// AddPersonalMetric adds a new personal metric
func (r *SupabaseEngagementRepository) AddPersonalMetric(ctx context.Context, metric *PersonalMetric) (*PersonalMetric, error) {
	if metric.ID == "" {
		metric.ID = uuid.New().String()
	}

	metric.TenantID = r.client.TenantID
	if metric.CreatedAt.IsZero() {
		metric.CreatedAt = time.Now()
	}

	err := r.client.Insert(ctx, "personal_metrics", metric)
	if err != nil {
		return nil, fmt.Errorf("failed to add personal metric: %w", err)
	}

	return metric, nil
}

// UpdatePersonalMetric updates an existing personal metric
func (r *SupabaseEngagementRepository) UpdatePersonalMetric(ctx context.Context, metric *PersonalMetric) (*PersonalMetric, error) {
	// Verify the metric exists and belongs to the tenant
	var metrics []PersonalMetric
	err := r.client.Query("personal_metrics").
		Select("*").
		Where("id", "eq", metric.ID).
		Where("tenant_id", "eq", r.client.TenantID).
		Execute(&metrics)

	if err != nil {
		return nil, fmt.Errorf("failed to verify personal metric: %w", err)
	}

	if len(metrics) == 0 {
		return nil, errors.New("personal metric not found")
	}

	// Keep original tenant ID and created date
	existing := metrics[0]
	metric.TenantID = existing.TenantID
	metric.CreatedAt = existing.CreatedAt

	err = r.client.Update(ctx, "personal_metrics", "id", metric.ID, metric)
	if err != nil {
		return nil, fmt.Errorf("failed to update personal metric: %w", err)
	}

	return metric, nil
}

// DeletePersonalMetric deletes a personal metric
func (r *SupabaseEngagementRepository) DeletePersonalMetric(ctx context.Context, tenantID, metricID string) error {
	// Verify the metric exists and belongs to the tenant
	var metrics []PersonalMetric
	err := r.client.Query("personal_metrics").
		Select("*").
		Where("id", "eq", metricID).
		Where("tenant_id", "eq", tenantID).
		Execute(&metrics)

	if err != nil {
		return fmt.Errorf("failed to verify personal metric: %w", err)
	}

	if len(metrics) == 0 {
		return errors.New("personal metric not found")
	}

	err = r.client.Delete(ctx, "personal_metrics", "id", metricID)
	if err != nil {
		return fmt.Errorf("failed to delete personal metric: %w", err)
	}

	return nil
}

// CompareMetrics compares personal metrics with competitor metrics
func (r *SupabaseEngagementRepository) CompareMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) (*ComparisonResult, error) {
	// Get personal metrics
	personalMetrics, err := r.GetPersonalMetrics(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get competitor metrics
	var competitorMetrics []struct {
		Likes          int     `json:"likes"`
		Shares         int     `json:"shares"`
		Comments       int     `json:"comments"`
		EngagementRate float64 `json:"engagement_rate"`
		AvgWatchTime   float64 `json:"avg_watch_time"`
	}

	err = r.client.Query("competitor_metrics").
		Select("likes", "shares", "comments", "engagement_rate", "avg_watch_time").
		Where("tenant_id", "eq", tenantID).
		Where("competitor_id", "eq", competitorID).
		Where("posted_at", "gte", startDate.Format(time.RFC3339)).
		Where("posted_at", "lte", endDate.Format(time.RFC3339)).
		Execute(&competitorMetrics)

	if err != nil {
		return nil, fmt.Errorf("failed to get competitor metrics: %w", err)
	}

	// Calculate aggregates for personal metrics
	personalAggregate := calculatePersonalAggregates(personalMetrics)

	// Calculate aggregates for competitor metrics
	competitorAggregate := calculateCompetitorAggregates(competitorMetrics)

	// Calculate ratios
	ratios := calculateRatios(competitorAggregate, personalAggregate)

	return &ComparisonResult{
		CompetitorMetrics: competitorAggregate,
		PersonalMetrics:   personalAggregate,
		Ratios:            ratios,
	}, nil
}

// calculatePersonalAggregates calculates aggregates for personal metrics
func calculatePersonalAggregates(metrics []PersonalMetric) MetricAggregate {
	var totalLikes, totalShares, totalComments int
	var sumEngagementRate, sumWatchTime float64

	for _, m := range metrics {
		totalLikes += m.Likes
		totalShares += m.Shares
		totalComments += m.Comments
		sumEngagementRate += m.EngagementRate
		sumWatchTime += m.AvgWatchTime
	}

	count := float64(len(metrics))
	var avgEngagementRate, avgWatchTime float64
	if count > 0 {
		avgEngagementRate = sumEngagementRate / count
		avgWatchTime = sumWatchTime / count
	}

	return MetricAggregate{
		TotalLikes:        totalLikes,
		TotalShares:       totalShares,
		TotalComments:     totalComments,
		AvgEngagementRate: avgEngagementRate,
		AvgWatchTime:      avgWatchTime,
	}
}

// calculateCompetitorAggregates calculates aggregates for competitor metrics
func calculateCompetitorAggregates(metrics []struct {
	Likes          int     `json:"likes"`
	Shares         int     `json:"shares"`
	Comments       int     `json:"comments"`
	EngagementRate float64 `json:"engagement_rate"`
	AvgWatchTime   float64 `json:"avg_watch_time"`
}) MetricAggregate {
	var totalLikes, totalShares, totalComments int
	var sumEngagementRate, sumWatchTime float64

	for _, m := range metrics {
		totalLikes += m.Likes
		totalShares += m.Shares
		totalComments += m.Comments
		sumEngagementRate += m.EngagementRate
		sumWatchTime += m.AvgWatchTime
	}

	count := float64(len(metrics))
	var avgEngagementRate, avgWatchTime float64
	if count > 0 {
		avgEngagementRate = sumEngagementRate / count
		avgWatchTime = sumWatchTime / count
	}

	return MetricAggregate{
		TotalLikes:        totalLikes,
		TotalShares:       totalShares,
		TotalComments:     totalComments,
		AvgEngagementRate: avgEngagementRate,
		AvgWatchTime:      avgWatchTime,
	}
}

// calculateRatios calculates ratios between competitor and personal metrics
func calculateRatios(competitor, personal MetricAggregate) ComparisonRatio {
	var likesRatio, sharesRatio, commentsRatio, engagementRateRatio, watchTimeRatio float64

	if personal.TotalLikes > 0 {
		likesRatio = float64(competitor.TotalLikes) / float64(personal.TotalLikes)
	}

	if personal.TotalShares > 0 {
		sharesRatio = float64(competitor.TotalShares) / float64(personal.TotalShares)
	}

	if personal.TotalComments > 0 {
		commentsRatio = float64(competitor.TotalComments) / float64(personal.TotalComments)
	}

	if personal.AvgEngagementRate > 0 {
		engagementRateRatio = competitor.AvgEngagementRate / personal.AvgEngagementRate
	}

	if personal.AvgWatchTime > 0 {
		watchTimeRatio = competitor.AvgWatchTime / personal.AvgWatchTime
	}

	return ComparisonRatio{
		LikesRatio:          likesRatio,
		SharesRatio:         sharesRatio,
		CommentsRatio:       commentsRatio,
		EngagementRateRatio: engagementRateRatio,
		WatchTimeRatio:      watchTimeRatio,
	}
}

// GetEngagementTrends retrieves engagement trends over time
func (r *SupabaseEngagementRepository) GetEngagementTrends(ctx context.Context, tenantID, period string, startDate, endDate time.Time) ([]EngagementTrend, error) {
	// Get all metrics for the period
	personalMetrics, err := r.GetPersonalMetrics(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Group by the requested period and calculate trends
	dateGroups := make(map[string][]PersonalMetric)
	for _, metric := range personalMetrics {
		dateKey := getDateKey(metric.PostedAt, period)
		dateGroups[dateKey] = append(dateGroups[dateKey], metric)
	}

	// Calculate trend for each group
	var trends []EngagementTrend
	for dateKey, metrics := range dateGroups {
		date, err := parseDataKey(dateKey, period)
		if err != nil {
			continue
		}

		// Calculate aggregate values for this period
		var totalLikes, totalShares, totalComments int
		var sumEngagementRate float64

		for _, m := range metrics {
			totalLikes += m.Likes
			totalShares += m.Shares
			totalComments += m.Comments
			sumEngagementRate += m.EngagementRate
		}

		count := float64(len(metrics))
		avgEngagementRate := 0.0
		if count > 0 {
			avgEngagementRate = sumEngagementRate / count
		}

		trend := EngagementTrend{
			Date:           date,
			EngagementRate: avgEngagementRate,
			Likes:          totalLikes,
			Shares:         totalShares,
			Comments:       totalComments,
			// ComparisonValue would be used to show week-over-week or month-over-month changes
			ComparisonValue: 0.0, // We would calculate this in a real implementation
		}

		trends = append(trends, trend)
	}

	return trends, nil
}

// getDateKey returns a string key for grouping metrics by time period
func getDateKey(date time.Time, period string) string {
	switch period {
	case "day":
		return date.Format("2006-01-02")
	case "week":
		year, week := date.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case "month":
		return date.Format("2006-01")
	default:
		return date.Format("2006-01-02")
	}
}

// parseDataKey converts a date key back to a time.Time
func parseDataKey(dateKey, period string) (time.Time, error) {
	switch period {
	case "day":
		return time.Parse("2006-01-02", dateKey)
	case "week":
		// Parse year-week format (e.g., "2023-W01")
		var year, week int
		_, err := fmt.Sscanf(dateKey, "%d-W%02d", &year, &week)
		if err != nil {
			return time.Time{}, err
		}
		// Calculate the first day of the week
		// This is a simplification - proper ISO week calculation is more complex
		jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		firstDay := jan1.AddDate(0, 0, (week-1)*7)
		return firstDay, nil
	case "month":
		return time.Parse("2006-01", dateKey)
	default:
		return time.Parse("2006-01-02", dateKey)
	}
}
