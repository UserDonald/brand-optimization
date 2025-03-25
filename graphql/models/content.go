package models

// ContentFormat represents a content format
type ContentFormat struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// FormatPerformance represents performance metrics for a content format
type FormatPerformance struct {
	FormatID       string  `json:"formatID"`
	EngagementRate float64 `json:"engagementRate"`
	ReachRate      float64 `json:"reachRate"`
	ConversionRate float64 `json:"conversionRate"`
}

// ScheduledPost represents a scheduled post
type ScheduledPost struct {
	ID            string `json:"id"`
	Content       string `json:"content"`
	ScheduledTime string `json:"scheduledTime"`
	Platform      string `json:"platform"`
	Format        string `json:"format"`
	Status        string `json:"status"`
}

// CreateContentFormatInput represents input for creating a new content format
type CreateContentFormatInput struct {
	TenantID    string `json:"tenantID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateContentFormatInput represents input for updating an existing content format
type UpdateContentFormatInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SchedulePostInput represents input for scheduling a new post
type SchedulePostInput struct {
	Content       string `json:"content"`
	ScheduledTime string `json:"scheduledTime"`
	Platform      string `json:"platform"`
	Format        string `json:"format"`
}

// UpdateScheduledPostInput represents input for updating a scheduled post
type UpdateScheduledPostInput struct {
	Content       string `json:"content"`
	Platform      string `json:"platform"`
	Format        string `json:"format"`
	Status        string `json:"status"`
	ScheduledTime string `json:"scheduledTime"`
}
