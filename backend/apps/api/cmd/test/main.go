package main

import (
	"api/internal/config"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "./config/config.yaml", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	cfg := config.Load(*configPath)
	fmt.Println(cfg)
	fmt.Println("-----------------------------")

	// Add a few default environmental attributes that always are included
	defaultAttrs := []slog.Attr{
		slog.String("service", "authService"),
		slog.String("node", "auth-node-1"),
		slog.String("environment", "production"),
	}

	// Create a Logger and apply the Default attrs to the handler
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}).WithAttrs(defaultAttrs)
	logger := slog.New(handler)

	logger.Debug("Starting authentication service")

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Mock a Request ID for tracing and UserID
		request_id, _ := uuid.NewV7()
		user_id := uuid.NewString()

		// Add Contextual information for this Function call,
		logger.Error("failed to login",
			slog.String("request_id", request_id.String()),
			slog.String("user_id", user_id),
			slog.Group("payload",
				slog.String("username", "Harry Potter"),
				slog.String("auth_method", "password"),
			),
		)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("the server crashed")
	}
}
