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

	// 🌐 Web framework and middleware
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	// 🔧 Environment and database
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"

	// 📊 Prometheus monitoring
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// 🗄️ Redis caching
	"github.com/redis/go-redis/v9"
	// 🔍 OpenTelemetry tracing - DISABLED
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// "go.opentelemetry.io/otel/sdk/resource"
	// sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// =============================================================================
// 📋 DATA STRUCTURES
// =============================================================================

// 🎯 Project represents a project
type Project struct {
	ID               int      `json:"id"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	ShortDescription string   `json:"short_description"`
	Type             string   `json:"type"`
	Modules          int      `json:"modules"`
	Icon             string   `json:"icon"`
	URL              string   `json:"url"`
	VideoURL         string   `json:"video_url,omitempty"`
	Technologies     []string `json:"technologies"`
	Active           bool     `json:"active"` // Controls visibility
}

// 📄 Content represents dynamic content from database
type Content struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// 👤 AboutData represents about page information
type AboutData struct {
	Description string `json:"description"`
	Highlights  []struct {
		Icon string `json:"icon"`
		Text string `json:"text"`
	} `json:"highlights"`
}

// 📞 ContactData represents contact information
type ContactData struct {
	Email        string `json:"email"`
	Location     string `json:"location"`
	LinkedIn     string `json:"linkedin"`
	GitHub       string `json:"github"`
	Availability string `json:"availability"`
}

// 🛠️ Skill represents a technical skill
type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Proficiency int    `json:"proficiency"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
}

// 💼 Experience represents work experience
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

// =============================================================================
// 🌍 GLOBAL VARIABLES
// =============================================================================

