package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/donaldnash/go-competitor/auth/pb"
	"github.com/donaldnash/go-competitor/auth/repository"
	"github.com/donaldnash/go-competitor/auth/server"
	"github.com/donaldnash/go-competitor/auth/service"
	"github.com/donaldnash/go-competitor/common/config"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up service port (for simplicity, we're just using the configured port + 1000)
	port := 9001
	if cfg.Port != "" {
		// In a real implementation, we would get the port from the configuration
		// For demonstration, we'll just use a fixed port
	}

	// Create repository
	repo, err := repository.NewSupabaseAuthRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create service
	svc := service.NewAuthService(repo)

	// Create server
	srv := server.NewAuthServer(svc)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register the server with the generated protobuf code
	pb.RegisterAuthServiceServer(grpcServer, srv)

	log.Printf("Starting auth service on port %d", port)

	// Set up HTTP server for health check
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "9001" // Default HTTP port
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

	log.Println("Shutting down auth service...")
	grpcServer.GracefulStop()
}
