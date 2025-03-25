package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/donaldnash/go-competitor/common/db"
	"github.com/donaldnash/go-competitor/engagement/repository"
	"github.com/google/uuid"
)

// AnalyticsRepository defines the interface for analytics data access
type AnalyticsRepository interface {
	// Predictive analytics
	GetPostingTimeRecommendations(ctx context.Context, tenantID string, dayOfWeek string) ([]PostingTimeRecommendation, error)
	GetContentFormatRecommendations(ctx context.Context, tenantID string) ([]ContentFormatRecommendation, error)

	// Performance predictions
	PredictEngagement(ctx context.Context, tenantID string, postTime time.Time, contentFormat string) (*EngagementPrediction, error)

	// Content analysis
	AnalyzeContentPerformance(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]ContentPerformance, error)

	// Performance tracking
	SaveRecommendation(ctx context.Context, rec *Recommendation) (*Recommendation, error)
	GetRecommendations(ctx context.Context, tenantID string, status string) ([]Recommendation, error)
	UpdateRecommendationStatus(ctx context.Context, tenantID, recID, status string) error
}

// PostingTimeRecommendation represents a recommended time to post content
type PostingTimeRecommendation struct {
	ID                      string    `json:"id"`
	TenantID                string    `json:"tenant_id"`
	DayOfWeek               string    `json:"day_of_week"`
	HourOfDay               int       `json:"hour_of_day"`
	PredictedEngagementRate float64   `json:"predicted_engagement_rate"`
	Confidence              float64   `json:"confidence"`
	CreatedAt               time.Time `json:"created_at"`
}

// ContentFormatRecommendation represents a recommended content format
type ContentFormatRecommendation struct {
	ID                      string    `json:"id"`
	TenantID                string    `json:"tenant_id"`
	Format                  string    `json:"format"`
	TargetAudience          string    `json:"target_audience"`
	PredictedEngagementRate float64   `json:"predicted_engagement_rate"`
	Confidence              float64   `json:"confidence"`
	CreatedAt               time.Time `json:"created_at"`
}

// EngagementPrediction represents a prediction of engagement for a specific post
type EngagementPrediction struct {
	ID                string    `json:"id"`
	TenantID          string    `json:"tenant_id"`
	PostTime          time.Time `json:"post_time"`
	ContentFormat     string    `json:"content_format"`
	PredictedLikes    int       `json:"predicted_likes"`
	PredictedShares   int       `json:"predicted_shares"`
	PredictedComments int       `json:"predicted_comments"`
	EngagementRate    float64   `json:"engagement_rate"`
	Confidence        float64   `json:"confidence"`
	CreatedAt         time.Time `json:"created_at"`
}

// ContentPerformance represents analytics about content performance
type ContentPerformance struct {
	Format            string  `json:"format"`
	TotalPosts        int     `json:"total_posts"`
	AvgEngagementRate float64 `json:"avg_engagement_rate"`
	AvgLikes          float64 `json:"avg_likes"`
	AvgShares         float64 `json:"avg_shares"`
	AvgComments       float64 `json:"avg_comments"`
	PerformanceScore  float64 `json:"performance_score"`
	PerformanceTrend  float64 `json:"performance_trend"` // positive = improving, negative = declining
}

