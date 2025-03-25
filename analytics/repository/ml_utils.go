package repository

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/donaldnash/go-competitor/engagement/repository"
)

// MLModel represents a simple machine learning model for predictions
type MLModel struct {
	// Model parameters
	weights map[string]float64
	biases  map[string]float64
}

// NewMLModel initializes a new ML model
func NewMLModel() *MLModel {
	return &MLModel{
		weights: make(map[string]float64),
		biases:  make(map[string]float64),
	}
}

// TrainTimeModel trains a model for time-based predictions
func (m *MLModel) TrainTimeModel(metrics []repository.PersonalMetric) {
	if len(metrics) == 0 {
		return
	}

	// Group metrics by day and hour
	dayHourEngagement := make(map[string]map[int][]float64)

	for _, metric := range metrics {
		day := metric.PostedAt.Weekday().String()
		hour := metric.PostedAt.Hour()

		if dayHourEngagement[day] == nil {
			dayHourEngagement[day] = make(map[int][]float64)
		}

		dayHourEngagement[day][hour] = append(dayHourEngagement[day][hour], metric.EngagementRate)
	}

	// Calculate weights based on average engagement rates
	for day, hourData := range dayHourEngagement {
		for hour, rates := range hourData {
			if len(rates) == 0 {
				continue
			}

			// Calculate average engagement rate
			sum := 0.0
			for _, rate := range rates {
				sum += rate
			}
			avg := sum / float64(len(rates))

			// Store as model weight
			key := timeKey(day, hour)
			m.weights[key] = avg

			// Calculate confidence based on sample size
			confidence := math.Min(0.5+(float64(len(rates))/20.0), 0.95)
			m.biases[key] = confidence
		}
	}
}

// TrainFormatModel trains a model for content format predictions
func (m *MLModel) TrainFormatModel(metrics []repository.PersonalMetric, formats map[string][]string) {
	if len(metrics) == 0 {
		return
	}

	// Group metrics by content format
	formatEngagement := make(map[string][]float64)

	for _, metric := range metrics {
		// In a real implementation, we would have content format information
		// attached to each metric. For this example, we'll use a simple lookup
		// from a provided mapping.
		format := "unknown"
		for f, postIDs := range formats {
			for _, id := range postIDs {
				if id == metric.PostID {
					format = f
					break
				}
			}
		}

		formatEngagement[format] = append(formatEngagement[format], metric.EngagementRate)
	}

	// Calculate weights based on average engagement rates
	for format, rates := range formatEngagement {
		if len(rates) == 0 {
			continue
		}

		// Calculate average engagement rate
		sum := 0.0
		for _, rate := range rates {
			sum += rate
		}
		avg := sum / float64(len(rates))

		// Store as model weight
		m.weights[format] = avg

		// Calculate confidence based on sample size
		confidence := math.Min(0.5+(float64(len(rates))/10.0), 0.95)
		m.biases[format] = confidence
	}
}

// PredictEngagementForTime predicts engagement rate for a given time
func (m *MLModel) PredictEngagementForTime(postTime time.Time) (float64, float64) {
	day := postTime.Weekday().String()
	hour := postTime.Hour()

	key := timeKey(day, hour)

	weight, exists := m.weights[key]
	if !exists {
		// Fall back to average for the day if specific hour isn't found
		dayTotal := 0.0
		dayCount := 0
		dayConfidence := 0.0

		for k, w := range m.weights {
			if strings.HasPrefix(k, day+":") {
				dayTotal += w
				dayCount++
				dayConfidence += m.biases[k]
			}
		}

		if dayCount > 0 {
			return dayTotal / float64(dayCount), dayConfidence / float64(dayCount)
		}

		// If no data for the day, return overall average
		return 0.05, 0.3 // Default values
	}

	return weight, m.biases[key]
}

// PredictEngagementForFormat predicts engagement rate for a given content format
func (m *MLModel) PredictEngagementForFormat(format string) (float64, float64) {
	weight, exists := m.weights[format]
	if !exists {
		// Return default if format isn't found
		return 0.05, 0.3 // Default values
	}

	return weight, m.biases[format]
}

// Helper function to create a key for day and hour
func timeKey(day string, hour int) string {
	return fmt.Sprintf("%s:%d", day, hour)
}
