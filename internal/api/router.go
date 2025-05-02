package api

import (
	"github.com/gin-gonic/gin"

	"github.com/open-move/intercord/internal/middleware"
)

func SetupRouter(
	authHandler *AuthHandler,
	teamHandler *TeamHandler,
	subscriptionHandler *SubscriptionHandler,
	channelHandler *ChannelHandler,
	notificationHandler *NotificationHandler,
	jwtMiddleware *middleware.JWTAuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/verify-email", authHandler.VerifyEmail)
		auth.POST("/request-reset-password", authHandler.RequestPasswordReset)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	api := router.Group("")
	api.Use(jwtMiddleware.AuthRequired())
	{

		teams := api.Group("/teams")
		{
			teams.GET("", teamHandler.GetTeams)
			teams.POST("", teamHandler.CreateTeam)
			teams.GET("/:id", teamHandler.GetTeam)
			teams.DELETE("/:id", teamHandler.DeleteTeam)
			teams.POST("/:id/invite", teamHandler.InviteToTeam)
			teams.POST("/:id/join", teamHandler.JoinTeam)
			teams.POST("/:id/leave", teamHandler.LeaveTeam)

			teams.GET("/:team_id/subscriptions", subscriptionHandler.GetTeamSubscriptions)

			teams.GET("/:team_id/channels", channelHandler.GetTeamChannels)
		}

		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("", subscriptionHandler.GetSubscriptions)
			subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
		}

		channels := api.Group("/channels")
		{
			channels.POST("", channelHandler.CreateChannel)
			channels.GET("", channelHandler.GetChannels)
			channels.GET("/:id", channelHandler.GetChannel)
			channels.PUT("/:id", channelHandler.UpdateChannel)
			channels.DELETE("/:id", channelHandler.DeleteChannel)
			channels.POST("/subscribe", channelHandler.SubscribeChannel)
			channels.POST("/unsubscribe", channelHandler.UnsubscribeChannel)
		}

		notifications := api.Group("/notifications")
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.GET("/:id", notificationHandler.GetNotification)
		}
	}

	return router
}
