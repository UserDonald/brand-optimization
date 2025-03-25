package service

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/scraper/repository"
	"github.com/robfig/cron/v3"
)

// Platforms supported by the scraper
var supportedPlatforms = map[string]PlatformInfo{
	"instagram": {
		Name:        "instagram",
		DisplayName: "Instagram",
		Description: "Meta's photo and video sharing social network",
		SupportedJobTypes: []repository.JobType{
			repository.JobTypeProfile,
			repository.JobTypePosts,
			repository.JobTypeEngagement,
			repository.JobTypeFollowers,
		},
		RateLimits: PlatformRateLimits{
			RequestsPerMinute: 30,
			RequestsPerHour:   500,
			RequestsPerDay:    5000,
		},
	},
	"twitter": {
		Name:        "twitter",
		DisplayName: "Twitter / X",
		Description: "Short-form microblogging social network",
		SupportedJobTypes: []repository.JobType{
			repository.JobTypeProfile,
			repository.JobTypePosts,
			repository.JobTypeEngagement,
			repository.JobTypeFollowers,
		},
		RateLimits: PlatformRateLimits{
			RequestsPerMinute: 50,
			RequestsPerHour:   1500,
			RequestsPerDay:    10000,
		},
	},
	"facebook": {
		Name:        "facebook",
		DisplayName: "Facebook",
		Description: "Meta's social networking platform",
		SupportedJobTypes: []repository.JobType{
			repository.JobTypeProfile,
			repository.JobTypePosts,
			repository.JobTypeEngagement,
		},
		RateLimits: PlatformRateLimits{
			RequestsPerMinute: 20,
			RequestsPerHour:   200,
			RequestsPerDay:    2000,
		},
	},
	"linkedin": {
		Name:        "linkedin",
		DisplayName: "LinkedIn",
		Description: "Professional networking and career development platform",
		SupportedJobTypes: []repository.JobType{
			repository.JobTypeProfile,
			repository.JobTypePosts,
			repository.JobTypeEngagement,
		},
		RateLimits: PlatformRateLimits{
			RequestsPerMinute: 10,
			RequestsPerHour:   100,
			RequestsPerDay:    1000,
		},
	},
	"tiktok": {
		Name:        "tiktok",
		DisplayName: "TikTok",
		Description: "Short-form video hosting service",
		SupportedJobTypes: []repository.JobType{
			repository.JobTypeProfile,
			repository.JobTypePosts,
			repository.JobTypeEngagement,
		},
		RateLimits: PlatformRateLimits{
			RequestsPerMinute: 15,
			RequestsPerHour:   150,
			RequestsPerDay:    1500,
		},
	},
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

// ScraperService provides business logic for scraper operations
type ScraperService struct {
	repo        repository.ScraperRepository
	scheduler   *cron.Cron
	platformAPI map[string]PlatformAPI
}

// PlatformAPI defines the interface for platform-specific scrapers
type PlatformAPI interface {
	GetProfile(ctx context.Context, targetID string) (map[string]interface{}, error)
	GetPosts(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error)
	GetEngagement(ctx context.Context, targetID, postID string) (map[string]interface{}, error)
	GetFollowers(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error)
	GetComments(ctx context.Context, targetID, postID string, count int) ([]map[string]interface{}, error)
}

// NewScraperService creates a new ScraperService
func NewScraperService(repo repository.ScraperRepository) *ScraperService {
	scheduler := cron.New(cron.WithSeconds())

	// Start the scheduler
	scheduler.Start()

	// Initialize platform APIs
	platformAPI := make(map[string]PlatformAPI)

	// For now, we'll just have placeholder implementations
	platformAPI["instagram"] = &InstagramAPI{}
	platformAPI["twitter"] = &TwitterAPI{}
	platformAPI["facebook"] = &FacebookAPI{}
	platformAPI["linkedin"] = &LinkedInAPI{}
	platformAPI["tiktok"] = &TikTokAPI{}

	return &ScraperService{
		repo:        repo,
		scheduler:   scheduler,
		platformAPI: platformAPI,
	}
}

// GetSupportedPlatforms returns a list of supported platforms
func (s *ScraperService) GetSupportedPlatforms(ctx context.Context) []PlatformInfo {
	platforms := make([]PlatformInfo, 0, len(supportedPlatforms))
	for _, platform := range supportedPlatforms {
		platforms = append(platforms, platform)
	}
	return platforms
}

// GetPlatformStatus returns the status of a platform
func (s *ScraperService) GetPlatformStatus(ctx context.Context, platform string) (*PlatformStatus, error) {
	info, exists := supportedPlatforms[platform]
	if !exists {
		return nil, fmt.Errorf("platform not supported: %s", platform)
	}

	// In a real implementation, we would check the actual status of the platform
	// For now, we'll just return a mock status
	return &PlatformStatus{
		Platform:      platform,
		Available:     true,
		StatusMessage: "Platform is operational",
		RateLimits:    info.RateLimits,
		LastChecked:   time.Now(),
	}, nil
}

// CreateScraperJob creates a new scraper job
func (s *ScraperService) CreateScraperJob(ctx context.Context, tenantID, platform, targetID string,
	jobType repository.JobType, schedule repository.ScraperSchedule, metadata map[string]string) (*repository.ScraperJob, error) {

	// Validate platform
	if _, exists := supportedPlatforms[platform]; !exists {
		return nil, fmt.Errorf("platform not supported: %s", platform)
	}

	// Create the job
	job := &repository.ScraperJob{
		TenantID:  tenantID,
		Platform:  platform,
		TargetID:  targetID,
		JobType:   jobType,
		Status:    repository.JobStatusPending,
		Schedule:  schedule,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Calculate next run time based on schedule
	nextRun, err := s.calculateNextRunTime(schedule)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule: %w", err)
	}
	job.NextRunAt = nextRun

	// If the job should run immediately, set its status to scheduled
	if nextRun.Before(time.Now().Add(5 * time.Minute)) {
		job.Status = repository.JobStatusScheduled
	}

	// Save to repository
	job, err = s.repo.CreateScraperJob(ctx, job)
	if err != nil {
		return nil, err
	}

	// Schedule the job
	s.scheduleJob(job)

	return job, nil
}

// GetScraperJob retrieves a specific scraper job
func (s *ScraperService) GetScraperJob(ctx context.Context, tenantID, jobID string) (*repository.ScraperJob, error) {
	return s.repo.GetScraperJob(ctx, tenantID, jobID)
}

// GetScraperJobs retrieves all scraper jobs matching the filters
func (s *ScraperService) GetScraperJobs(ctx context.Context, tenantID, platform string,
	jobType repository.JobType, status repository.JobStatus) ([]repository.ScraperJob, error) {

	return s.repo.GetScraperJobs(ctx, tenantID, platform, jobType, status)
}

// CancelScraperJob cancels a scraper job
func (s *ScraperService) CancelScraperJob(ctx context.Context, tenantID, jobID string) (*repository.ScraperJob, error) {
	job, err := s.repo.GetScraperJob(ctx, tenantID, jobID)
	if err != nil {
		return nil, err
	}

	// Only pending or scheduled jobs can be cancelled
	if job.Status != repository.JobStatusPending && job.Status != repository.JobStatusScheduled {
		return nil, fmt.Errorf("job cannot be cancelled in its current state: %s", job.Status.String())
	}

	job.Status = repository.JobStatusCancelled
	job.UpdatedAt = time.Now()

	return s.repo.UpdateScraperJob(ctx, job)
}

// DeleteScraperJob deletes a scraper job
func (s *ScraperService) DeleteScraperJob(ctx context.Context, tenantID, jobID string) error {
	return s.repo.DeleteScraperJob(ctx, tenantID, jobID)
}

// GetScrapedData retrieves scraped data for a specific job
func (s *ScraperService) GetScrapedData(ctx context.Context, tenantID, jobID string,
	startDate, endDate time.Time) ([]repository.ScrapedDataItem, error) {

	return s.repo.GetScrapedData(ctx, tenantID, jobID, startDate, endDate)
}

// calculateNextRunTime calculates the next run time based on the schedule
func (s *ScraperService) calculateNextRunTime(schedule repository.ScraperSchedule) (time.Time, error) {
	// If there's a cron expression, use it
	if schedule.CronExpression != "" {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(schedule.CronExpression)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid cron expression: %w", err)
		}
		return schedule.Next(time.Now()), nil
	}

	// Otherwise, use the frequency
	now := time.Now()
	switch schedule.Frequency {
	case repository.FrequencyOnce:
		// If start date is provided and it's in the future, use it
		if !schedule.StartDate.IsZero() && schedule.StartDate.After(now) {
			return schedule.StartDate, nil
		}
		// Otherwise, run immediately
		return now, nil
	case repository.FrequencyHourly:
		return now.Add(1 * time.Hour), nil
	case repository.FrequencyDaily:
		return now.Add(24 * time.Hour), nil
	case repository.FrequencyWeekly:
		return now.Add(7 * 24 * time.Hour), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported frequency: %s", schedule.Frequency.String())
	}
}

