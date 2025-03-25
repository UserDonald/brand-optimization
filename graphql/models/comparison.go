package models

// ComparisonResult represents the full comparison result
type ComparisonResult struct {
	Competitor *CompetitorComparison `json:"competitor"`
	Personal   *PersonalComparison   `json:"personal"`
	Ratios     *ComparisonRatios     `json:"ratios"`
}
