package server

import (
	"context"
	"fmt"

	"github.com/donaldnash/go-competitor/notification/pb"
	"github.com/donaldnash/go-competitor/notification/repository"
	"github.com/donaldnash/go-competitor/notification/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NotificationServer is the gRPC server implementation for the notification service
type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
	service service.NotificationService
}

// NewNotificationServer creates a new notification gRPC server
func NewNotificationServer(svc service.NotificationService) (*NotificationServer, error) {
	if svc == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	return &NotificationServer{
		service: svc,
	}, nil
}

// CreateNotification creates a new notification
func (s *NotificationServer) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.Notification, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if req.Message == "" {
		return nil, status.Error(codes.InvalidArgument, "message is required")
	}

	// Call service to create notification
	notification, err := s.service.CreateNotification(ctx, req.UserId, req.Type, req.Title, req.Message, req.Priority, req.Metadata)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create notification: %v", err)
	}

	// Convert to protobuf message
	return convertNotificationToProto(notification), nil
}

// GetNotifications retrieves notifications based on filters
func (s *NotificationServer) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.NotificationsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service to get notifications
	notifications, err := s.service.GetNotifications(ctx, req.TenantId, req.UserId, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get notifications: %v", err)
	}

	// Convert to protobuf response
	response := &pb.NotificationsResponse{
		Notifications: make([]*pb.Notification, 0, len(notifications)),
	}

	for _, notification := range notifications {
		response.Notifications = append(response.Notifications, convertNotificationToProto(&notification))
	}

	return response, nil
}

// MarkNotificationAsRead marks a notification as read
func (s *NotificationServer) MarkNotificationAsRead(ctx context.Context, req *pb.NotificationStatusRequest) (*pb.UpdateStatusResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.NotificationId == "" {
		return nil, status.Error(codes.InvalidArgument, "notification ID is required")
	}

	// Call service to mark notification as read
	err := s.service.MarkNotificationAsRead(ctx, req.TenantId, req.NotificationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark notification as read: %v", err)
	}

	return &pb.UpdateStatusResponse{
		Success: true,
		Message: "Notification marked as read",
	}, nil
}

// ArchiveNotification archives a notification
func (s *NotificationServer) ArchiveNotification(ctx context.Context, req *pb.NotificationStatusRequest) (*pb.UpdateStatusResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.NotificationId == "" {
		return nil, status.Error(codes.InvalidArgument, "notification ID is required")
	}

	// Call service to archive notification
	err := s.service.ArchiveNotification(ctx, req.TenantId, req.NotificationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to archive notification: %v", err)
	}

	return &pb.UpdateStatusResponse{
		Success: true,
		Message: "Notification archived",
	}, nil
}

// DeleteNotification deletes a notification
func (s *NotificationServer) DeleteNotification(ctx context.Context, req *pb.NotificationStatusRequest) (*pb.UpdateStatusResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.NotificationId == "" {
		return nil, status.Error(codes.InvalidArgument, "notification ID is required")
	}

	// Call service to delete notification
	err := s.service.DeleteNotification(ctx, req.TenantId, req.NotificationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete notification: %v", err)
	}

	return &pb.UpdateStatusResponse{
		Success: true,
		Message: "Notification deleted",
	}, nil
}

// CreateAlertThreshold creates a new alert threshold
func (s *NotificationServer) CreateAlertThreshold(ctx context.Context, req *pb.CreateAlertThresholdRequest) (*pb.AlertThreshold, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.MetricType == "" {
		return nil, status.Error(codes.InvalidArgument, "metric type is required")
	}

	if req.ComparisonType == "" {
		return nil, status.Error(codes.InvalidArgument, "comparison type is required")
	}

	// Call service to create alert threshold
	threshold, err := s.service.CreateAlertThreshold(ctx, req.UserId, req.Name, req.MetricType, req.ComparisonType, req.Value, req.Percentage, req.Period)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create alert threshold: %v", err)
	}

	// Convert to protobuf message
	return convertAlertThresholdToProto(threshold), nil
}

