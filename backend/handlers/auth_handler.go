package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"backend/models"
	"gorm.io/gorm"
)

const (
	authCookieName = "auth_token"
	cookieExpiry   = 72 * time.Hour // 3 days
)

type AuthHandler struct {
	DB *gorm.DB
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate auth token (in real implementation, use JWT or similar)
	token := generateSecureToken()
	
	c.SetCookie(authCookieName, token, int(cookieExpiry.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"admin":   user.Admin,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(authCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func generateSecureToken() string {
	// In production, use a proper token generation method like JWT
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic("failed to generate token")
	}
	return base64.URLEncoding.EncodeToString(token)
}
