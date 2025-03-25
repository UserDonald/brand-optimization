package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/common/db"
	"github.com/google/uuid"
)

// ScraperRepository defines the interface for scraper data access
type ScraperRepository interface {
	// Job management
	GetScraperJobs(ctx context.Context, tenantID, platform string, jobType JobType, status JobStatus) ([]ScraperJob, error)
	GetScraperJob(ctx context.Context, tenantID, jobID string) (*ScraperJob, error)
	CreateScraperJob(ctx context.Context, job *ScraperJob) (*ScraperJob, error)
	UpdateScraperJob(ctx context.Context, job *ScraperJob) (*ScraperJob, error)
	DeleteScraperJob(ctx context.Context, tenantID, jobID string) error

	// Data management
	GetScrapedData(ctx context.Context, tenantID, jobID string, startDate, endDate time.Time) ([]ScrapedDataItem, error)
	SaveScrapedData(ctx context.Context, tenantID string, data []ScrapedDataItem) (int, error)
}

// JobType represents the type of scraper job
type JobType int

const (
	JobTypeUnspecified JobType = iota
	JobTypeProfile
	JobTypePosts
	JobTypeEngagement
	JobTypeComments
	JobTypeFollowers
)

// String returns the string representation of JobType
func (j JobType) String() string {
	switch j {
	case JobTypeProfile:
		return "profile"
	case JobTypePosts:
		return "posts"
	case JobTypeEngagement:
		return "engagement"
	case JobTypeComments:
		return "comments"
	case JobTypeFollowers:
		return "followers"
	default:
		return "unspecified"
	}
}

// JobStatus represents the status of a scraper job
type JobStatus int

const (
	JobStatusUnspecified JobStatus = iota
	JobStatusPending
	JobStatusScheduled
	JobStatusRunning
	JobStatusCompleted
	JobStatusFailed
	JobStatusCancelled
)

// String returns the string representation of JobStatus
func (s JobStatus) String() string {
	switch s {
	case JobStatusPending:
		return "pending"
	case JobStatusScheduled:
		return "scheduled"
	case JobStatusRunning:
		return "running"
	case JobStatusCompleted:
		return "completed"
	case JobStatusFailed:
		return "failed"
	case JobStatusCancelled:
		return "cancelled"
	default:
		return "unspecified"
	}
}

// ScheduleFrequency represents the frequency of a scheduled job
type ScheduleFrequency int

const (
	FrequencyUnspecified ScheduleFrequency = iota
	FrequencyOnce
	FrequencyHourly
	FrequencyDaily
	FrequencyWeekly
)

// String returns the string representation of ScheduleFrequency
func (f ScheduleFrequency) String() string {
	switch f {
	case FrequencyOnce:
		return "once"
	case FrequencyHourly:
		return "hourly"
	case FrequencyDaily:
		return "daily"
	case FrequencyWeekly:
		return "weekly"
	default:
		return "unspecified"
	}
}

// DataType represents the type of scraped data
type DataType int

const (
	DataTypeUnspecified DataType = iota
	DataTypeProfile
	DataTypePost
	DataTypeStory
	DataTypeComment
	DataTypeFollower
)

// String returns the string representation of DataType
func (d DataType) String() string {
	switch d {
	case DataTypeProfile:
		return "profile"
	case DataTypePost:
		return "post"
	case DataTypeStory:
		return "story"
	case DataTypeComment:
		return "comment"
	case DataTypeFollower:
		return "follower"
	default:
		return "unspecified"
	}
}

// ScraperSchedule represents a schedule for a scraper job
type ScraperSchedule struct {
	CronExpression string            `json:"cron_expression"`
	Frequency      ScheduleFrequency `json:"frequency"`
	StartDate      time.Time         `json:"start_date"`
	EndDate        time.Time         `json:"end_date"`
}

