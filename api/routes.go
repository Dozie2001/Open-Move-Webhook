package api

import (
	"github.com/Dozie2001/Open-Move-Webhook/internal/handlers"
	"github.com/Dozie2001/Open-Move-Webhook/internal/middleware"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func BuildRoutesHandler() *gin.Engine {
	r := gin.New()

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(middleware.RateLimitMiddleware())

	r.GET("/health", handlers.HealthHandler)

	// Subscription routes
	subscriptionRoutes := r.Group("/api/v1/subscriptions")
	{
		subscriptionRoutes.Use(middleware.CacheMiddleware())
		subscriptionRoutes.POST("", handlers.CreateSubscription)
		subscriptionRoutes.GET("", handlers.ListSubscriptions)
		subscriptionRoutes.DELETE("/:id", handlers.DeleteSubscription)
		subscriptionRoutes.PATCH("/:id", handlers.UpdateSubscription)
	}

	return r
}