package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
)

// run handles the actual execution of the program. It's separated from main to allow for:
// 1. Proper error handling (since main can't return errors)
// 2. Easier testing (by injecting dependencies)
// 3. Better resource cleanup
//
// Parameters:
// - ctx: Controls the lifecycle of the program
// - w: Where to write output (usually os.Stdout)
// - args: Command line arguments (usually os.Args)
// - getenv: Function to read environment variables (usually os.Getenv)
func run(ctx context.Context, w io.Writer, args []string, getenv func(string) string) error {
	// Create a new context that will be cancelled when the program receives an interrupt signal
	// This enables graceful shutdown when someone presses Ctrl+C
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	// Ensure we call cancel() when the function returns to clean up resources
	defer cancel()

	// Getting env file locally
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	// Get the API key from the environment
	apiKey := getenv("TFL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("TFL_API_KEY environment variable is required")
	}

	// Initialize the client
	client := tfl.NewClient(apiKey)

	// Create server mux and add routes
	mux := http.NewServeMux()
	addRoutes(mux, client)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server
	go func() {
		log.Printf("Server active and listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Wait for interrupt
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error shutting down http server: %w", err)
	}

	return nil
}

// main is the entry point of the program.
// It's kept minimal and delegates actual work to run()
func main() {
	// Create the base context for the program
	ctx := context.Background()

	// Call run with the standard OS outputs and handle any errors
	if err := run(ctx, os.Stdout, os.Args, os.Getenv); err != nil {
		// If there's an error, write it to stderr and exit with error code 1
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
