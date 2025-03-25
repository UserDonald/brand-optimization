package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/donaldnash/go-competitor/audience/pb"
	"github.com/donaldnash/go-competitor/audience/repository"
	"github.com/donaldnash/go-competitor/audience/server"
	"github.com/donaldnash/go-competitor/audience/service"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config holds the server configuration
type Config struct {
	Port         string `envconfig:"PORT" default:"50053"`
	SupabaseURL  string `envconfig:"SUPABASE_URL" required:"true"`
	SupabaseKey  string `envconfig:"SUPABASE_ANON_KEY" required:"true"`
	ServiceRole  string `envconfig:"SUPABASE_SERVICE_ROLE" required:"true"`
	TenantHeader string `envconfig:"TENANT_HEADER" default:"X-Tenant-ID"`
}

func main() {
	// Load configuration
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up listener
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Server starting on port %s", cfg.Port)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Set up the repository, service, and server
	// We use a blank tenant ID for initialization; the actual tenant ID will be
	// supplied in each request context
	repo, err := repository.NewSupabaseAudienceRepository("")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	audienceService, err := service.NewAudienceService(repo)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	audienceServer := server.NewAudienceServer(audienceService)
	pb.RegisterAudienceServiceServer(grpcServer, audienceServer)

	// Register reflection service for development tools
	reflection.Register(grpcServer)

	// Set up HTTP server for health check
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "9006" // Default HTTP port
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"UP"}`))
		})

		httpServer := &http.Server{
			Addr:    ":" + httpPort,
			Handler: mux,
		}

		log.Printf("HTTP health check server listening on port %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Set up graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Shutting down server...")
		cancel()
		grpcServer.GracefulStop()
	}()

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
