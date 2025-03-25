package repository

import (
	"context"
	"errors"
	"time"

	"github.com/donaldnash/go-competitor/common/db"
)

// AuthRepository defines the interface for auth data access
type AuthRepository interface {
	// User authentication
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User, password string) (*User, error)
	ValidatePassword(ctx context.Context, userID string, password string) (bool, error)

	// Organization/Tenant management
	CreateOrganization(ctx context.Context, org *Organization) (*Organization, error)
	GetOrganization(ctx context.Context, orgID string) (*Organization, error)
	UpdateOrganization(ctx context.Context, org *Organization) (*Organization, error)
	DeleteOrganization(ctx context.Context, orgID string) error

	// User management
	GetUser(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	InviteUser(ctx context.Context, email, orgID, role string) (string, error)
	ListOrganizationUsers(ctx context.Context, orgID string) ([]User, error)

	// Token management
	CreateToken(ctx context.Context, userID string) (*Token, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	InvalidateToken(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*Token, error)
}

// User represents a user entity
type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	OrganizationID string    `json:"organization_id"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Organization represents an organization (tenant) entity
type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	AccountOwner string    `json:"account_owner"`
	Tier         string    `json:"tier"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Token represents an authentication token
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	ExpiresAt      int64  `json:"exp"`
}

// SupabaseAuthRepository implements AuthRepository using Supabase
type SupabaseAuthRepository struct {
	client *db.SupabaseClient
}

