package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

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
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	// Get the API key from the environment
	apiKey := getenv("TFL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("TFL_API_KEY environment variable is required")
	}

	// Initialize the client
	client := tfl.NewClient(apiKey)

	predictions, err := client.GetVictoriaPredictions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Victoria line predictions: %w", err)
	}

	fmt.Fprintln(w, "Predictions fetched:")
	for _, prediction := range predictions {
		fmt.Fprintf(w, "%+v\n", prediction)
	}

	// TODO: Add server setup and main loop here

	http.HandleFunc("/", handler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()

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
