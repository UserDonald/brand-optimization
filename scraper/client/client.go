package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/scraper/pb"
	"github.com/donaldnash/go-competitor/scraper/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScraperClient defines the interface for client communication with the scraper service
type ScraperClient interface {
	// Job management
	CreateScraperJob(ctx context.Context, tenantID, platform, targetID string, jobType repository.JobType, schedule repository.ScraperSchedule, metadata map[string]string) (*repository.ScraperJob, error)
	GetScraperJob(ctx context.Context, tenantID, jobID string) (*repository.ScraperJob, error)
	ListScraperJobs(ctx context.Context, tenantID, platform string, jobType repository.JobType, status repository.JobStatus) ([]repository.ScraperJob, error)
	CancelScraperJob(ctx context.Context, tenantID, jobID string) (*repository.ScraperJob, error)
	DeleteScraperJob(ctx context.Context, tenantID, jobID string) error

	// Platform operations
	ListSupportedPlatforms(ctx context.Context, tenantID string) ([]PlatformInfo, error)
	GetPlatformStatus(ctx context.Context, tenantID, platform string) (*PlatformStatus, error)

	// Scraper results
	GetScrapedData(ctx context.Context, tenantID, jobID string, startDate, endDate time.Time) ([]repository.ScrapedDataItem, error)

	// Close closes the client connection
	Close() error
}

// PlatformInfo contains information about a supported platform
type PlatformInfo struct {
	Name              string
	DisplayName       string
	Description       string
	SupportedJobTypes []repository.JobType
	RateLimits        PlatformRateLimits
}

// PlatformRateLimits contains rate limit information for a platform
type PlatformRateLimits struct {
	RequestsPerMinute int
	RequestsPerHour   int
	RequestsPerDay    int
	AvailableRequests int
	ResetAt           time.Time
}

// PlatformStatus contains the status of a platform
type PlatformStatus struct {
	Platform      string
	Available     bool
	StatusMessage string
	RateLimits    PlatformRateLimits
	LastChecked   time.Time
}

// GRPCScraperClient implements ScraperClient using gRPC
type GRPCScraperClient struct {
	conn   *grpc.ClientConn
	client pb.ScraperServiceClient
}

// NewGRPCScraperClient creates a new GRPCScraperClient
func NewGRPCScraperClient(serverAddr string) (*GRPCScraperClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to scraper service: %w", err)
	}

	client := pb.NewScraperServiceClient(conn)
	return &GRPCScraperClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *GRPCScraperClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateScraperJob creates a new scraper job
func (c *GRPCScraperClient) CreateScraperJob(ctx context.Context, tenantID, platform, targetID string,
	jobType repository.JobType, schedule repository.ScraperSchedule, metadata map[string]string) (*repository.ScraperJob, error) {

	// Convert job type to protobuf format
	protoJobType := convertJobTypeToProto(jobType)

	// Convert schedule to protobuf format
	protoSchedule := &pb.ScraperSchedule{
		CronExpression: schedule.CronExpression,
		Frequency:      convertScheduleFrequencyToProto(schedule.Frequency),
	}

	// Add timestamps if present
	if !schedule.StartDate.IsZero() {
		protoSchedule.StartDate = timestamppb.New(schedule.StartDate)
	}

	if !schedule.EndDate.IsZero() {
		protoSchedule.EndDate = timestamppb.New(schedule.EndDate)
	}

	// Create the request
	req := &pb.CreateScraperJobRequest{
		TenantId: tenantID,
		Platform: platform,
		TargetId: targetID,
		JobType:  protoJobType,
		Schedule: protoSchedule,
		Metadata: metadata,
	}

	// Call the service
	resp, err := c.client.CreateScraperJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create scraper job: %w", err)
	}

	// Convert the response to repository format
	return convertJobFromProto(resp), nil
}

// GetScraperJob retrieves a specific scraper job
func (c *GRPCScraperClient) GetScraperJob(ctx context.Context, tenantID, jobID string) (*repository.ScraperJob, error) {
	req := &pb.GetScraperJobRequest{
		TenantId: tenantID,
		JobId:    jobID,
	}

	resp, err := c.client.GetScraperJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get scraper job: %w", err)
	}

	return convertJobFromProto(resp), nil
}

