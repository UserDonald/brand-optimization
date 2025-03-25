package models

// Competitor represents a competitor being tracked
type Competitor struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	CreatedAt string `json:"createdAt"`
}

// CompetitorMetric represents a single metric record for a competitor
type CompetitorMetric struct {
	ID               string  `json:"id"`
	CompetitorID     string  `json:"competitorID"`
	PostID           string  `json:"postID"`
	Likes            int     `json:"likes"`
	Shares           int     `json:"shares"`
	Comments         int     `json:"comments"`
	ClickThroughRate float64 `json:"clickThroughRate"`
	AvgWatchTime     float64 `json:"avgWatchTime"`
	EngagementRate   float64 `json:"engagementRate"`
	PostedAt         string  `json:"postedAt"`
}

// CompetitorComparison represents metrics for a competitor
type CompetitorComparison struct {
	Metrics    []*CompetitorMetric `json:"metrics"`
	Aggregates *MetricAggregates   `json:"aggregates"`
}

// AddCompetitorInput represents input for adding a new competitor
type AddCompetitorInput struct {
	TenantID string `json:"tenantID"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

// UpdateCompetitorInput represents input for updating an existing competitor
type UpdateCompetitorInput struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
}
