package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/donaldnash/go-competitor/common/config"
	"github.com/donaldnash/go-competitor/engagement/pb"
	"github.com/donaldnash/go-competitor/engagement/repository"
	"github.com/donaldnash/go-competitor/engagement/server"
	"github.com/donaldnash/go-competitor/engagement/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up service port
	port := 9004 // Default port

	// Create repository using a system tenant ID for service-level access
	// In a real implementation, the tenant ID would come from authentication middleware
	repo, err := repository.NewSupabaseEngagementRepository("system")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create service
	svc := service.NewEngagementService(repo)

	// Create server
	srv := server.NewEngagementServer(svc)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register the server with the generated protobuf code
	pb.RegisterEngagementServiceServer(grpcServer, srv)

	// Register reflection service for development/debugging
	reflection.Register(grpcServer)

	// Set up HTTP server for health check
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "9004" // Use the same port
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

	// For now, we'll just log that we're starting the server
	log.Printf("Starting engagement service on port %d (config env: %s)", port, cfg.Environment)

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down engagement service...")
	grpcServer.GracefulStop()
}