var (
	// 🗄️ Database and Redis connections
	db  *sql.DB
	rdb *redis.Client

	// 🚦 Rate limiting
	rateLimitMap   = make(map[string][]time.Time)
	rateLimitMutex sync.RWMutex

	// 📊 Prometheus metrics
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

	// 🔍 Validation patterns
	urlPattern   = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// =============================================================================
// 🚀 INITIALIZATION
// =============================================================================

func init() {
	// 📊 Register Prometheus metrics
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// =============================================================================
// 🗄️ CACHE HELPERS (Redis Operations)
// =============================================================================

// getFromCache retrieves data from Redis cache with type safety
func getFromCache[T any](ctx context.Context, key string) (T, bool) {
	var result T
	if rdb == nil {
		log.Printf("❌ DEBUG: Cache miss - Redis not available for key: %s", key)
		return result, false
	}

	log.Printf("🔍 DEBUG: Attempting to get from cache - key: %s", key)
	cached, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Printf("❌ DEBUG: Cache miss - Redis error for key %s: %v", key, err)
		return result, false
	}

	if err := json.Unmarshal([]byte(cached), &result); err != nil {
		log.Printf("❌ DEBUG: Cache miss - JSON unmarshal error for key %s: %v", key, err)
		return result, false
	}

	log.Printf("✅ DEBUG: Cache hit - Retrieved data for key: %s", key)
	return result, true
}

// setCache stores data in Redis cache with duration
func setCache(ctx context.Context, key string, data interface{}, duration time.Duration) {
	if rdb == nil {
		return
	}

	if dataBytes, err := json.Marshal(data); err == nil {
		rdb.Set(ctx, key, dataBytes, duration)
	}
}

// clearCache removes data from Redis cache
func clearCache(ctx context.Context, key string) {
	if rdb != nil {
		rdb.Del(ctx, key)
	}
}

// =============================================================================
// 🗃️ DATABASE HELPERS (with tracing)
// =============================================================================

// queryRowWithTracing executes a single row query with OpenTelemetry tracing - DISABLED
func queryRowWithTracing(ctx context.Context, tracerName, queryName, query string, dest ...interface{}) error {
	// tracer := otel.Tracer(tracerName)
	// _, span := tracer.Start(ctx, fmt.Sprintf("db.query.%s", queryName))
	// defer span.End()

	err := db.QueryRowContext(ctx, query, dest...).Scan(dest...)
	// if err != nil {
	// 	span.RecordError(err)
	// }
	return err
}

// execWithTracing executes a database statement with OpenTelemetry tracing - DISABLED
func execWithTracing(ctx context.Context, tracerName, queryName, query string, args ...interface{}) (sql.Result, error) {
	// tracer := otel.Tracer(tracerName)
	// _, span := tracer.Start(ctx, fmt.Sprintf("db.exec.%s", queryName))
	// defer span.End()

	result, err := db.ExecContext(ctx, query, args...)
	// if err != nil {
	// 	span.RecordError(err)
	// }
	return result, err
}

// =============================================================================
// 🌐 HTTP RESPONSE HELPERS
// =============================================================================

// respondWithETag sends response with ETag for caching
func respondWithETag(c *gin.Context, data interface{}, status int) {
	etag := generateETag(data)
	if c.GetHeader("If-None-Match") == etag {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("ETag", etag)
	c.JSON(status, data)
}

// respondWithError sends standardized error response
func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// respondWithSuccess sends standardized success response
func respondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// =============================================================================
// 🔧 UTILITY FUNCTIONS
// =============================================================================

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// generateETag creates MD5 hash for ETag
func generateETag(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	hash := md5.Sum(jsonData)
	return fmt.Sprintf("\"%x\"", hash)
}

// getContentValue retrieves content from database with fallback
func getContentValue(ctx context.Context, key, field, defaultValue string) string {
	var value string
	query := fmt.Sprintf("SELECT value->>'%s' FROM content WHERE key = $1", field)
	err := db.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}

// validateURL validates URL format
func validateURL(url string) bool {
	return urlPattern.MatchString(url)
}

// validateEmail validates email format
func validateEmail(email string) bool {
	return emailPattern.MatchString(email)
}

// sanitizeString removes potentially dangerous content
func sanitizeString(input string) string {
	dangerous := []string{"<script>", "</script>", "javascript:", "onload=", "onerror="}
	result := input
	for _, danger := range dangerous {
		result = strings.ReplaceAll(result, danger, "")
	}
	return result
}

// isValidIP validates IP address format
func isValidIP(ip string) bool {
	if ip == "localhost" || ip == "127.0.0.1" || ip == "::1" {
		return true
	}

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

// countActiveProjects counts active projects
func countActiveProjects(projects []Project) int {
	count := 0
	for _, p := range projects {
		if p.Active {
			count++
		}
	}
	return count
}

// countInactiveProjects counts inactive projects
func countInactiveProjects(projects []Project) int {
	count := 0
	for _, p := range projects {
		if !p.Active {
			count++
		}
	}
	return count
}

// =============================================================================
// 🔍 OPEN TELEMETRY SETUP
// =============================================================================

// initTracer sets up OpenTelemetry tracing - DISABLED
func initTracer() (interface{}, error) {
	// exporter, err := otlptracegrpc.New(context.Background(),
	// 	otlptracegrpc.WithEndpoint("localhost:4317"),
	// 	otlptracegrpc.WithInsecure(),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	// }

	// res, err := resource.New(context.Background(),
	// 	resource.WithAttributes(
	// 		semconv.ServiceName("bruno-api"),
	// 		semconv.ServiceVersion("1.0.0"),
	// 	),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create resource: %w", err)
	// }

	// tp := sdktrace.NewTracerProvider(
	// 	sdktrace.WithBatcher(exporter),
	// 	sdktrace.WithResource(res),
	// )

	// otel.SetTracerProvider(tp)
	// return tp, nil
	return nil, nil
}

// =============================================================================
// 🗄️ DATABASE & REDIS INITIALIZATION
// =============================================================================

// initDB initializes database connection
func initDB() {
	log.Printf("🔍 DEBUG: Initializing database connection...")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// 🔧 Fallback to individual environment variables
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "bruno_user")
		dbPassword := getEnv("DB_PASSWORD", "bruno_password")
		dbName := getEnv("DB_NAME", "bruno_site")

		log.Printf("🔍 DEBUG: Using individual DB config - Host: %s, Port: %s, User: %s, DB: %s",
			dbHost, dbPort, dbUser, dbName)

		dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
	} else {
		log.Printf("🔍 DEBUG: Using DATABASE_URL environment variable")
	}

	var err error
	log.Printf("🔍 DEBUG: Opening database connection...")
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Printf("🔍 DEBUG: Pinging database...")
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}

	log.Println("✅ Successfully connected to database")
}