// ListScraperJobs retrieves all scraper jobs matching the filters
func (c *GRPCScraperClient) ListScraperJobs(ctx context.Context, tenantID, platform string,
	jobType repository.JobType, status repository.JobStatus) ([]repository.ScraperJob, error) {

	req := &pb.ListScraperJobsRequest{
		TenantId: tenantID,
		Platform: platform,
		JobType:  convertJobTypeToProto(jobType),
		Status:   convertJobStatusToProto(status),
	}

	resp, err := c.client.ListScraperJobs(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list scraper jobs: %w", err)
	}

	// Convert the response to repository format
	jobs := make([]repository.ScraperJob, len(resp.Jobs))
	for i, job := range resp.Jobs {
		jobs[i] = *convertJobFromProto(job)
	}

	return jobs, nil
}

// CancelScraperJob cancels a scraper job
func (c *GRPCScraperClient) CancelScraperJob(ctx context.Context, tenantID, jobID string) (*repository.ScraperJob, error) {
	req := &pb.CancelScraperJobRequest{
		TenantId: tenantID,
		JobId:    jobID,
	}

	resp, err := c.client.CancelScraperJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel scraper job: %w", err)
	}

	return convertJobFromProto(resp), nil
}

// DeleteScraperJob deletes a scraper job
func (c *GRPCScraperClient) DeleteScraperJob(ctx context.Context, tenantID, jobID string) error {
	req := &pb.DeleteScraperJobRequest{
		TenantId: tenantID,
		JobId:    jobID,
	}

	_, err := c.client.DeleteScraperJob(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete scraper job: %w", err)
	}

	return nil
}

// ListSupportedPlatforms retrieves all supported platforms
func (c *GRPCScraperClient) ListSupportedPlatforms(ctx context.Context, tenantID string) ([]PlatformInfo, error) {
	req := &pb.ListSupportedPlatformsRequest{
		TenantId: tenantID,
	}

	resp, err := c.client.ListSupportedPlatforms(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list supported platforms: %w", err)
	}

	// Convert the response to repository format
	platforms := make([]PlatformInfo, len(resp.Platforms))
	for i, platform := range resp.Platforms {
		platforms[i] = PlatformInfo{
			Name:              platform.Name,
			DisplayName:       platform.DisplayName,
			Description:       platform.Description,
			SupportedJobTypes: convertJobTypesFromProto(platform.SupportedJobTypes),
			RateLimits: PlatformRateLimits{
				RequestsPerMinute: int(platform.RateLimits.RequestsPerMinute),
				RequestsPerHour:   int(platform.RateLimits.RequestsPerHour),
				RequestsPerDay:    int(platform.RateLimits.RequestsPerDay),
				AvailableRequests: int(platform.RateLimits.AvailableRequests),
			},
		}

		if platform.RateLimits.ResetAt != nil {
			platforms[i].RateLimits.ResetAt = platform.RateLimits.ResetAt.AsTime()
		}
	}

	return platforms, nil
}

// GetPlatformStatus retrieves the status of a platform
func (c *GRPCScraperClient) GetPlatformStatus(ctx context.Context, tenantID, platform string) (*PlatformStatus, error) {
	req := &pb.GetPlatformStatusRequest{
		TenantId: tenantID,
		Platform: platform,
	}

	resp, err := c.client.GetPlatformStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform status: %w", err)
	}

	// Convert the response to repository format
	status := &PlatformStatus{
		Platform:      resp.Platform,
		Available:     resp.Available,
		StatusMessage: resp.StatusMessage,
		LastChecked:   resp.LastChecked.AsTime(),
		RateLimits: PlatformRateLimits{
			RequestsPerMinute: int(resp.RateLimits.RequestsPerMinute),
			RequestsPerHour:   int(resp.RateLimits.RequestsPerHour),
			RequestsPerDay:    int(resp.RateLimits.RequestsPerDay),
			AvailableRequests: int(resp.RateLimits.AvailableRequests),
		},
	}

	if resp.RateLimits.ResetAt != nil {
		status.RateLimits.ResetAt = resp.RateLimits.ResetAt.AsTime()
	}

	return status, nil
}

// GetScrapedData retrieves scraped data for a specific job
func (c *GRPCScraperClient) GetScrapedData(ctx context.Context, tenantID, jobID string,
	startDate, endDate time.Time) ([]repository.ScrapedDataItem, error) {

	req := &pb.GetScrapedDataRequest{
		TenantId: tenantID,
		JobId:    jobID,
	}

	if !startDate.IsZero() {
		req.StartDate = timestamppb.New(startDate)
	}

	if !endDate.IsZero() {
		req.EndDate = timestamppb.New(endDate)
	}

	resp, err := c.client.GetScrapedData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get scraped data: %w", err)
	}

	// Convert the response to repository format
	items := make([]repository.ScrapedDataItem, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = *convertDataItemFromProto(item)
	}

	return items, nil
}

