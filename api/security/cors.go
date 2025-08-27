package security

import (
	"fmt"
	"net/http"
	"strings"

	"bruno-api/utils"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           string
}

// CORSMiddleware creates a secure CORS middleware
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 🛡️ SECURITY: Validate origin
		if !isOriginAllowed(origin, config.AllowedOrigins) {
			utils.SecureLog.Warning("CORS: Blocked request from unauthorized origin", map[string]interface{}{
				"origin": origin,
				"ip":     c.ClientIP(),
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			})

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "CORS: Origin not allowed",
				"message": "Access denied from this origin",
			})
			c.Abort()
			return
		}

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		c.Header("Access-Control-Max-Age", config.MaxAge)

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			c.Abort()
			return
		}

		// Log successful CORS requests in development
		if gin.Mode() == gin.DebugMode {
			utils.SecureLog.Info("CORS: Request allowed", map[string]interface{}{
				"origin": origin,
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			})
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	// 🚨 SECURITY: Never allow wildcard origins in production
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			// Only allow wildcard in development
			if gin.Mode() == gin.DebugMode {
				return true
			}
			utils.SecureLog.Error("CORS: Wildcard origin detected in production", nil, map[string]interface{}{
				"origin": origin,
				"mode":   gin.Mode(),
			})
			return false
		}

		if allowed == origin {
			return true
		}
	}
	return false
}

// ValidateCORSConfig validates CORS configuration for security
func ValidateCORSConfig(config CORSConfig) error {
	// Check for wildcard origins
	for _, origin := range config.AllowedOrigins {
		if origin == "*" {
			return fmt.Errorf("CORS: Wildcard origin '*' is not allowed for security reasons")
		}
	}

	// Validate origins format
	for _, origin := range config.AllowedOrigins {
		if !isValidOrigin(origin) {
			return fmt.Errorf("CORS: Invalid origin format: %s", origin)
		}
	}

	// Validate methods
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true, "OPTIONS": true,
		"HEAD": true, "PATCH": true,
	}

	for _, method := range config.AllowedMethods {
		if !validMethods[method] {
			return fmt.Errorf("CORS: Invalid method: %s", method)
		}
	}

	return nil
}

// isValidOrigin validates origin format
func isValidOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	// Must start with http:// or https://
	if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
		return false
	}

	// Must have a valid domain
	parts := strings.Split(strings.TrimPrefix(origin, "http://"), "://")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]
	if domain == "" || strings.Contains(domain, " ") {
		return false
	}

	return true
}

// GetCORSConfigFromEnv creates CORS config from environment variables
func GetCORSConfigFromEnv() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{}, // Will be set by main.go
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           "12h",
	}
}
