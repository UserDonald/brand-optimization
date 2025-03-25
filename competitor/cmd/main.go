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
	"github.com/donaldnash/go-competitor/competitor/pb"
	"github.com/donaldnash/go-competitor/competitor/repository"
	"github.com/donaldnash/go-competitor/competitor/server"
	"github.com/donaldnash/go-competitor/competitor/service"
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
	port := 9003 // Default port
	if cfg.Port != "" {
		// In a real implementation, we would get the port from the configuration
	}

	// Create repository (using a temporary tenant ID for demonstration)
	// In a real implementation, we would use different tenant IDs based on the request
	repo, err := repository.NewSupabaseCompetitorRepository("system")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create service
	svc := service.NewCompetitorService(repo)

	// Create server
	srv := server.NewCompetitorServer(svc)

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register the server with the generated protobuf code
	pb.RegisterCompetitorServiceServer(grpcServer, srv)

	// Register reflection service for development/debugging
	reflection.Register(grpcServer)

	log.Printf("Starting competitor service on port %d", port)

	// Set up HTTP server for health check
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "9003" // Default HTTP port
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

	log.Println("Shutting down competitor service...")
	grpcServer.GracefulStop()
}
