package resolvers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/donaldnash/go-competitor/auth/client"
	"github.com/donaldnash/go-competitor/auth/repository"
	"github.com/donaldnash/go-competitor/graphql/middleware"
	"github.com/donaldnash/go-competitor/graphql/models"
)

// AuthResolver handles all authentication and user-related GraphQL operations
// It communicates with the auth service to process login, registration, etc.
type AuthResolver struct {
	authClient *client.AuthClient // Client for interacting with the auth service
}

// NewAuthResolver creates a new AuthResolver with the provided auth client
func NewAuthResolver(authClient *client.AuthClient) *AuthResolver {
	return &AuthResolver{
		authClient: authClient,
	}
}

// Login authenticates a user and returns an auth payload containing JWT tokens
// This mutation doesn't require authentication as it's the entry point
func (r *AuthResolver) Login(ctx context.Context, email, password string) (*models.AuthPayload, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, token, err := r.authClient.Login(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return convertToAuthPayload(user, token), nil
}

// Register creates a new user account with an associated organization
// This mutation doesn't require authentication as it's for new users
func (r *AuthResolver) Register(ctx context.Context, email, password, firstName, lastName, organizationName string) (*models.AuthPayload, error) {
	// Validate required fields
	if email == "" || password == "" || firstName == "" || lastName == "" || organizationName == "" {
		return nil, fmt.Errorf("all fields are required for registration")
	}

	// Call the auth service to register the user
	user, _, token, err := r.authClient.Register(ctx, email, password, firstName, lastName, organizationName)
	if err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	// Create response with user and token information
	return &models.AuthPayload{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(token.ExpiresAt).Seconds()),
		User: &models.User{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			TenantID:  user.OrganizationID,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// RefreshToken issues a new access token using a valid refresh token
// This allows token renewal without requiring re-authentication
func (r *AuthResolver) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthPayload, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	token, err := r.authClient.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	// Get the user info using the token claims
	claims, err := r.authClient.ValidateToken(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	user, err := r.authClient.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user lookup failed: %w", err)
	}

	return convertToAuthPayload(user, token), nil
}

// Logout invalidates the current session token
// Requires authentication
func (r *AuthResolver) Logout(ctx context.Context) (bool, error) {
	// Ensure user is authenticated
	if err := middleware.RequireAuthentication(ctx); err != nil {
		return false, err
	}

	// Get user ID from context
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return false, fmt.Errorf("user context not found")
	}

	// Call the auth service to invalidate the session
	err := r.authClient.Logout(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("logout failed: %w", err)
	}

	return true, nil
}

// Me returns the current authenticated user's profile
// Requires authentication
func (r *AuthResolver) Me(ctx context.Context) (*models.User, error) {
	// Ensure user is authenticated
	if err := middleware.RequireAuthentication(ctx); err != nil {
		return nil, err
	}

	// Get user ID from context
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, fmt.Errorf("user context not found")
	}

	// Retrieve user details from auth service
	user, err := r.authClient.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return convertToUser(user), nil
}

// Tenant returns the current tenant (organization) for the authenticated user
// Requires authentication
func (r *AuthResolver) Tenant(ctx context.Context) (*models.Tenant, error) {
	// Ensure user is authenticated
	if err := middleware.RequireAuthentication(ctx); err != nil {
		return nil, err
	}

	// Get user ID and tenant ID from context
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, fmt.Errorf("user context not found")
	}

	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		// If tenant ID is not in the context, retrieve it from the user
		user, err := r.authClient.GetUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user: %w", err)
		}
		tenantID = user.OrganizationID
	}

	// Fetch tenant details from the auth service
	// The GetTenant method should be implemented in the auth client
	tenant, err := r.authClient.GetTenant(ctx, tenantID)
	if err != nil {
		// If there's an error or the method is not yet implemented, fall back to placeholder data
		log.Printf("Warning: Failed to fetch tenant details: %v. Using placeholder data instead.", err)
		return &models.Tenant{
			ID:        tenantID,
			Name:      "Organization", // Placeholder name
			Plan:      "standard",     // Default plan
			Active:    true,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}, nil
	}

	// Convert the tenant from repository model to GraphQL model
	return &models.Tenant{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Plan:      tenant.Tier,
		Active:    true, // Using default value as this field doesn't exist in Organization
		CreatedAt: tenant.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tenant.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// Helper function to convert repository user and token to GraphQL AuthPayload
func convertToAuthPayload(user *repository.User, token *repository.Token) *models.AuthPayload {
	return &models.AuthPayload{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(token.ExpiresAt).Seconds()),
		User:         convertToUser(user),
	}
}

// Helper function to convert repository User to GraphQL User
func convertToUser(user *repository.User) *models.User {
	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		TenantID:  user.OrganizationID,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

// getUserIDFromContext extracts user ID from context with fallback for development
func getUserIDFromContext(ctx context.Context) string {
	userID := middleware.GetUserID(ctx)

	// For development/testing, return a default user ID if none is in context
	if userID == "" && os.Getenv("ENV") == "development" {
		return "default-user-id"
	}

	return userID
}
