package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Log environment configuration in debug mode
	debugLogEnvironment()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create and configure the server
	s, err := CreateVoiceNotifyServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start the server with stdio transport
	go func() {
		log.Println("Starting voice notify MCP server...")
		debugLog("Server starting with PID: %d", os.Getpid())
		if err := server.ServeStdio(s); err != nil {
			log.Printf("Server error: %v", err)
			debugLog("Server error details: %+v", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received")
	debugLog("Received signal, initiating graceful shutdown")
	log.Println("Shutting down voice notify MCP server...")
}
