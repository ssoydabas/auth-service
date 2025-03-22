package main

// @title Authentication Service API
// @version 1.0
// @description A comprehensive authentication service providing user management, authentication, and authorization capabilities.
// @termsOfService http://swagger.io/terms/

// @contact.name Sertan Soydabas
// @contact.email ssoydabas41@gmail.com
// @contact.url https://github.com/ssoydabas

// @license.name MIT
// @license.url https://github.com/ssoydabas/auth-service/LICENSE

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"

// @schemes http https
// @produce application/json
// @consumes application/json

// @tag.name Accounts
// @tag.description Account management operations including registration, login, profile updates, and password management
// @tag.docs.url https://docs.example.com/accounts
// @tag.docs.description Extended documentation for the account operations

// @tag.name Authentication
// @tag.description Authentication operations including token management and validation
// @tag.docs.url https://docs.example.com/authentication
// @tag.docs.description Detailed information about the authentication process

// @tag.name Authorization
// @tag.description Authorization operations including role management and permissions
// @tag.docs.url https://docs.example.com/authorization
// @tag.docs.description Complete guide to authorization mechanisms

// @x-logo {"url": "https://your-logo-url.com", "backgroundColor": "#FFFFFF", "altText": "API Logo"}

// @externalDocs.description Find out more about our API
// @externalDocs.url https://docs.example.com
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
	_ "github.com/ssoydabas/auth-service/docs"
	"github.com/ssoydabas/auth-service/pkg/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	e.Use(middleware.ErrorHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
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
