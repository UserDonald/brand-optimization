package models

// PostingTimeRecommendation represents a recommended posting time
type PostingTimeRecommendation struct {
	DayOfWeek               string  `json:"dayOfWeek"`
	TimeOfDay               string  `json:"timeOfDay"`
	PredictedEngagementRate float64 `json:"predictedEngagementRate"`
	Confidence              float64 `json:"confidence"`
}

// ContentFormatRecommendation represents a recommended content format
type ContentFormatRecommendation struct {
	Format                  string  `json:"format"`
	PredictedEngagementRate float64 `json:"predictedEngagementRate"`
	TargetAudience          string  `json:"targetAudience"`
	Confidence              float64 `json:"confidence"`
}
