package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// 🌐 Web framework and middleware
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	// 🔧 Environment and database
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	// 📊 Prometheus monitoring
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// 🗄️ Redis caching
	"github.com/redis/go-redis/v9"

	// 🔒 Security package
	"bruno-api/security"

	// 🤖 LLM services
	"bruno-api/services"

	// 🔐 Secure logging utilities
	"bruno-api/utils"
)

// =============================================================================
// 📋 GLOBAL VARIABLES
// =============================================================================

var (
	db          *sql.DB
	redisClient *redis.Client
	secConfig   security.SecurityConfig
	llmService  *services.LLMService
)

// =============================================================================
// 🚀 MAIN APPLICATION
// =============================================================================

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database connection
	if err := initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Redis connection
	if err := initRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize security configuration
	initSecurityConfig()

	// Initialize LLM service
	initLLMService()

	// Initialize OpenTelemetry (if enabled)
	initTracing()

	// Setup Gin router
	router := setupRouter()

	// Get port from environment
	port := getEnv("PORT", "8080")

	// Start server
	utils.SecureLog.Info("Server starting", map[string]interface{}{
		"port": port,
	})
	if err := router.Run(":" + port); err != nil {
		utils.SecureLog.Error("Failed to start server", err, nil)
		log.Fatalf("Failed to start server: %v", err)
	}
}

// =============================================================================
// 🔧 INITIALIZATION FUNCTIONS
// =============================================================================

func initDatabase() error {
	// Construct connection string from individual environment variables
	host := getEnv("DATABASE_HOST", "localhost")
	port := getEnv("DATABASE_PORT", "5432")
	user := getEnv("DATABASE_USER", "postgres")
	password := getEnv("PGPASSWORD", "secure-password")
	dbname := getEnv("DATABASE_NAME", "bruno_site")

	// 🔒 SECURITY: SSL mode configuration
	// Default to 'require' for production, 'disable' for development
	env := getEnv("APP_ENV", "development")
	sslMode := getEnv("DATABASE_SSL_MODE", "")

	if sslMode == "" {
		switch env {
		case "production":
			sslMode = "require" // 🔒 Production must use SSL
		case "staging":
			sslMode = "require" // 🔒 Staging should use SSL
		default:
			sslMode = "disable" // Development can use plain text for local testing
		}
	}

	// Validate SSL mode
	validSSLModes := []string{"disable", "require", "verify-ca", "verify-full"}
	validMode := false
	for _, mode := range validSSLModes {
		if sslMode == mode {
			validMode = true
			break
		}
	}
	if !validMode {
		return fmt.Errorf("invalid DATABASE_SSL_MODE: %s. Valid modes: %v", sslMode, validSSLModes)
	}

	// URL-encode the password to handle special characters
	encodedPassword := url.QueryEscape(password)
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, encodedPassword, host, port, dbname, sslMode)

	// Log SSL configuration (without sensitive data)
	utils.SecureLog.Info("Database connection configuration", map[string]interface{}{
		"host":     host,
		"port":     port,
		"database": dbname,
		"user":     user,
		"ssl_mode": sslMode,
		"env":      env,
	})

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Retry connection with exponential backoff
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err != nil {
			utils.SecureLog.LogDatabaseConnection(i+1, maxRetries, host, port, dbname, err)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * 2 * time.Second)
			}
		} else {
			utils.SecureLog.LogDatabaseConnection(i+1, maxRetries, host, port, dbname, nil)
			return nil
		}
	}

	return fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

func initRedis() error {
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	redisClient = redis.NewClient(opts)

	// Extract host and port for logging (without credentials)
	host := "localhost"
	port := "6379"
	if strings.HasPrefix(redisURL, "redis://") {
		parts := strings.Split(strings.TrimPrefix(redisURL, "redis://"), "@")
		if len(parts) > 1 {
			hostPort := strings.Split(parts[1], "/")[0]
			hostPortParts := strings.Split(hostPort, ":")
			if len(hostPortParts) >= 2 {
				host = hostPortParts[0]
				port = hostPortParts[1]
			}
		}
	}

	// Retry connection with exponential backoff
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := redisClient.Ping(ctx).Err(); err != nil {
			cancel()
			utils.SecureLog.LogRedisConnection(i+1, maxRetries, host, port, err)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * 2 * time.Second)
			}
		} else {
			cancel()
			utils.SecureLog.LogRedisConnection(i+1, maxRetries, host, port, nil)
			return nil
		}
	}

	return fmt.Errorf("failed to connect to Redis after %d attempts", maxRetries)
}

