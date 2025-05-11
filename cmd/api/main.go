package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrshabel/sgbank/internal/config"
	"github.com/mrshabel/sgbank/internal/db"
	"github.com/mrshabel/sgbank/internal/handlers"
	log "github.com/mrshabel/sgbank/internal/logger"
	"github.com/mrshabel/sgbank/internal/repository"
)

func main() {
	// setup configurations
	cfg := config.New()
	logger := log.New(cfg.Env)

	// setup database
	db, err := db.New(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to connect to db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// http server
	router := gin.Default()
	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	// register middlewares

	// create repositories
	userRepo := repository.NewUserRepository(db, logger)

	// create handlers
	userHandler := handlers.NewUserHandler(userRepo, logger)

	// register handlers here
	handlers.RegisterPingHandler(router, logger)
	handlers.RegisterUserHandlers(userHandler, router, logger)

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
