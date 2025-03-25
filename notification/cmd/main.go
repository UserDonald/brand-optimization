package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/donaldnash/go-competitor/notification/pb"
	"github.com/donaldnash/go-competitor/notification/repository"
	"github.com/donaldnash/go-competitor/notification/server"
	"github.com/donaldnash/go-competitor/notification/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 50053, "The server port")
	flag.Parse()

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("notification service starting on port %d...", *port)

	// Get tenant ID from environment variable
	// In a real production environment, this would be handled more robustly
	tenantID := os.Getenv("TENANT_ID")
	if tenantID == "" {
		// For development, use a default tenant ID
		tenantID = "default-tenant"
		log.Printf("TENANT_ID not set, using default: %s", tenantID)
	}

	// Create repository
	repo, err := repository.NewSupabaseNotificationRepository(tenantID)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	// Create service
	svc, err := service.NewNotificationService(repo)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create server
	srv, err := server.NewNotificationServer(svc)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register service
	// This line will be uncommented once we generate the protobuf code
	pb.RegisterNotificationServiceServer(grpcServer, srv)

	// Register reflection service (useful for grpcurl and other gRPC debug tools)
	reflection.Register(grpcServer)

	// Set up HTTP server for health check
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "9002" // Default HTTP port for notification service
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

	// Handle graceful shutdown
	// We create a context that we could use to propagate cancellation
	// signals to any ongoing operations if needed
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channel to listen for interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %v, initiating graceful shutdown", sig)
		cancel()
		grpcServer.GracefulStop()
	}()

	// Start server
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("notification service stopped")
}
