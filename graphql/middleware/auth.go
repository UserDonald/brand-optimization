package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/donaldnash/go-competitor/auth/client"
)

// contextKey is a private type for context keys to prevent collisions
type contextKey string

// String returns the string representation of the context key
func (c contextKey) String() string {
	return string(c)
}

// Auth context keys used to store authentication information in the request context
const (
	UserIDKey          = contextKey("user_id")          // The authenticated user's ID
	TenantIDKey        = contextKey("tenant_id")        // The user's tenant (organization) ID
	UserRoleKey        = contextKey("user_role")        // The user's role within their tenant
	IsAuthenticatedKey = contextKey("is_authenticated") // Whether the user is authenticated
)

// AuthMiddleware creates middleware for JWT authentication and tenant context population
// It validates tokens using the auth service and adds user information to the request context
func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for introspection queries and options requests
			// GraphQL introspection is used by tools like GraphiQL and should remain accessible
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No auth header provided - continue as unauthenticated
				// This allows public operations to proceed, but protected operations
				// will be rejected by the resolvers
				log.Println("Request without auth token received")
				next.ServeHTTP(w, r)
				return
			}

			// Check if header has Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token with the auth service
			claims, err := authClient.ValidateToken(r.Context(), token)
			if err != nil {
				// Token validation failed
				log.Printf("Token validation failed: %v", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract tenant context for multi-tenancy
			if claims.OrganizationID == "" {
				log.Println("Warning: Token missing organization/tenant ID")
			}

			// Add auth information to request context
			// This will be used by resolvers to enforce tenant boundaries
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, TenantIDKey, claims.OrganizationID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			ctx = context.WithValue(ctx, IsAuthenticatedKey, true)

			// Log successful authentication
			log.Printf("Authenticated user: %s, tenant: %s, role: %s",
				claims.UserID, claims.OrganizationID, claims.Role)

			// Continue with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the context if available
// Returns empty string if not authenticated
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		return userID
	}
	return ""
}

// GetTenantID extracts the tenant ID from the context if available
// Returns empty string if not authenticated or tenant ID not present
func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok && tenantID != "" {
		return tenantID
	}
	return ""
}

// GetUserRole extracts the user role from the context if available
// Returns empty string if not authenticated or role not present
func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok && role != "" {
		return role
	}
	return ""
}

// IsAuthenticated checks if the user is authenticated based on the context
// This can be used by resolvers to enforce authentication requirements
func IsAuthenticated(ctx context.Context) bool {
	if isAuth, ok := ctx.Value(IsAuthenticatedKey).(bool); ok {
		return isAuth
	}
	return false
}

// RequireAuthentication is a helper that returns an error if the user is not authenticated
// Use this at the beginning of resolver methods that require authentication
func RequireAuthentication(ctx context.Context) error {
	if !IsAuthenticated(ctx) {
		return fmt.Errorf("authentication required")
	}
	return nil
}

// RequireRole is a helper that returns an error if the user doesn't have the required role
// Use this at the beginning of resolver methods that require specific roles
func RequireRole(ctx context.Context, requiredRole string) error {
	if err := RequireAuthentication(ctx); err != nil {
		return err
	}

	role := GetUserRole(ctx)
	if role != requiredRole {
		return fmt.Errorf("required role '%s' not found, user has role '%s'", requiredRole, role)
	}

	return nil
}
