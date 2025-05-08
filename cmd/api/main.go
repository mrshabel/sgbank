package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrshabel/shbank/internal/config"
	"github.com/mrshabel/shbank/internal/handlers"
	log "github.com/mrshabel/shbank/internal/logger"
)

func main() {
	router := gin.Default()
	cfg := config.New()

	logger := log.New(cfg.Env)

	// http server
	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	// register middlewares

	// register handlers here
	handlers.RegisterPingHandler(router, logger)

	// start server in background
	go func() {
		logger.Info("Server is starting on " + cfg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// monitor signal interrupts
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	// timeout graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	<-sigCh
	logger.Warn("Received interrupt. Shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	logger.Info("Server shutdown complete")
	os.Exit(0)
}