// initRedis initializes Redis connection
func initRedis() {
	log.Printf("🔍 DEBUG: Initializing Redis connection...")

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		// 🔧 Fallback to individual environment variables
		redisHost := getEnv("REDIS_HOST", "localhost")
		redisPort := getEnv("REDIS_PORT", "6379")
		redisPassword := getEnv("REDIS_PASSWORD", "")

		log.Printf("🔍 DEBUG: Using individual Redis config - Host: %s, Port: %s, Password: %s",
			redisHost, redisPort, func() string {
				if redisPassword != "" {
					return "***"
				} else {
					return "none"
				}
			}())

		if redisPassword != "" {
			redisURL = fmt.Sprintf("redis://:%s@%s:%s", redisPassword, redisHost, redisPort)
		} else {
			redisURL = fmt.Sprintf("redis://%s:%s", redisHost, redisPort)
		}
	} else {
		log.Printf("🔍 DEBUG: Using REDIS_URL environment variable")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("⚠️ Failed to parse Redis URL: %v, continuing without Redis", err)
		rdb = nil
		return
	}

	rdb = redis.NewClient(opt)

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️ Failed to connect to Redis: %v, continuing without Redis", err)
		rdb = nil
		return
	}

	log.Println("✅ Successfully connected to Redis")
}

// =============================================================================
// 🚦 MIDDLEWARE
// =============================================================================

// rateLimitMiddleware implements rate limiting (100 requests per minute)
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		rateLimitMutex.Lock()
		now := time.Now()

		// 🧹 Clean old requests (older than 1 minute)
		var validRequests []time.Time
		for _, reqTime := range rateLimitMap[clientIP] {
			if now.Sub(reqTime) < time.Minute {
				validRequests = append(validRequests, reqTime)
			}
		}
		rateLimitMap[clientIP] = validRequests

		// 🚫 Check rate limit (100 requests per minute)
		if len(rateLimitMap[clientIP]) >= 100 {
			rateLimitMutex.Unlock()
			respondWithError(c, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		// ➕ Add current request
		rateLimitMap[clientIP] = append(rateLimitMap[clientIP], now)
		rateLimitMutex.Unlock()

		c.Next()
	}
}

// prometheusMiddleware records metrics for all requests
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

// =============================================================================
// 🎯 API HANDLERS
// =============================================================================

// getProjects returns active projects
func getProjects(c *gin.Context) {
	// tracer := otel.Tracer("bruno-api")
	// ctx, span := tracer.Start(c.Request.Context(), "getProjects")
	// defer span.End()
	ctx := c.Request.Context()

	log.Printf("🔍 DEBUG: getProjects called - User-Agent: %s, Remote-Addr: %s", c.GetHeader("User-Agent"), c.ClientIP())

	// 🗄️ Try cache first
	if cached, found := getFromCache[[]Project](ctx, "projects"); found {
		log.Printf("✅ DEBUG: Cache hit for projects, returning %d projects", len(cached))
		// span.AddEvent("cache_hit")
		respondWithETag(c, cached, http.StatusOK)
		return
	}

	log.Printf("❌ DEBUG: Cache miss for projects, querying database")

	// span.AddEvent("cache_miss")

	// 🗃️ Query database
	query := `
		SELECT id, title, description, description as short_description, type, modules, github_url, live_url, video_url, technologies, active
		FROM projects 
		WHERE active = true
		ORDER BY "order" ASC, id ASC
	`
	log.Printf("🔍 Executing query: %s", query)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("❌ DEBUG: Database query failed: %v", err)
		// span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}
	defer rows.Close()

	log.Printf("✅ DEBUG: Database query successful, processing rows")

	var projects []Project
	rowCount := 0
	for rows.Next() {
		var p Project
		var githubURL, liveURL sql.NullString
		var videoURL sql.NullString
		var technologies pq.StringArray
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ShortDescription, &p.Type, &p.Modules, &githubURL, &liveURL, &videoURL, &technologies, &p.Active); err != nil {
			continue
		}
		// 🔗 Use live_url if available, otherwise use github_url
		if liveURL.Valid {
			p.URL = liveURL.String
		} else if githubURL.Valid {
			p.URL = githubURL.String
		}
		if videoURL.Valid {
			p.VideoURL = videoURL.String
		}
		// Set technologies
		p.Technologies = []string(technologies)
		projects = append(projects, p)
		rowCount++
	}

	log.Printf("✅ DEBUG: Processed %d projects from database", rowCount)

	// 🗄️ Cache the result
	log.Printf("💾 DEBUG: Caching %d projects for 5 minutes", len(projects))
	setCache(ctx, "projects", projects, 5*time.Minute)

	// span.SetAttributes(semconv.DBStatement("SELECT projects"))
	log.Printf("✅ DEBUG: Returning %d projects to client", len(projects))
	respondWithETag(c, projects, http.StatusOK)
}