// GetAlertThresholds retrieves alert thresholds based on filters
func (s *NotificationServer) GetAlertThresholds(ctx context.Context, req *pb.GetAlertThresholdsRequest) (*pb.AlertThresholdsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service to get alert thresholds
	thresholds, err := s.service.GetAlertThresholds(ctx, req.TenantId, req.MetricType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get alert thresholds: %v", err)
	}

	// Convert to protobuf response
	response := &pb.AlertThresholdsResponse{
		Thresholds: make([]*pb.AlertThreshold, 0, len(thresholds)),
	}

	for _, threshold := range thresholds {
		response.Thresholds = append(response.Thresholds, convertAlertThresholdToProto(&threshold))
	}

	return response, nil
}

// UpdateAlertThreshold updates an alert threshold
func (s *NotificationServer) UpdateAlertThreshold(ctx context.Context, req *pb.UpdateAlertThresholdRequest) (*pb.AlertThreshold, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.ThresholdId == "" {
		return nil, status.Error(codes.InvalidArgument, "threshold ID is required")
	}

	// Call service to update alert threshold
	threshold, err := s.service.UpdateAlertThreshold(ctx, req.ThresholdId, req.TenantId, req.UserId, req.Name, req.MetricType, req.ComparisonType, req.Value, req.Percentage, req.Period, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update alert threshold: %v", err)
	}

	// Convert to protobuf message
	return convertAlertThresholdToProto(threshold), nil
}

// DeleteAlertThreshold deletes an alert threshold
func (s *NotificationServer) DeleteAlertThreshold(ctx context.Context, req *pb.DeleteAlertThresholdRequest) (*pb.UpdateStatusResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.ThresholdId == "" {
		return nil, status.Error(codes.InvalidArgument, "threshold ID is required")
	}

	// Call service to delete alert threshold
	err := s.service.DeleteAlertThreshold(ctx, req.TenantId, req.ThresholdId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete alert threshold: %v", err)
	}

	return &pb.UpdateStatusResponse{
		Success: true,
		Message: "Alert threshold deleted",
	}, nil
}

// CreateScheduledReport creates a new scheduled report
func (s *NotificationServer) CreateScheduledReport(ctx context.Context, req *pb.CreateScheduledReportRequest) (*pb.ScheduledReport, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.ReportType == "" {
		return nil, status.Error(codes.InvalidArgument, "report type is required")
	}

	// Call service to create scheduled report
	report, err := s.service.CreateScheduledReport(ctx, req.UserId, req.Name, req.Description, req.ReportType, req.Schedule, req.Filters, req.DeliveryType, req.EmailAddresses)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create scheduled report: %v", err)
	}

	// Convert to protobuf message
	return convertScheduledReportToProto(report), nil
}

// GetScheduledReports retrieves scheduled reports
func (s *NotificationServer) GetScheduledReports(ctx context.Context, req *pb.GetScheduledReportsRequest) (*pb.ScheduledReportsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service to get scheduled reports
	reports, err := s.service.GetScheduledReports(ctx, req.TenantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get scheduled reports: %v", err)
	}

	// Convert to protobuf response
	response := &pb.ScheduledReportsResponse{
		Reports: make([]*pb.ScheduledReport, 0, len(reports)),
	}

	for _, report := range reports {
		response.Reports = append(response.Reports, convertScheduledReportToProto(&report))
	}

	return response, nil
}

// UpdateScheduledReport updates a scheduled report
func (s *NotificationServer) UpdateScheduledReport(ctx context.Context, req *pb.UpdateScheduledReportRequest) (*pb.ScheduledReport, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.ReportId == "" {
		return nil, status.Error(codes.InvalidArgument, "report ID is required")
	}

	// Call service to update scheduled report
	report, err := s.service.UpdateScheduledReport(
		ctx,
		req.ReportId,
		req.TenantId,
		req.UserId,
		req.Name,
		req.Description,
		req.ReportType,
		req.Schedule,
		req.Filters,
		req.DeliveryType,
		req.EmailAddresses,
		req.Status,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update scheduled report: %v", err)
	}

	// Convert to protobuf message
	return convertScheduledReportToProto(report), nil
}

