# Auth Service

The Auth Service is responsible for user authentication, authorization, and tenant management in the Strategic Brand Optimization Platform. It provides JWT-based authentication and integrates with Supabase for secure user management.

## Features

- **User Authentication**: Email/password authentication using Supabase Auth
- **JWT Token Management**: Generation and validation of access and refresh tokens
- **Tenant Isolation**: Association of users with tenants (organizations)
- **Role-Based Access Control**: Role management within tenant boundaries
- **User Management**: User registration, profile updates, and deactivation
- **Health Monitoring**: Health check endpoint for service monitoring

## Service Architecture

```
auth/
├── cmd/               # Entry point for the service
├── server/            # gRPC server implementation
├── client/            # gRPC client for other services to use
├── pb/                # Protocol Buffers definitions
├── service/           # Auth business logic
└── repository/        # Data access layer with Supabase
```

## API Endpoints

### gRPC Service

The Auth service exposes a gRPC API defined in `auth/pb/auth.proto`.

Key operations:
- `RegisterUser`: Register a new user
- `Login`: Authenticate a user
- `ValidateToken`: Validate and decode a JWT token
- `RefreshToken`: Generate new tokens using a refresh token
- `GetUserProfile`: Retrieve user profile information
- `UpdateUserProfile`: Update user profile

### HTTP Endpoints

- **Health Check**: `/health` - Returns health status of the service

## Technical Details

### Authentication Flow

1. **Registration**:
   - User calls `RegisterUser` with email, password, and tenant details
   - Service creates user in Supabase Auth
   - Service creates tenant record and associates user with it
   - JWT tokens with tenant context are returned

2. **Login**:
   - User calls `Login` with email and password
   - Service validates credentials with Supabase Auth
   - JWT tokens with tenant context are returned
   
3. **Token Validation**:
   - Services call `ValidateToken` to verify token authenticity
   - Service decodes token and returns user claims (ID, tenant, role)

4. **Token Refresh**:
   - User calls `RefreshToken` with an expired access token and a valid refresh token
   - Service validates the refresh token and issues new tokens

### JWT Token Structure

```json
{
  "sub": "user-uuid",
  "iss": "go-competitor-auth",
  "exp": 1617278374,
  "iat": 1617274774,
  "tenant_id": "org-uuid",
  "role": "admin"
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to listen on | `9001` |
| `SUPABASE_URL` | Supabase instance URL | - |
| `SUPABASE_ANON_KEY` | Supabase anon key | - |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | - |
| `JWT_SECRET` | Secret for signing JWT tokens | - |
| `ACCESS_TOKEN_EXPIRY` | Access token expiry in seconds | `3600` (1 hour) |
| `REFRESH_TOKEN_EXPIRY` | Refresh token expiry in seconds | `604800` (7 days) |

## Usage Examples

### Starting the Service

The service can be started directly:

```bash
go run auth/cmd/main.go
```

Or via Docker:

```bash
docker build -t auth-service -f auth/Dockerfile .
docker run -p 9001:9001 auth-service
```

### Health Check

```bash
curl http://localhost:9001/health
```

Example response:
```json
{"status":"UP"}
```

### Accessing via gRPC Client

To use the Auth service from other services or applications:

```go
import "github.com/donaldnash/go-competitor/auth/client"

func main() {
    // Create auth client
    authClient, err := client.NewAuthClient("localhost:9001")
    if err != nil {
        log.Fatalf("Failed to create auth client: %v", err)
    }
    
    // Validate token
    claims, err := authClient.ValidateToken(ctx, token)
    if err != nil {
        log.Fatalf("Token validation failed: %v", err)
    }
    
    // Use claims
    log.Printf("User ID: %s, Tenant: %s", claims.UserID, claims.OrganizationID)
}
```

## Error Handling

Common error scenarios:

- **Invalid Credentials**: Returned when login credentials are incorrect
- **User Already Exists**: Returned when trying to register with an existing email
- **Invalid Token**: Returned when a token is expired or malformed
- **Authorization Error**: Returned when a user doesn't have permission for an operation

## Dependencies

- **Supabase Auth**: For user authentication and management
- **JWT-Go**: For JWT token generation and validation
- **gRPC**: For service API
- **PostgreSQL**: Via Supabase for user and tenant data storage 