package resolvers

import (
	"context"
	"time"

	analyticsClient "github.com/donaldnash/go-competitor/analytics/client"
	"github.com/donaldnash/go-competitor/analytics/pb"
)

// AnalyticsResolver handles analytics-related GraphQL queries
type AnalyticsResolver struct {
	client *analyticsClient.AnalyticsClient
}

// NewAnalyticsResolver creates a new analytics resolver
func NewAnalyticsResolver(client *analyticsClient.AnalyticsClient) *AnalyticsResolver {
	return &AnalyticsResolver{
		client: client,
	}
}

// PostingTimeRecommendation represents a recommended posting time
type PostingTimeRecommendation struct {
	DayOfWeek               string
	TimeOfDay               string
	PredictedEngagementRate float64
	Confidence              float64
}

// ContentFormatRecommendation represents a recommended content format
type ContentFormatRecommendation struct {
	Format                  string
	PredictedEngagementRate float64
	TargetAudience          string
	Confidence              float64
}

// GetRecommendedPostingTimes returns recommended posting times based on historical data
func (r *AnalyticsResolver) GetRecommendedPostingTimes(ctx context.Context, tenantID, dayOfWeek string) ([]PostingTimeRecommendation, error) {
	// Call analytics service
	response, err := r.client.GetPostingTimeRecommendations(ctx, tenantID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL type
	recommendations := make([]PostingTimeRecommendation, 0, len(response.Recommendations))
	for _, rec := range response.Recommendations {
		// Format the time of day as "HH:00" (e.g., "14:00")
		timeOfDay := formatHourToTimeString(int(rec.HourOfDay))

		recommendations = append(recommendations, PostingTimeRecommendation{
			DayOfWeek:               rec.DayOfWeek,
			TimeOfDay:               timeOfDay,
			PredictedEngagementRate: rec.PredictedEngagementRate,
			Confidence:              rec.Confidence,
		})
	}

	return recommendations, nil
}

// GetRecommendedContentFormats returns recommended content formats based on historical data
func (r *AnalyticsResolver) GetRecommendedContentFormats(ctx context.Context, tenantID string) ([]ContentFormatRecommendation, error) {
	// Call analytics service
	response, err := r.client.GetContentFormatRecommendations(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL type
	recommendations := make([]ContentFormatRecommendation, 0, len(response.Recommendations))
	for _, rec := range response.Recommendations {
		recommendations = append(recommendations, ContentFormatRecommendation{
			Format:                  rec.Format,
			PredictedEngagementRate: rec.PredictedEngagementRate,
			TargetAudience:          rec.TargetAudience,
			Confidence:              rec.Confidence,
		})
	}

	return recommendations, nil
}

// PredictPostEngagement predicts engagement metrics for a potential post
func (r *AnalyticsResolver) PredictPostEngagement(ctx context.Context, tenantID, contentFormat, scheduledTime string) (*pb.EngagementPrediction, error) {
	// Parse scheduled time
	postTime, err := time.Parse(time.RFC3339, scheduledTime)
	if err != nil {
		return nil, err
	}

	// Call analytics service
	return r.client.PredictEngagement(ctx, tenantID, postTime, contentFormat)
}

// GetContentPerformanceAnalysis returns performance analysis for different content types
func (r *AnalyticsResolver) GetContentPerformanceAnalysis(ctx context.Context, tenantID, startDate, endDate string) (*pb.ContentPerformanceResponse, error) {
	// Parse date range
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, err
	}

	// Call analytics service
	return r.client.AnalyzeContentPerformance(ctx, tenantID, start, end)
}

// Helper to format hour to time string
func formatHourToTimeString(hour int) string {
	return time.Date(0, 0, 0, hour, 0, 0, 0, time.UTC).Format("15:04")
}