// DeleteScheduledReport deletes a scheduled report
func (s *NotificationServer) DeleteScheduledReport(ctx context.Context, req *pb.DeleteScheduledReportRequest) (*pb.UpdateStatusResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.ReportId == "" {
		return nil, status.Error(codes.InvalidArgument, "report ID is required")
	}

	// Call service to delete scheduled report
	err := s.service.DeleteScheduledReport(ctx, req.TenantId, req.ReportId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete scheduled report: %v", err)
	}

	return &pb.UpdateStatusResponse{
		Success: true,
		Message: "Scheduled report deleted",
	}, nil
}

// CheckAlertThresholds checks for triggered alert thresholds
func (s *NotificationServer) CheckAlertThresholds(ctx context.Context, req *pb.CheckAlertThresholdsRequest) (*pb.NotificationsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service to check alert thresholds
	notifications, err := s.service.CheckAlertThresholds(ctx, req.TenantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check alert thresholds: %v", err)
	}

	// Convert to protobuf response
	response := &pb.NotificationsResponse{
		Notifications: make([]*pb.Notification, 0, len(notifications)),
	}

	for _, notification := range notifications {
		response.Notifications = append(response.Notifications, convertNotificationToProto(&notification))
	}

	return response, nil
}

// ProcessScheduledReports processes scheduled reports
func (s *NotificationServer) ProcessScheduledReports(ctx context.Context, req *pb.ProcessScheduledReportsRequest) (*pb.NotificationsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Call service to process scheduled reports
	notifications, err := s.service.ProcessScheduledReports(ctx, req.TenantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process scheduled reports: %v", err)
	}

	// Convert to protobuf response
	response := &pb.NotificationsResponse{
		Notifications: make([]*pb.Notification, 0, len(notifications)),
	}

	for _, notification := range notifications {
		response.Notifications = append(response.Notifications, convertNotificationToProto(&notification))
	}

	return response, nil
}

// Helper functions to convert between domain models and protobuf messages

func convertNotificationToProto(notification *repository.Notification) *pb.Notification {
	return &pb.Notification{
		Id:        notification.ID,
		TenantId:  notification.TenantID,
		UserId:    notification.UserID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		Priority:  notification.Priority,
		Status:    notification.Status,
		Metadata:  notification.Metadata,
		CreatedAt: timestamppb.New(notification.CreatedAt),
		UpdatedAt: timestamppb.New(notification.UpdatedAt),
	}
}

func convertAlertThresholdToProto(threshold *repository.AlertThreshold) *pb.AlertThreshold {
	return &pb.AlertThreshold{
		Id:             threshold.ID,
		TenantId:       threshold.TenantID,
		UserId:         threshold.UserID,
		Name:           threshold.Name,
		MetricType:     threshold.MetricType,
		ComparisonType: threshold.ComparisonType,
		Value:          threshold.Value,
		Percentage:     threshold.Percentage,
		Period:         threshold.Period,
		Status:         threshold.Status,
		CreatedAt:      timestamppb.New(threshold.CreatedAt),
		UpdatedAt:      timestamppb.New(threshold.UpdatedAt),
	}
}

func convertScheduledReportToProto(report *repository.ScheduledReport) *pb.ScheduledReport {
	return &pb.ScheduledReport{
		Id:             report.ID,
		TenantId:       report.TenantID,
		UserId:         report.UserID,
		Name:           report.Name,
		Description:    report.Description,
		ReportType:     report.ReportType,
		Schedule:       report.Schedule,
		Filters:        report.Filters,
		DeliveryType:   report.DeliveryType,
		EmailAddresses: report.EmailAddresses,
		Status:         report.Status,
		LastRunAt:      timestamppb.New(report.LastRunAt),
		NextRunAt:      timestamppb.New(report.NextRunAt),
		CreatedAt:      timestamppb.New(report.CreatedAt),
		UpdatedAt:      timestamppb.New(report.UpdatedAt),
	}
}
