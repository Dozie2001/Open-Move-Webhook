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

	r.GET("/health", handlers.HealthHandler)

	// r.POST("/api/auth/verify", handlers.HandleAuth)
	// r.POST("/api/auth/refresh", AuthMiddleware(), handlers.HandleTokenRefresh)

	// // OAuth routes
	// oauthRoutes := r.Group("/oauth")

	// oauthRoutes.GET("/initialize", handlers.InitalizeOAuthSignIn)
	// oauthRoutes.GET("/callback", handlers.HandleOAuthCallBack)

	// All other API routes should be mounted on this route group
	// apiRoutes := r.Group("/api")

	// // mount the API routes auth middleware
	// apiRoutes.Use(AuthMiddleware())

	api := r.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/signup", handlers.SignUp)
		auth.POST("/login", handlers.Login)
		auth.POST("/refresh", handlers.RefreshAccessToken)
	}

	zkauth := api.Group("/auth/zklogin")
	{ // missing middleware -verifygoogletoken
		zkauth.POST("/register", handlers.ZKRegister)
		zkauth.GET("/salt", middleware.AuthMiddleware(), handlers.ZkSalt)
		zkauth.POST("/refresh", handlers.ZKRefreshAccessToken)
	}
	// Protected route
	api.GET("/me", middleware.AuthMiddleware(), handlers.Me)

	return r
}