// scheduleJob schedules a job to be executed
func (s *ScraperService) scheduleJob(job *repository.ScraperJob) {
	// In a real implementation, we would use the scheduler to run the job
	// For now, we'll just log that the job is scheduled
	fmt.Printf("Job %s scheduled to run at %s\n", job.ID, job.NextRunAt.Format(time.RFC3339))
}

// Placeholder implementations for platform APIs

// InstagramAPI is a placeholder implementation for Instagram
type InstagramAPI struct{}

func (a *InstagramAPI) GetProfile(ctx context.Context, targetID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": targetID, "platform": "instagram"}, nil
}

func (a *InstagramAPI) GetPosts(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "post1", "platform": "instagram"}}, nil
}

func (a *InstagramAPI) GetEngagement(ctx context.Context, targetID, postID string) (map[string]interface{}, error) {
	return map[string]interface{}{"likes": 100, "comments": 20}, nil
}

func (a *InstagramAPI) GetFollowers(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "user1", "platform": "instagram"}}, nil
}

func (a *InstagramAPI) GetComments(ctx context.Context, targetID, postID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "comment1", "text": "Great post!"}}, nil
}

// TwitterAPI is a placeholder implementation for Twitter
type TwitterAPI struct{}

func (a *TwitterAPI) GetProfile(ctx context.Context, targetID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": targetID, "platform": "twitter"}, nil
}

