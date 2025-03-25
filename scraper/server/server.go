package server

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/scraper/pb"
	"github.com/donaldnash/go-competitor/scraper/repository"
	"github.com/donaldnash/go-competitor/scraper/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScraperServer implements the gRPC server for scraper service
type ScraperServer struct {
	pb.UnimplementedScraperServiceServer
	service *service.ScraperService
}

// NewScraperServer creates a new ScraperServer
func NewScraperServer(service *service.ScraperService) *ScraperServer {
	return &ScraperServer{
		service: service,
	}
}

// CreateScraperJob handles the CreateScraperJob RPC call
func (s *ScraperServer) CreateScraperJob(ctx context.Context, req *pb.CreateScraperJobRequest) (*pb.ScraperJob, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.Platform == "" {
		return nil, status.Error(codes.InvalidArgument, "platform is required")
	}

	if req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "target ID is required")
	}

	if req.JobType == pb.ScraperJobType_JOB_TYPE_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "job type is required")
	}

	// Convert job type from protobuf to repository
	jobType := convertJobTypeFromProto(req.JobType)

	// Convert schedule from protobuf to repository
	schedule := repository.ScraperSchedule{}
	if req.Schedule != nil {
		schedule = repository.ScraperSchedule{
			CronExpression: req.Schedule.CronExpression,
			Frequency:      convertScheduleFrequencyFromProto(req.Schedule.Frequency),
		}

		// Convert timestamps if present
		if req.Schedule.StartDate != nil {
			schedule.StartDate = req.Schedule.StartDate.AsTime()
		}

		if req.Schedule.EndDate != nil {
			schedule.EndDate = req.Schedule.EndDate.AsTime()
		}
	}

	// Create the job
	job, err := s.service.CreateScraperJob(ctx, req.TenantId, req.Platform, req.TargetId, jobType, schedule, req.Metadata)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert the job to protobuf format
	return convertJobToProto(job), nil
}

