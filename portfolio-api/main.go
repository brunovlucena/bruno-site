package main

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	// Web framework and middleware
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	// Environment and database
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	// Prometheus monitoring
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Redis caching
	"github.com/redis/go-redis/v9"

	// OpenTelemetry tracing - Core packages
	"go.opentelemetry.io/otel"                                        // Global tracer provider
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc" // gRPC exporter for traces
	"go.opentelemetry.io/otel/sdk/resource"                           // Resource information (service name, version)
	sdktrace "go.opentelemetry.io/otel/sdk/trace"                     // Trace provider implementation
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"                // Semantic conventions (standard attributes)

	// OpenTelemetry Gin middleware for automatic HTTP tracing
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Project struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Modules     int    `json:"modules"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	Active      bool   `json:"active"` // Controls project visibility
}

type Content struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AboutData struct {
	Description string `json:"description"`
	Highlights  []struct {
		Icon string `json:"icon"`
		Text string `json:"text"`
	} `json:"highlights"`
}

type ContactData struct {
	Email        string `json:"email"`
	Location     string `json:"location"`
	LinkedIn     string `json:"linkedin"`
	GitHub       string `json:"github"`
	Availability string `json:"availability"`
}

var db *sql.DB
var rdb *redis.Client

// Rate limiting map
var rateLimitMap = make(map[string][]time.Time)
var rateLimitMutex sync.Mutex

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// initTracer initializes OpenTelemetry tracing with gRPC exporter
// This function sets up the complete tracing infrastructure for the application
func initTracer() (*sdktrace.TracerProvider, error) {
	// Create OTLP gRPC exporter (more efficient than HTTP)
	// gRPC provides better performance, reliability, and connection handling
	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint("localhost:4317"), // OpenTelemetry Collector gRPC endpoint
		otlptracegrpc.WithInsecure(),                 // Disable TLS for local development
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information
	// This metadata helps identify the service in tracing backends
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("portfolio-api"), // Service name for identification
			semconv.ServiceVersion("1.0.0"),      // Service version for tracking
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider with batched exporter
	// Batched exporter improves performance by sending traces in batches
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // Batch traces before sending
		sdktrace.WithResource(res),     // Attach service metadata
	)

	// Set global trace provider so it can be accessed throughout the application
	otel.SetTracerProvider(tp)

	return tp, nil
}

// Validation patterns
var (
	urlPattern   = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize OpenTelemetry tracing infrastructure
	// This sets up the global tracer provider and connects to the collector
	tp, err := initTracer()
	if err != nil {
		log.Printf("Failed to initialize tracer: %v", err)
	} else {
		// Ensure graceful shutdown of tracing when the application exits
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	// Initialize database connection
	initDB()
	defer db.Close()

	// Initialize Redis connection
	initRedis()
	defer rdb.Close()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	r := gin.Default()

	// Add compression middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Add Prometheus middleware globally
	r.Use(prometheusMiddleware())

	// Add OpenTelemetry Gin middleware for automatic HTTP request tracing
	// This automatically creates spans for all HTTP requests with timing and metadata
	r.Use(otelgin.Middleware("portfolio-api"))

	// Configure CORS with stricter settings
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = false
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))

	// Health check endpoint with minimal information
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes with rate limiting
	api := r.Group("/api/v1")
	{
		api.GET("/projects", rateLimitMiddleware(), getProjects)
		api.GET("/about", rateLimitMiddleware(), getAbout)
		api.GET("/contact", rateLimitMiddleware(), getContact)
		api.GET("/content/skills", rateLimitMiddleware(), getSkills)
		api.GET("/content/experience", rateLimitMiddleware(), getExperience)
		api.POST("/analytics/visit", rateLimitMiddleware(), trackVisit)
	}

	// Admin routes for project management
	admin := r.Group("/admin")
	{
		admin.GET("/projects", rateLimitMiddleware(), getAllProjects) // Get all projects (including inactive)
		admin.PUT("/projects/:id/activate", rateLimitMiddleware(), activateProject)
		admin.PUT("/projects/:id/deactivate", rateLimitMiddleware(), deactivateProject)
		admin.GET("/projects/stats", rateLimitMiddleware(), getProjectStats)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting portfolio API server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() {
	// Get database connection string from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback to individual environment variables
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "portfolio_user")
		dbPassword := getEnv("DB_PASSWORD", "portfolio_password")
		dbName := getEnv("DB_NAME", "portfolio")

		dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")
}

func initRedis() {
	// Get Redis connection string from environment
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}

	rdb = redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Successfully connected to Redis")
}

func getProjects(c *gin.Context) {
	// Get tracer instance for this service
	// This creates a tracer that will be used to create spans for this function
	tracer := otel.Tracer("portfolio-api")

	// Start a new span for this function call
	// The span will track the entire duration of this function execution
	ctx, span := tracer.Start(c.Request.Context(), "getProjects")
	defer span.End() // Ensure the span is ended when the function returns

	// Try to get from cache first (Redis)
	cached, err := rdb.Get(ctx, "projects").Result()
	if err == nil {
		var projects []Project
		if err := json.Unmarshal([]byte(cached), &projects); err == nil {
			// Add event to span indicating cache hit for debugging
			span.AddEvent("cache_hit")

			// Add ETag for HTTP caching
			etag := generateETag(projects)
			if c.GetHeader("If-None-Match") == etag {
				c.Status(http.StatusNotModified)
				return
			}
			c.Header("ETag", etag)
			c.JSON(http.StatusOK, projects)
			return
		}
	}

	// Add event to span indicating cache miss for debugging
	span.AddEvent("cache_miss")

	// Create a child span for database query operation
	// This allows us to track database performance separately
	_, dbSpan := tracer.Start(ctx, "db.query.projects")
	rows, err := db.Query(`
		SELECT id, title, description, type, modules, github_url, live_url, active
		FROM projects 
		WHERE active = true
		ORDER BY "order" ASC, id ASC
	`)
	dbSpan.End() // End the database span

	if err != nil {
		// Record the error in the span for debugging and alerting
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	defer rows.Close()

	// Process database results
	var projects []Project
	for rows.Next() {
		var p Project
		var githubURL, liveURL sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Type, &p.Modules, &githubURL, &liveURL, &p.Active); err != nil {
			continue
		}
		// Use live_url if available, otherwise use github_url
		if liveURL.Valid {
			p.URL = liveURL.String
		} else if githubURL.Valid {
			p.URL = githubURL.String
		}
		projects = append(projects, p)
	}

	// Create a child span for Redis cache operation
	// This allows us to track cache performance separately
	_, cacheSpan := tracer.Start(ctx, "redis.set.projects")
	if data, err := json.Marshal(projects); err == nil {
		rdb.Set(ctx, "projects", data, 300000000000) // 5 minutes in nanoseconds
	}
	cacheSpan.End() // End the cache span

	// Add database statement attribute to the span for debugging
	span.SetAttributes(semconv.DBStatement("SELECT projects"))
	c.JSON(http.StatusOK, projects)
}

// getAllProjects returns all projects (including inactive ones) for admin management
func getAllProjects(c *gin.Context) {
	// Get tracer for this function
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "getAllProjects")
	defer span.End()

	// Query database for all projects (including inactive)
	_, dbSpan := tracer.Start(ctx, "db.query.all_projects")
	rows, err := db.Query(`
		SELECT id, title, description, type, modules, github_url, live_url, active
		FROM projects 
		ORDER BY "order" ASC, id ASC
	`)
	dbSpan.End()

	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var githubURL, liveURL sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Type, &p.Modules, &githubURL, &liveURL, &p.Active); err != nil {
			continue
		}
		// Use live_url if available, otherwise use github_url
		if liveURL.Valid {
			p.URL = liveURL.String
		} else if githubURL.Valid {
			p.URL = githubURL.String
		}
		projects = append(projects, p)
	}

	span.SetAttributes(semconv.DBStatement("SELECT all projects"))
	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"total":    len(projects),
		"active":   countActiveProjects(projects),
		"inactive": countInactiveProjects(projects),
	})
}

// activateProject activates a project by setting active = true
func activateProject(c *gin.Context) {
	// Get tracer for this function
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "activateProject")
	defer span.End()

	// Get project ID from URL parameter
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	// Update project to active
	_, dbSpan := tracer.Start(ctx, "db.update.activate_project")
	result, err := db.Exec("UPDATE projects SET active = true WHERE id = $1", projectID)
	dbSpan.End()

	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate project"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Clear cache to ensure fresh data
	rdb.Del(ctx, "projects")

	span.AddEvent("project_activated")
	c.JSON(http.StatusOK, gin.H{
		"message":    "Project activated successfully",
		"project_id": projectID,
	})
}

// deactivateProject deactivates a project by setting active = false
func deactivateProject(c *gin.Context) {
	// Get tracer for this function
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "deactivateProject")
	defer span.End()

	// Get project ID from URL parameter
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	// Update project to inactive
	_, dbSpan := tracer.Start(ctx, "db.update.deactivate_project")
	result, err := db.Exec("UPDATE projects SET active = false WHERE id = $1", projectID)
	dbSpan.End()

	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate project"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Clear cache to ensure fresh data
	rdb.Del(ctx, "projects")

	span.AddEvent("project_deactivated")
	c.JSON(http.StatusOK, gin.H{
		"message":    "Project deactivated successfully",
		"project_id": projectID,
	})
}

// getProjectStats returns statistics about projects
func getProjectStats(c *gin.Context) {
	// Get tracer for this function
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "getProjectStats")
	defer span.End()

	// Query database for project statistics
	_, dbSpan := tracer.Start(ctx, "db.query.project_stats")
	var total, active, inactive int
	err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&total)
	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project stats"})
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM projects WHERE active = true").Scan(&active)
	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active project count"})
		return
	}

	inactive = total - active
	dbSpan.End()

	c.JSON(http.StatusOK, gin.H{
		"total":             total,
		"active":            active,
		"inactive":          inactive,
		"active_percentage": float64(active) / float64(total) * 100,
	})
}

// Helper functions for counting projects
func countActiveProjects(projects []Project) int {
	count := 0
	for _, p := range projects {
		if p.Active {
			count++
		}
	}
	return count
}

func countInactiveProjects(projects []Project) int {
	count := 0
	for _, p := range projects {
		if !p.Active {
			count++
		}
	}
	return count
}

func getAbout(c *gin.Context) {
	// Try to get from cache first
	ctx := context.Background()
	cached, err := rdb.Get(ctx, "about").Result()
	if err == nil {
		var aboutData AboutData
		if err := json.Unmarshal([]byte(cached), &aboutData); err == nil {
			c.JSON(http.StatusOK, aboutData)
			return
		}
	}

	// Query database for about content
	var description string
	err = db.QueryRow("SELECT value->>'description' FROM content WHERE key = 'about'").Scan(&description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch about data"})
		return
	}

	// Query highlights from content table
	rows, err := db.Query("SELECT value->>'highlights' FROM content WHERE key = 'about'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch highlights"})
		return
	}
	defer rows.Close()

	var highlights []struct {
		Icon string `json:"icon"`
		Text string `json:"text"`
	}

	// For now, return empty highlights since we need to parse JSON
	// TODO: Implement proper JSON parsing for highlights

	aboutData := AboutData{
		Description: description,
		Highlights:  highlights,
	}

	// Cache the result for 10 minutes
	if data, err := json.Marshal(aboutData); err == nil {
		rdb.Set(ctx, "about", data, 600000000000) // 10 minutes in nanoseconds
	}

	c.JSON(http.StatusOK, aboutData)
}

func getContact(c *gin.Context) {
	// Try to get from cache first
	ctx := context.Background()
	cached, err := rdb.Get(ctx, "contact").Result()
	if err == nil {
		var contactData ContactData
		if err := json.Unmarshal([]byte(cached), &contactData); err == nil {
			c.JSON(http.StatusOK, contactData)
			return
		}
	}

	// Query database for contact information
	var email, location, linkedin, github, availability string

	err = db.QueryRow("SELECT value->>'email' FROM content WHERE key = 'contact'").Scan(&email)
	if err != nil {
		email = "bruno.lucena@example.com"
	}

	err = db.QueryRow("SELECT value->>'location' FROM content WHERE key = 'contact'").Scan(&location)
	if err != nil {
		location = "Brazil"
	}

	err = db.QueryRow("SELECT value->>'linkedin' FROM content WHERE key = 'contact'").Scan(&linkedin)
	if err != nil {
		linkedin = "https://www.linkedin.com/in/bvlucena"
	}

	err = db.QueryRow("SELECT value->>'github' FROM content WHERE key = 'contact'").Scan(&github)
	if err != nil {
		github = "https://github.com/brunovlucena"
	}

	err = db.QueryRow("SELECT value->>'availability' FROM content WHERE key = 'contact'").Scan(&availability)
	if err != nil {
		availability = "Open to new opportunities in SRE, DevSecOps, and AI Engineering roles."
	}

	contactData := ContactData{
		Email:        email,
		Location:     location,
		LinkedIn:     linkedin,
		GitHub:       github,
		Availability: availability,
	}

	// Cache the result for 10 minutes
	if data, err := json.Marshal(contactData); err == nil {
		rdb.Set(ctx, "contact", data, 600000000000) // 10 minutes in nanoseconds
	}

	c.JSON(http.StatusOK, contactData)
}

func trackVisit(c *gin.Context) {
	var visit struct {
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		Referrer  string `json:"referrer"`
	}

	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Input validation and sanitization
	if visit.IP == "" {
		visit.IP = c.ClientIP()
	}

	// Validate IP format (basic validation)
	if !isValidIP(visit.IP) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address format"})
		return
	}

	// Sanitize user agent and referrer
	visit.UserAgent = sanitizeString(visit.UserAgent)
	visit.Referrer = sanitizeString(visit.Referrer)

	// Validate referrer URL if provided
	if visit.Referrer != "" && !validateURL(visit.Referrer) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referrer URL"})
		return
	}

	// Store visit in database
	_, err := db.Exec(`
		INSERT INTO visitors (ip, user_agent, first_visit, last_visit, visit_count)
		VALUES ($1, $2, NOW(), NOW(), 1)
		ON CONFLICT (ip) DO UPDATE SET
			last_visit = NOW(),
			visit_count = visitors.visit_count + 1
	`, visit.IP, visit.UserAgent)

	if err != nil {
		log.Printf("Failed to store analytics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track visit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Rate limiting middleware
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		rateLimitMutex.Lock()
		now := time.Now()

		// Clean old requests (older than 1 minute)
		var validRequests []time.Time
		for _, reqTime := range rateLimitMap[clientIP] {
			if now.Sub(reqTime) < time.Minute {
				validRequests = append(validRequests, reqTime)
			}
		}
		rateLimitMap[clientIP] = validRequests

		// Check rate limit (100 requests per minute)
		if len(rateLimitMap[clientIP]) >= 100 {
			rateLimitMutex.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		// Add current request
		rateLimitMap[clientIP] = append(rateLimitMap[clientIP], now)
		rateLimitMutex.Unlock()

		c.Next()
	}
}

// Input validation functions
func validateURL(url string) bool {
	return urlPattern.MatchString(url)
}

func validateEmail(email string) bool {
	return emailPattern.MatchString(email)
}

func sanitizeString(input string) string {
	// Remove potentially dangerous characters
	dangerous := []string{"<script>", "</script>", "javascript:", "onload=", "onerror="}
	result := input
	for _, danger := range dangerous {
		result = strings.ReplaceAll(result, danger, "")
	}
	return result
}

func isValidIP(ip string) bool {
	// Basic IP validation
	if ip == "localhost" || ip == "127.0.0.1" || ip == "::1" {
		return true
	}

	// Check for IPv4 format
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		for _, part := range parts {
			if len(part) == 0 || len(part) > 3 {
				return false
			}
			for _, char := range part {
				if char < '0' || char > '9' {
					return false
				}
			}
			num := 0
			for _, char := range part {
				num = num*10 + int(char-'0')
			}
			if num > 255 {
				return false
			}
			if part[0] == '0' && len(part) > 1 {
				return false
			}
		}
		return true
	}

	return false
}

// Prometheus middleware to record metrics for all requests
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			httpRequestDuration.WithLabelValues(path, method, fmt.Sprint(c.Writer.Status())).Observe(v)
		}))
		c.Next()
		timer.ObserveDuration()

		httpRequestsTotal.WithLabelValues(path, method, fmt.Sprint(c.Writer.Status())).Inc()
	}
}

func generateETag(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	hash := md5.Sum(jsonData)
	return fmt.Sprintf("\"%x\"", hash)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Skill represents a technical skill
type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Proficiency int    `json:"proficiency"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
}