// Helper functions for type conversions

// convertJobTypeToProto converts a job type from repository to protobuf format
func convertJobTypeToProto(jobType repository.JobType) pb.ScraperJobType {
	switch jobType {
	case repository.JobTypeProfile:
		return pb.ScraperJobType_JOB_TYPE_PROFILE
	case repository.JobTypePosts:
		return pb.ScraperJobType_JOB_TYPE_POSTS
	case repository.JobTypeEngagement:
		return pb.ScraperJobType_JOB_TYPE_ENGAGEMENT
	case repository.JobTypeComments:
		return pb.ScraperJobType_JOB_TYPE_COMMENTS
	case repository.JobTypeFollowers:
		return pb.ScraperJobType_JOB_TYPE_FOLLOWERS
	default:
		return pb.ScraperJobType_JOB_TYPE_UNSPECIFIED
	}
}

// convertJobTypeFromProto converts a job type from protobuf to repository format
func convertJobTypeFromProto(jobType pb.ScraperJobType) repository.JobType {
	switch jobType {
	case pb.ScraperJobType_JOB_TYPE_PROFILE:
		return repository.JobTypeProfile
	case pb.ScraperJobType_JOB_TYPE_POSTS:
		return repository.JobTypePosts
	case pb.ScraperJobType_JOB_TYPE_ENGAGEMENT:
		return repository.JobTypeEngagement
	case pb.ScraperJobType_JOB_TYPE_COMMENTS:
		return repository.JobTypeComments
	case pb.ScraperJobType_JOB_TYPE_FOLLOWERS:
		return repository.JobTypeFollowers
	default:
		return repository.JobTypeUnspecified
	}
}

// convertJobTypesFromProto converts a slice of job types from protobuf to repository format
func convertJobTypesFromProto(jobTypes []pb.ScraperJobType) []repository.JobType {
	repoJobTypes := make([]repository.JobType, len(jobTypes))
	for i, jobType := range jobTypes {
		repoJobTypes[i] = convertJobTypeFromProto(jobType)
	}
	return repoJobTypes
}

// convertJobStatusToProto converts a job status from repository to protobuf format
func convertJobStatusToProto(status repository.JobStatus) pb.ScraperJobStatus {
	switch status {
	case repository.JobStatusPending:
		return pb.ScraperJobStatus_JOB_STATUS_PENDING
	case repository.JobStatusScheduled:
		return pb.ScraperJobStatus_JOB_STATUS_SCHEDULED
	case repository.JobStatusRunning:
		return pb.ScraperJobStatus_JOB_STATUS_RUNNING
	case repository.JobStatusCompleted:
		return pb.ScraperJobStatus_JOB_STATUS_COMPLETED
	case repository.JobStatusFailed:
		return pb.ScraperJobStatus_JOB_STATUS_FAILED
	case repository.JobStatusCancelled:
		return pb.ScraperJobStatus_JOB_STATUS_CANCELLED
	default:
		return pb.ScraperJobStatus_JOB_STATUS_UNSPECIFIED
	}
}

// convertJobStatusFromProto converts a job status from protobuf to repository format
func convertJobStatusFromProto(status pb.ScraperJobStatus) repository.JobStatus {
	switch status {
	case pb.ScraperJobStatus_JOB_STATUS_PENDING:
		return repository.JobStatusPending
	case pb.ScraperJobStatus_JOB_STATUS_SCHEDULED:
		return repository.JobStatusScheduled
	case pb.ScraperJobStatus_JOB_STATUS_RUNNING:
		return repository.JobStatusRunning
	case pb.ScraperJobStatus_JOB_STATUS_COMPLETED:
		return repository.JobStatusCompleted
	case pb.ScraperJobStatus_JOB_STATUS_FAILED:
		return repository.JobStatusFailed
	case pb.ScraperJobStatus_JOB_STATUS_CANCELLED:
		return repository.JobStatusCancelled
	default:
		return repository.JobStatusUnspecified
	}
}

// convertScheduleFrequencyToProto converts a schedule frequency from repository to protobuf format
func convertScheduleFrequencyToProto(frequency repository.ScheduleFrequency) pb.ScheduleFrequency {
	switch frequency {
	case repository.FrequencyOnce:
		return pb.ScheduleFrequency_FREQUENCY_ONCE
	case repository.FrequencyHourly:
		return pb.ScheduleFrequency_FREQUENCY_HOURLY
	case repository.FrequencyDaily:
		return pb.ScheduleFrequency_FREQUENCY_DAILY
	case repository.FrequencyWeekly:
		return pb.ScheduleFrequency_FREQUENCY_WEEKLY
	default:
		return pb.ScheduleFrequency_FREQUENCY_UNSPECIFIED
	}
}