// NewSupabaseAuthRepository creates a new SupabaseAuthRepository
func NewSupabaseAuthRepository() (*SupabaseAuthRepository, error) {
	// For auth service, we use a system tenant ID since it needs to access data across tenants
	client, err := db.NewSupabaseClient("system")
	if err != nil {
		return nil, err
	}

	return &SupabaseAuthRepository{
		client: client,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (r *SupabaseAuthRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// In a real implementation, this would query Supabase Auth or the users table
	// For demonstration purposes, we'll return a placeholder error
	return nil, errors.New("user not found")
}

// CreateUser creates a new user
func (r *SupabaseAuthRepository) CreateUser(ctx context.Context, user *User, password string) (*User, error) {
	// In a real implementation, this would create a user in Supabase Auth
	// and store additional user data in the users table
	// For demonstration purposes, we'll return a placeholder implementation
	user.ID = "user-id"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user, nil
}

// ValidatePassword validates a user's password
func (r *SupabaseAuthRepository) ValidatePassword(ctx context.Context, userID string, password string) (bool, error) {
	// In a real implementation, this would validate the password with Supabase Auth
	// For demonstration purposes, we'll return true
	return true, nil
}

// CreateOrganization creates a new organization
func (r *SupabaseAuthRepository) CreateOrganization(ctx context.Context, org *Organization) (*Organization, error) {
	// In a real implementation, this would create an organization in the organizations table
	// For demonstration purposes, we'll return a placeholder implementation
	org.ID = "org-id"
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()
	return org, nil
}

// GetOrganization retrieves an organization by ID
func (r *SupabaseAuthRepository) GetOrganization(ctx context.Context, orgID string) (*Organization, error) {
	// In a real implementation, this would query the organizations table using Supabase
	// For development purposes, we'll return a dummy organization
	if orgID == "" {
		return nil, errors.New("organization ID cannot be empty")
	}

	// First check if we should try to query the actual database
	var orgs []Organization
	err := r.client.Query("organizations").
		Select("*").
		Where("id", "eq", orgID).
		Execute(&orgs)

	// If we have results from the database, return the first one
	if err == nil && len(orgs) > 0 {
		return &orgs[0], nil
	}

	// Otherwise, return a dummy organization with the provided ID
	// This is useful for development when the database might not be set up
	now := time.Now()
	return &Organization{
		ID:           orgID,
		Name:         "Organization " + orgID,
		AccountOwner: "default-owner",
		Tier:         "standard",
		CreatedAt:    now.Add(-24 * time.Hour), // Created yesterday
		UpdatedAt:    now,
	}, nil
}

// UpdateOrganization updates an organization
func (r *SupabaseAuthRepository) UpdateOrganization(ctx context.Context, org *Organization) (*Organization, error) {
	// In a real implementation, this would update the organization in the organizations table
	// For demonstration purposes, we'll return a placeholder implementation
	org.UpdatedAt = time.Now()
	return org, nil
}

// DeleteOrganization deletes an organization
func (r *SupabaseAuthRepository) DeleteOrganization(ctx context.Context, orgID string) error {
	// In a real implementation, this would delete the organization from the organizations table
	// For demonstration purposes, we'll return nil
	return nil
}

// GetUser retrieves a user by ID
func (r *SupabaseAuthRepository) GetUser(ctx context.Context, userID string) (*User, error) {
	// In a real implementation, this would query the users table using Supabase
	// For development purposes, we'll return a dummy user
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	// First check if we should try to query the actual database
	var users []User
	err := r.client.Query("users").
		Select("*").
		Where("id", "eq", userID).
		Execute(&users)

	// If we have results from the database, return the first one
	if err == nil && len(users) > 0 {
		return &users[0], nil
	}

	// Otherwise, return a dummy user with the provided ID
	// This is useful for development when the database might not be set up
	now := time.Now()
	return &User{
		ID:             userID,
		Email:          "user@example.com",
		FirstName:      "Test",
		LastName:       "User",
		OrganizationID: "org-" + userID, // Organization ID based on user ID
		Role:           "user",
		CreatedAt:      now.Add(-24 * time.Hour), // Created yesterday
		UpdatedAt:      now,
	}, nil
}

// UpdateUser updates a user
func (r *SupabaseAuthRepository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	// In a real implementation, this would update the user in the users table
	// For demonstration purposes, we'll return a placeholder implementation
	user.UpdatedAt = time.Now()
	return user, nil
}

// InviteUser invites a user to an organization
func (r *SupabaseAuthRepository) InviteUser(ctx context.Context, email, orgID, role string) (string, error) {
	// In a real implementation, this would create an invitation in the invitations table
	// and send an email to the user
	// For demonstration purposes, we'll return a placeholder invite ID
	return "invite-id", nil
}

// ListOrganizationUsers lists all users in an organization
func (r *SupabaseAuthRepository) ListOrganizationUsers(ctx context.Context, orgID string) ([]User, error) {
	// In a real implementation, this would query the users table for users with the given organization ID
	// For demonstration purposes, we'll return a placeholder empty slice
	return []User{}, nil
}

// CreateToken creates a new token for a user
func (r *SupabaseAuthRepository) CreateToken(ctx context.Context, userID string) (*Token, error) {
	// In a real implementation, this would create a JWT token with the user's claims
	// For demonstration purposes, we'll return a placeholder implementation
	return &Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

// ValidateToken validates a token and returns its claims
func (r *SupabaseAuthRepository) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	// In a real implementation, this would validate the JWT token and extract its claims
	// For demonstration purposes, we'll return a placeholder implementation
	return &TokenClaims{
		UserID:         "user-id",
		OrganizationID: "org-id",
		Role:           "admin",
		ExpiresAt:      time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

// InvalidateToken invalidates a token
func (r *SupabaseAuthRepository) InvalidateToken(ctx context.Context, token string) error {
	// In a real implementation, this would add the token to a blacklist or revoke it
	// For demonstration purposes, we'll return nil
	return nil
}

// RefreshToken refreshes a token
func (r *SupabaseAuthRepository) RefreshToken(ctx context.Context, refreshToken string) (*Token, error) {
	// In a real implementation, this would validate the refresh token and create a new access token
	// For demonstration purposes, we'll return a placeholder implementation
	return &Token{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}