// Experience represents work experience
type Experience struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Company      string   `json:"company"`
	StartDate    string   `json:"start_date"`
	EndDate      *string  `json:"end_date"`
	Current      bool     `json:"current"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	Order        int      `json:"order"`
}

func getSkills(c *gin.Context) {
	// Get tracer instance for this service
	tracer := otel.Tracer("portfolio-api")

	// Start a new span for this function call
	ctx, span := tracer.Start(c.Request.Context(), "getSkills")
	defer span.End()

	// Try to get from cache first (Redis)
	cached, err := rdb.Get(ctx, "skills").Result()
	if err == nil {
		var skills []Skill
		if err := json.Unmarshal([]byte(cached), &skills); err == nil {
			span.AddEvent("cache_hit")
			etag := generateETag(skills)
			c.Header("ETag", etag)
			c.JSON(http.StatusOK, skills)
			return
		}
	}

	// Query database for skills
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, category, proficiency, icon, "order" 
		FROM skills 
		ORDER BY "order", name
	`)
	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var skill Skill
		err := rows.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Proficiency, &skill.Icon, &skill.Order)
		if err != nil {
			span.RecordError(err)
			continue
		}
		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan skills"})
		return
	}

	// Cache the result in Redis (5 minutes)
	if skillsData, err := json.Marshal(skills); err == nil {
		rdb.Set(ctx, "skills", skillsData, 5*time.Minute)
	}

	etag := generateETag(skills)
	c.Header("ETag", etag)
	c.JSON(http.StatusOK, skills)
}