// Recommendation represents an actionable recommendation
type Recommendation struct {
	ID                  string    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	Type                string    `json:"type"` // "posting_time", "content_format", etc.
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	ExpectedImprovement float64   `json:"expected_improvement"`
	Status              string    `json:"status"` // "pending", "applied", "dismissed"
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// SupabaseAnalyticsRepository implements AnalyticsRepository using Supabase
type SupabaseAnalyticsRepository struct {
	client *db.SupabaseClient
	model  *MLModel
}

// NewSupabaseAnalyticsRepository creates a new SupabaseAnalyticsRepository
func NewSupabaseAnalyticsRepository(tenantID string) (*SupabaseAnalyticsRepository, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	client, err := db.NewSupabaseClient(tenantID)
	if err != nil {
		return nil, err
	}

	repo := &SupabaseAnalyticsRepository{
		client: client,
		model:  NewMLModel(),
	}

	// Asynchronously initialize ML model in the background
	go func() {
		ctx := context.Background()
		metrics, err := repo.getHistoricalData(ctx, tenantID)
		if err != nil {
			log.Printf("Error loading historical data for ML model: %v", err)
			return
		}

		// Get content formats (in a real implementation, this would come from the content service)
		formats := repo.getContentFormats(ctx, tenantID)

		// Train models
		repo.model.TrainTimeModel(metrics)
		repo.model.TrainFormatModel(metrics, formats)

		log.Printf("ML model training completed with %d data points", len(metrics))
	}()

	return repo, nil
}

// getHistoricalData retrieves historical engagement data to use for predictions
func (r *SupabaseAnalyticsRepository) getHistoricalData(ctx context.Context, tenantID string) ([]repository.PersonalMetric, error) {
	// Get personal metrics for the last 90 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, -3, 0) // 3 months ago

	var metrics []repository.PersonalMetric
	err := r.client.Query("personal_metrics").
		Select("*").
		Where("tenant_id", "eq", tenantID).
		Where("posted_at", "gte", startDate.Format(time.RFC3339)).
		Where("posted_at", "lte", endDate.Format(time.RFC3339)).
		Execute(&metrics)

	if err != nil {
		return nil, fmt.Errorf("failed to get historical metrics: %w", err)
	}

	return metrics, nil
}

// getContentFormats gets a mapping of content formats to post IDs
// In a real implementation, this would fetch from the content service
func (r *SupabaseAnalyticsRepository) getContentFormats(ctx context.Context, tenantID string) map[string][]string {
	// This is a stub implementation for demonstration
	// In a real application, we would fetch this data from another service or database
	formats := map[string][]string{
		"video":       {"post1", "post4", "post7"},
		"image":       {"post2", "post5", "post8"},
		"text":        {"post3", "post6", "post9"},
		"infographic": {"post10", "post11"},
		"poll":        {"post12", "post13"},
	}

	return formats
}

// GetPostingTimeRecommendations recommends optimal posting times
func (r *SupabaseAnalyticsRepository) GetPostingTimeRecommendations(ctx context.Context, tenantID string, dayOfWeek string) ([]PostingTimeRecommendation, error) {
	// Get historical data
	historicalData, err := r.getHistoricalData(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Group metrics by day of week and hour of day
	timeStats := make(map[string]map[int][]repository.PersonalMetric)
	for _, metric := range historicalData {
		day := metric.PostedAt.Weekday().String()
		hour := metric.PostedAt.Hour()

		if timeStats[day] == nil {
			timeStats[day] = make(map[int][]repository.PersonalMetric)
		}

		timeStats[day][hour] = append(timeStats[day][hour], metric)
	}

	// If a specific day is requested, filter results
	var days []string
	if dayOfWeek != "" {
		days = []string{dayOfWeek}
	} else {
		days = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	}

	// Calculate recommendations
	var recommendations []PostingTimeRecommendation

	// For each day, find the best posting hours
	for _, day := range days {
		hourlyStats := timeStats[day]
		if len(hourlyStats) == 0 {
			continue
		}

		// Calculate engagement rate for each hour
		type hourlyEngagement struct {
			Hour           int
			EngagementRate float64
			PostCount      int
		}

		var hourlyEngagements []hourlyEngagement

		for hour, metrics := range hourlyStats {
			if len(metrics) == 0 {
				continue
			}

			// Calculate average engagement rate
			var totalEngagementRate float64
			for _, metric := range metrics {
				totalEngagementRate += metric.EngagementRate
			}
			avgEngagementRate := totalEngagementRate / float64(len(metrics))

			hourlyEngagements = append(hourlyEngagements, hourlyEngagement{
				Hour:           hour,
				EngagementRate: avgEngagementRate,
				PostCount:      len(metrics),
			})
		}

		// Sort by engagement rate and pick top 3
		// In a real implementation, we would use a proper sorting algorithm
		// For simplicity, we'll just iterate and find the top 3
		type topHour struct {
			Hour           int
			EngagementRate float64
			Confidence     float64
		}

		var topHours []topHour

		for _, he := range hourlyEngagements {
			// Simple confidence calculation based on post count
			// In a real implementation, we would use more sophisticated statistics
			confidence := 0.5
			if he.PostCount >= 5 {
				confidence = 0.7
			}
			if he.PostCount >= 10 {
				confidence = 0.9
			}

			// Keep the top 3 hours
			if len(topHours) < 3 {
				topHours = append(topHours, topHour{
					Hour:           he.Hour,
					EngagementRate: he.EngagementRate,
					Confidence:     confidence,
				})
			} else {
				// Replace the lowest engagement rate if this one is higher
				minIndex := 0
				minRate := topHours[0].EngagementRate

				for i, th := range topHours {
					if th.EngagementRate < minRate {
						minIndex = i
						minRate = th.EngagementRate
					}
				}

				if he.EngagementRate > minRate {
					topHours[minIndex] = topHour{
						Hour:           he.Hour,
						EngagementRate: he.EngagementRate,
						Confidence:     confidence,
					}
				}
			}
		}

		// Convert to recommendations
		for _, th := range topHours {
			recommendations = append(recommendations, PostingTimeRecommendation{
				ID:                      uuid.New().String(),
				TenantID:                tenantID,
				DayOfWeek:               day,
				HourOfDay:               th.Hour,
				PredictedEngagementRate: th.EngagementRate,
				Confidence:              th.Confidence,
				CreatedAt:               time.Now(),
			})
		}
	}

	return recommendations, nil
}

// GetContentFormatRecommendations recommends optimal content formats
func (r *SupabaseAnalyticsRepository) GetContentFormatRecommendations(ctx context.Context, tenantID string) ([]ContentFormatRecommendation, error) {
	// In a real implementation, we would analyze historical content performance
	// For this demo, we'll create some sample recommendations

	// Get content formats
	var contentFormats []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := r.client.Query("content_formats").
		Select("id", "name", "description").
		Where("tenant_id", "eq", tenantID).
		Execute(&contentFormats)

	if err != nil {
		return nil, fmt.Errorf("failed to get content formats: %w", err)
	}

	// Get audience segments
	var audienceSegments []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err = r.client.Query("audience_segments").
		Select("id", "name", "description").
		Where("tenant_id", "eq", tenantID).
		Execute(&audienceSegments)

	if err != nil {
		return nil, fmt.Errorf("failed to get audience segments: %w", err)
	}

	// For each content format, generate a recommendation
	var recommendations []ContentFormatRecommendation

	// If no content formats or audience segments, return empty
	if len(contentFormats) == 0 || len(audienceSegments) == 0 {
		return recommendations, nil
	}

	for _, format := range contentFormats {
		// In a real implementation, we would analyze which segment engages most with this format
		// For this demo, we'll just randomly assign a segment
		segmentIndex := 0 // In a real implementation, use proper assignment
		segment := audienceSegments[segmentIndex]

		recommendations = append(recommendations, ContentFormatRecommendation{
			ID:                      uuid.New().String(),
			TenantID:                tenantID,
			Format:                  format.Name,
			TargetAudience:          segment.Name,
			PredictedEngagementRate: 0.08, // Sample value, would be calculated from historical data
			Confidence:              0.7,  // Sample value, would be calculated from historical data
			CreatedAt:               time.Now(),
		})
	}

	return recommendations, nil
}

// PredictEngagement predicts engagement metrics for a potential post
func (r *SupabaseAnalyticsRepository) PredictEngagement(ctx context.Context, tenantID string, postTime time.Time, contentFormat string) (*EngagementPrediction, error) {
	// Get engagement rate predictions
	timeEngagement, timeConfidence := r.model.PredictEngagementForTime(postTime)
	formatEngagement, formatConfidence := r.model.PredictEngagementForFormat(contentFormat)

	// Combine predictions (weighted average based on confidence)
	totalConfidence := timeConfidence + formatConfidence
	combinedEngagement := ((timeEngagement * timeConfidence) + (formatEngagement * formatConfidence)) / totalConfidence

	// Scale confidence based on both factors
	confidence := (timeConfidence + formatConfidence) / 2.0

	// Calculate predicted metrics
	// This is a simplified approach; in a real system, we would use more sophisticated models
	averageFollowers := 5000.0 // Assume average follower count
	impressionRate := 0.3      // Assume 30% of followers see a post

	// Estimate metrics
	estimatedImpressions := averageFollowers * impressionRate
	predictedLikes := int(estimatedImpressions * combinedEngagement * 0.6)    // Assume 60% of engaged users like
	predictedShares := int(estimatedImpressions * combinedEngagement * 0.1)   // Assume 10% of engaged users share
	predictedComments := int(estimatedImpressions * combinedEngagement * 0.2) // Assume 20% of engaged users comment

	// Create prediction
	prediction := &EngagementPrediction{
		ID:                uuid.New().String(),
		TenantID:          tenantID,
		PostTime:          postTime,
		ContentFormat:     contentFormat,
		PredictedLikes:    predictedLikes,
		PredictedShares:   predictedShares,
		PredictedComments: predictedComments,
		EngagementRate:    combinedEngagement,
		Confidence:        confidence,
		CreatedAt:         time.Now(),
	}

	return prediction, nil
}

// AnalyzeContentPerformance calculates content performance
func (r *SupabaseAnalyticsRepository) AnalyzeContentPerformance(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]ContentPerformance, error) {
	// In a real implementation, we would analyze content performance across different formats
	// For this demo, we'll simulate content performance data

	// Get personal metrics in the date range
	var metrics []repository.PersonalMetric
	err := r.client.Query("personal_metrics").
		Select("*").
		Where("tenant_id", "eq", tenantID).
		Where("posted_at", "gte", startDate.Format(time.RFC3339)).
		Where("posted_at", "lte", endDate.Format(time.RFC3339)).
		Execute(&metrics)

	if err != nil {
		return nil, fmt.Errorf("failed to get metrics for content analysis: %w", err)
	}

	// Get content posts with formats
	var posts []struct {
		ID     string `json:"id"`
		Format string `json:"format"`
		PostID string `json:"post_id"`
	}

	err = r.client.Query("content_posts").
		Select("id", "format", "post_id").
		Where("tenant_id", "eq", tenantID).
		Execute(&posts)

	if err != nil {
		return nil, fmt.Errorf("failed to get content posts: %w", err)
	}

	// Create a map from post_id to format
	postFormats := make(map[string]string)
	for _, post := range posts {
		postFormats[post.PostID] = post.Format
	}

	// Group metrics by format
	formatMetrics := make(map[string][]repository.PersonalMetric)

	for _, metric := range metrics {
		format, ok := postFormats[metric.PostID]
		if !ok {
			// If no format is found, use "unknown"
			format = "unknown"
		}

		formatMetrics[format] = append(formatMetrics[format], metric)
	}

	// Calculate performance for each format
	var performances []ContentPerformance

	for format, metrics := range formatMetrics {
		if len(metrics) == 0 {
			continue
		}

		var totalLikes, totalShares, totalComments int
		var totalEngagementRate float64

		for _, metric := range metrics {
			totalLikes += metric.Likes
			totalShares += metric.Shares
			totalComments += metric.Comments
			totalEngagementRate += metric.EngagementRate
		}

		count := float64(len(metrics))
		avgLikes := float64(totalLikes) / count
		avgShares := float64(totalShares) / count
		avgComments := float64(totalComments) / count
		avgEngagementRate := totalEngagementRate / count

		// Calculate performance score (simplified)
		performanceScore := avgEngagementRate * 100

		// Calculate trend (simplified)
		// In a real implementation, we would compare with previous period
		performanceTrend := 0.0

		performances = append(performances, ContentPerformance{
			Format:            format,
			TotalPosts:        len(metrics),
			AvgEngagementRate: avgEngagementRate,
			AvgLikes:          avgLikes,
			AvgShares:         avgShares,
			AvgComments:       avgComments,
			PerformanceScore:  performanceScore,
			PerformanceTrend:  performanceTrend,
		})
	}

	return performances, nil
}

// SaveRecommendation saves a new recommendation
func (r *SupabaseAnalyticsRepository) SaveRecommendation(ctx context.Context, rec *Recommendation) (*Recommendation, error) {
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}

	rec.TenantID = r.client.TenantID
	rec.CreatedAt = time.Now()
	rec.UpdatedAt = time.Now()

	if rec.Status == "" {
		rec.Status = "pending"
	}

	err := r.client.Insert(ctx, "recommendations", rec)
	if err != nil {
		return nil, fmt.Errorf("failed to save recommendation: %w", err)
	}

	return rec, nil
}

