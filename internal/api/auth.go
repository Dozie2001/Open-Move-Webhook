package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-move/intercord/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
	baseURL     string
}

func NewAuthHandler(userService *services.UserService, baseURL string) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		baseURL:     baseURL,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input services.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input: " + err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), input, h.baseURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input: " + err.Error()})
		return
	}

	auth, err := h.userService.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, auth)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Token is required"})
		return
	}

	err := h.userService.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Email verified successfully"})
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input: " + err.Error()})
		return
	}

	err := h.userService.RequestPasswordReset(c.Request.Context(), input.Email, h.baseURL)
	if err != nil {

		c.JSON(http.StatusOK, SuccessResponse{Message: "If your email is registered, you will receive a password reset link"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "If your email is registered, you will receive a password reset link"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input services.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input: " + err.Error()})
		return
	}

	err := h.userService.ResetPassword(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Password reset successfully"})
}