func (a *TwitterAPI) GetPosts(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "tweet1", "platform": "twitter"}}, nil
}

func (a *TwitterAPI) GetEngagement(ctx context.Context, targetID, postID string) (map[string]interface{}, error) {
	return map[string]interface{}{"likes": 50, "retweets": 10}, nil
}

func (a *TwitterAPI) GetFollowers(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "user1", "platform": "twitter"}}, nil
}

func (a *TwitterAPI) GetComments(ctx context.Context, targetID, postID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "reply1", "text": "Interesting!"}}, nil
}

// FacebookAPI is a placeholder implementation for Facebook
type FacebookAPI struct{}

func (a *FacebookAPI) GetProfile(ctx context.Context, targetID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": targetID, "platform": "facebook"}, nil
}

func (a *FacebookAPI) GetPosts(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "post1", "platform": "facebook"}}, nil
}

func (a *FacebookAPI) GetEngagement(ctx context.Context, targetID, postID string) (map[string]interface{}, error) {
	return map[string]interface{}{"likes": 200, "shares": 30}, nil
}

func (a *FacebookAPI) GetFollowers(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "user1", "platform": "facebook"}}, nil
}

func (a *FacebookAPI) GetComments(ctx context.Context, targetID, postID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "comment1", "text": "Awesome!"}}, nil
}

// LinkedInAPI is a placeholder implementation for LinkedIn
type LinkedInAPI struct{}

func (a *LinkedInAPI) GetProfile(ctx context.Context, targetID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": targetID, "platform": "linkedin"}, nil
}

func (a *LinkedInAPI) GetPosts(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "post1", "platform": "linkedin"}}, nil
}

func (a *LinkedInAPI) GetEngagement(ctx context.Context, targetID, postID string) (map[string]interface{}, error) {
	return map[string]interface{}{"likes": 150, "comments": 25}, nil
}

func (a *LinkedInAPI) GetFollowers(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "user1", "platform": "linkedin"}}, nil
}

func (a *LinkedInAPI) GetComments(ctx context.Context, targetID, postID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "comment1", "text": "Insightful!"}}, nil
}

// TikTokAPI is a placeholder implementation for TikTok
type TikTokAPI struct{}

func (a *TikTokAPI) GetProfile(ctx context.Context, targetID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": targetID, "platform": "tiktok"}, nil
}

func (a *TikTokAPI) GetPosts(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "video1", "platform": "tiktok"}}, nil
}

func (a *TikTokAPI) GetEngagement(ctx context.Context, targetID, postID string) (map[string]interface{}, error) {
	return map[string]interface{}{"likes": 500, "shares": 100}, nil
}

func (a *TikTokAPI) GetFollowers(ctx context.Context, targetID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "user1", "platform": "tiktok"}}, nil
}

func (a *TikTokAPI) GetComments(ctx context.Context, targetID, postID string, count int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": "comment1", "text": "Amazing video!"}}, nil
}
