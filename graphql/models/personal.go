package models

// PersonalMetric represents a single metric record for the tenant's own brand
type PersonalMetric struct {
	PostID           string  `json:"postID"`
	Likes            int     `json:"likes"`
	Shares           int     `json:"shares"`
	Comments         int     `json:"comments"`
	ClickThroughRate float64 `json:"clickThroughRate"`
	AvgWatchTime     float64 `json:"avgWatchTime"`
	EngagementRate   float64 `json:"engagementRate"`
	PostedAt         string  `json:"postedAt"`
}

// PersonalComparison represents metrics for the tenant's own brand
type PersonalComparison struct {
	Metrics    []*PersonalMetric `json:"metrics"`
	Aggregates *MetricAggregates `json:"aggregates"`
}

// UpdatePersonalDataInput represents input for updating personal data
type UpdatePersonalDataInput struct {
	PostID  string        `json:"postID"`
	Metrics *MetricsInput `json:"metrics"`
}

// MetricsInput represents input for updating metrics
type MetricsInput struct {
	Likes            int     `json:"likes"`
	Shares           int     `json:"shares"`
	Comments         int     `json:"comments"`
	ClickThroughRate float64 `json:"clickThroughRate"`
	AvgWatchTime     float64 `json:"avgWatchTime"`
}