// getAllProjects returns all projects (including inactive) for admin
func getAllProjects(c *gin.Context) {
	// tracer := otel.Tracer("bruno-api")
	// ctx, span := tracer.Start(c.Request.Context(), "getAllProjects")
	// defer span.End()
	ctx := c.Request.Context()

	query := `
		SELECT id, title, description, description as short_description, type, modules, github_url, live_url, video_url, technologies, active
		FROM projects 
		ORDER BY "order" ASC, id ASC
	`
	log.Printf("🔍 getAllProjects: Executing query: %s", query)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("❌ getAllProjects: Database query failed: %v", err)
		// span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var githubURL, liveURL sql.NullString
		var videoURL sql.NullString
		var technologies pq.StringArray
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ShortDescription, &p.Type, &p.Modules, &githubURL, &liveURL, &videoURL, &technologies, &p.Active); err != nil {
			continue
		}
		if liveURL.Valid {
			p.URL = liveURL.String
		} else if githubURL.Valid {
			p.URL = githubURL.String
		}
		if videoURL.Valid {
			p.VideoURL = videoURL.String
		}
		// Set technologies
		p.Technologies = []string(technologies)
		projects = append(projects, p)
	}

	log.Printf("🔍 getAllProjects: Found %d projects", len(projects))

	// span.SetAttributes(semconv.DBStatement("SELECT all projects"))
	respondWithSuccess(c, gin.H{
		"projects": projects,
		"total":    len(projects),
		"active":   countActiveProjects(projects),
		"inactive": countInactiveProjects(projects),
	})
}

// activateProject activates a project
func activateProject(c *gin.Context) {
	// tracer := otel.Tracer("bruno-api")
	// ctx, span := tracer.Start(c.Request.Context(), "activateProject")
	// defer span.End()
	ctx := c.Request.Context()

	projectID := c.Param("id")
	if projectID == "" {
		respondWithError(c, http.StatusBadRequest, "Project ID is required")
		return
	}

	result, err := execWithTracing(ctx, "bruno-api", "activate_project", "UPDATE projects SET active = true WHERE id = $1", projectID)
	if err != nil {
		// span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to activate project")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(c, http.StatusNotFound, "Project not found")
		return
	}

	clearCache(ctx, "projects")
	// span.AddEvent("project_activated")
	respondWithSuccess(c, gin.H{
		"message":    "Project activated successfully",
		"project_id": projectID,
	})
}

// deactivateProject deactivates a project
func deactivateProject(c *gin.Context) {
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "deactivateProject")
	defer span.End()

	projectID := c.Param("id")
	if projectID == "" {
		respondWithError(c, http.StatusBadRequest, "Project ID is required")
		return
	}

	result, err := execWithTracing(ctx, "portfolio-api", "deactivate_project", "UPDATE projects SET active = false WHERE id = $1", projectID)
	if err != nil {
		span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to deactivate project")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(c, http.StatusNotFound, "Project not found")
		return
	}

	clearCache(ctx, "projects")
	span.AddEvent("project_deactivated")
	respondWithSuccess(c, gin.H{
		"message":    "Project deactivated successfully",
		"project_id": projectID,
	})
}

