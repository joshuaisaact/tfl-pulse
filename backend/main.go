package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
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

	// Example of using getenv (I'll add more environment variables later)
	apiKey := getenv("TFL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("TFL_API_KEY environment variable is required")
	}

	// TODO: Add server setup and main loop here

	return nil // For now, just return nil (no error)
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
