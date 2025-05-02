package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/open-move/intercord/internal/api"
	"github.com/open-move/intercord/internal/config"
	"github.com/open-move/intercord/internal/database"
	"github.com/open-move/intercord/internal/middleware"
	"github.com/open-move/intercord/internal/services"
)

func main() {

	cfg := config.Load()

	db := database.New(&cfg.Database)
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + cfg.Server.Port
	}

	emailService := services.NewEmailService(&cfg.Email)
	userService := services.NewUserService(db, &cfg.JWT, emailService)
	teamService := services.NewTeamService(db, emailService)
	subscriptionService := services.NewSubscriptionService(db, teamService)
	channelService := services.NewChannelService(db, teamService)
	notificationService := services.NewNotificationService(db, teamService)

	jwtMiddleware := middleware.NewJWTAuthMiddleware(&cfg.JWT)

	authHandler := api.NewAuthHandler(userService, baseURL)
	teamHandler := api.NewTeamHandler(teamService, baseURL)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService)
	channelHandler := api.NewChannelHandler(channelService)
	notificationHandler := api.NewNotificationHandler(notificationService)

	router := api.SetupRouter(
		authHandler,
		teamHandler,
		subscriptionHandler,
		channelHandler,
		notificationHandler,
		jwtMiddleware,
	)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