// getProjectStats returns project statistics
func getProjectStats(c *gin.Context) {
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "getProjectStats")
	defer span.End()

	var total, active int
	err := queryRowWithTracing(ctx, "portfolio-api", "project_stats_total", "SELECT COUNT(*) FROM projects", &total)
	if err != nil {
		span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to get project stats")
		return
	}

	err = queryRowWithTracing(ctx, "portfolio-api", "project_stats_active", "SELECT COUNT(*) FROM projects WHERE active = true", &active)
	if err != nil {
		span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to get active project count")
		return
	}

	inactive := total - active
	respondWithSuccess(c, gin.H{
		"total":             total,
		"active":            active,
		"inactive":          inactive,
		"active_percentage": float64(active) / float64(total) * 100,
	})
}

// getAbout returns about page data
func getAbout(c *gin.Context) {
	ctx := context.Background()

	// 🗄️ Try cache first
	if cached, found := getFromCache[AboutData](ctx, "about"); found {
		respondWithSuccess(c, cached)
		return
	}

	var description string
	err := db.QueryRowContext(ctx, "SELECT value->>'description' FROM content WHERE key = 'about'").Scan(&description)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch about data")
		return
	}

	aboutData := AboutData{
		Description: description,
		Highlights: []struct {
			Icon string `json:"icon"`
			Text string `json:"text"`
		}{}, // TODO: Implement proper JSON parsing for highlights
	}

	setCache(ctx, "about", aboutData, 10*time.Minute)
	respondWithSuccess(c, aboutData)
}

// getContact returns contact information
func getContact(c *gin.Context) {
	ctx := context.Background()

	// 🗄️ Try cache first
	if cached, found := getFromCache[ContactData](ctx, "contact"); found {
		respondWithSuccess(c, cached)
		return
	}

	contactData := ContactData{
		Email:        getContentValue(ctx, "contact", "email", "bruno@lucena.cloud"),
		Location:     getContentValue(ctx, "contact", "location", "Brazil"),
		LinkedIn:     getContentValue(ctx, "contact", "linkedin", "https://www.linkedin.com/in/bvlucena"),
		GitHub:       getContentValue(ctx, "contact", "github", "https://github.com/brunovlucena"),
		Availability: getContentValue(ctx, "contact", "availability", "Open to new opportunities in SRE, DevSecOps, and AI Engineering roles."),
	}

	setCache(ctx, "contact", contactData, 10*time.Minute)
	respondWithSuccess(c, contactData)
}

// getSkills returns technical skills
func getSkills(c *gin.Context) {
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "getSkills")
	defer span.End()

	// 🗄️ Try cache first
	if cached, found := getFromCache[[]Skill](ctx, "skills"); found {
		span.AddEvent("cache_hit")
		respondWithETag(c, cached, http.StatusOK)
		return
	}

	rows, err := db.QueryContext(ctx, `
		SELECT id, name, category, proficiency, icon, "order" 
		FROM skills 
		ORDER BY "order", name
	`)
	if err != nil {
		span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch skills")
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
		respondWithError(c, http.StatusInternalServerError, "Failed to scan skills")
		return
	}

	setCache(ctx, "skills", skills, 5*time.Minute)
	respondWithETag(c, skills, http.StatusOK)
}

// getExperience returns work experience
func getExperience(c *gin.Context) {
	tracer := otel.Tracer("portfolio-api")
	ctx, span := tracer.Start(c.Request.Context(), "getExperience")
	defer span.End()

	// 🗄️ Try cache first
	if cached, found := getFromCache[[]Experience](ctx, "experience"); found {
		span.AddEvent("cache_hit")
		respondWithETag(c, cached, http.StatusOK)
		return
	}

	rows, err := db.QueryContext(ctx, `
		SELECT id, title, company, start_date, end_date, current, description, technologies, "order" 
		FROM experience 
		ORDER BY "order" DESC, start_date DESC
	`)
	if err != nil {
		span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch experience")
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

		if err := json.Unmarshal(technologiesJSON, &exp.Technologies); err != nil {
			exp.Technologies = []string{}
		}

		experience = append(experience, exp)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		respondWithError(c, http.StatusInternalServerError, "Failed to scan experience")
		return
	}

	setCache(ctx, "experience", experience, 5*time.Minute)
	respondWithETag(c, experience, http.StatusOK)
}

