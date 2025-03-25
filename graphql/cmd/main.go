package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/donaldnash/go-competitor/graphql/server"
)

func main() {
	log.Println("Starting GraphQL server...")

	// Get port from environment variable
	port := os.Getenv("GRAPHQL_PORT")
	if port == "" {
		port = "8080"
	}

	// Create the GraphQL server
	srv, err := server.NewGraphQLServer()
	if err != nil {
		log.Fatalf("Failed to create GraphQL server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("GraphQL server listening on :%s", port)
		if err := srv.Start(":" + port); err != nil {
			log.Fatalf("Failed to start GraphQL server: %v", err)
		}
	}()

	// Wait for interruption signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
}
