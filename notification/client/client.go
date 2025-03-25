package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/notification/pb"
	"github.com/donaldnash/go-competitor/notification/repository"
	"github.com/donaldnash/go-competitor/notification/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NotificationClient is the client interface for the notification service
type NotificationClient interface {
	// Notification management
	CreateNotification(ctx context.Context, userID, notificationType, title, message, priority, metadata string) (*repository.Notification, error)
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

	// Close the connection
	Close() error
}

// grpcNotificationClient implements NotificationClient interface using gRPC
type grpcNotificationClient struct {
	conn   *grpc.ClientConn
	client pb.NotificationServiceClient
	// We'll keep the local service option for easier testing
	service service.NotificationService
}

// NewGRPCNotificationClient creates a new gRPC notification client
func NewGRPCNotificationClient(serverAddr string) (NotificationClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}

	client := pb.NewNotificationServiceClient(conn)
	return &grpcNotificationClient{
		conn:   conn,
		client: client,
	}, nil
}

// NewLocalNotificationClient creates a client that directly calls the local service
// This is useful for development or when running all services in one process
func NewLocalNotificationClient(svc service.NotificationService) NotificationClient {
	return &grpcNotificationClient{
		service: svc,
	}
}

// Close closes the client connection
func (c *grpcNotificationClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateNotification creates a new notification
func (c *grpcNotificationClient) CreateNotification(ctx context.Context, userID, notificationType, title, message, priority, metadata string) (*repository.Notification, error) {
	if c.service != nil {
		return c.service.CreateNotification(ctx, userID, notificationType, title, message, priority, metadata)
	}

	// Use gRPC client
	req := &pb.CreateNotificationRequest{
		TenantId: "default", // We should get this from context/config
		UserId:   userID,
		Type:     notificationType,
		Title:    title,
		Message:  message,
		Priority: priority,
		Metadata: metadata,
	}

	resp, err := c.client.CreateNotification(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Convert proto notification to repository notification
	return &repository.Notification{
		ID:        resp.Id,
		TenantID:  resp.TenantId,
		UserID:    resp.UserId,
		Type:      resp.Type,
		Title:     resp.Title,
		Message:   resp.Message,
		Priority:  resp.Priority,
		Status:    resp.Status,
		Metadata:  resp.Metadata,
		CreatedAt: resp.CreatedAt.AsTime(),
		UpdatedAt: resp.UpdatedAt.AsTime(),
	}, nil
}

// GetNotifications retrieves notifications based on filters
func (c *grpcNotificationClient) GetNotifications(ctx context.Context, tenantID, userID, status string) ([]repository.Notification, error) {
	if c.service != nil {
		return c.service.GetNotifications(ctx, tenantID, userID, status)
	}

	// Use gRPC client
	req := &pb.GetNotificationsRequest{
		TenantId: tenantID,
		UserId:   userID,
		Status:   status,
	}

	resp, err := c.client.GetNotifications(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	// Convert proto notifications to repository notifications
	notifications := make([]repository.Notification, len(resp.Notifications))
	for i, n := range resp.Notifications {
		notifications[i] = repository.Notification{
			ID:        n.Id,
			TenantID:  n.TenantId,
			UserID:    n.UserId,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			Priority:  n.Priority,
			Status:    n.Status,
			Metadata:  n.Metadata,
			CreatedAt: n.CreatedAt.AsTime(),
			UpdatedAt: n.UpdatedAt.AsTime(),
		}
	}

	return notifications, nil
}

// MarkNotificationAsRead marks a notification as read
func (c *grpcNotificationClient) MarkNotificationAsRead(ctx context.Context, tenantID, notificationID string) error {
	if c.service != nil {
		return c.service.MarkNotificationAsRead(ctx, tenantID, notificationID)
	}

	// Use gRPC client
	req := &pb.NotificationStatusRequest{
		TenantId:       tenantID,
		NotificationId: notificationID,
	}

	_, err := c.client.MarkNotificationAsRead(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}

// ArchiveNotification archives a notification
func (c *grpcNotificationClient) ArchiveNotification(ctx context.Context, tenantID, notificationID string) error {
	if c.service != nil {
		return c.service.ArchiveNotification(ctx, tenantID, notificationID)
	}

	// Use gRPC client
	req := &pb.NotificationStatusRequest{
		TenantId:       tenantID,
		NotificationId: notificationID,
	}

	_, err := c.client.ArchiveNotification(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to archive notification: %w", err)
	}

	return nil
}

// DeleteNotification deletes a notification
func (c *grpcNotificationClient) DeleteNotification(ctx context.Context, tenantID, notificationID string) error {
	if c.service != nil {
		return c.service.DeleteNotification(ctx, tenantID, notificationID)
	}

	// Use gRPC client
	req := &pb.NotificationStatusRequest{
		TenantId:       tenantID,
		NotificationId: notificationID,
	}

	_, err := c.client.DeleteNotification(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// CreateAlertThreshold creates a new alert threshold
func (c *grpcNotificationClient) CreateAlertThreshold(ctx context.Context, userID, name, metricType, comparisonType string, value float64, percentage bool, period string) (*repository.AlertThreshold, error) {
	if c.service != nil {
		return c.service.CreateAlertThreshold(ctx, userID, name, metricType, comparisonType, value, percentage, period)
	}

	// Use gRPC client
	req := &pb.CreateAlertThresholdRequest{
		TenantId:       "default", // We should get this from context/config
		UserId:         userID,
		Name:           name,
		MetricType:     metricType,
		ComparisonType: comparisonType,
		Value:          value,
		Percentage:     percentage,
		Period:         period,
	}

	resp, err := c.client.CreateAlertThreshold(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert threshold: %w", err)
	}

	// Convert proto threshold to repository threshold
	return &repository.AlertThreshold{
		ID:             resp.Id,
		TenantID:       resp.TenantId,
		UserID:         resp.UserId,
		Name:           resp.Name,
		MetricType:     resp.MetricType,
		ComparisonType: resp.ComparisonType,
		Value:          resp.Value,
		Percentage:     resp.Percentage,
		Period:         resp.Period,
		Status:         resp.Status,
		CreatedAt:      resp.CreatedAt.AsTime(),
		UpdatedAt:      resp.UpdatedAt.AsTime(),
	}, nil
}

// GetAlertThresholds retrieves alert thresholds based on filters
func (c *grpcNotificationClient) GetAlertThresholds(ctx context.Context, tenantID, metricType string) ([]repository.AlertThreshold, error) {
	if c.service != nil {
		return c.service.GetAlertThresholds(ctx, tenantID, metricType)
	}

	// Use gRPC client
	req := &pb.GetAlertThresholdsRequest{
		TenantId:   tenantID,
		MetricType: metricType,
	}

	resp, err := c.client.GetAlertThresholds(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert thresholds: %w", err)
	}

	// Convert proto thresholds to repository thresholds
	thresholds := make([]repository.AlertThreshold, len(resp.Thresholds))
	for i, t := range resp.Thresholds {
		thresholds[i] = repository.AlertThreshold{
			ID:             t.Id,
			TenantID:       t.TenantId,
			UserID:         t.UserId,
			Name:           t.Name,
			MetricType:     t.MetricType,
			ComparisonType: t.ComparisonType,
			Value:          t.Value,
			Percentage:     t.Percentage,
			Period:         t.Period,
			Status:         t.Status,
			CreatedAt:      t.CreatedAt.AsTime(),
			UpdatedAt:      t.UpdatedAt.AsTime(),
		}
	}

	return thresholds, nil
}

// UpdateAlertThreshold updates an alert threshold
func (c *grpcNotificationClient) UpdateAlertThreshold(ctx context.Context, thresholdID, tenantID, userID, name, metricType, comparisonType string, value float64, percentage bool, period, status string) (*repository.AlertThreshold, error) {
	if c.service != nil {
		return c.service.UpdateAlertThreshold(ctx, thresholdID, tenantID, userID, name, metricType, comparisonType, value, percentage, period, status)
	}

	// Use gRPC client
	req := &pb.UpdateAlertThresholdRequest{
		ThresholdId:    thresholdID,
		TenantId:       tenantID,
		UserId:         userID,
		Name:           name,
		MetricType:     metricType,
		ComparisonType: comparisonType,
		Value:          value,
		Percentage:     percentage,
		Period:         period,
		Status:         status,
	}

	resp, err := c.client.UpdateAlertThreshold(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update alert threshold: %w", err)
	}

	// Convert proto threshold to repository threshold
	return &repository.AlertThreshold{
		ID:             resp.Id,
		TenantID:       resp.TenantId,
		UserID:         resp.UserId,
		Name:           resp.Name,
		MetricType:     resp.MetricType,
		ComparisonType: resp.ComparisonType,
		Value:          resp.Value,
		Percentage:     resp.Percentage,
		Period:         resp.Period,
		Status:         resp.Status,
		CreatedAt:      resp.CreatedAt.AsTime(),
		UpdatedAt:      resp.UpdatedAt.AsTime(),
	}, nil
}

// DeleteAlertThreshold deletes an alert threshold
func (c *grpcNotificationClient) DeleteAlertThreshold(ctx context.Context, tenantID, thresholdID string) error {
	if c.service != nil {
		return c.service.DeleteAlertThreshold(ctx, tenantID, thresholdID)
	}

	// Use gRPC client
	req := &pb.DeleteAlertThresholdRequest{
		TenantId:    tenantID,
		ThresholdId: thresholdID,
	}

	_, err := c.client.DeleteAlertThreshold(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete alert threshold: %w", err)
	}

	return nil
}

// CreateScheduledReport creates a new scheduled report
func (c *grpcNotificationClient) CreateScheduledReport(ctx context.Context, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses string) (*repository.ScheduledReport, error) {
	if c.service != nil {
		return c.service.CreateScheduledReport(ctx, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses)
	}

	// Use gRPC client
	req := &pb.CreateScheduledReportRequest{
		TenantId:       "default", // We should get this from context/config
		UserId:         userID,
		Name:           name,
		Description:    description,
		ReportType:     reportType,
		Schedule:       schedule,
		Filters:        filters,
		DeliveryType:   deliveryType,
		EmailAddresses: emailAddresses,
	}

	resp, err := c.client.CreateScheduledReport(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduled report: %w", err)
	}

	// Convert proto report to repository report
	return convertScheduledReportFromProto(resp), nil
}

// GetScheduledReports retrieves scheduled reports
func (c *grpcNotificationClient) GetScheduledReports(ctx context.Context, tenantID string) ([]repository.ScheduledReport, error) {
	if c.service != nil {
		return c.service.GetScheduledReports(ctx, tenantID)
	}

	// Use gRPC client
	req := &pb.GetScheduledReportsRequest{
		TenantId: tenantID,
	}

	resp, err := c.client.GetScheduledReports(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled reports: %w", err)
	}

	// Convert proto reports to repository reports
	reports := make([]repository.ScheduledReport, len(resp.Reports))
	for i, r := range resp.Reports {
		reports[i] = *convertScheduledReportFromProto(r)
	}

	return reports, nil
}

// UpdateScheduledReport updates a scheduled report
func (c *grpcNotificationClient) UpdateScheduledReport(ctx context.Context, reportID, tenantID, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses, status string) (*repository.ScheduledReport, error) {
	if c.service != nil {
		return c.service.UpdateScheduledReport(ctx, reportID, tenantID, userID, name, description, reportType, schedule, filters, deliveryType, emailAddresses, status)
	}

	// Use gRPC client
	req := &pb.UpdateScheduledReportRequest{
		ReportId:       reportID,
		TenantId:       tenantID,
		UserId:         userID,
		Name:           name,
		Description:    description,
		ReportType:     reportType,
		Schedule:       schedule,
		Filters:        filters,
		DeliveryType:   deliveryType,
		EmailAddresses: emailAddresses,
		Status:         status,
	}

	resp, err := c.client.UpdateScheduledReport(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update scheduled report: %w", err)
	}

	// Convert proto report to repository report
	return convertScheduledReportFromProto(resp), nil
}

// DeleteScheduledReport deletes a scheduled report
func (c *grpcNotificationClient) DeleteScheduledReport(ctx context.Context, tenantID, reportID string) error {
	if c.service != nil {
		return c.service.DeleteScheduledReport(ctx, tenantID, reportID)
	}

	// Use gRPC client
	req := &pb.DeleteScheduledReportRequest{
		TenantId: tenantID,
		ReportId: reportID,
	}

	_, err := c.client.DeleteScheduledReport(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete scheduled report: %w", err)
	}

	return nil
}

// CheckAlertThresholds checks for triggered alert thresholds
func (c *grpcNotificationClient) CheckAlertThresholds(ctx context.Context, tenantID string) ([]repository.Notification, error) {
	if c.service != nil {
		return c.service.CheckAlertThresholds(ctx, tenantID)
	}

	// Use gRPC client
	req := &pb.CheckAlertThresholdsRequest{
		TenantId: tenantID,
	}

	resp, err := c.client.CheckAlertThresholds(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check alert thresholds: %w", err)
	}

	// Convert proto notifications to repository notifications
	notifications := make([]repository.Notification, len(resp.Notifications))
	for i, n := range resp.Notifications {
		notifications[i] = repository.Notification{
			ID:        n.Id,
			TenantID:  n.TenantId,
			UserID:    n.UserId,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			Priority:  n.Priority,
			Status:    n.Status,
			Metadata:  n.Metadata,
			CreatedAt: n.CreatedAt.AsTime(),
			UpdatedAt: n.UpdatedAt.AsTime(),
		}
	}

	return notifications, nil
}

// ProcessScheduledReports processes scheduled reports
func (c *grpcNotificationClient) ProcessScheduledReports(ctx context.Context, tenantID string) ([]repository.Notification, error) {
	if c.service != nil {
		return c.service.ProcessScheduledReports(ctx, tenantID)
	}

	// Use gRPC client
	req := &pb.ProcessScheduledReportsRequest{
		TenantId: tenantID,
	}

	resp, err := c.client.ProcessScheduledReports(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process scheduled reports: %w", err)
	}

	// Convert proto notifications to repository notifications
	notifications := make([]repository.Notification, len(resp.Notifications))
	for i, n := range resp.Notifications {
		notifications[i] = repository.Notification{
			ID:        n.Id,
			TenantID:  n.TenantId,
			UserID:    n.UserId,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			Priority:  n.Priority,
			Status:    n.Status,
			Metadata:  n.Metadata,
			CreatedAt: n.CreatedAt.AsTime(),
			UpdatedAt: n.UpdatedAt.AsTime(),
		}
	}

	return notifications, nil
}

// Helper functions

// convertScheduledReportFromProto converts a proto scheduled report to repository format
func convertScheduledReportFromProto(report *pb.ScheduledReport) *repository.ScheduledReport {
	var lastRunAt, nextRunAt time.Time
	if report.LastRunAt != nil {
		lastRunAt = report.LastRunAt.AsTime()
	}
	if report.NextRunAt != nil {
		nextRunAt = report.NextRunAt.AsTime()
	}

	return &repository.ScheduledReport{
		ID:             report.Id,
		TenantID:       report.TenantId,
		UserID:         report.UserId,
		Name:           report.Name,
		Description:    report.Description,
		ReportType:     report.ReportType,
		Schedule:       report.Schedule,
		Filters:        report.Filters,
		DeliveryType:   report.DeliveryType,
		EmailAddresses: report.EmailAddresses,
		Status:         report.Status,
		LastRunAt:      lastRunAt,
		NextRunAt:      nextRunAt,
		CreatedAt:      report.CreatedAt.AsTime(),
		UpdatedAt:      report.UpdatedAt.AsTime(),
	}
}
