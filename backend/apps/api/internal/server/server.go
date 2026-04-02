package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"

	"api/internal/database"
	"api/internal/middleware"
)

type Server struct {
	port int

	db          database.Service
	logger      *slog.Logger
	rateLimiter *middleware.RateLimiter
	deps        *Dependencies
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db:     database.New(),
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		rateLimiter: middleware.NewRateLimiter(
			redis.NewClient(&redis.Options{
				Addr: envOrDefault("REDIS_ADDR", "localhost:6379"),
			}),
			100,
			time.Minute,
		),
	}
	NewServer.deps = NewDependencies(NewServer.logger)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
