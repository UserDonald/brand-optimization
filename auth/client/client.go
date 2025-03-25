package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donaldnash/go-competitor/auth/pb"
	"github.com/donaldnash/go-competitor/auth/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient is the client for the Auth service
type AuthClient struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

// NewAuthClient creates a new AuthClient
func NewAuthClient(serverAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)

	return &AuthClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Login authenticates a user
func (c *AuthClient) Login(ctx context.Context, email, password string) (*repository.User, *repository.Token, error) {
	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to login: %w", err)
	}

	user := &repository.User{
		ID:             resp.User.Id,
		Email:          resp.User.Email,
		FirstName:      resp.User.FirstName,
		LastName:       resp.User.LastName,
		OrganizationID: resp.User.TenantId,
		Role:           resp.User.Role,
	}

	if resp.User.CreatedAt != nil {
		user.CreatedAt = resp.User.CreatedAt.AsTime()
	}

	if resp.User.UpdatedAt != nil {
		user.UpdatedAt = resp.User.UpdatedAt.AsTime()
	}

	expiresAt := time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	token := &repository.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
	}

	return user, token, nil
}

// Register registers a new user and organization
func (c *AuthClient) Register(ctx context.Context, email, password, firstName, lastName, orgName string) (*repository.User, *repository.Organization, *repository.Token, error) {
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Email:            email,
		Password:         password,
		FirstName:        firstName,
		LastName:         lastName,
		OrganizationName: orgName,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to register: %w", err)
	}

	user := &repository.User{
		ID:             resp.User.Id,
		Email:          resp.User.Email,
		FirstName:      resp.User.FirstName,
		LastName:       resp.User.LastName,
		OrganizationID: resp.User.TenantId,
		Role:           resp.User.Role,
	}

	if resp.User.CreatedAt != nil {
		user.CreatedAt = resp.User.CreatedAt.AsTime()
	}

	if resp.User.UpdatedAt != nil {
		user.UpdatedAt = resp.User.UpdatedAt.AsTime()
	}

	org := &repository.Organization{
		ID:   resp.Tenant.Id,
		Name: resp.Tenant.Name,
		Tier: resp.Tenant.Plan,
	}

	if resp.Tenant.CreatedAt != nil {
		org.CreatedAt = resp.Tenant.CreatedAt.AsTime()
	}

	if resp.Tenant.UpdatedAt != nil {
		org.UpdatedAt = resp.Tenant.UpdatedAt.AsTime()
	}

	expiresAt := time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	token := &repository.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
	}

	return user, org, token, nil
}

// Logout logs out a user
func (c *AuthClient) Logout(ctx context.Context, accessToken string) error {
	_, err := c.client.Logout(ctx, &pb.LogoutRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

// RefreshToken refreshes a token
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*repository.Token, error) {
	resp, err := c.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	token := &repository.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
	}

	return token, nil
}

// ValidateToken validates a token
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*repository.TokenClaims, error) {
	resp, err := c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		AccessToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	claims := &repository.TokenClaims{
		UserID:         resp.UserId,
		OrganizationID: resp.TenantId,
		Role:           resp.Role,
	}

	return claims, nil
}

// GetUser retrieves a user by ID
func (c *AuthClient) GetUser(ctx context.Context, userID string) (*repository.User, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := &repository.User{
		ID:             resp.Id,
		Email:          resp.Email,
		FirstName:      resp.FirstName,
		LastName:       resp.LastName,
		OrganizationID: resp.TenantId,
		Role:           resp.Role,
	}

	if resp.CreatedAt != nil {
		user.CreatedAt = resp.CreatedAt.AsTime()
	}

	if resp.UpdatedAt != nil {
		user.UpdatedAt = resp.UpdatedAt.AsTime()
	}

	return user, nil
}

// CreateOrganization creates a new organization
func (c *AuthClient) CreateOrganization(ctx context.Context, name, accountOwnerID, tier string) (*repository.Organization, error) {
	resp, err := c.client.CreateOrganization(ctx, &pb.CreateOrganizationRequest{
		Name:           name,
		AccountOwnerId: accountOwnerID,
		Plan:           tier,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	org := &repository.Organization{
		ID:   resp.Id,
		Name: resp.Name,
		Tier: resp.Plan,
	}

	if resp.CreatedAt != nil {
		org.CreatedAt = resp.CreatedAt.AsTime()
	}

	if resp.UpdatedAt != nil {
		org.UpdatedAt = resp.UpdatedAt.AsTime()
	}

	return org, nil
}

// CreateTenant creates a new tenant (organization)
func (c *AuthClient) CreateTenant(ctx context.Context, name, plan string, metadata map[string]string) (*repository.Organization, error) {
	resp, err := c.client.CreateTenant(ctx, &pb.CreateTenantRequest{
		Name:     name,
		Plan:     plan,
		Metadata: metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	org := &repository.Organization{
		ID:   resp.Id,
		Name: resp.Name,
		Tier: resp.Plan,
	}

	if resp.CreatedAt != nil {
		org.CreatedAt = resp.CreatedAt.AsTime()
	}

	if resp.UpdatedAt != nil {
		org.UpdatedAt = resp.UpdatedAt.AsTime()
	}

	return org, nil
}

// HasPermission checks if a user has a permission
func (c *AuthClient) HasPermission(ctx context.Context, userID, permission, resourceID string) (bool, error) {
	resp, err := c.client.HasPermission(ctx, &pb.HasPermissionRequest{
		UserId:     userID,
		Permission: permission,
		ResourceId: resourceID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return resp.HasPermission, nil
}

// GetUserPermissions gets all permissions for a user
func (c *AuthClient) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	resp, err := c.client.GetUserPermissions(ctx, &pb.GetUserPermissionsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	return resp.Permissions, nil
}

// GetTenant retrieves tenant details by ID
func (c *AuthClient) GetTenant(ctx context.Context, tenantID string) (*repository.Organization, error) {
	resp, err := c.client.GetTenant(ctx, &pb.GetTenantRequest{
		TenantId: tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Create organization with the fields we're sure about
	org := &repository.Organization{
		ID:   resp.Id,
		Name: resp.Name,
		Tier: resp.Plan,
	}

	if resp.CreatedAt != nil {
		org.CreatedAt = resp.CreatedAt.AsTime()
	}

	if resp.UpdatedAt != nil {
		org.UpdatedAt = resp.UpdatedAt.AsTime()
	}

	return org, nil
}
