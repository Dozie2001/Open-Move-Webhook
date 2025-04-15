package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/Dozie2001/Open-Move-Webhook/internal/utils"
	"github.com/google/uuid"

	"github.com/Dozie2001/Open-Move-Webhook/internal/db"
	"github.com/Dozie2001/Open-Move-Webhook/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// signup with good old email and password.
func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	user := models.User{
		Email:        body.Email,
		PasswordHash: utils.SQLNullString(string(hashedPassword)),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

// login user with email and password
// TODO: send otp to email for verification
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // 15 minutes
	})
	accessTokenString, _ := accessToken.SignedString(jwtSecret)

	refreshToken := uuid.NewString()
	user.RefreshToken = utils.SQLNullString(refreshToken)
	user.TokenExpiry = utils.SQLNullTime(time.Now().Add(7 * 24 * time.Hour)) // 7 days
	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshToken,
	})
}

// refresh token to get the access token
func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	var user models.User
	if err := db.DB.Where("refresh_token = ?", body.RefreshToken).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if !user.TokenExpiry.Valid || user.TokenExpiry.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Generate new access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	accessTokenString, _ := accessToken.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessTokenString,
	})
}

// get all of the user information
// TODO: hide sensitive data in the payload
func Me(c *gin.Context) {
	userId := c.GetString("user_id")

	var user models.User
	if err := db.DB.First(&user, "id = ?", userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func ZkLogin(c *gin.Context) {

}

func ZkSalt(c *gin.Context) {

}
