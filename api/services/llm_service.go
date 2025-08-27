package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"bruno-api/utils"
)

// LLMService handles communication with Ollama
type LLMService struct {
	ollamaURL      string
	model          string
	contextBuilder *ContextBuilder
	httpClient     *http.Client
}

// ChatRequest represents an incoming chat request
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
	Context string `json:"context,omitempty"`
}

// ChatResponse represents the response from the chatbot
type ChatResponse struct {
	Response  string   `json:"response"`
	Sources   []string `json:"sources,omitempty"`
	Model     string   `json:"model"`
	Timestamp string   `json:"timestamp"`
}

// OllamaRequest represents request format for Ollama Chat API
type OllamaRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaResponse represents response format from Ollama Chat API
type OllamaResponse struct {
	Message OllamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewLLMService creates a new LLM service
func NewLLMService(db *sql.DB) *LLMService {
	service := &LLMService{
		ollamaURL:      getEnv("OLLAMA_URL", "http://192.168.0.3:11434"),
		model:          getEnv("GEMMA_MODEL", "gemma3n:e4b"),
		contextBuilder: NewContextBuilder(db),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	utils.SecureLog.Info("LLM Service initialized", map[string]interface{}{
		"ollama_url": service.ollamaURL,
		"model":      service.model,
		"timeout":    service.httpClient.Timeout.String(),
	})

	// Test connection on startup
	go service.testConnectionOnStartup()

	return service
}

// testConnectionOnStartup tests the Ollama connection in background
func (llm *LLMService) testConnectionOnStartup() {
	utils.SecureLog.Info("Testing Ollama connection on startup")

	// Wait a bit for the service to fully start
	time.Sleep(2 * time.Second)

	if err := llm.HealthCheck(); err != nil {
		utils.SecureLog.Error("Ollama connection test failed", err, map[string]interface{}{
			"ollama_url": llm.ollamaURL,
			"model":      llm.model,
		})
		utils.SecureLog.Info("Troubleshooting tips: Check if Ollama is running, verify network connectivity, check model availability, verify firewall settings")
	} else {
		utils.SecureLog.Info("Ollama connection test successful")
	}
}

// ProcessChat handles a chat request and returns an AI response
func (llm *LLMService) ProcessChat(request ChatRequest) (*ChatResponse, error) {
	startTime := time.Now()
	requestID := fmt.Sprintf("chat_%d", startTime.UnixNano())

	utils.SecureLog.Info("Starting chat processing", map[string]interface{}{
		"request_id":     requestID,
		"message_length": len(request.Message),
		"model":          llm.model,
	})

	// Build context from PostgreSQL data
	utils.SecureLog.Info("Building context from database", map[string]interface{}{
		"request_id": requestID,
	})
	context, err := llm.contextBuilder.BuildContext(request.Message)
	if err != nil {
		utils.SecureLog.Error("Context building failed", err, map[string]interface{}{
			"request_id":   requestID,
			"db_connected": llm.contextBuilder.db != nil,
		})
		return nil, fmt.Errorf("failed to build context: %v", err)
	}
	utils.SecureLog.Info("Context built successfully", map[string]interface{}{
		"request_id":     requestID,
		"context_length": len(context),
	})

	// Generate response using Ollama
	utils.SecureLog.Info("Calling Ollama API", map[string]interface{}{
		"request_id": requestID,
	})
	if err := llm.HealthCheck(); err != nil {
		utils.SecureLog.Warning("Ollama health check failed before API call", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
		})
	}

	response, err := llm.callOllama(context, requestID)

	if err != nil {
		utils.SecureLog.Error("Ollama API call failed", err, map[string]interface{}{
			"request_id": requestID,
		})
		return nil, fmt.Errorf("LLM request failed: %v", err)
	}

	// Create response
	chatResponse := &ChatResponse{
		Response:  response,
		Model:     llm.model,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Sources:   []string{"PostgreSQL Database"},
	}

	duration := time.Since(startTime)
	utils.SecureLog.Info("Chat processing completed", map[string]interface{}{
		"request_id":      requestID,
		"duration":        duration.String(),
		"response_length": len(response),
		"model":           llm.model,
	})

	return chatResponse, nil
}

// callOllama sends request to Ollama API with enhanced logging
func (llm *LLMService) callOllama(prompt string, requestID string) (string, error) {
	utils.SecureLog.Info("Preparing Ollama request", map[string]interface{}{
		"request_id":    requestID,
		"model":         llm.model,
		"prompt_length": len(prompt),
		"timeout":       llm.httpClient.Timeout.String(),
	})

	requestBody := OllamaRequest{
		Model: llm.model,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "You are a fact-based assistant. NEVER use greetings, introductions, or pleasantries. Answer questions immediately with facts only. Maximum 2 sentences. Start directly with the answer.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		utils.SecureLog.Error("Failed to marshal request", err, map[string]interface{}{
			"request_id": requestID,
		})
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}
	utils.SecureLog.Info("Request payload prepared", map[string]interface{}{
		"request_id":   requestID,
		"payload_size": len(jsonData),
	})

	// Log request details (without sensitive data)
	utils.SecureLog.Info("Sending HTTP POST request", map[string]interface{}{
		"request_id": requestID,
		"timeout":    llm.httpClient.Timeout.String(),
	})

	startTime := time.Now()
	resp, err := llm.httpClient.Post(
		fmt.Sprintf("%s/api/chat", llm.ollamaURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	requestDuration := time.Since(startTime)

	if err != nil {
		utils.SecureLog.Error("HTTP request failed", err, map[string]interface{}{
			"request_id": requestID,
			"duration":   requestDuration.String(),
		})
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	utils.SecureLog.Info("Received response", map[string]interface{}{
		"request_id":  requestID,
		"duration":    requestDuration.String(),
		"status_code": resp.StatusCode,
		"status":      resp.Status,
	})

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		utils.SecureLog.Error("Ollama API error response", fmt.Errorf("status %d: %s", resp.StatusCode, string(body)), map[string]interface{}{
			"request_id":  requestID,
			"status_code": resp.StatusCode,
			"model":       llm.model,
		})
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SecureLog.Error("Failed to read response body", err, map[string]interface{}{
			"request_id": requestID,
		})
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	utils.SecureLog.Info("Response body read", map[string]interface{}{
		"request_id": requestID,
		"body_size":  len(body),
	})

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		utils.SecureLog.Error("Failed to unmarshal response", err, map[string]interface{}{
			"request_id": requestID,
		})
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	response := strings.TrimSpace(ollamaResp.Message.Content)
	response = strings.TrimSpace(response)

	utils.SecureLog.Info("Ollama response processed successfully", map[string]interface{}{
		"request_id":      requestID,
		"response_length": len(response),
		"model":           llm.model,
		"total_time":      requestDuration.String(),
	})

	return response, nil
}

// HealthCheck checks if Ollama service is available with enhanced logging
func (llm *LLMService) HealthCheck() error {
	utils.SecureLog.Info("Starting Ollama health check", map[string]interface{}{
		"timeout": llm.httpClient.Timeout.String(),
	})

	startTime := time.Now()
	resp, err := llm.httpClient.Get(fmt.Sprintf("%s/api/tags", llm.ollamaURL))
	duration := time.Since(startTime)

	if err != nil {
		utils.SecureLog.Error("Health check failed", err, map[string]interface{}{
			"duration": duration.String(),
		})
		return fmt.Errorf("ollama health check failed: %v", err)
	}
	defer resp.Body.Close()

	utils.SecureLog.Info("Health check response received", map[string]interface{}{
		"duration":    duration.String(),
		"status_code": resp.StatusCode,
		"status":      resp.Status,
	})

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		utils.SecureLog.Error("Health check failed with status", fmt.Errorf("status %d: %s", resp.StatusCode, string(body)), map[string]interface{}{
			"status_code": resp.StatusCode,
		})
		return fmt.Errorf("ollama health check failed with status: %d", resp.StatusCode)
	}

	// Try to parse the response to get model information
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SecureLog.Warning("Health check succeeded but failed to read response", map[string]interface{}{
			"error": err.Error(),
		})
		utils.SecureLog.Info("Ollama is responding (status 200)")
		return nil
	}

	// Parse models list
	var modelsResponse struct {
		Models []struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &modelsResponse); err != nil {
		utils.SecureLog.Warning("Health check succeeded but failed to parse models", map[string]interface{}{
			"error": err.Error(),
		})
		utils.SecureLog.Info("Ollama is responding (status 200)")
		return nil
	}

	utils.SecureLog.Info("Ollama health check successful", map[string]interface{}{
		"available_models": len(modelsResponse.Models),
	})

	// Check if our model is available
	modelFound := false
	for _, model := range modelsResponse.Models {
		if model.Name == llm.model {
			modelFound = true
			utils.SecureLog.Info("Required model found", map[string]interface{}{
				"model":      model.Name,
				"size_bytes": model.Size,
			})
			break
		}
	}

	if !modelFound {
		utils.SecureLog.Warning("Required model not found in available models", map[string]interface{}{
			"required_model":   llm.model,
			"available_models": len(modelsResponse.Models),
		})
	}

	return nil
}

// Helper function to get environment variables
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