func initSecurityConfig() {
	// 🛡️ SECURITY: Environment-specific CORS configuration
	env := getEnv("APP_ENV", "development")
	var allowedOrigins []string

	switch env {
	case "production":
		// 🚨 CRITICAL: Production must have specific origins, never wildcards
		origins := getEnv("ALLOWED_ORIGINS", "")
		if origins == "" || origins == "*" {
			utils.SecureLog.Error("CRITICAL: ALLOWED_ORIGINS not configured for production", nil, map[string]interface{}{
				"env":     env,
				"message": "Production environment requires explicit ALLOWED_ORIGINS configuration",
			})
			// Fallback to secure defaults - adjust these for your actual domains
			allowedOrigins = []string{
				"https://yourdomain.com",
				"https://www.yourdomain.com",
			}
		} else {
			allowedOrigins = strings.Split(origins, ",")
		}
	case "staging":
		// Staging environment - more permissive but still controlled
		origins := getEnv("ALLOWED_ORIGINS", "")
		if origins == "" {
			allowedOrigins = []string{
				"https://staging.yourdomain.com",
				"http://localhost:3000", // For testing
			}
		} else {
			allowedOrigins = strings.Split(origins, ",")
		}
	default:
		// Development environment - localhost only
		origins := getEnv("ALLOWED_ORIGINS", "")
		if origins == "" || origins == "*" {
			allowedOrigins = []string{
				"http://localhost:3000",
				"http://localhost:8080",
				"http://127.0.0.1:3000",
				"http://127.0.0.1:8080",
			}
		} else {
			allowedOrigins = strings.Split(origins, ",")
		}
	}

	// Initialize security configuration from environment variables
	secConfig = security.SecurityConfig{
		EnableMetricsAuth: getEnv("ENABLE_METRICS_AUTH", "true") == "true",
		MetricsUsername:   getEnv("METRICS_USERNAME", "admin"),
		MetricsPassword:   getEnv("METRICS_PASSWORD", "secure_password_change_me"),
		AllowedOrigins:    allowedOrigins,
		EnableCSP:         getEnv("ENABLE_CSP", "true") == "true",
		CSPPolicy:         getEnv("CSP_POLICY", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;"),
	}

	utils.SecureLog.Info("Security configuration initialized", map[string]interface{}{
		"env":                   env,
		"allowed_origins_count": len(allowedOrigins),
		"cors_configured":       true,
	})
}

func initLLMService() {
	llmService = services.NewLLMService(db)

	// Test LLM service health
	if err := llmService.HealthCheck(); err != nil {
		utils.SecureLog.Warning("LLM service health check failed", map[string]interface{}{
			"error": err.Error(),
		})
		utils.SecureLog.Info("Make sure Ollama is running and the model is available")
	} else {
		utils.SecureLog.Info("LLM service initialized and healthy")
	}
}

func initTracing() {
	// OpenTelemetry initialization (currently disabled)
	// This can be enabled when needed for distributed tracing
	utils.SecureLog.Info("OpenTelemetry tracing disabled")
}

// =============================================================================
// 🌐 ROUTER SETUP
// =============================================================================

func setupRouter() *gin.Engine {
	// Set Gin mode
	if getEnv("GIN_MODE", "release") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(requestLogger())
	router.Use(errorHandler())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(security.EnhancedSecurityHeaders(secConfig))
	router.Use(security.SQLInjectionProtectionMiddleware())
	router.Use(security.RateLimitMiddleware())
	// 🛡️ SECURITY: Use custom CORS middleware with validation
	corsConfig := security.CORSConfig{
		AllowedOrigins:   secConfig.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           "12h",
	}

	// Validate CORS configuration
	if err := security.ValidateCORSConfig(corsConfig); err != nil {
		utils.SecureLog.Error("CORS configuration validation failed", err, map[string]interface{}{
			"config": corsConfig,
		})
		log.Fatalf("CORS configuration error: %v", err)
	}

	router.Use(security.CORSMiddleware(corsConfig))

	// Health check endpoint
	router.GET("/health", healthCheck)

	// Prometheus metrics endpoint (secured)
	router.GET("/metrics", security.MetricsAuthMiddleware(secConfig), gin.WrapH(promhttp.Handler()))

	// API routes (v1)
	api := router.Group("/api/v1")
	{
		// Projects
		api.GET("/projects", getProjects)
		api.GET("/projects/:id", getProject)
		api.POST("/projects", createProject)
		api.PUT("/projects/:id", updateProject)
		api.DELETE("/projects/:id", deleteProject)

		// Skills
		api.GET("/skills", getSkills)
		api.GET("/skills/:id", getSkill)
		api.POST("/skills", createSkill)
		api.PUT("/skills/:id", updateSkill)
		api.DELETE("/skills/:id", deleteSkill)

		// Experiences
		api.GET("/experiences", getExperiences)
		api.GET("/experiences/:id", getExperience)
		api.POST("/experiences", createExperience)
		api.PUT("/experiences/:id", updateExperience)
		api.DELETE("/experiences/:id", deleteExperience)

		// Content
		api.GET("/content", getContent)
		api.GET("/content/:type", getContentByType)
		api.POST("/content", createContent)
		api.PUT("/content/:id", updateContent)
		api.DELETE("/content/:id", deleteContent)

		// About
		api.GET("/about", getAbout)
		api.PUT("/about", updateAbout)

		// Contact
		api.GET("/contact", getContact)
		api.PUT("/contact", updateContact)

		// 🤖 AI Chat endpoint
		api.POST("/chat", handleChat)
		api.GET("/chat/health", handleChatHealth)

		// 📊 Analytics endpoint
		api.POST("/analytics/track", handleAnalyticsTrack)
	}

	// Legacy API routes (for frontend compatibility)
	legacyApi := router.Group("/api")
	{
		// Projects
		legacyApi.GET("/projects", getProjects)
		legacyApi.GET("/projects/:id", getProject)
		legacyApi.POST("/projects", createProject)
		legacyApi.PUT("/projects/:id", updateProject)
		legacyApi.DELETE("/projects/:id", deleteProject)

		// Skills
		legacyApi.GET("/skills", getSkills)
		legacyApi.GET("/skills/:id", getSkill)
		legacyApi.POST("/skills", createSkill)
		legacyApi.PUT("/skills/:id", updateSkill)
		legacyApi.DELETE("/skills/:id", deleteSkill)

		// Experiences
		legacyApi.GET("/experiences", getExperiences)
		legacyApi.GET("/experiences/:id", getExperience)
		legacyApi.POST("/experiences", createExperience)
		legacyApi.PUT("/experiences/:id", updateExperience)
		legacyApi.DELETE("/experiences/:id", deleteExperience)

		// Content
		legacyApi.GET("/content", getContent)
		legacyApi.GET("/content/:type", getContentByType)
		legacyApi.POST("/content", createContent)
		legacyApi.PUT("/content/:id", updateContent)
		legacyApi.DELETE("/content/:id", deleteContent)

		// About
		legacyApi.GET("/about", getAbout)
		legacyApi.PUT("/about", updateAbout)

		// Contact
		legacyApi.GET("/contact", getContact)
		legacyApi.PUT("/contact", updateContact)

		// 🤖 AI Chat endpoint
		legacyApi.POST("/chat", handleChat)
		legacyApi.GET("/chat/health", handleChatHealth)

		// 📊 Analytics endpoint
		legacyApi.POST("/analytics/track", handleAnalyticsTrack)
	}

	return router
}

// =============================================================================
// 🛠️ UTILITY FUNCTIONS
// =============================================================================

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

// Helper function to truncate strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "bruno-api",
	})
}

