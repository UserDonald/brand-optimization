package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/donaldnash/go-competitor/common/db"
	"github.com/google/uuid"
)

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	// Notification management
	CreateNotification(ctx context.Context, notification *Notification) (*Notification, error)
	GetNotifications(ctx context.Context, tenantID string, userID string, status string) ([]Notification, error)
	UpdateNotificationStatus(ctx context.Context, tenantID, notificationID, status string) error
	DeleteNotification(ctx context.Context, tenantID, notificationID string) error

	// Alert threshold management
	CreateAlertThreshold(ctx context.Context, threshold *AlertThreshold) (*AlertThreshold, error)
	GetAlertThresholds(ctx context.Context, tenantID string, metricType string) ([]AlertThreshold, error)
	UpdateAlertThreshold(ctx context.Context, threshold *AlertThreshold) (*AlertThreshold, error)
	DeleteAlertThreshold(ctx context.Context, tenantID, thresholdID string) error

	// Scheduled report management
	CreateScheduledReport(ctx context.Context, report *ScheduledReport) (*ScheduledReport, error)
	GetScheduledReports(ctx context.Context, tenantID string) ([]ScheduledReport, error)
	UpdateScheduledReport(ctx context.Context, report *ScheduledReport) (*ScheduledReport, error)
	DeleteScheduledReport(ctx context.Context, tenantID, reportID string) error

	// Check for triggered alerts
	CheckAlertThresholds(ctx context.Context, tenantID string) ([]Notification, error)

	// Process scheduled reports
	ProcessScheduledReports(ctx context.Context, tenantID string) ([]Notification, error)
}

