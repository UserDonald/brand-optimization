# Authentication Service

This service handles user authentication, authorization, and tenant management for the Strategic Brand Optimization Platform.

## Purpose

The Auth Service:
- Authenticates users via email/password
- Issues JWT tokens for authenticated sessions
- Manages tenant (organization) data 
- Controls permissions and access rights
- Provides user profile management

## Architecture

The service follows a clean architecture pattern:

- `cmd/` - Entry point and application setup
- `server/` - gRPC server implementation
- `client/` - gRPC client for other services
- `pb/` - Protocol Buffers definitions
- `service/` - Business logic
- `repository/` - Data access layer

## API

The service exposes the following gRPC endpoints:

### Authentication
- `Login` - Authenticates a user and returns tokens
- `Register` - Creates a new user and organization
- `Logout` - Invalidates the current session
- `RefreshToken` - Issues new tokens using a refresh token

### User Management
- `GetUser` - Retrieves user details
- `CreateUser` - Creates a new user for an existing organization
- `UpdateUser` - Updates user details
- `DeleteUser` - Deactivates a user

### Tenant Management
- `GetTenant` - Retrieves tenant details
- `CreateTenant` - Creates a new tenant
- `UpdateTenant` - Updates tenant details
- `DeleteTenant` - Deactivates a tenant

### Authorization
- `ValidateToken` - Validates a JWT token
- `HasPermission` - Checks if a user has specific permissions
- `GetUserPermissions` - Lists all permissions for a user

## Development

### Running Locally

```bash
cd cmd
go run main.go
```

The server will start on port 9001.

### Environment Variables

- `SUPABASE_URL` - Supabase URL
- `SUPABASE_ANON_KEY` - Supabase anonymous key
- `SUPABASE_SERVICE_ROLE` - Supabase service role key
- `JWT_SECRET` - Secret for signing JWT tokens

## Usage from Other Services

Other services should use the provided client to interact with the auth service:

```go
import "github.com/donaldnash/go-competitor/auth/client"

// Create client
authClient, err := client.NewAuthClient("localhost:9001")
if err != nil {
    // Handle error
}
defer authClient.Close()

// Validate token
claims, err := authClient.ValidateToken(ctx, token)
if err != nil {
    // Token is invalid
}

// Get user details
user, err := authClient.GetUser(ctx, userID)
if err != nil {
    // Handle error
}
```

## Testing

During development, the service provides dummy data when Supabase is not configured, making it easier to test other services that depend on authentication. 