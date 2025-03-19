package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"backend/models"
	"gorm.io/gorm"
)

const (
	authCookieName = "auth_token"
	cookieExpiry   = 72 * time.Hour // 3 days
	jwtSecret      = "your-256-bit-secret" // In production, use environment variable
)

type JWTClaims struct {
	UserID uint `json:"userId"`
	Admin  bool `json:"admin"`
	jwt.RegisteredClaims
}

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

	// Generate JWT token
	expirationTime := time.Now().Add(cookieExpiry)
	claims := &JWTClaims{
		UserID: user.ID,
		Admin:  user.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate CSRF token
	csrfToken := generateSecureToken()
	
	// Set auth cookie
	c.SetCookie(authCookieName, tokenString, int(cookieExpiry.Seconds()), "/", "", true, true)
	// Set CSRF cookie
	c.SetCookie("csrf_token", csrfToken, int(cookieExpiry.Seconds()), "/", "", false, true)
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Login successful",
		"admin":      user.Admin,
		"csrf_token": csrfToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(authCookieName, "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify JWT token
		tokenString, err := c.Cookie(authCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check if token is about to expire
		if time.Until(claims.ExpiresAt.Time) < 1*time.Hour {
			// Refresh token
			newToken, err := h.refreshToken(claims)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token refresh failed"})
				return
			}
			c.SetCookie(authCookieName, newToken, int(cookieExpiry.Seconds()), "/", "", true, true)
		}

		// Get user from database
		var user models.User
		if err := h.DB.First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Check if user is admin
		if !user.Admin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}

func (h *AuthHandler) refreshToken(claims *JWTClaims) (string, error) {
	// Generate new token with same claims but new expiration
	expirationTime := time.Now().Add(cookieExpiry)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func generateSecureToken() string {
	// In production, use a proper token generation method like JWT
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic("failed to generate token")
	}
	return base64.URLEncoding.EncodeToString(token)
}
