package api

import (
	"os"

	"github.com/Dozie2001/Open-Move-Webhook/internal/handlers"
	"github.com/Dozie2001/Open-Move-Webhook/internal/middleware"

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

	// All other API routes should be mounted on this route group
	apiRoutes := r.Group("/api/v1")

	// Subscription routes

	apiRoutes.POST("/subscription", handlers.CreateSubscription)
	apiRoutes.GET("/subscription", handlers.ListSubscriptions)
	apiRoutes.DELETE("/subscriptions/:id", handlers.DeleteSubscription)
	apiRoutes.PATCH("/subscriptions/:id", handlers.UpdateSubscription)

	// subscriptionRoutes.Use(middleware.CacheMiddleware())

	return r
}
