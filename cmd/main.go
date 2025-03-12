package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"log"

	"github.com/ssoydabas/auth-service/internal/repository"
	"github.com/ssoydabas/auth-service/internal/service"
	"github.com/ssoydabas/auth-service/internal/transport/http/handler"
	"github.com/ssoydabas/auth-service/pkg/config"
	"github.com/ssoydabas/auth-service/pkg/postgres"

	"github.com/labstack/echo/v4"
)

const ShutdownTimeout = 5 * time.Second

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.ConnectPQ(*cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	e := echo.New()
	apiVersion := "/v1"
	apiPrefix := e.Group("/api" + apiVersion)

	handler.NewAccountHandler(service.NewAccountService(repository.NewAccountRepository(db))).AddRoutes(apiPrefix)

	// Graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := e.Start(":" + strconv.Itoa(cfg.Port)); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Fatalf("Failed to start server: %v", err)
	case <-shutdownChan:
		log.Println("Received shutdown signal, shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}

		log.Println("Server shutdown successfully")
	}
}
