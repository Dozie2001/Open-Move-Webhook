package handlers

import (
	"net/http"

	"github.com/Dozie2001/Open-Move-Webhook/internal/db"
	"github.com/Dozie2001/Open-Move-Webhook/internal/models"
	"github.com/Dozie2001/Open-Move-Webhook/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type CreateSubscriptionRequest struct {
	WebhookURL     string         `json:"webhook_url" binding:"required,url"`
	Secret         string         `json:"secret"`
	EventType      string         `json:"event_type" binding:"required"`
	FilterCriteria datatypes.JSON `json:"filter_criteria"`
}

func CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	webhook := &models.Webhook{
		Url:    req.WebhookURL,
		Secret: req.Secret,
		Status: true,
	}
	result := db.DB.Where("url = ?", req.WebhookURL).First(webhook)
	if result.Error != nil {
		if err := db.DB.Create(webhook).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to create webhook")
			return
		}
	}

	// Create subscription
	subscription := &models.Subscription{
		WebhookId:      webhook.Id,
		EventType:      req.EventType,
		FilterCriteria: req.FilterCriteria,
	}

	if err := db.DB.Create(subscription).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create subscription")
		return
	}

	response.Success(c, http.StatusCreated, "Subscription created successfully", subscription)
}

func ListSubscriptions(c *gin.Context) {
	var subscriptions []models.Subscription

	result := db.DB.Find(&subscriptions)
	if result.Error != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch subscriptions")
		return
	}

	response.Success(c, http.StatusOK, "Subscriptions retrieved successfully", subscriptions)
}


func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	result := db.DB.Delete(&models.Subscription{}, "id = ?", id)
	if result.Error != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete subscription")
		return
	}

	if result.RowsAffected == 0 {
		response.Error(c, http.StatusNotFound, "Subscription not found")
		return
	}

	response.Success(c, http.StatusNoContent, "Subscription deleted successfully", nil)
}

func UpdateSubscription(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		EventType      string         `json:"event_type"`
		FilterCriteria datatypes.JSON `json:"filter_criteria"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var subscription models.Subscription
	if err := db.DB.First(&subscription, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Subscription not found")
		return
	}

	// Update fields
	if req.EventType != "" {
		subscription.EventType = req.EventType
	}
	if req.FilterCriteria != nil {
		subscription.FilterCriteria = req.FilterCriteria
	}

	if err := db.DB.Save(&subscription).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update subscription")
		return
	}

	response.Success(c, http.StatusOK, "Subscription updated successfully", subscription)
}