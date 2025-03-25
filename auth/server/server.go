package server

import (
	"context"
	"time"

	"github.com/donaldnash/go-competitor/auth/pb"
	"github.com/donaldnash/go-competitor/auth/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthServer implements the AuthService gRPC server
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	service *service.AuthService
}

// NewAuthServer creates a new AuthServer
func NewAuthServer(service *service.AuthService) *AuthServer {
	return &AuthServer{
		service: service,
	}
}

// Login handles the Login RPC call
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Call the service
	user, token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Calculate expiration time in seconds
	expiresIn := int32(time.Until(token.ExpiresAt).Seconds())

	// Convert to protobuf response
	return &pb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User: &pb.User{
			Id:        user.ID,
			TenantId:  user.OrganizationID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Active:    true,
		},
	}, nil
}

// Register handles the Register RPC call
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first_name is required")
	}

	if req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "last_name is required")
	}

	if req.OrganizationName == "" {
		return nil, status.Error(codes.InvalidArgument, "organization_name is required")
	}

	// Call the service
	user, org, token, err := s.service.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName, req.OrganizationName)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	// Calculate expiration time in seconds
	expiresIn := int32(time.Until(token.ExpiresAt).Seconds())

	// Convert to protobuf response
	return &pb.RegisterResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User: &pb.User{
			Id:        user.ID,
			TenantId:  user.OrganizationID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Active:    true,
		},
		Tenant: &pb.Tenant{
			Id:       org.ID,
			Name:     org.Name,
			Plan:     org.Tier,
			Active:   true,
			Metadata: make(map[string]string),
		},
	}, nil
}

// Logout handles the Logout RPC call
func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	// Extract token from request
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	// For now, we'll just return success
	// In a real implementation, we would extract the user ID from the token and call the service

	return &emptypb.Empty{}, nil
}

// RefreshToken handles the RefreshToken RPC call
func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Validate request
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// Call the service
	token, err := s.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Calculate expiration time in seconds
	expiresIn := int32(time.Until(token.ExpiresAt).Seconds())

	// Convert to protobuf response
	return &pb.RefreshTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}

// GetUser handles the GetUser RPC call
func (s *AuthServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Call the service
	user, err := s.service.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Convert to protobuf response
	return &pb.User{
		Id:        user.ID,
		TenantId:  user.OrganizationID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Active:    true,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// CreateTenant handles the CreateTenant RPC call
func (s *AuthServer) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.Tenant, error) {
	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// Use service to create organization (tenant)
	// We'll use a default owner ID for now (in a real implementation, this would come from auth context)
	ownerID := "system"

	org, err := s.service.CreateOrganization(ctx, req.Name, ownerID, req.Plan)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
	return &pb.Tenant{
		Id:        org.ID,
		Name:      org.Name,
		Plan:      org.Tier,
		Active:    true,
		Metadata:  make(map[string]string),
		CreatedAt: timestamppb.New(org.CreatedAt),
		UpdatedAt: timestamppb.New(org.UpdatedAt),
	}, nil
}

// CreateOrganization handles the CreateOrganization RPC call
func (s *AuthServer) CreateOrganization(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.Tenant, error) {
	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.AccountOwnerId == "" {
		return nil, status.Error(codes.InvalidArgument, "account_owner_id is required")
	}

	// Call the service
	org, err := s.service.CreateOrganization(ctx, req.Name, req.AccountOwnerId, req.Plan)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
	return &pb.Tenant{
		Id:        org.ID,
		Name:      org.Name,
		Plan:      org.Tier,
		Active:    true,
		Metadata:  make(map[string]string),
		CreatedAt: timestamppb.New(org.CreatedAt),
		UpdatedAt: timestamppb.New(org.UpdatedAt),
	}, nil
}

// ValidateToken handles the ValidateToken RPC call
func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// Validate request
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Call the service
	claims, err := s.service.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Convert to protobuf response
	return &pb.ValidateTokenResponse{
		Valid:    true,
		UserId:   claims.UserID,
		TenantId: claims.OrganizationID,
		Role:     claims.Role,
	}, nil
}

// Implement stub methods for the remaining methods to satisfy the interface

// CreateUser handles the CreateUser RPC call
func (s *AuthServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

// UpdateUser handles the UpdateUser RPC call
func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

// DeleteUser handles the DeleteUser RPC call
func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, status.Error(codes.Unimplemented, "method not implemented")
}

// GetTenant handles the GetTenant RPC call
func (s *AuthServer) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.Tenant, error) {
	// Validate request
	if req.TenantId == "" {
		return nil, status.Error(codes.InvalidArgument, "tenant_id is required")
	}

	// Call the service
	org, err := s.service.GetOrganization(ctx, req.TenantId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Convert to protobuf response
	return &pb.Tenant{
		Id:        org.ID,
		Name:      org.Name,
		Plan:      org.Tier,
		Active:    true, // Default to active since this field isn't in the Organization model
		Metadata:  make(map[string]string),
		CreatedAt: timestamppb.New(org.CreatedAt),
		UpdatedAt: timestamppb.New(org.UpdatedAt),
	}, nil
}

// ListTenants handles the ListTenants RPC call
func (s *AuthServer) ListTenants(ctx context.Context, req *pb.ListTenantsRequest) (*pb.ListTenantsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

// UpdateTenant handles the UpdateTenant RPC call
func (s *AuthServer) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.Tenant, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

// DeleteTenant handles the DeleteTenant RPC call
func (s *AuthServer) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, status.Error(codes.Unimplemented, "method not implemented")
}

// HasPermission handles the HasPermission RPC call
func (s *AuthServer) HasPermission(ctx context.Context, req *pb.HasPermissionRequest) (*pb.HasPermissionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

// GetUserPermissions handles the GetUserPermissions RPC call
func (s *AuthServer) GetUserPermissions(ctx context.Context, req *pb.GetUserPermissionsRequest) (*pb.GetUserPermissionsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}
