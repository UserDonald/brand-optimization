package service

import (
	"context"
	"errors"
	"time"

	"github.com/donaldnash/go-competitor/notification/repository"
)

// NotificationService defines the interface for notification service operations
type NotificationService interface {
	// Notification management
	CreateNotification(ctx context.Context, userID, notificationType, title, message, priority string, metadata string) (*repository.Notification, error)
	GetNotifications(ctx context.Context, tenantID, userID, status string) ([]repository.Notification, error)
	MarkNotificationAsRead(ctx context.Context, tenantID, notificationID string) error
	ArchiveNotification(ctx context.Context, tenantID, notificationID string) error
	DeleteNotification(ctx context.Context, tenantID, notificationID string) error

	// Alert threshold management
	CreateAlertThreshold(ctx context.Context, userID, name, metricType, comparisonType string, value float64, percentage bool, period string) (*repository.AlertThreshold, error)
	GetAlertThresholds(ctx context.Context, tenantID, metricType string) ([]repository.AlertThreshold, error)
	UpdateAlertThreshold(ctx context.Context, thresholdID, tenantID, userID, name, metricType, comparisonType string, value float64, percentage bool, period, status string) (*repository.AlertThreshold, error)
	DeleteAlertThreshold(ctx context.Context, tenantID, thresholdID string) error

	// Scheduled report management
	CreateScheduledReport(ctx context.Context, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses string) (*repository.ScheduledReport, error)
	GetScheduledReports(ctx context.Context, tenantID string) ([]repository.ScheduledReport, error)
	UpdateScheduledReport(ctx context.Context, reportID, tenantID, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses, status string) (*repository.ScheduledReport, error)
	DeleteScheduledReport(ctx context.Context, tenantID, reportID string) error

	// Background processes
	CheckAlertThresholds(ctx context.Context, tenantID string) ([]repository.Notification, error)
	ProcessScheduledReports(ctx context.Context, tenantID string) ([]repository.Notification, error)
}