// trackVisit tracks visitor analytics
func trackVisit(c *gin.Context) {
	var visit struct {
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		Referrer  string `json:"referrer"`
	}

	if err := c.ShouldBindJSON(&visit); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if visit.IP == "" {
		visit.IP = c.ClientIP()
	}

	if !isValidIP(visit.IP) {
		respondWithError(c, http.StatusBadRequest, "Invalid IP address format")
		return
	}

	visit.UserAgent = sanitizeString(visit.UserAgent)
	visit.Referrer = sanitizeString(visit.Referrer)

	if visit.Referrer != "" && !validateURL(visit.Referrer) {
		respondWithError(c, http.StatusBadRequest, "Invalid referrer URL")
		return
	}

	_, err := db.ExecContext(c.Request.Context(), `
		INSERT INTO visitors (ip, user_agent, first_visit, last_visit, visit_count)
		VALUES ($1, $2, NOW(), NOW(), 1)
		ON CONFLICT (ip) DO UPDATE SET
			last_visit = NOW(),
			visit_count = visitors.visit_count + 1
	`, visit.IP, visit.UserAgent)

	if err != nil {
		log.Printf("❌ Failed to store analytics: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to track visit")
		return
	}

	respondWithSuccess(c, gin.H{"status": "success"})
}

// =============================================================================
// 🚀 MAIN APPLICATION
// =============================================================================

func main() {
	// 🔧 Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ No .env file found, using environment variables")
	}

	// 🔍 Initialize OpenTelemetry tracing - DISABLED
	// tp, err := initTracer()
	// if err != nil {
	// 	log.Printf("⚠️ Failed to initialize tracer: %v", err)
	// } else {
	// 	defer func() {
	// 		if err := tp.Shutdown(context.Background()); err != nil {
	// 			log.Printf("❌ Error shutting down tracer provider: %v", err)
	// 		}
	// 	}()
	// }

	// 🗄️ Initialize database and Redis
	initDB()
	defer db.Close()

	initRedis()
	if rdb != nil {
		defer rdb.Close()
	}

	// 🌐 Set up Gin router
	// Enable debug mode for verbose logging
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 🔧 Add middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(prometheusMiddleware())
	// r.Use(otelgin.Middleware("portfolio-api")) // DISABLED

	// 🌍 Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080", "*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "User-Agent", "Referer"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))

	// 🔍 Add debug middleware for all requests
	r.Use(func(c *gin.Context) {
		log.Printf("🔍 DEBUG: Incoming request - Method: %s, Path: %s, IP: %s, User-Agent: %s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP(), c.GetHeader("User-Agent"))
		c.Next()
		log.Printf("✅ DEBUG: Request completed - Status: %d, Size: %d",
			c.Writer.Status(), c.Writer.Size())
	})

	// 🏥 Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		log.Printf("🏥 DEBUG: Health check requested from %s", c.ClientIP())
		respondWithSuccess(c, gin.H{"status": "ok"})
	})

	// 📊 Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 🎯 API routes with rate limiting
	v1 := r.Group("/v1")
	{
		v1.GET("/projects", rateLimitMiddleware(), getProjects)
		v1.GET("/about", rateLimitMiddleware(), getAbout)
		v1.GET("/contact", rateLimitMiddleware(), getContact)
		v1.GET("/content/skills", rateLimitMiddleware(), getSkills)
		v1.GET("/content/experience", rateLimitMiddleware(), getExperience)
		v1.POST("/analytics/visit", rateLimitMiddleware(), trackVisit)
	}

	// 👑 Admin routes for project management
	admin := r.Group("/admin")
	{
		admin.GET("/projects", rateLimitMiddleware(), getAllProjects)
		admin.PUT("/projects/:id/activate", rateLimitMiddleware(), activateProject)
		admin.PUT("/projects/:id/deactivate", rateLimitMiddleware(), deactivateProject)
		admin.GET("/projects/stats", rateLimitMiddleware(), getProjectStats)
	}

	// 🚀 Start server
	port := getEnv("PORT", "8080")
	log.Printf("🚀 Starting Bruno API server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
