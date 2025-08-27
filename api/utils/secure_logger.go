package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// 🔒 SecureLogger provides secure logging with sensitive data sanitization
type SecureLogger struct {
	level string
}

// LogLevel represents the logging level
type LogLevel string

const (
	DEBUG   LogLevel = "DEBUG"
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
)

// NewSecureLogger creates a new secure logger instance
func NewSecureLogger(level string) *SecureLogger {
	return &SecureLogger{
		level: level,
	}
}

// sanitizeString removes or masks sensitive information from strings
func (sl *SecureLogger) sanitizeString(input string) string {
	if input == "" {
		return input
	}

	// 🔐 Sanitize database connection strings
	dbPattern := regexp.MustCompile(`postgresql://[^:]+:[^@]+@[^/]+/[^?]*`)
	input = dbPattern.ReplaceAllString(input, "postgresql://***:***@***/***")

	// 🔐 Sanitize Redis URLs
	redisPattern := regexp.MustCompile(`redis://[^:]+:[^@]+@[^/]+`)
	input = redisPattern.ReplaceAllString(input, "redis://***:***@***")

	// 🔐 Sanitize HTTP URLs with credentials
	httpPattern := regexp.MustCompile(`https?://[^:]+:[^@]+@[^/]+`)
	input = httpPattern.ReplaceAllString(input, "***://***:***@***")

	// 🔐 Sanitize passwords in environment variables
	envPattern := regexp.MustCompile(`(PASSWORD|SECRET|KEY|TOKEN)=[^,\s]+`)
	input = envPattern.ReplaceAllString(input, "$1=***")

	// 🔐 Sanitize API keys
	apiKeyPattern := regexp.MustCompile(`[a-zA-Z0-9]{32,}`)
	input = apiKeyPattern.ReplaceAllString(input, "***")

	// 🔐 Sanitize IP addresses (keep localhost)
	ipPattern := regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	if !strings.Contains(input, "localhost") && !strings.Contains(input, "127.0.0.1") {
		input = ipPattern.ReplaceAllString(input, "***.***.***.***")
	}

	return input
}

// sanitizeError sanitizes error messages to remove sensitive information
func (sl *SecureLogger) sanitizeError(err error) string {
	if err == nil {
		return ""
	}

	errorMsg := err.Error()

	// 🔐 Remove connection strings from error messages
	errorMsg = sl.sanitizeString(errorMsg)

	// 🔐 Remove stack traces that might contain sensitive data
	if strings.Contains(errorMsg, "goroutine") {
		lines := strings.Split(errorMsg, "\n")
		var sanitizedLines []string
		for _, line := range lines {
			if !strings.Contains(line, "goroutine") && !strings.Contains(line, "runtime.") {
				sanitizedLines = append(sanitizedLines, line)
			}
		}
		errorMsg = strings.Join(sanitizedLines, "\n")
	}

	return errorMsg
}

// formatMessage creates a structured log message
func (sl *SecureLogger) formatMessage(level LogLevel, message string, fields map[string]interface{}) string {
	timestamp := time.Now().Format("2006-01-02T15:04:05.000Z07:00")

	// 🔐 Sanitize the main message
	sanitizedMessage := sl.sanitizeString(message)

	// 🔐 Sanitize field values
	sanitizedFields := make(map[string]interface{})
	for key, value := range fields {
		if str, ok := value.(string); ok {
			sanitizedFields[key] = sl.sanitizeString(str)
		} else {
			sanitizedFields[key] = value
		}
	}

	// Create structured log entry
	logEntry := fmt.Sprintf("[%s] %s: %s", timestamp, level, sanitizedMessage)

	// Add fields if present
	if len(sanitizedFields) > 0 {
		fieldStr := ""
		for key, value := range sanitizedFields {
			if fieldStr != "" {
				fieldStr += ", "
			}
			fieldStr += fmt.Sprintf("%s=%v", key, value)
		}
		logEntry += fmt.Sprintf(" | %s", fieldStr)
	}

	return logEntry
}

// Debug logs debug level messages
func (sl *SecureLogger) Debug(message string, fields ...map[string]interface{}) {
	if sl.level == "DEBUG" {
		var fieldMap map[string]interface{}
		if len(fields) > 0 {
			fieldMap = fields[0]
		}
		log.Println(sl.formatMessage(DEBUG, message, fieldMap))
	}
}

// Info logs info level messages
func (sl *SecureLogger) Info(message string, fields ...map[string]interface{}) {
	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}
	log.Println(sl.formatMessage(INFO, message, fieldMap))
}

// Warning logs warning level messages
func (sl *SecureLogger) Warning(message string, fields ...map[string]interface{}) {
	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}
	log.Println(sl.formatMessage(WARNING, message, fieldMap))
}

// Error logs error level messages with sanitized error details
func (sl *SecureLogger) Error(message string, err error, fields ...map[string]interface{}) {
	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}

	// 🔐 Sanitize error details
	if err != nil {
		sanitizedErr := sl.sanitizeError(err)
		if fieldMap == nil {
			fieldMap = make(map[string]interface{})
		}
		fieldMap["error"] = sanitizedErr
	}

	log.Println(sl.formatMessage(ERROR, message, fieldMap))
}

// 🔐 Secure connection logging utilities

// LogDatabaseConnection logs database connection attempts without sensitive data
func (sl *SecureLogger) LogDatabaseConnection(attempt, maxAttempts int, host, port, dbname string, err error) {
	fields := map[string]interface{}{
		"attempt":      attempt,
		"max_attempts": maxAttempts,
		"host":         host,
		"port":         port,
		"database":     dbname,
	}

	if err != nil {
		sl.Error("Database connection attempt failed", err, fields)
	} else {
		sl.Info("Database connection successful", fields)
	}
}

// LogRedisConnection logs Redis connection attempts without sensitive data
func (sl *SecureLogger) LogRedisConnection(attempt, maxAttempts int, host, port string, err error) {
	fields := map[string]interface{}{
		"attempt":      attempt,
		"max_attempts": maxAttempts,
		"host":         host,
		"port":         port,
	}

	if err != nil {
		sl.Error("Redis connection attempt failed", err, fields)
	} else {
		sl.Info("Redis connection successful", fields)
	}
}

// LogAPICall logs API calls without sensitive request/response data
func (sl *SecureLogger) LogAPICall(method, endpoint, status string, duration time.Duration, requestID string) {
	fields := map[string]interface{}{
		"method":     method,
		"endpoint":   endpoint,
		"status":     status,
		"duration":   duration.String(),
		"request_id": requestID,
	}

	sl.Info("API call completed", fields)
}

// LogServiceHealth logs service health checks without sensitive details
func (sl *SecureLogger) LogServiceHealth(serviceName, status string, duration time.Duration) {
	fields := map[string]interface{}{
		"service":  serviceName,
		"status":   status,
		"duration": duration.String(),
	}

	sl.Info("Service health check", fields)
}

// 🔐 Global secure logger instance
var SecureLog = NewSecureLogger("INFO")

// 🔐 Convenience functions for backward compatibility
func LogInfo(message string, fields ...map[string]interface{}) {
	SecureLog.Info(message, fields...)
}

func LogError(message string, err error, fields ...map[string]interface{}) {
	SecureLog.Error(message, err, fields...)
}

func LogWarning(message string, fields ...map[string]interface{}) {
	SecureLog.Warning(message, fields...)
}

func LogDebug(message string, fields ...map[string]interface{}) {
	SecureLog.Debug(message, fields...)
}
