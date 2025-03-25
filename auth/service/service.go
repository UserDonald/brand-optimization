package service

import (
	"context"
	"errors"

	"github.com/donaldnash/go-competitor/auth/repository"
)

// AuthService provides business logic for authentication and authorization
type AuthService struct {
	repo repository.AuthRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, email, password string) (*repository.User, *repository.Token, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	// Validate password
	valid, err := s.repo.ValidatePassword(ctx, user.ID, password)
	if err != nil {
		return nil, nil, err
	}

	if !valid {
		return nil, nil, errors.New("invalid credentials")
	}

	// Create token
	token, err := s.repo.CreateToken(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

// Register creates a new user and organization
func (s *AuthService) Register(ctx context.Context, email, password, firstName, lastName, orgName string) (*repository.User, *repository.Organization, *repository.Token, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, nil, nil, errors.New("user already exists")
	}

	// Create organization
	org := &repository.Organization{
		Name:         orgName,
		AccountOwner: email,
		Tier:         "standard", // Default tier
	}

	createdOrg, err := s.repo.CreateOrganization(ctx, org)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create user
	user := &repository.User{
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		OrganizationID: createdOrg.ID,
		Role:           "admin", // First user is an admin
	}

	createdUser, err := s.repo.CreateUser(ctx, user, password)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create token
	token, err := s.repo.CreateToken(ctx, createdUser.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	return createdUser, createdOrg, token, nil
}

// Logout logs out a user
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	// Nothing to do for now, in a real implementation we might invalidate tokens
	return nil
}

// RefreshToken refreshes a token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*repository.Token, error) {
	return s.repo.RefreshToken(ctx, refreshToken)
}

// CreateOrganization creates a new organization
func (s *AuthService) CreateOrganization(ctx context.Context, name, accountOwnerID, tier string) (*repository.Organization, error) {
	org := &repository.Organization{
		Name:         name,
		AccountOwner: accountOwnerID,
		Tier:         tier,
	}

	return s.repo.CreateOrganization(ctx, org)
}

// GetOrganization retrieves an organization by ID
func (s *AuthService) GetOrganization(ctx context.Context, orgID string) (*repository.Organization, error) {
	return s.repo.GetOrganization(ctx, orgID)
}

// UpdateOrganization updates an organization
func (s *AuthService) UpdateOrganization(ctx context.Context, orgID, name, tier string) (*repository.Organization, error) {
	// First get the existing organization
	org, err := s.repo.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if name != "" {
		org.Name = name
	}
	if tier != "" {
		org.Tier = tier
	}

	return s.repo.UpdateOrganization(ctx, org)
}

// DeleteOrganization deletes an organization
func (s *AuthService) DeleteOrganization(ctx context.Context, orgID string) error {
	return s.repo.DeleteOrganization(ctx, orgID)
}

// ValidateToken validates a token and returns its claims
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*repository.TokenClaims, error) {
	return s.repo.ValidateToken(ctx, token)
}

// GetUser retrieves a user by ID
func (s *AuthService) GetUser(ctx context.Context, userID string) (*repository.User, error) {
	return s.repo.GetUser(ctx, userID)
}

// UpdateUser updates a user
func (s *AuthService) UpdateUser(ctx context.Context, userID, firstName, lastName, role string) (*repository.User, error) {
	// First get the existing user
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	if role != "" {
		user.Role = role
	}

	return s.repo.UpdateUser(ctx, user)
}

// InviteUser invites a user to an organization
func (s *AuthService) InviteUser(ctx context.Context, email, orgID, role string) (string, error) {
	return s.repo.InviteUser(ctx, email, orgID, role)
}

// ListOrganizationUsers lists all users in an organization
func (s *AuthService) ListOrganizationUsers(ctx context.Context, orgID string) ([]repository.User, error) {
	return s.repo.ListOrganizationUsers(ctx, orgID)
}
