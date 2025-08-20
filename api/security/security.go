package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// =============================================================================
// 🔒 SECURITY CONSTANTS AND CONFIGURATION
// =============================================================================

const (
	// Rate limiting constants
	MaxRequestsPerMinute = 100
	MaxRequestsPerHour   = 1000

	// Input validation constants
	MaxTitleLength       = 255
	MaxDescriptionLength = 5000
	MaxURLLength         = 2048
	MaxEmailLength       = 254

	// Token constants
	TokenLength = 32
	TokenExpiry = 24 * time.Hour
)

// SecurityConfig holds security configuration
type SecurityConfig struct {
	EnableMetricsAuth bool
	MetricsUsername   string
	MetricsPassword   string
	AllowedOrigins    []string
	EnableCSP         bool
	CSPPolicy         string
}

// =============================================================================
// 🔍 INPUT VALIDATION AND SANITIZATION
// =============================================================================

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult holds validation results
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// SanitizeString removes potentially dangerous characters and normalizes input
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// Remove null bytes and control characters
	re := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	input = re.ReplaceAllString(input, "")

	// HTML escape to prevent XSS
	input = html.EscapeString(input)

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}

// ValidateAndSanitizeTitle validates and sanitizes project/skill titles
func ValidateAndSanitizeTitle(title string) (string, *ValidationError) {
	if title == "" {
		return "", &ValidationError{
			Field:   "title",
			Message: "Title is required",
		}
	}

	// Check for potentially dangerous patterns before sanitization
	dangerousPatterns := []string{
		"<script", "javascript:", "onload=", "onerror=", "onclick=",
		"<iframe", "<object", "<embed", "data:text/html",
	}

	lowerTitle := strings.ToLower(title)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerTitle, pattern) {
			return "", &ValidationError{
				Field:   "title",
				Message: "Title contains potentially dangerous content",
			}
		}
	}

	// Sanitize input
	title = SanitizeString(title)

	// Check length
	if len(title) > MaxTitleLength {
		return "", &ValidationError{
			Field:   "title",
			Message: fmt.Sprintf("Title must be %d characters or less", MaxTitleLength),
		}
	}

	return title, nil
}

// ValidateAndSanitizeDescription validates and sanitizes descriptions
func ValidateAndSanitizeDescription(description string) (string, *ValidationError) {
	if description == "" {
		return "", &ValidationError{
			Field:   "description",
			Message: "Description is required",
		}
	}

	// Sanitize input
	description = SanitizeString(description)

	// Check length
	if len(description) > MaxDescriptionLength {
		return "", &ValidationError{
			Field:   "description",
			Message: fmt.Sprintf("Description must be %d characters or less", MaxDescriptionLength),
		}
	}

	return description, nil
}

// ValidateAndSanitizeURL validates and sanitizes URLs
func ValidateAndSanitizeURL(urlStr, fieldName string) (string, *ValidationError) {
	if urlStr == "" {
		return "", nil // URLs can be optional
	}

	// Sanitize input
	urlStr = SanitizeString(urlStr)

	// Check length
	if len(urlStr) > MaxURLLength {
		return "", &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("URL must be %d characters or less", MaxURLLength),
		}
	}

	// Basic URL validation
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return "", &ValidationError{
			Field:   fieldName,
			Message: "URL must start with http:// or https://",
		}
	}

	// Check for potentially dangerous URLs
	dangerousPatterns := []string{
		"javascript:", "data:text/html", "vbscript:", "file://",
	}

	lowerURL := strings.ToLower(urlStr)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerURL, pattern) {
			return "", &ValidationError{
				Field:   fieldName,
				Message: "URL contains potentially dangerous content",
			}
		}
	}

	return urlStr, nil
}

// ValidateAndSanitizeEmail validates and sanitizes email addresses
func ValidateAndSanitizeEmail(email string) (string, *ValidationError) {
	if email == "" {
		return "", &ValidationError{
			Field:   "email",
			Message: "Email is required",
		}
	}

	// Sanitize input
	email = SanitizeString(email)

	// Check length
	if len(email) > MaxEmailLength {
		return "", &ValidationError{
			Field:   "email",
			Message: fmt.Sprintf("Email must be %d characters or less", MaxEmailLength),
		}
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "", &ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		}
	}

	return email, nil
}

// ValidateInteger validates integer parameters
func ValidateInteger(value string, fieldName string, min, max int) (int, *ValidationError) {
	if value == "" {
		return 0, &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s is required", fieldName),
		}
	}

	// Check if it's a valid integer
	intValue, err := parseInt(value)
	if err != nil {
		return 0, &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s must be a valid integer", fieldName),
		}
	}

	// Check range
	if intValue < min || intValue > max {
		return 0, &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s must be between %d and %d", fieldName, min, max),
		}
	}

	return intValue, nil
}

// =============================================================================
// 🔐 AUTHENTICATION AND AUTHORIZATION
// =============================================================================

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashPassword creates a bcrypt hash of a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SecureCompare performs a constant-time comparison
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// =============================================================================
// 🛡️ SECURITY MIDDLEWARE
// =============================================================================

// MetricsAuthMiddleware provides authentication for metrics endpoint
func MetricsAuthMiddleware(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.EnableMetricsAuth {
			c.Next()
			return
		}

		// Check for Basic Auth
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Metrics"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Validate credentials
		if username != config.MetricsUsername || !SecureCompare(password, config.MetricsPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// EnhancedSecurityHeaders adds comprehensive security headers
func EnhancedSecurityHeaders(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Content Security Policy
		if config.EnableCSP {
			c.Header("Content-Security-Policy", config.CSPPolicy)
		}

		// HSTS (HTTP Strict Transport Security)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting with Redis (placeholder)
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement Redis-based rate limiting
		// For now, just pass through
		c.Next()
	}
}

// SQLInjectionProtectionMiddleware adds additional SQL injection protection
func SQLInjectionProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check query parameters for SQL injection patterns
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsSQLInjectionPattern(value) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input detected"})
					c.Abort()
					return
				}
			}
		}

		// Check path parameters
		for _, param := range c.Params {
			if containsSQLInjectionPattern(param.Value) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input detected"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// =============================================================================
// 🛠️ UTILITY FUNCTIONS
// =============================================================================

// containsSQLInjectionPattern checks for common SQL injection patterns
func containsSQLInjectionPattern(input string) bool {
	if input == "" {
		return false
	}

	// Convert to lowercase for case-insensitive matching
	lowerInput := strings.ToLower(input)

	// Common SQL injection patterns
	patterns := []string{
		"union select", "union all select", "drop table", "delete from",
		"insert into", "update set", "alter table", "create table",
		"exec(", "execute(", "xp_", "sp_", "--", "/*", "*/",
		"waitfor delay", "benchmark(", "sleep(", "load_file(",
		"into outfile", "into dumpfile", "information_schema",
		"update ", "set ", "where ", "select ", "from ",
	}

	for _, pattern := range patterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	return false
}

// parseInt safely parses an integer string
func parseInt(s string) (int, error) {
	// Remove any whitespace
	s = strings.TrimSpace(s)

	// Check for non-numeric characters
	for _, char := range s {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("non-numeric character found")
		}
	}

	// Use strconv.Atoi for actual parsing
	return strconv.Atoi(s)
}
