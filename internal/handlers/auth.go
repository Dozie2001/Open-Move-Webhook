package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/Dozie2001/Open-Move-Webhook/internal/utils"
	"github.com/google/uuid"

	"database/sql"

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
func RefreshAccessToken(c *gin.Context) {
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

func ZKRegister(c *gin.Context) {
	var body struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Salt    string `json:"salt"`
		Address string `json:"address"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user models.User
	result := db.DB.First(&user, "sub = ?", body.Sub)

	if result.Error == nil {
		salt, _ := utils.Decrypt(user.Salt.String)
		access, refresh, _ := utils.GenerateTokens(body.Sub, body.Email)
		c.JSON(http.StatusOK, gin.H{
			"email":   user.Email,
			"address": user.SuiAddress.String,
			"salt":    salt,
			"access":  access,
			"refresh": refresh,
		})
		return
	}

	encryptedSalt, err := utils.Encrypt(body.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
		return
	}

	user = models.User{
		Email:      body.Email,
		Sub:        sql.NullString{String: body.Sub, Valid: true},
		Salt:       sql.NullString{String: encryptedSalt, Valid: true},
		SuiAddress: sql.NullString{String: body.Address, Valid: true},
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
		return
	}

	access, refresh, _ := utils.GenerateTokens(body.Sub, body.Email)

	c.JSON(http.StatusOK, gin.H{
		"email":   user.Email,
		"address": user.SuiAddress.String,
		"salt":    body.Salt,
		"access":  access,
		"refresh": refresh,
	})
}

func ZkSalt(c *gin.Context) {
	sub := c.Query("sub")
	var user models.User
	result := db.DB.First(&user, "sub = ?", sub)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	decrypted, err := utils.Decrypt(user.Salt.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decryption failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"salt": decrypted})
}

func ZKRefreshAccessToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := utils.VerifyRefreshToken(body.RefreshToken)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)

	var user models.User
	result := db.DB.First(&user, "sub = ?", sub)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	access, refresh, _ := utils.GenerateTokens(user.Sub.String, user.Email)
	c.JSON(http.StatusOK, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}