// =============================================================================
// 🤖 CHAT HANDLERS
// =============================================================================

func handleChat(c *gin.Context) {
	startTime := time.Now()
	requestID := fmt.Sprintf("chat_handler_%d", startTime.UnixNano())

	log.Printf("🤖 [%s] Chat request received", requestID)
	log.Printf("   📍 Remote IP: %s", c.ClientIP())
	log.Printf("   📍 User Agent: %s", c.GetHeader("User-Agent"))
	log.Printf("   📍 Content-Type: %s", c.GetHeader("Content-Type"))
	log.Printf("   📍 Content-Length: %s", c.GetHeader("Content-Length"))

	var request services.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("❌ [%s] JSON binding failed: %v", requestID, err)
		log.Printf("   📄 Request body: %s", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ [%s] JSON binding successful", requestID)
	log.Printf("   📝 Message: %s", truncateString(request.Message, 100))
	log.Printf("   📝 Context: %s", truncateString(request.Context, 50))

	// Validate message is not empty
	if strings.TrimSpace(request.Message) == "" {
		log.Printf("❌ [%s] Empty message received", requestID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Message cannot be empty",
		})
		return
	}

	log.Printf("🔄 [%s] Processing chat request...", requestID)

	// Process chat request
	response, err := llmService.ProcessChat(request)
	if err != nil {
		log.Printf("❌ [%s] Chat processing error: %v", requestID, err)
		log.Printf("   🔍 Error type: %T", err)
		log.Printf("   🔍 Full error details: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process chat request",
			"details": err.Error(),
		})
		return
	}

	duration := time.Since(startTime)
	log.Printf("✅ [%s] Chat request completed successfully in %v", requestID, duration)
	log.Printf("   📤 Response length: %d chars", len(response.Response))
	log.Printf("   🎯 Model used: %s", response.Model)

	c.JSON(http.StatusOK, response)
}

func handleChatHealth(c *gin.Context) {
	if err := llmService.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "unhealthy",
			"error":     err.Error(),
			"timestamp": time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"provider":  "ollama",
		"model":     getEnv("GEMMA_MODEL", "gemma3n:e4b"),
		"timestamp": time.Now().UTC(),
	})
}