// GetScraperJob handles the GetScraperJob RPC call
func (s *ScraperServer) GetScraperJob(ctx context.Context, req *pb.GetScraperJobRequest) (*pb.ScraperJob, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.JobId == "" {
		return nil, status.Error(codes.InvalidArgument, "job ID is required")
	}

	job, err := s.service.GetScraperJob(ctx, req.TenantId, req.JobId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertJobToProto(job), nil
}

// ListScraperJobs handles the ListScraperJobs RPC call
func (s *ScraperServer) ListScraperJobs(ctx context.Context, req *pb.ListScraperJobsRequest) (*pb.ListScraperJobsResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	// Convert job type and status from protobuf to repository
	jobType := convertJobTypeFromProto(req.JobType)
	jobStatus := convertJobStatusFromProto(req.Status)

	jobs, err := s.service.GetScraperJobs(ctx, req.TenantId, req.Platform, jobType, jobStatus)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert jobs to protobuf format
	protoJobs := make([]*pb.ScraperJob, len(jobs))
	for i, job := range jobs {
		protoJobs[i] = convertJobToProto(&job)
	}

	return &pb.ListScraperJobsResponse{
		Jobs: protoJobs,
	}, nil
}

// CancelScraperJob handles the CancelScraperJob RPC call
func (s *ScraperServer) CancelScraperJob(ctx context.Context, req *pb.CancelScraperJobRequest) (*pb.ScraperJob, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.JobId == "" {
		return nil, status.Error(codes.InvalidArgument, "job ID is required")
	}

	job, err := s.service.CancelScraperJob(ctx, req.TenantId, req.JobId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertJobToProto(job), nil
}

// DeleteScraperJob handles the DeleteScraperJob RPC call
func (s *ScraperServer) DeleteScraperJob(ctx context.Context, req *pb.DeleteScraperJobRequest) (*emptypb.Empty, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.JobId == "" {
		return nil, status.Error(codes.InvalidArgument, "job ID is required")
	}

	err := s.service.DeleteScraperJob(ctx, req.TenantId, req.JobId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// ListSupportedPlatforms handles the ListSupportedPlatforms RPC call
func (s *ScraperServer) ListSupportedPlatforms(ctx context.Context, req *pb.ListSupportedPlatformsRequest) (*pb.ListSupportedPlatformsResponse, error) {
	platforms := s.service.GetSupportedPlatforms(ctx)

	// Convert platforms to protobuf format
	protoPlatforms := make([]*pb.PlatformInfo, len(platforms))
	for i, platform := range platforms {
		protoPlatforms[i] = &pb.PlatformInfo{
			Name:              platform.Name,
			DisplayName:       platform.DisplayName,
			Description:       platform.Description,
			SupportedJobTypes: convertJobTypesToProto(platform.SupportedJobTypes),
			RateLimits: &pb.PlatformRateLimits{
				RequestsPerMinute: int32(platform.RateLimits.RequestsPerMinute),
				RequestsPerHour:   int32(platform.RateLimits.RequestsPerHour),
				RequestsPerDay:    int32(platform.RateLimits.RequestsPerDay),
				AvailableRequests: int32(platform.RateLimits.AvailableRequests),
			},
		}

		if !platform.RateLimits.ResetAt.IsZero() {
			protoPlatforms[i].RateLimits.ResetAt = timestamppb.New(platform.RateLimits.ResetAt)
		}
	}

	return &pb.ListSupportedPlatformsResponse{
		Platforms: protoPlatforms,
	}, nil
}

// GetPlatformStatus handles the GetPlatformStatus RPC call
func (s *ScraperServer) GetPlatformStatus(ctx context.Context, req *pb.GetPlatformStatusRequest) (*pb.PlatformStatus, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.Platform == "" {
		return nil, status.Error(codes.InvalidArgument, "platform is required")
	}

	platformStatus, err := s.service.GetPlatformStatus(ctx, req.Platform)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert status to protobuf format
	protoStatus := &pb.PlatformStatus{
		Platform:      platformStatus.Platform,
		Available:     platformStatus.Available,
		StatusMessage: platformStatus.StatusMessage,
		LastChecked:   timestamppb.New(platformStatus.LastChecked),
		RateLimits: &pb.PlatformRateLimits{
			RequestsPerMinute: int32(platformStatus.RateLimits.RequestsPerMinute),
			RequestsPerHour:   int32(platformStatus.RateLimits.RequestsPerHour),
			RequestsPerDay:    int32(platformStatus.RateLimits.RequestsPerDay),
			AvailableRequests: int32(platformStatus.RateLimits.AvailableRequests),
		},
	}

	if !platformStatus.RateLimits.ResetAt.IsZero() {
		protoStatus.RateLimits.ResetAt = timestamppb.New(platformStatus.RateLimits.ResetAt)
	}

	return protoStatus, nil
}

// GetScrapedData handles the GetScrapedData RPC call
func (s *ScraperServer) GetScrapedData(ctx context.Context, req *pb.GetScrapedDataRequest) (*pb.GetScrapedDataResponse, error) {
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant ID is required")
	}

	if req.JobId == "" {
		return nil, status.Error(codes.InvalidArgument, "job ID is required")
	}

	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate = req.StartDate.AsTime()
	}

	if req.EndDate != nil {
		endDate = req.EndDate.AsTime()
	}

	items, err := s.service.GetScrapedData(ctx, req.TenantId, req.JobId, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert items to protobuf format
	protoItems := make([]*pb.ScrapedDataItem, len(items))
	for i, item := range items {
		protoItems[i] = convertDataItemToProto(&item)
	}

	return &pb.GetScrapedDataResponse{
		Items: protoItems,
	}, nil
}

// Helper functions for type conversions

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

// convertJobTypesToProto converts a slice of job types from repository to protobuf format
func convertJobTypesToProto(jobTypes []repository.JobType) []pb.ScraperJobType {
	protoJobTypes := make([]pb.ScraperJobType, len(jobTypes))
	for i, jobType := range jobTypes {
		protoJobTypes[i] = convertJobTypeToProto(jobType)
	}
	return protoJobTypes
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

// convertDataTypeToProto converts a data type from repository to protobuf format
func convertDataTypeToProto(dataType repository.DataType) pb.ScraperDataType {
	switch dataType {
	case repository.DataTypeProfile:
		return pb.ScraperDataType_DATA_TYPE_PROFILE
	case repository.DataTypePost:
		return pb.ScraperDataType_DATA_TYPE_POST
	case repository.DataTypeStory:
		return pb.ScraperDataType_DATA_TYPE_STORY
	case repository.DataTypeComment:
		return pb.ScraperDataType_DATA_TYPE_COMMENT
	case repository.DataTypeFollower:
		return pb.ScraperDataType_DATA_TYPE_FOLLOWER
	default:
		return pb.ScraperDataType_DATA_TYPE_UNSPECIFIED
	}
}

// convertJobToProto converts a job from repository to protobuf format
func convertJobToProto(job *repository.ScraperJob) *pb.ScraperJob {
	if job == nil {
		return nil
	}

	protoJob := &pb.ScraperJob{
		Id:       job.ID,
		TenantId: job.TenantID,
		Platform: job.Platform,
		TargetId: job.TargetID,
		JobType:  convertJobTypeToProto(job.JobType),
		Status:   convertJobStatusToProto(job.Status),
		Schedule: &pb.ScraperSchedule{
			CronExpression: job.Schedule.CronExpression,
			Frequency:      convertScheduleFrequencyToProto(job.Schedule.Frequency),
		},
		LastError: job.LastError,
		RunCount:  int32(job.RunCount),
		Metadata:  job.Metadata,
	}

	// Convert timestamps if present
	if !job.LastRunAt.IsZero() {
		protoJob.LastRunAt = timestamppb.New(job.LastRunAt)
	}

	if !job.NextRunAt.IsZero() {
		protoJob.NextRunAt = timestamppb.New(job.NextRunAt)
	}

	if !job.Schedule.StartDate.IsZero() {
		protoJob.Schedule.StartDate = timestamppb.New(job.Schedule.StartDate)
	}

	if !job.Schedule.EndDate.IsZero() {
		protoJob.Schedule.EndDate = timestamppb.New(job.Schedule.EndDate)
	}

	if !job.CreatedAt.IsZero() {
		protoJob.CreatedAt = timestamppb.New(job.CreatedAt)
	}

	if !job.UpdatedAt.IsZero() {
		protoJob.UpdatedAt = timestamppb.New(job.UpdatedAt)
	}

	return protoJob
}

// convertDataItemToProto converts a data item from repository to protobuf format
func convertDataItemToProto(item *repository.ScrapedDataItem) *pb.ScrapedDataItem {
	if item == nil {
		return nil
	}

	protoItem := &pb.ScrapedDataItem{
		Id:                item.ID,
		JobId:             item.JobID,
		TenantId:          item.TenantID,
		Platform:          item.Platform,
		TargetId:          item.TargetID,
		PostId:            item.PostID,
		DataType:          convertDataTypeToProto(item.DataType),
		Likes:             int32(item.Likes),
		Shares:            int32(item.Shares),
		Comments:          int32(item.Comments),
		ClickThroughRate:  item.CTR,
		AvgWatchTime:      item.AvgWatchTime,
		EngagementRate:    item.EngagementRate,
		ContentType:       item.ContentType,
		ContentUrl:        item.ContentURL,
		ContentAttributes: item.ContentAttributes,
	}

	// Convert timestamps if present
	if !item.PostedAt.IsZero() {
		protoItem.PostedAt = timestamppb.New(item.PostedAt)
	}

	if !item.ScrapedAt.IsZero() {
		protoItem.ScrapedAt = timestamppb.New(item.ScrapedAt)
	}

	if !item.CreatedAt.IsZero() {
		protoItem.CreatedAt = timestamppb.New(item.CreatedAt)
	}

	return protoItem
}
