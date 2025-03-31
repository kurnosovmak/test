package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kurnosovmak/test/internal/config"
	"github.com/kurnosovmak/test/pkg/logger"

	"go.uber.org/zap"
)

var (
	envPath = flag.String("config", ".yaml", "Path to the environment file")
)

func init() {
	flag.Parse()
}

func main() {

	// Initialize logger
	logger.Init()
	defer logger.Info("Shutting down the application")

	// Load configuration with env file path
	cfg, err := config.LoadConfig(*envPath)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.String("env path", *envPath), zap.Error(err))
	}

	// Log successful startup
	logger.Info("Starting the application",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// Create server address
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// Basic health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy"))
	})
	// Start the server
	logger.Info("Server is starting", zap.String("address", addr))
	srv := &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	// Start server in a goroutine so it doesn't block
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown the server
	logger.Info("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
