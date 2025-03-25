package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/donaldnash/go-competitor/engagement/repository"
)

// EngagementService provides business logic for engagement operations
type EngagementService struct {
	repo repository.EngagementRepository
}

// NewEngagementService creates a new EngagementService
func NewEngagementService(repo repository.EngagementRepository) *EngagementService {
	return &EngagementService{
		repo: repo,
	}
}

// GetPersonalMetrics retrieves personal metrics for a date range
func (s *EngagementService) GetPersonalMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]repository.PersonalMetric, error) {
	return s.repo.GetPersonalMetrics(ctx, tenantID, startDate, endDate)
}

// AddPersonalMetric adds a new personal metric
func (s *EngagementService) AddPersonalMetric(ctx context.Context, tenantID string, metric *repository.PersonalMetric) (*repository.PersonalMetric, error) {
	// Ensure tenant ID is set correctly
	metric.TenantID = tenantID
	return s.repo.AddPersonalMetric(ctx, metric)
}

// UpdatePersonalMetric updates an existing personal metric
func (s *EngagementService) UpdatePersonalMetric(ctx context.Context, tenantID string, metric *repository.PersonalMetric) (*repository.PersonalMetric, error) {
	// Ensure tenant ID is set correctly
	metric.TenantID = tenantID
	return s.repo.UpdatePersonalMetric(ctx, metric)
}

// DeletePersonalMetric deletes a personal metric
func (s *EngagementService) DeletePersonalMetric(ctx context.Context, tenantID, metricID string) error {
	return s.repo.DeletePersonalMetric(ctx, tenantID, metricID)
}

// CompareMetrics compares personal metrics with competitor metrics
func (s *EngagementService) CompareMetrics(ctx context.Context, tenantID, competitorID string, startDate, endDate time.Time) (*repository.ComparisonResult, error) {
	return s.repo.CompareMetrics(ctx, tenantID, competitorID, startDate, endDate)
}

// GetEngagementTrends retrieves engagement trends over time
func (s *EngagementService) GetEngagementTrends(ctx context.Context, tenantID, period string, startDate, endDate time.Time) ([]repository.EngagementTrend, error) {
	// Validate the period
	switch period {
	case "day", "week", "month":
		// Valid period
	default:
		// Default to day if period is invalid
		period = "day"
	}

	return s.repo.GetEngagementTrends(ctx, tenantID, period, startDate, endDate)
}