// ScraperJob represents a scraper job
type ScraperJob struct {
	ID        string            `json:"id"`
	TenantID  string            `json:"tenant_id"`
	Platform  string            `json:"platform"`
	TargetID  string            `json:"target_id"`
	JobType   JobType           `json:"job_type"`
	Status    JobStatus         `json:"status"`
	Schedule  ScraperSchedule   `json:"schedule"`
	LastError string            `json:"last_error"`
	RunCount  int               `json:"run_count"`
	LastRunAt time.Time         `json:"last_run_at"`
	NextRunAt time.Time         `json:"next_run_at"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ScrapedDataItem represents a scraped data item
type ScrapedDataItem struct {
	ID                string            `json:"id"`
	JobID             string            `json:"job_id"`
	TenantID          string            `json:"tenant_id"`
	Platform          string            `json:"platform"`
	TargetID          string            `json:"target_id"`
	PostID            string            `json:"post_id"`
	DataType          DataType          `json:"data_type"`
	PostedAt          time.Time         `json:"posted_at"`
	Likes             int               `json:"likes"`
	Shares            int               `json:"shares"`
	Comments          int               `json:"comments"`
	CTR               float64           `json:"ctr"`
	AvgWatchTime      float64           `json:"avg_watch_time"`
	EngagementRate    float64           `json:"engagement_rate"`
	ContentType       string            `json:"content_type"`
	ContentURL        string            `json:"content_url"`
	ContentAttributes map[string]string `json:"content_attributes"`
	ScrapedAt         time.Time         `json:"scraped_at"`
	CreatedAt         time.Time         `json:"created_at"`
}

// SupabaseScraperRepository implements ScraperRepository using Supabase
type SupabaseScraperRepository struct {
	client *db.SupabaseClient
}

// NewSupabaseScraperRepository creates a new SupabaseScraperRepository
func NewSupabaseScraperRepository(tenantID string) (*SupabaseScraperRepository, error) {
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	client, err := db.NewSupabaseClient(tenantID)
	if err != nil {
		return nil, err
	}

	return &SupabaseScraperRepository{
		client: client,
	}, nil
}

// GetScraperJobs retrieves all scraper jobs matching the filters
func (r *SupabaseScraperRepository) GetScraperJobs(ctx context.Context, tenantID, platform string, jobType JobType, status JobStatus) ([]ScraperJob, error) {
	query := r.client.Query("scraper_jobs").Select("*")

	// Apply filters if provided
	if platform != "" {
		query = query.Where("platform", "eq", platform)
	}

	if jobType != JobTypeUnspecified {
		query = query.Where("job_type", "eq", jobType.String())
	}

	if status != JobStatusUnspecified {
		query = query.Where("status", "eq", status.String())
	}

	var jobs []ScraperJob
	err := query.Execute(&jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to get scraper jobs: %w", err)
	}
	return jobs, nil
}

// GetScraperJob retrieves a specific scraper job
func (r *SupabaseScraperRepository) GetScraperJob(ctx context.Context, tenantID, jobID string) (*ScraperJob, error) {
	var jobs []ScraperJob
	err := r.client.Query("scraper_jobs").
		Select("*").
		Where("id", "eq", jobID).
		Execute(&jobs)

	if err != nil {
		return nil, fmt.Errorf("failed to get scraper job: %w", err)
	}

	if len(jobs) == 0 {
		return nil, errors.New("scraper job not found")
	}

	return &jobs[0], nil
}

// CreateScraperJob creates a new scraper job
func (r *SupabaseScraperRepository) CreateScraperJob(ctx context.Context, job *ScraperJob) (*ScraperJob, error) {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}

	job.TenantID = r.client.TenantID
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
	job.UpdatedAt = time.Now()

	err := r.client.Insert(ctx, "scraper_jobs", job)
	if err != nil {
		return nil, fmt.Errorf("failed to create scraper job: %w", err)
	}

	return job, nil
}

// UpdateScraperJob updates an existing scraper job
func (r *SupabaseScraperRepository) UpdateScraperJob(ctx context.Context, job *ScraperJob) (*ScraperJob, error) {
	// Make sure the job belongs to the tenant
	existing, err := r.GetScraperJob(ctx, r.client.TenantID, job.ID)
	if err != nil {
		return nil, err
	}

	// Keep original tenant ID and created date
	job.TenantID = existing.TenantID
	job.CreatedAt = existing.CreatedAt
	job.UpdatedAt = time.Now()

	err = r.client.Update(ctx, "scraper_jobs", "id", job.ID, job)
	if err != nil {
		return nil, fmt.Errorf("failed to update scraper job: %w", err)
	}

	return job, nil
}

// DeleteScraperJob deletes a scraper job
func (r *SupabaseScraperRepository) DeleteScraperJob(ctx context.Context, tenantID, jobID string) error {
	// Verify the job exists and belongs to the tenant
	_, err := r.GetScraperJob(ctx, tenantID, jobID)
	if err != nil {
		return err
	}

	err = r.client.Delete(ctx, "scraper_jobs", "id", jobID)
	if err != nil {
		return fmt.Errorf("failed to delete scraper job: %w", err)
	}

	return nil
}

// GetScrapedData retrieves scraped data for a specific job within a date range
func (r *SupabaseScraperRepository) GetScrapedData(ctx context.Context, tenantID, jobID string, startDate, endDate time.Time) ([]ScrapedDataItem, error) {
	// First verify the job exists and belongs to the tenant
	_, err := r.GetScraperJob(ctx, tenantID, jobID)
	if err != nil {
		return nil, err
	}

	var items []ScrapedDataItem
	query := r.client.Query("scraped_data").
		Select("*").
		Where("job_id", "eq", jobID)

	if !startDate.IsZero() {
		query = query.Where("scraped_at", "gte", startDate.Format(time.RFC3339))
	}

	if !endDate.IsZero() {
		query = query.Where("scraped_at", "lte", endDate.Format(time.RFC3339))
	}

	err = query.Order("scraped_at", false).Execute(&items)

	if err != nil {
		return nil, fmt.Errorf("failed to get scraped data: %w", err)
	}

	return items, nil
}

// SaveScrapedData saves scraped data items
func (r *SupabaseScraperRepository) SaveScrapedData(ctx context.Context, tenantID string, data []ScrapedDataItem) (int, error) {
	// Insert each data item
	for i := range data {
		// Set IDs and ensure tenant ID is set
		if data[i].ID == "" {
			data[i].ID = uuid.New().String()
		}
		data[i].TenantID = tenantID

		if data[i].CreatedAt.IsZero() {
			data[i].CreatedAt = time.Now()
		}

		err := r.client.Insert(ctx, "scraped_data", data[i])
		if err != nil {
			return i, fmt.Errorf("failed to save data item at index %d: %w", i, err)
		}
	}

	return len(data), nil
}