func getExperience(c *gin.Context) {
	// Get tracer instance for this service
	tracer := otel.Tracer("portfolio-api")

	// Start a new span for this function call
	ctx, span := tracer.Start(c.Request.Context(), "getExperience")
	defer span.End()

	// Try to get from cache first (Redis)
	cached, err := rdb.Get(ctx, "experience").Result()
	if err == nil {
		var experience []Experience
		if err := json.Unmarshal([]byte(cached), &experience); err == nil {
			span.AddEvent("cache_hit")
			etag := generateETag(experience)
			c.Header("ETag", etag)
			c.JSON(http.StatusOK, experience)
			return
		}
	}

	// Query database for experience
	rows, err := db.QueryContext(ctx, `
		SELECT id, title, company, start_date, end_date, current, description, technologies, "order" 
		FROM experience 
		ORDER BY "order" DESC, start_date DESC
	`)
	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch experience"})
		return
	}
	defer rows.Close()

	var experience []Experience
	for rows.Next() {
		var exp Experience
		var technologiesJSON []byte
		err := rows.Scan(&exp.ID, &exp.Title, &exp.Company, &exp.StartDate, &exp.EndDate, &exp.Current, &exp.Description, &technologiesJSON, &exp.Order)
		if err != nil {
			span.RecordError(err)
			continue
		}

		// Parse technologies JSON array
		if err := json.Unmarshal(technologiesJSON, &exp.Technologies); err != nil {
			exp.Technologies = []string{}
		}

		experience = append(experience, exp)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan experience"})
		return
	}

	// Cache the result in Redis (5 minutes)
	if experienceData, err := json.Marshal(experience); err == nil {
		rdb.Set(ctx, "experience", experienceData, 5*time.Minute)
	}

	etag := generateETag(experience)
	c.Header("ETag", etag)
	c.JSON(http.StatusOK, experience)
}
