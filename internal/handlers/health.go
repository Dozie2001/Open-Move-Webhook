package handlers

import (
	"github.com/Dozie2001/Open-Move-Webhook/pkg/response"
	"net/http"
	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	response.Success(c, http.StatusOK, "Open Move Webhook API", nil)
}