// convertScheduleFrequencyFromProto converts a schedule frequency from protobuf to repository format
func convertScheduleFrequencyFromProto(frequency pb.ScheduleFrequency) repository.ScheduleFrequency {
	switch frequency {
	case pb.ScheduleFrequency_FREQUENCY_ONCE:
		return repository.FrequencyOnce
	case pb.ScheduleFrequency_FREQUENCY_HOURLY:
		return repository.FrequencyHourly
	case pb.ScheduleFrequency_FREQUENCY_DAILY:
		return repository.FrequencyDaily
	case pb.ScheduleFrequency_FREQUENCY_WEEKLY:
		return repository.FrequencyWeekly
	default:
		return repository.FrequencyUnspecified
	}
}

// convertDataTypeFromProto converts a data type from protobuf to repository format
func convertDataTypeFromProto(dataType pb.ScraperDataType) repository.DataType {
	switch dataType {
	case pb.ScraperDataType_DATA_TYPE_PROFILE:
		return repository.DataTypeProfile
	case pb.ScraperDataType_DATA_TYPE_POST:
		return repository.DataTypePost
	case pb.ScraperDataType_DATA_TYPE_STORY:
		return repository.DataTypeStory
	case pb.ScraperDataType_DATA_TYPE_COMMENT:
		return repository.DataTypeComment
	case pb.ScraperDataType_DATA_TYPE_FOLLOWER:
		return repository.DataTypeFollower
	default:
		return repository.DataTypeUnspecified
	}
}

// convertJobFromProto converts a job from protobuf to repository format
func convertJobFromProto(job *pb.ScraperJob) *repository.ScraperJob {
	if job == nil {
		return nil
	}

	repoJob := &repository.ScraperJob{
		ID:        job.Id,
		TenantID:  job.TenantId,
		Platform:  job.Platform,
		TargetID:  job.TargetId,
		JobType:   convertJobTypeFromProto(job.JobType),
		Status:    convertJobStatusFromProto(job.Status),
		LastError: job.LastError,
		RunCount:  int(job.RunCount),
		Metadata:  job.Metadata,
	}

	if job.Schedule != nil {
		repoJob.Schedule = repository.ScraperSchedule{
			CronExpression: job.Schedule.CronExpression,
			Frequency:      convertScheduleFrequencyFromProto(job.Schedule.Frequency),
		}

		if job.Schedule.StartDate != nil {
			repoJob.Schedule.StartDate = job.Schedule.StartDate.AsTime()
		}

		if job.Schedule.EndDate != nil {
			repoJob.Schedule.EndDate = job.Schedule.EndDate.AsTime()
		}
	}

	if job.LastRunAt != nil {
		repoJob.LastRunAt = job.LastRunAt.AsTime()
	}

	if job.NextRunAt != nil {
		repoJob.NextRunAt = job.NextRunAt.AsTime()
	}

	if job.CreatedAt != nil {
		repoJob.CreatedAt = job.CreatedAt.AsTime()
	}

	if job.UpdatedAt != nil {
		repoJob.UpdatedAt = job.UpdatedAt.AsTime()
	}

	return repoJob
}

// convertDataItemFromProto converts a data item from protobuf to repository format
func convertDataItemFromProto(item *pb.ScrapedDataItem) *repository.ScrapedDataItem {
	if item == nil {
		return nil
	}

	repoItem := &repository.ScrapedDataItem{
		ID:                item.Id,
		JobID:             item.JobId,
		TenantID:          item.TenantId,
		Platform:          item.Platform,
		TargetID:          item.TargetId,
		PostID:            item.PostId,
		DataType:          convertDataTypeFromProto(item.DataType),
		Likes:             int(item.Likes),
		Shares:            int(item.Shares),
		Comments:          int(item.Comments),
		CTR:               item.ClickThroughRate,
		AvgWatchTime:      item.AvgWatchTime,
		EngagementRate:    item.EngagementRate,
		ContentType:       item.ContentType,
		ContentURL:        item.ContentUrl,
		ContentAttributes: item.ContentAttributes,
	}

	if item.PostedAt != nil {
		repoItem.PostedAt = item.PostedAt.AsTime()
	}

	if item.ScrapedAt != nil {
		repoItem.ScrapedAt = item.ScrapedAt.AsTime()
	}

	if item.CreatedAt != nil {
		repoItem.CreatedAt = item.CreatedAt.AsTime()
	}

	return repoItem
}
