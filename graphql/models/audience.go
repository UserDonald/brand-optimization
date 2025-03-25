package models

// AudienceSegment represents an audience segment
type AudienceSegment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SegmentMetric represents metrics for an audience segment
type SegmentMetric struct {
	SegmentID         string  `json:"segmentID"`
	Size              int     `json:"size"`
	EngagementRate    float64 `json:"engagementRate"`
	ContentPreference string  `json:"contentPreference"`
}

// CreateAudienceSegmentInput represents input for creating a new audience segment
type CreateAudienceSegmentInput struct {
	TenantID    string `json:"tenantID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// UpdateAudienceSegmentInput represents input for updating an existing audience segment
type UpdateAudienceSegmentInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