// Notification represents a single notification to a user
type Notification struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // "alert", "report", "system"
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Priority  string    `json:"priority"` // "high", "medium", "low"
	Status    string    `json:"status"`   // "unread", "read", "archived"
	Metadata  string    `json:"metadata"` // JSON string with additional data
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AlertThreshold represents a threshold for triggering an alert
type AlertThreshold struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	MetricType     string    `json:"metric_type"`     // "engagement_rate", "likes", "shares", etc.
	ComparisonType string    `json:"comparison_type"` // "above", "below", "change_by"
	Value          float64   `json:"value"`           // The threshold value
	Percentage     bool      `json:"percentage"`      // If true, value is a percentage
	Period         string    `json:"period"`          // "daily", "weekly", "monthly"
	Status         string    `json:"status"`          // "active", "inactive"
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ScheduledReport represents a scheduled report configuration
type ScheduledReport struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ReportType     string    `json:"report_type"`     // "competitor_analysis", "content_performance", etc.
	Schedule       string    `json:"schedule"`        // "daily", "weekly", "monthly", custom cron string
	Filters        string    `json:"filters"`         // JSON string with report filters
	DeliveryType   string    `json:"delivery_type"`   // "email", "in_app", "both"
	EmailAddresses string    `json:"email_addresses"` // Comma-separated list of email addresses
	Status         string    `json:"status"`          // "active", "inactive"
	LastRunAt      time.Time `json:"last_run_at"`
	NextRunAt      time.Time `json:"next_run_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SupabaseNotificationRepository implements NotificationRepository using Supabase
type SupabaseNotificationRepository struct {
	client *db.SupabaseClient
}

// NewSupabaseNotificationRepository creates a new SupabaseNotificationRepository
func NewSupabaseNotificationRepository(tenantID string) (*SupabaseNotificationRepository, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	client, err := db.NewSupabaseClient(tenantID)
	if err != nil {
		return nil, err
	}

	return &SupabaseNotificationRepository{
		client: client,
	}, nil
}

// CreateNotification creates a new notification
func (r *SupabaseNotificationRepository) CreateNotification(ctx context.Context, notification *Notification) (*Notification, error) {
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}

	notification.TenantID = r.client.TenantID
	now := time.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now

	if notification.Status == "" {
		notification.Status = "unread"
	}

	err := r.client.Insert(ctx, "notifications", notification)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

// GetNotifications retrieves notifications based on filters
func (r *SupabaseNotificationRepository) GetNotifications(ctx context.Context, tenantID string, userID string, status string) ([]Notification, error) {
	var notifications []Notification
	query := r.client.Query("notifications").Select("*")

	// Add tenant filter
	query = query.Where("tenant_id", "eq", tenantID)

	// Add user filter if provided
	if userID != "" {
		query = query.Where("user_id", "eq", userID)
	}

	// Add status filter if provided
	if status != "" {
		query = query.Where("status", "eq", status)
	}

	// Order by creation time, newest first
	query = query.Order("created_at", true)

	err := query.Execute(&notifications)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}

// UpdateNotificationStatus updates a notification's status
func (r *SupabaseNotificationRepository) UpdateNotificationStatus(ctx context.Context, tenantID, notificationID, status string) error {
	// Verify the notification exists and belongs to the tenant
	var notifications []Notification
	err := r.client.Query("notifications").
		Select("*").
		Where("id", "eq", notificationID).
		Where("tenant_id", "eq", tenantID).
		Execute(&notifications)

	if err != nil {
		return fmt.Errorf("failed to verify notification: %w", err)
	}

	if len(notifications) == 0 {
		return errors.New("notification not found")
	}

	// Update the status
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	err = r.client.Update(ctx, "notifications", "id", notificationID, updateData)
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

// DeleteNotification deletes a notification
func (r *SupabaseNotificationRepository) DeleteNotification(ctx context.Context, tenantID, notificationID string) error {
	// Verify the notification exists and belongs to the tenant
	var notifications []Notification
	err := r.client.Query("notifications").
		Select("*").
		Where("id", "eq", notificationID).
		Where("tenant_id", "eq", tenantID).
		Execute(&notifications)

	if err != nil {
		return fmt.Errorf("failed to verify notification: %w", err)
	}

	if len(notifications) == 0 {
		return errors.New("notification not found")
	}

	// First, try to perform a hard delete
	err = r.client.Delete(ctx, "notifications", "id", notificationID)

	// If the delete operation is not supported or fails, fall back to soft delete by updating status
	if err != nil {
		log.Printf("Warning: Hard delete failed: %v. Falling back to soft delete.", err)

		updateData := map[string]interface{}{
			"status":     "deleted",
			"updated_at": time.Now(),
		}

		err = r.client.Update(ctx, "notifications", "id", notificationID, updateData)
		if err != nil {
			return fmt.Errorf("failed to soft delete notification: %w", err)
		}
	}

	return nil
}

// CreateAlertThreshold creates a new alert threshold
func (r *SupabaseNotificationRepository) CreateAlertThreshold(ctx context.Context, threshold *AlertThreshold) (*AlertThreshold, error) {
	if threshold.ID == "" {
		threshold.ID = uuid.New().String()
	}

	threshold.TenantID = r.client.TenantID
	now := time.Now()
	threshold.CreatedAt = now
	threshold.UpdatedAt = now

	if threshold.Status == "" {
		threshold.Status = "active"
	}

	err := r.client.Insert(ctx, "alert_thresholds", threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert threshold: %w", err)
	}

	return threshold, nil
}

// GetAlertThresholds retrieves alert thresholds based on filters
func (r *SupabaseNotificationRepository) GetAlertThresholds(ctx context.Context, tenantID string, metricType string) ([]AlertThreshold, error) {
	var thresholds []AlertThreshold
	query := r.client.Query("alert_thresholds").Select("*")

	// Add tenant filter
	query = query.Where("tenant_id", "eq", tenantID)

	// Add metric type filter if provided
	if metricType != "" {
		query = query.Where("metric_type", "eq", metricType)
	}

	// Only get active thresholds
	query = query.Where("status", "eq", "active")

	err := query.Execute(&thresholds)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert thresholds: %w", err)
	}

	return thresholds, nil
}

// UpdateAlertThreshold updates an alert threshold
func (r *SupabaseNotificationRepository) UpdateAlertThreshold(ctx context.Context, threshold *AlertThreshold) (*AlertThreshold, error) {
	// Verify the threshold exists and belongs to the tenant
	var thresholds []AlertThreshold
	err := r.client.Query("alert_thresholds").
		Select("*").
		Where("id", "eq", threshold.ID).
		Where("tenant_id", "eq", threshold.TenantID).
		Execute(&thresholds)

	if err != nil {
		return nil, fmt.Errorf("failed to verify alert threshold: %w", err)
	}

	if len(thresholds) == 0 {
		return nil, errors.New("alert threshold not found")
	}

	// Update the threshold
	threshold.UpdatedAt = time.Now()

	err = r.client.Update(ctx, "alert_thresholds", "id", threshold.ID, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to update alert threshold: %w", err)
	}

	return threshold, nil
}

// DeleteAlertThreshold deletes an alert threshold
func (r *SupabaseNotificationRepository) DeleteAlertThreshold(ctx context.Context, tenantID, thresholdID string) error {
	// Verify the threshold exists and belongs to the tenant
	var thresholds []AlertThreshold
	err := r.client.Query("alert_thresholds").
		Select("*").
		Where("id", "eq", thresholdID).
		Where("tenant_id", "eq", tenantID).
		Execute(&thresholds)

	if err != nil {
		return fmt.Errorf("failed to verify alert threshold: %w", err)
	}

	if len(thresholds) == 0 {
		return errors.New("alert threshold not found")
	}

	// Update status to inactive instead of deleting
	updateData := map[string]interface{}{
		"status":     "inactive",
		"updated_at": time.Now(),
	}

	err = r.client.Update(ctx, "alert_thresholds", "id", thresholdID, updateData)
	if err != nil {
		return fmt.Errorf("failed to delete alert threshold: %w", err)
	}

	return nil
}

// CreateScheduledReport creates a new scheduled report
func (r *SupabaseNotificationRepository) CreateScheduledReport(ctx context.Context, report *ScheduledReport) (*ScheduledReport, error) {
	if report.ID == "" {
		report.ID = uuid.New().String()
	}

	report.TenantID = r.client.TenantID
	now := time.Now()
	report.CreatedAt = now
	report.UpdatedAt = now

	if report.Status == "" {
		report.Status = "active"
	}

	// Calculate next run time
	nextRun, err := calculateNextRunTime(report.Schedule, now)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate next run time: %w", err)
	}
	report.NextRunAt = nextRun

	err = r.client.Insert(ctx, "scheduled_reports", report)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduled report: %w", err)
	}

	return report, nil
}

// GetScheduledReports retrieves scheduled reports
func (r *SupabaseNotificationRepository) GetScheduledReports(ctx context.Context, tenantID string) ([]ScheduledReport, error) {
	var reports []ScheduledReport
	query := r.client.Query("scheduled_reports").Select("*")

	// Add tenant filter
	query = query.Where("tenant_id", "eq", tenantID)

	err := query.Execute(&reports)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled reports: %w", err)
	}

	return reports, nil
}

// UpdateScheduledReport updates a scheduled report
func (r *SupabaseNotificationRepository) UpdateScheduledReport(ctx context.Context, report *ScheduledReport) (*ScheduledReport, error) {
	// Verify the report exists and belongs to the tenant
	var reports []ScheduledReport
	err := r.client.Query("scheduled_reports").
		Select("*").
		Where("id", "eq", report.ID).
		Where("tenant_id", "eq", report.TenantID).
		Execute(&reports)

	if err != nil {
		return nil, fmt.Errorf("failed to verify scheduled report: %w", err)
	}

	if len(reports) == 0 {
		return nil, errors.New("scheduled report not found")
	}

	// Update the report
	report.UpdatedAt = time.Now()

	// If the schedule has changed, recalculate next run time
	if report.Schedule != reports[0].Schedule {
		nextRun, err := calculateNextRunTime(report.Schedule, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to calculate next run time: %w", err)
		}
		report.NextRunAt = nextRun
	}

	err = r.client.Update(ctx, "scheduled_reports", "id", report.ID, report)
	if err != nil {
		return nil, fmt.Errorf("failed to update scheduled report: %w", err)
	}

	return report, nil
}

// DeleteScheduledReport deletes a scheduled report
func (r *SupabaseNotificationRepository) DeleteScheduledReport(ctx context.Context, tenantID, reportID string) error {
	// Verify the report exists and belongs to the tenant
	var reports []ScheduledReport
	err := r.client.Query("scheduled_reports").
		Select("*").
		Where("id", "eq", reportID).
		Where("tenant_id", "eq", tenantID).
		Execute(&reports)

	if err != nil {
		return fmt.Errorf("failed to verify scheduled report: %w", err)
	}

	if len(reports) == 0 {
		return errors.New("scheduled report not found")
	}

	// Update status to inactive instead of deleting
	updateData := map[string]interface{}{
		"status":     "inactive",
		"updated_at": time.Now(),
	}

	err = r.client.Update(ctx, "scheduled_reports", "id", reportID, updateData)
	if err != nil {
		return fmt.Errorf("failed to delete scheduled report: %w", err)
	}

	return nil
}

// CheckAlertThresholds checks all active alert thresholds and creates notifications if triggered
func (r *SupabaseNotificationRepository) CheckAlertThresholds(ctx context.Context, tenantID string) ([]Notification, error) {
	// Get all active alert thresholds
	thresholds, err := r.GetAlertThresholds(ctx, tenantID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get alert thresholds: %w", err)
	}

	var notifications []Notification

	// For each threshold, check if it's triggered
	for _, threshold := range thresholds {
		isTriggered, metricValue, err := r.checkThreshold(ctx, threshold)
		if err != nil {
			// Log error but continue with other thresholds
			fmt.Printf("Error checking threshold %s: %v\n", threshold.ID, err)
			continue
		}

		// If triggered, create a notification
		if isTriggered {
			notification := &Notification{
				TenantID:  tenantID,
				UserID:    threshold.UserID,
				Type:      "alert",
				Title:     fmt.Sprintf("Alert: %s", threshold.Name),
				Message:   generateAlertMessage(threshold, metricValue),
				Priority:  "high",
				Status:    "unread",
				Metadata:  fmt.Sprintf(`{"threshold_id":"%s","metric_type":"%s","value":%f}`, threshold.ID, threshold.MetricType, metricValue),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			createdNotification, err := r.CreateNotification(ctx, notification)
			if err != nil {
				fmt.Printf("Error creating notification for threshold %s: %v\n", threshold.ID, err)
				continue
			}

			notifications = append(notifications, *createdNotification)
		}
	}

	return notifications, nil
}

// ProcessScheduledReports processes all scheduled reports that are due to run
func (r *SupabaseNotificationRepository) ProcessScheduledReports(ctx context.Context, tenantID string) ([]Notification, error) {
	// Get all active scheduled reports
	reports, err := r.GetScheduledReports(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled reports: %w", err)
	}

	now := time.Now()
	var notifications []Notification

	// For each report, check if it's due to run
	for _, report := range reports {
		// Skip reports that are inactive
		if report.Status != "active" {
			continue
		}

		// Skip reports that aren't due yet
		if now.Before(report.NextRunAt) {
			continue
		}

		// Generate report (in a real implementation, this would generate the actual report)
		// For now, we'll just create a notification
		notification := &Notification{
			TenantID:  tenantID,
			UserID:    report.UserID,
			Type:      "report",
			Title:     fmt.Sprintf("Report: %s", report.Name),
			Message:   fmt.Sprintf("Your scheduled report '%s' is ready.", report.Name),
			Priority:  "medium",
			Status:    "unread",
			Metadata:  fmt.Sprintf(`{"report_id":"%s","report_type":"%s"}`, report.ID, report.ReportType),
			CreatedAt: now,
			UpdatedAt: now,
		}

		createdNotification, err := r.CreateNotification(ctx, notification)
		if err != nil {
			fmt.Printf("Error creating notification for report %s: %v\n", report.ID, err)
			continue
		}

		notifications = append(notifications, *createdNotification)

		// Update the report's last run time and calculate next run time
		report.LastRunAt = now
		nextRun, err := calculateNextRunTime(report.Schedule, now)
		if err != nil {
			fmt.Printf("Error calculating next run time for report %s: %v\n", report.ID, err)
			continue
		}

		report.NextRunAt = nextRun
		report.UpdatedAt = now

		_, err = r.UpdateScheduledReport(ctx, &report)
		if err != nil {
			fmt.Printf("Error updating report %s after processing: %v\n", report.ID, err)
		}
	}

	return notifications, nil
}

// Helper functions

// checkThreshold checks if an alert threshold is triggered
func (r *SupabaseNotificationRepository) checkThreshold(ctx context.Context, threshold AlertThreshold) (bool, float64, error) {
	// This would be implemented based on your specific metrics data
	// For now, we'll return a simulated result
	metricValue := 0.0

	// Simulate fetching current metric value
	switch threshold.MetricType {
	case "engagement_rate":
		metricValue = 0.05 // 5% engagement rate
	case "likes":
		metricValue = 150.0
	case "shares":
		metricValue = 25.0
	case "comments":
		metricValue = 35.0
	default:
		return false, 0, fmt.Errorf("unsupported metric type: %s", threshold.MetricType)
	}

	// Check if threshold is triggered
	var isTriggered bool

	switch threshold.ComparisonType {
	case "above":
		isTriggered = metricValue > threshold.Value
	case "below":
		isTriggered = metricValue < threshold.Value
	case "change_by":
		// This would involve comparing with historical data
		// For simplicity, we'll just simulate a change
		isTriggered = true
	default:
		return false, 0, fmt.Errorf("unsupported comparison type: %s", threshold.ComparisonType)
	}

	return isTriggered, metricValue, nil
}

// generateAlertMessage generates a user-friendly message for an alert
func generateAlertMessage(threshold AlertThreshold, currentValue float64) string {
	var comparisonStr string

	switch threshold.ComparisonType {
	case "above":
		comparisonStr = "above"
	case "below":
		comparisonStr = "below"
	case "change_by":
		if currentValue > threshold.Value {
			comparisonStr = "increased by"
		} else {
			comparisonStr = "decreased by"
		}
	}

	valueStr := fmt.Sprintf("%.2f", threshold.Value)
	if threshold.Percentage {
		valueStr = fmt.Sprintf("%.2f%%", threshold.Value*100)
	}

	currentValueStr := fmt.Sprintf("%.2f", currentValue)
	if threshold.Percentage {
		currentValueStr = fmt.Sprintf("%.2f%%", currentValue*100)
	}

	return fmt.Sprintf("Your %s is %s threshold of %s. Current value: %s.",
		formatMetricType(threshold.MetricType),
		comparisonStr,
		valueStr,
		currentValueStr)
}

// formatMetricType returns a user-friendly name for a metric type
func formatMetricType(metricType string) string {
	switch metricType {
	case "engagement_rate":
		return "engagement rate"
	case "likes":
		return "like count"
	case "shares":
		return "share count"
	case "comments":
		return "comment count"
	default:
		return metricType
	}
}

// calculateNextRunTime calculates the next run time based on a schedule
func calculateNextRunTime(schedule string, from time.Time) (time.Time, error) {
	switch schedule {
	case "daily":
		// Next day, same time
		return from.AddDate(0, 0, 1), nil
	case "weekly":
		// Next week, same day and time
		return from.AddDate(0, 0, 7), nil
	case "monthly":
		// Next month, same day and time
		return from.AddDate(0, 1, 0), nil
	default:
		// For custom schedules, we would parse cron syntax
		// For simplicity, we'll treat any other value as daily
		return from.AddDate(0, 0, 1), nil
	}
}
