package resolvers

import (
	"context"
)

// RootResolver combines all service resolvers
type RootResolver struct {
	AuthResolver       *AuthResolver
	CompetitorResolver *CompetitorResolver
	AudienceResolver   *AudienceResolver
	ContentResolver    *ContentResolver
	AnalyticsResolver  *AnalyticsResolver
}

// NewRootResolver creates a new RootResolver
func NewRootResolver(
	authResolver *AuthResolver,
	competitorResolver *CompetitorResolver,
	audienceResolver *AudienceResolver,
	contentResolver *ContentResolver,
	analyticsResolver *AnalyticsResolver,
) *RootResolver {
	return &RootResolver{
		AuthResolver:       authResolver,
		CompetitorResolver: competitorResolver,
		AudienceResolver:   audienceResolver,
		ContentResolver:    contentResolver,
		AnalyticsResolver:  analyticsResolver,
	}
}

// These methods implement the root resolver operations

// Query handlers

// GetCompetitors handles the competitors query
func (r *RootResolver) GetCompetitors(ctx context.Context, tenantID string) (interface{}, error) {
	return r.CompetitorResolver.GetCompetitors(ctx, tenantID)
}

// GetCompetitor handles the competitor query
func (r *RootResolver) GetCompetitor(ctx context.Context, tenantID, id string) (interface{}, error) {
	return r.CompetitorResolver.GetCompetitor(ctx, tenantID, id)
}

// GetAudienceSegments handles the audienceSegments query
func (r *RootResolver) GetAudienceSegments(ctx context.Context, tenantID string) (interface{}, error) {
	return r.AudienceResolver.GetAudienceSegments(ctx, tenantID)
}

// GetAudienceSegment handles the audienceSegment query
func (r *RootResolver) GetAudienceSegment(ctx context.Context, tenantID, id string) (interface{}, error) {
	return r.AudienceResolver.GetAudienceSegment(ctx, tenantID, id)
}

// GetContentFormats handles the contentFormats query
func (r *RootResolver) GetContentFormats(ctx context.Context, tenantID string) (interface{}, error) {
	return r.ContentResolver.GetContentFormats(ctx, tenantID)
}

// GetContentFormat handles the contentFormat query
func (r *RootResolver) GetContentFormat(ctx context.Context, tenantID, id string) (interface{}, error) {
	return r.ContentResolver.GetContentFormat(ctx, tenantID, id)
}

// GetScheduledPosts handles the scheduledPosts query
func (r *RootResolver) GetScheduledPosts(ctx context.Context, tenantID string) (interface{}, error) {
	return r.ContentResolver.GetScheduledPosts(ctx, tenantID)
}

// GetScheduledPost handles the scheduledPost query
func (r *RootResolver) GetScheduledPost(ctx context.Context, tenantID, id string) (interface{}, error) {
	return r.ContentResolver.GetScheduledPost(ctx, tenantID, id)
}

// GetRecommendedPostingTimes handles the recommendedPostingTimes query
func (r *RootResolver) GetRecommendedPostingTimes(ctx context.Context, tenantID, dayOfWeek string) (interface{}, error) {
	return r.AnalyticsResolver.GetRecommendedPostingTimes(ctx, tenantID, dayOfWeek)
}

// GetRecommendedContentFormats handles the recommendedContentFormats query
func (r *RootResolver) GetRecommendedContentFormats(ctx context.Context, tenantID string) (interface{}, error) {
	return r.AnalyticsResolver.GetRecommendedContentFormats(ctx, tenantID)
}

// PredictPostEngagement handles the predictPostEngagement query
func (r *RootResolver) PredictPostEngagement(ctx context.Context, tenantID, contentFormat, scheduledTime string) (interface{}, error) {
	return r.AnalyticsResolver.PredictPostEngagement(ctx, tenantID, contentFormat, scheduledTime)
}

// GetContentPerformanceAnalysis handles the contentPerformanceAnalysis query
func (r *RootResolver) GetContentPerformanceAnalysis(ctx context.Context, tenantID, startDate, endDate string) (interface{}, error) {
	return r.AnalyticsResolver.GetContentPerformanceAnalysis(ctx, tenantID, startDate, endDate)
}

// Auth Query handlers

// Me handles the me query
func (r *RootResolver) Me(ctx context.Context) (interface{}, error) {
	return r.AuthResolver.Me(ctx)
}

// Tenant handles the tenant query
func (r *RootResolver) Tenant(ctx context.Context) (interface{}, error) {
	return r.AuthResolver.Tenant(ctx)
}

// Auth Mutation handlers

// Login handles the login mutation
func (r *RootResolver) Login(ctx context.Context, args struct {
	Email    string
	Password string
}) (interface{}, error) {
	return r.AuthResolver.Login(ctx, args.Email, args.Password)
}

// Register handles the register mutation
func (r *RootResolver) Register(ctx context.Context, args struct {
	Email            string
	Password         string
	FirstName        string
	LastName         string
	OrganizationName string
}) (interface{}, error) {
	return r.AuthResolver.Register(ctx, args.Email, args.Password, args.FirstName, args.LastName, args.OrganizationName)
}

// RefreshToken handles the refreshToken mutation
func (r *RootResolver) RefreshToken(ctx context.Context, args struct {
	RefreshToken string
}) (interface{}, error) {
	return r.AuthResolver.RefreshToken(ctx, args.RefreshToken)
}

// Logout handles the logout mutation
func (r *RootResolver) Logout(ctx context.Context) (bool, error) {
	return r.AuthResolver.Logout(ctx)
}

// Health returns a simple health check status
func (r *RootResolver) Health() string {
	return "OK"
}

// Ping returns a simple ping response
func (r *RootResolver) Ping() string {
	return "pong"
}
