package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"backend/models"
	"gorm.io/gorm"
)

const (
	authCookieName = "auth_token"
	cookieExpiry   = 72 * time.Hour // 3 days
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	return []byte(secret)
}

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
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=50"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8,max=50"`
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex = regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,50}$`)
)

func sanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

func (h *AuthHandler) Register(c *gin.Context) {
	// Verify CSRF token
	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
		return
	}

	cookieCsrfToken, err := c.Cookie("csrf_token")
	if err != nil || csrfToken != cookieCsrfToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize inputs
	req.Email = sanitizeInput(req.Email)
	req.Password = sanitizeInput(req.Password)
	req.ConfirmPassword = sanitizeInput(req.ConfirmPassword)

	// Validate email format
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate password format
	if !passwordRegex.MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password format"})
		return
	}

	// Check password match
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Create new user
	newUser := models.User{
		Email:    req.Email,
		Admin:    false, // Default to non-admin
	}

	// Hash password
	if err := newUser.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Save user to database
	if err := h.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Verify CSRF token
	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
		return
	}

	cookieCsrfToken, err := c.Cookie("csrf_token")
	if err != nil || csrfToken != cookieCsrfToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize inputs
	req.Email = sanitizeInput(req.Email)
	req.Password = sanitizeInput(req.Password)

	// Validate email format
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate password format
	if !passwordRegex.MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password format"})
		return
	}

	// Use parameterized query to prevent SQL injection
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
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate new CSRF token
	newCsrfToken := generateSecureToken()
	
	// Set auth cookie with proper domain and secure settings
	domain := ""
	if os.Getenv("ENV") == "production" {
		domain = os.Getenv("DOMAIN")
	}
	c.SetCookie(authCookieName, tokenString, int(cookieExpiry.Seconds()), "/", domain, os.Getenv("ENV") == "production", true)
	// Set CSRF cookie
	c.SetCookie("csrf_token", newCsrfToken, int(cookieExpiry.Seconds()), "/", domain, os.Getenv("ENV") == "production", true)
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Login successful",
		"admin":      user.Admin,
		"email":      user.Email,
		"csrf_token": newCsrfToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(authCookieName, "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify JWT token
		tokenString, err := c.Cookie(authCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
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

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}

func (h *AuthHandler) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First verify auth using standard middleware
		h.AuthMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// Then check if user is admin
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if !user.(models.User).Admin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
	}
}

func (h *AuthHandler) refreshToken(claims *JWTClaims) (string, error) {
	// Generate new token with same claims but new expiration
	expirationTime := time.Now().Add(cookieExpiry)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required,min=8,max=50"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=50"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Get user from context
	user := c.MustGet("user").(models.User)

	// Verify CSRF token
	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
		return
	}

	cookieCsrfToken, err := c.Cookie("csrf_token")
	if err != nil || csrfToken != cookieCsrfToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize inputs
	req.CurrentPassword = sanitizeInput(req.CurrentPassword)
	req.NewPassword = sanitizeInput(req.NewPassword)

	// Validate new password format
	if !passwordRegex.MatchString(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid new password format"})
		return
	}

	// Check current password
	if err := user.CheckPassword(req.CurrentPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	if err := user.HashPassword(req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password in database
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Logout user by clearing auth cookie
	c.SetCookie(authCookieName, "", -1, "/", "", true, true)
	
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func generateSecureToken() string {
	// Generate cryptographically secure random token
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		log.Printf("ERROR: Failed to generate secure token: %v", err)
		panic("failed to generate secure token")
	}
	// Use URL-safe base64 encoding without padding
	return base64.RawURLEncoding.EncodeToString(token)
}

func (h *AuthHandler) GetCSRFToken(c *gin.Context) {
	// Generate new CSRF token
	csrfToken := generateSecureToken()
	
	// Set CSRF cookie
	c.SetCookie("csrf_token", csrfToken, int(cookieExpiry.Seconds()), "/", "", false, true)
	
	c.JSON(http.StatusOK, gin.H{
		"csrf_token": csrfToken,
	})
}

func (h *AuthHandler) CheckAuth(c *gin.Context) {
	// Get auth token from cookie
	tokenString, err := c.Cookie(authCookieName)
	if err != nil {
		log.Printf("INFO: CheckAuth - No auth cookie found")
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"error":         "No auth cookie found",
		})
		return
	}

	// Verify token
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		log.Printf("WARN: CheckAuth - Invalid token: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"error":         "Invalid token: " + err.Error(),
		})
		return
	}

	if !token.Valid {
		log.Printf("INFO: CheckAuth - Expired token")
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"error":         "Expired token",
		})
		return
	}

	// Get user from database
	var user models.User
	if err := h.DB.First(&user, claims.UserID).Error; err != nil {
		log.Printf("ERROR: CheckAuth - User not found: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"error":         "User not found: " + err.Error(),
		})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"admin":         user.Admin,
		"userId":        user.ID,
		"email":         user.Email,
	})
}