// notificationService implements the NotificationService interface
type notificationService struct {
	repo repository.NotificationRepository
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(repo repository.NotificationRepository) (NotificationService, error) {
	if repo == nil {
		return nil, errors.New("repository is required")
	}

	return &notificationService{
		repo: repo,
	}, nil
}

// CreateNotification creates a new notification
func (s *notificationService) CreateNotification(ctx context.Context, userID, notificationType, title, message, priority, metadata string) (*repository.Notification, error) {
	// Validate inputs
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if notificationType == "" {
		return nil, errors.New("notification type is required")
	}

	if title == "" {
		return nil, errors.New("title is required")
	}

	if message == "" {
		return nil, errors.New("message is required")
	}

	// Default priority to medium if not provided
	if priority == "" {
		priority = "medium"
	}

	// Create notification object
	notification := &repository.Notification{
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Priority:  priority,
		Metadata:  metadata,
		Status:    "unread",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Call repository to save
	return s.repo.CreateNotification(ctx, notification)
}

// GetNotifications retrieves notifications based on filters
func (s *notificationService) GetNotifications(ctx context.Context, tenantID, userID, status string) ([]repository.Notification, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.GetNotifications(ctx, tenantID, userID, status)
}

// MarkNotificationAsRead marks a notification as read
func (s *notificationService) MarkNotificationAsRead(ctx context.Context, tenantID, notificationID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if notificationID == "" {
		return errors.New("notification ID is required")
	}

	return s.repo.UpdateNotificationStatus(ctx, tenantID, notificationID, "read")
}

// ArchiveNotification archives a notification
func (s *notificationService) ArchiveNotification(ctx context.Context, tenantID, notificationID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if notificationID == "" {
		return errors.New("notification ID is required")
	}

	return s.repo.UpdateNotificationStatus(ctx, tenantID, notificationID, "archived")
}

// DeleteNotification deletes a notification
func (s *notificationService) DeleteNotification(ctx context.Context, tenantID, notificationID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if notificationID == "" {
		return errors.New("notification ID is required")
	}

	return s.repo.DeleteNotification(ctx, tenantID, notificationID)
}

// CreateAlertThreshold creates a new alert threshold
func (s *notificationService) CreateAlertThreshold(ctx context.Context, userID, name, metricType, comparisonType string, value float64, percentage bool, period string) (*repository.AlertThreshold, error) {
	// Validate inputs
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if metricType == "" {
		return nil, errors.New("metric type is required")
	}

	if comparisonType == "" {
		return nil, errors.New("comparison type is required")
	}

	if period == "" {
		period = "daily" // Default to daily if not specified
	}

	// Create threshold object
	threshold := &repository.AlertThreshold{
		UserID:         userID,
		Name:           name,
		MetricType:     metricType,
		ComparisonType: comparisonType,
		Value:          value,
		Percentage:     percentage,
		Period:         period,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Call repository to save
	return s.repo.CreateAlertThreshold(ctx, threshold)
}

// GetAlertThresholds retrieves alert thresholds based on filters
func (s *notificationService) GetAlertThresholds(ctx context.Context, tenantID, metricType string) ([]repository.AlertThreshold, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.GetAlertThresholds(ctx, tenantID, metricType)
}

// UpdateAlertThreshold updates an alert threshold
func (s *notificationService) UpdateAlertThreshold(ctx context.Context, thresholdID, tenantID, userID, name, metricType, comparisonType string, value float64, percentage bool, period, status string) (*repository.AlertThreshold, error) {
	// Validate inputs
	if thresholdID == "" {
		return nil, errors.New("threshold ID is required")
	}

	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Get existing threshold to preserve any fields not being updated
	thresholds, err := s.repo.GetAlertThresholds(ctx, tenantID, "")
	if err != nil {
		return nil, err
	}

	var existingThreshold *repository.AlertThreshold
	for _, t := range thresholds {
		if t.ID == thresholdID {
			existingThreshold = &t
			break
		}
	}

	if existingThreshold == nil {
		return nil, errors.New("threshold not found")
	}

	// Update fields if provided
	if userID != "" {
		existingThreshold.UserID = userID
	}
	if name != "" {
		existingThreshold.Name = name
	}
	if metricType != "" {
		existingThreshold.MetricType = metricType
	}
	if comparisonType != "" {
		existingThreshold.ComparisonType = comparisonType
	}
	if value != 0 {
		existingThreshold.Value = value
	}
	existingThreshold.Percentage = percentage
	if period != "" {
		existingThreshold.Period = period
	}
	if status != "" {
		existingThreshold.Status = status
	}

	existingThreshold.UpdatedAt = time.Now()

	// Call repository to update
	return s.repo.UpdateAlertThreshold(ctx, existingThreshold)
}

// DeleteAlertThreshold deletes an alert threshold
func (s *notificationService) DeleteAlertThreshold(ctx context.Context, tenantID, thresholdID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if thresholdID == "" {
		return errors.New("threshold ID is required")
	}

	return s.repo.DeleteAlertThreshold(ctx, tenantID, thresholdID)
}

// CreateScheduledReport creates a new scheduled report
func (s *notificationService) CreateScheduledReport(ctx context.Context, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses string) (*repository.ScheduledReport, error) {
	// Validate inputs
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if name == "" {
		return nil, errors.New("name is required")
	}

	if reportType == "" {
		return nil, errors.New("report type is required")
	}

	if schedule == "" {
		schedule = "weekly" // Default to weekly if not specified
	}

	if deliveryType == "" {
		deliveryType = "in_app" // Default to in-app if not specified
	}

	// Create report object
	report := &repository.ScheduledReport{
		UserID:         userID,
		Name:           name,
		Description:    description,
		ReportType:     reportType,
		Schedule:       schedule,
		Filters:        filters,
		DeliveryType:   deliveryType,
		EmailAddresses: emailAddresses,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Call repository to save
	return s.repo.CreateScheduledReport(ctx, report)
}

// GetScheduledReports retrieves scheduled reports
func (s *notificationService) GetScheduledReports(ctx context.Context, tenantID string) ([]repository.ScheduledReport, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.GetScheduledReports(ctx, tenantID)
}

// UpdateScheduledReport updates a scheduled report
func (s *notificationService) UpdateScheduledReport(ctx context.Context, reportID, tenantID, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses, status string) (*repository.ScheduledReport, error) {
	// Validate inputs
	if reportID == "" {
		return nil, errors.New("report ID is required")
	}

	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Get existing reports
	reports, err := s.repo.GetScheduledReports(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var existingReport *repository.ScheduledReport
	for _, r := range reports {
		if r.ID == reportID {
			existingReport = &r
			break
		}
	}

	if existingReport == nil {
		return nil, errors.New("report not found")
	}

	// Update fields if provided
	if userID != "" {
		existingReport.UserID = userID
	}
	if name != "" {
		existingReport.Name = name
	}
	if description != "" {
		existingReport.Description = description
	}
	if reportType != "" {
		existingReport.ReportType = reportType
	}
	if schedule != "" {
		existingReport.Schedule = schedule
	}
	if filters != "" {
		existingReport.Filters = filters
	}
	if deliveryType != "" {
		existingReport.DeliveryType = deliveryType
	}
	if emailAddresses != "" {
		existingReport.EmailAddresses = emailAddresses
	}
	if status != "" {
		existingReport.Status = status
	}

	existingReport.UpdatedAt = time.Now()

	// Call repository to update
	return s.repo.UpdateScheduledReport(ctx, existingReport)
}

// DeleteScheduledReport deletes a scheduled report
func (s *notificationService) DeleteScheduledReport(ctx context.Context, tenantID, reportID string) error {
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	if reportID == "" {
		return errors.New("report ID is required")
	}

	return s.repo.DeleteScheduledReport(ctx, tenantID, reportID)
}

// CheckAlertThresholds checks all active alert thresholds and creates notifications if triggered
func (s *notificationService) CheckAlertThresholds(ctx context.Context, tenantID string) ([]repository.Notification, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.CheckAlertThresholds(ctx, tenantID)
}

// ProcessScheduledReports processes all scheduled reports that are due to run
func (s *notificationService) ProcessScheduledReports(ctx context.Context, tenantID string) ([]repository.Notification, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.ProcessScheduledReports(ctx, tenantID)
}
