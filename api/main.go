package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// 🌐 Web framework and middleware
	"github.com/gin-contrib/cors"
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
	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// =============================================================================
// 🔧 INITIALIZATION FUNCTIONS
// =============================================================================

func initDatabase() error {
	connStr := getEnv("DATABASE_URL", "postgresql://postgres:secure-password@localhost:5432/bruno_site?sslmode=disable")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Database connected successfully")
	return nil
}

func initRedis() error {
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	redisClient = redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Println("✅ Redis connected successfully")
	return nil
}

func initSecurityConfig() {
	// Initialize security configuration from environment variables
	secConfig = security.SecurityConfig{
		EnableMetricsAuth: getEnv("ENABLE_METRICS_AUTH", "true") == "true",
		MetricsUsername:   getEnv("METRICS_USERNAME", "admin"),
		MetricsPassword:   getEnv("METRICS_PASSWORD", "secure_password_change_me"),
		AllowedOrigins:    strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ","),
		EnableCSP:         getEnv("ENABLE_CSP", "true") == "true",
		CSPPolicy:         getEnv("CSP_POLICY", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;"),
	}

	log.Println("🔒 Security configuration initialized")
}

func initLLMService() {
	llmService = services.NewLLMService(db)

	// Test LLM service health
	if err := llmService.HealthCheck(); err != nil {
		log.Printf("⚠️ LLM service health check failed: %v", err)
		log.Println("💡 Make sure Ollama is running and the model is available")
	} else {
		log.Println("🤖 LLM service initialized and healthy")
	}
}

func initTracing() {
	// OpenTelemetry initialization (currently disabled)
	// This can be enabled when needed for distributed tracing
	log.Println("ℹ️  OpenTelemetry tracing disabled")
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     secConfig.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
