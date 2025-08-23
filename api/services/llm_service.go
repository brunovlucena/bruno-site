package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

	log.Printf("🤖 LLM Service initialized - Model: %s", service.model)
	return service
}

// ProcessChat handles a chat request and returns an AI response
func (llm *LLMService) ProcessChat(request ChatRequest) (*ChatResponse, error) {
	log.Printf("💬 Processing chat request: %s", request.Message)

	// Build context from PostgreSQL data
	context, err := llm.contextBuilder.BuildContext(request.Message)
	if err != nil {
		log.Printf("⚠️ Error building context: %v", err)
		return nil, fmt.Errorf("failed to build context: %v", err)
	}

	// Generate response using Ollama
	response, err := llm.callOllama(context)

	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %v", err)
	}

	// Create response
	chatResponse := &ChatResponse{
		Response:  response,
		Model:     llm.model,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Sources:   []string{"PostgreSQL Database"}, // Could be enhanced to show specific tables used
	}

	log.Printf("✅ Chat response generated successfully")
	return chatResponse, nil
}

// callOllama sends request to Ollama API
func (llm *LLMService) callOllama(prompt string) (string, error) {
	log.Printf("🦙 Calling Ollama API at %s", llm.ollamaURL)

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
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := llm.httpClient.Post(
		fmt.Sprintf("%s/api/chat", llm.ollamaURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	response := strings.TrimSpace(ollamaResp.Message.Content)

	// Clean up any double spaces and trim
	response = strings.TrimSpace(response)

	return response, nil
}

// HealthCheck checks if Ollama service is available
func (llm *LLMService) HealthCheck() error {
	resp, err := llm.httpClient.Get(fmt.Sprintf("%s/api/tags", llm.ollamaURL))
	if err != nil {
		return fmt.Errorf("ollama health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama health check failed with status: %d", resp.StatusCode)
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