// GetRecommendations retrieves recommendations
func (r *SupabaseAnalyticsRepository) GetRecommendations(ctx context.Context, tenantID string, status string) ([]Recommendation, error) {
	var recommendations []Recommendation

	query := r.client.Query("recommendations").
		Select("*").
		Where("tenant_id", "eq", tenantID)

	if status != "" {
		query = query.Where("status", "eq", status)
	}

	err := query.Order("created_at", true).Execute(&recommendations)

	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	return recommendations, nil
}

// UpdateRecommendationStatus updates a recommendation's status
func (r *SupabaseAnalyticsRepository) UpdateRecommendationStatus(ctx context.Context, tenantID, recID, status string) error {
	// First verify the recommendation exists and belongs to the tenant
	var recs []Recommendation
	err := r.client.Query("recommendations").
		Select("*").
		Where("id", "eq", recID).
		Where("tenant_id", "eq", tenantID).
		Execute(&recs)

	if err != nil {
		return fmt.Errorf("failed to verify recommendation: %w", err)
	}

	if len(recs) == 0 {
		return errors.New("recommendation not found")
	}

	// Update the status
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	err = r.client.Update(ctx, "recommendations", "id", recID, updateData)
	if err != nil {
		return fmt.Errorf("failed to update recommendation status: %w", err)
	}

	return nil
}
