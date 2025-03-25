package models

// Common types used across multiple services

// DateRange represents a date range for querying data
type DateRange struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// MetricAggregates represents aggregated metrics over a time period
type MetricAggregates struct {
	TotalLikes        int     `json:"totalLikes"`
	TotalShares       int     `json:"totalShares"`
	TotalComments     int     `json:"totalComments"`
	AvgEngagementRate float64 `json:"avgEngagementRate"`
	AvgWatchTime      float64 `json:"avgWatchTime"`
}

// ComparisonRatios represents the ratios of metrics between a competitor and the client's own brand
type ComparisonRatios struct {
	LikesRatio          float64 `json:"likesRatio"`
	SharesRatio         float64 `json:"sharesRatio"`
	CommentsRatio       float64 `json:"commentsRatio"`
	EngagementRateRatio float64 `json:"engagementRateRatio"`
	WatchTimeRatio      float64 `json:"watchTimeRatio"`
}
