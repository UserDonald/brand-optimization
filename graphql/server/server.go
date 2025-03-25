package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/donaldnash/go-competitor/auth/client"
	"github.com/donaldnash/go-competitor/graphql/middleware"
	"github.com/donaldnash/go-competitor/graphql/resolvers"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// GraphQLServer handles the GraphQL API requests and delegates to the appropriate service
// It serves as the API gateway for the entire platform
type GraphQLServer struct {
	schema     *graphql.Schema    // Parsed GraphQL schema
	router     http.Handler       // HTTP router for handling requests
	authClient *client.AuthClient // Client for auth service communication
}

// NewGraphQLServer creates a new GraphQLServer instance and initializes all service clients
// and resolvers needed to handle GraphQL queries and mutations
func NewGraphQLServer() (*GraphQLServer, error) {
	// Read GraphQL schema file
	schema, err := os.ReadFile("../schema.graphql")
	if err != nil {
		return nil, fmt.Errorf("failed to read GraphQL schema: %w", err)
	}

	// Get service URLs from environment variables, with sensible defaults for local development
	serviceURLs := getServiceURLs()

	// Initialize clients for all microservices
	// Start with auth client which is required for authentication
	authClient, err := client.NewAuthClient(serviceURLs["auth"])
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %w", err)
	}

	// Initialize all resolvers with their respective service clients
	// For each resolver, we pass the appropriate service client to handle domain-specific operations
	resolvers, err := initializeResolvers(authClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resolvers: %w", err)
	}

	// Parse the GraphQL schema with the root resolver
	parsedSchema := graphql.MustParseSchema(string(schema), resolvers)

	// Set up HTTP router with middleware chain
	router := setupRouter(parsedSchema, authClient)

	return &GraphQLServer{
		schema:     parsedSchema,
		router:     router,
		authClient: authClient,
	}, nil
}

// getServiceURLs retrieves all service URLs from environment variables
// or falls back to sensible defaults for local development
func getServiceURLs() map[string]string {
	urls := make(map[string]string)

	// Define all services and their environment variable names
	services := map[string]string{
		"auth":         "AUTH_SERVICE_URL",
		"notification": "NOTIFICATION_SERVICE_URL",
		"competitor":   "COMPETITOR_SERVICE_URL",
		"engagement":   "ENGAGEMENT_SERVICE_URL",
		"content":      "CONTENT_SERVICE_URL",
		"audience":     "AUDIENCE_SERVICE_URL",
		"analytics":    "ANALYTICS_SERVICE_URL",
		"scraper":      "SCRAPER_SERVICE_URL",
	}

	// Default ports for local development
	defaultPorts := map[string]string{
		"auth":         "9001",
		"notification": "9002",
		"competitor":   "9003",
		"engagement":   "9004",
		"content":      "9005",
		"audience":     "9006",
		"analytics":    "9007",
		"scraper":      "9008",
	}

	// Get URLs from environment or use defaults
	for service, envVar := range services {
		url := os.Getenv(envVar)
		if url == "" {
			url = "localhost:" + defaultPorts[service]
		}
		urls[service] = url
	}

	return urls
}

// initializeResolvers creates and initializes all resolvers needed for the GraphQL API
func initializeResolvers(authClient *client.AuthClient) (*resolvers.RootResolver, error) {
	// Initialize the auth resolver with the auth client
	authResolver := resolvers.NewAuthResolver(authClient)

	// Initialize service clients and their respective resolvers
	// Each resolver is given its service client to handle domain-specific operations

	// Use empty resolver implementations for now
	// These will be replaced with actual implementations as we connect to each service
	// The client initialization code is ready to be used when the services are available

	// Note: Client initialization is commented out until the services are fully implemented
	// and their client packages are available

	/*
		competitorClient, err := competitor.NewCompetitorClient(serviceURLs["competitor"])
		if err != nil {
			return nil, fmt.Errorf("failed to create competitor client: %w", err)
		}

		audienceClient, err := audience.NewAudienceClient(serviceURLs["audience"])
		if err != nil {
			return nil, fmt.Errorf("failed to create audience client: %w", err)
		}

		contentClient, err := content.NewContentClient(serviceURLs["content"])
		if err != nil {
			return nil, fmt.Errorf("failed to create content client: %w", err)
		}

		analyticsClient, err := analytics.NewAnalyticsClient(serviceURLs["analytics"])
		if err != nil {
			return nil, fmt.Errorf("failed to create analytics client: %w", err)
		}

		notificationClient, err := notification.NewNotificationClient(serviceURLs["notification"])
		if err != nil {
			return nil, fmt.Errorf("failed to create notification client: %w", err)
		}

		engagementClient, err := engagement.NewEngagementClient(serviceURLs["engagement"])
		if err != nil {
			return nil, fmt.Errorf("failed to create engagement client: %w", err)
		}

		scraperClient, err := scraper.NewScraperClient(serviceURLs["scraper"])
		if err != nil {
			return nil, fmt.Errorf("failed to create scraper client: %w", err)
		}
	*/

	// Using stub implementations for now
	competitorResolver := &resolvers.CompetitorResolver{}
	audienceResolver := &resolvers.AudienceResolver{}
	contentResolver := &resolvers.ContentResolver{}
	analyticsResolver := &resolvers.AnalyticsResolver{}

	// Create a root resolver that combines all resolvers
	rootResolver := &resolvers.RootResolver{
		AuthResolver:       authResolver,
		CompetitorResolver: competitorResolver,
		AudienceResolver:   audienceResolver,
		ContentResolver:    contentResolver,
		AnalyticsResolver:  analyticsResolver,
	}

	return rootResolver, nil
}

// setupRouter configures the HTTP router with all necessary middleware
func setupRouter(schema *graphql.Schema, authClient *client.AuthClient) http.Handler {
	router := http.NewServeMux()

	// Create the base GraphQL handler
	graphqlHandler := &relay.Handler{Schema: schema}

	// Add authentication middleware
	authHandler := middleware.AuthMiddleware(authClient)(graphqlHandler)

	// Add CORS middleware for browser support
	corsHandler := addCORS(authHandler)

	// Register the handler for the /query endpoint
	router.Handle("/query", corsHandler)

	// Add health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP"}`))
	})

	return router
}

// ServeHTTP implements the http.Handler interface
// This allows the GraphQLServer to be used directly with http.ListenAndServe
func (s *GraphQLServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Start begins listening for HTTP requests on the specified address
func (s *GraphQLServer) Start(addr string) error {
	log.Printf("GraphQL server listening on %s", addr)
	return http.ListenAndServe(addr, s)
}

// Execute executes a GraphQL query programmatically
// This is useful for internal service-to-service communication
func (s *GraphQLServer) Execute(ctx context.Context, query string) *graphql.Response {
	return s.schema.Exec(ctx, query, "", nil)
}

// addCORS adds Cross-Origin Resource Sharing headers to enable browser access
func addCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		h.ServeHTTP(w, r)
	})
}