// GetEngagementInsights generates insights from engagement data
func (s *EngagementService) GetEngagementInsights(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]EngagementInsight, error) {
	// Get trends
	dayTrends, err := s.GetEngagementTrends(ctx, tenantID, "day", startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get weekly trends
	weekTrends, err := s.GetEngagementTrends(ctx, tenantID, "week", startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get personal metrics
	metrics, err := s.GetPersonalMetrics(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Generate insights
	insights := generateInsightsFromData(dayTrends, weekTrends, metrics)

	return insights, nil
}

// EngagementInsight represents an insight derived from engagement data
type EngagementInsight struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	InsightType string            `json:"insight_type"` // "trend", "anomaly", "recommendation"
	Confidence  float64           `json:"confidence"`
	Metadata    map[string]string `json:"metadata"`
}

// generateInsightsFromData analyzes trends and metrics to generate insights
func generateInsightsFromData(dayTrends []repository.EngagementTrend, weekTrends []repository.EngagementTrend, metrics []repository.PersonalMetric) []EngagementInsight {
	insights := []EngagementInsight{}

	// Example: Detect day of week pattern
	if len(dayTrends) >= 7 {
		// Group by day of week
		dayOfWeekEngagement := make(map[string]float64)
		dayOfWeekCount := make(map[string]int)

		for _, trend := range dayTrends {
			dayOfWeek := trend.Date.Weekday().String()
			dayOfWeekEngagement[dayOfWeek] += trend.EngagementRate
			dayOfWeekCount[dayOfWeek]++
		}

		// Find best and worst days
		var bestDay, worstDay string
		var bestRate, worstRate float64

		for day, total := range dayOfWeekEngagement {
			count := dayOfWeekCount[day]
			if count == 0 {
				continue
			}

			avgRate := total / float64(count)

			if bestDay == "" || avgRate > bestRate {
				bestDay = day
				bestRate = avgRate
			}

			if worstDay == "" || avgRate < worstRate {
				worstDay = day
				worstRate = avgRate
			}
		}

		// Only create insight if there's a significant difference
		if bestRate > worstRate*1.2 { // At least 20% better
			insights = append(insights, EngagementInsight{
				Title:       "Day of Week Pattern Detected",
				Description: bestDay + " posts receive higher engagement than " + worstDay + " posts.",
				InsightType: "trend",
				Confidence:  0.85,
				Metadata: map[string]string{
					"best_day":       bestDay,
					"best_day_rate":  formatFloat(bestRate),
					"worst_day":      worstDay,
					"worst_day_rate": formatFloat(worstRate),
					"percent_higher": formatFloat((bestRate - worstRate) / worstRate * 100),
				},
			})
		}
	}

	// Example: Detect growth or decline trend
	if len(weekTrends) >= 2 {
		firstWeek := weekTrends[len(weekTrends)-1]
		lastWeek := weekTrends[0]

		changePercent := (lastWeek.EngagementRate - firstWeek.EngagementRate) / firstWeek.EngagementRate * 100

		if changePercent >= 10 { // At least 10% growth
			insights = append(insights, EngagementInsight{
				Title:       "Engagement Growth Trend",
				Description: "Your engagement rate has increased by " + formatFloat(changePercent) + "% over the analyzed period.",
				InsightType: "trend",
				Confidence:  0.9,
				Metadata: map[string]string{
					"start_date":     firstWeek.Date.Format("2006-01-02"),
					"end_date":       lastWeek.Date.Format("2006-01-02"),
					"start_rate":     formatFloat(firstWeek.EngagementRate),
					"end_rate":       formatFloat(lastWeek.EngagementRate),
					"percent_change": formatFloat(changePercent),
				},
			})
		} else if changePercent <= -10 { // At least 10% decline
			insights = append(insights, EngagementInsight{
				Title:       "Engagement Decline Alert",
				Description: "Your engagement rate has decreased by " + formatFloat(-changePercent) + "% over the analyzed period.",
				InsightType: "anomaly",
				Confidence:  0.9,
				Metadata: map[string]string{
					"start_date":     firstWeek.Date.Format("2006-01-02"),
					"end_date":       lastWeek.Date.Format("2006-01-02"),
					"start_rate":     formatFloat(firstWeek.EngagementRate),
					"end_rate":       formatFloat(lastWeek.EngagementRate),
					"percent_change": formatFloat(changePercent),
				},
			})
		}
	}

	// Example: Content volume recommendation
	if len(metrics) > 0 {
		daysInPeriod := endDate(metrics).Sub(startDate(metrics)).Hours() / 24
		postsPerDay := float64(len(metrics)) / daysInPeriod

		if postsPerDay < 0.5 { // Less than 1 post every 2 days
			insights = append(insights, EngagementInsight{
				Title:       "Posting Frequency Recommendation",
				Description: "Your posting frequency is lower than recommended. Consider increasing to at least 1 post per day.",
				InsightType: "recommendation",
				Confidence:  0.8,
				Metadata: map[string]string{
					"current_posts_per_day": formatFloat(postsPerDay),
					"recommended_minimum":   "1",
					"days_analyzed":         formatFloat(daysInPeriod),
					"total_posts":           formatInt(len(metrics)),
				},
			})
		}
	}

	return insights
}

// Helper functions
func startDate(metrics []repository.PersonalMetric) time.Time {
	if len(metrics) == 0 {
		return time.Now()
	}

	earliest := metrics[0].PostedAt
	for _, m := range metrics {
		if m.PostedAt.Before(earliest) {
			earliest = m.PostedAt
		}
	}
	return earliest
}

func endDate(metrics []repository.PersonalMetric) time.Time {
	if len(metrics) == 0 {
		return time.Now()
	}

	latest := metrics[0].PostedAt
	for _, m := range metrics {
		if m.PostedAt.After(latest) {
			latest = m.PostedAt
		}
	}
	return latest
}

func formatFloat(value float64) string {
	return fmt.Sprintf("%.1f", value)
}

func formatInt(value int) string {
	return strconv.Itoa(value)
}
