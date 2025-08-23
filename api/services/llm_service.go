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

// OllamaRequest represents request format for Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents response format from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// NewLLMService creates a new LLM service
func NewLLMService(db *sql.DB) *LLMService {
	service := &LLMService{
		ollamaURL:      getEnv("OLLAMA_URL", "http://192.168.0.3:11434"),
		model:          getEnv("GEMMA_MODEL", "gemma2:2b"),
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
		// Continue with basic context
		context = llm.getFallbackContext()
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
		Model:  llm.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := llm.httpClient.Post(
		fmt.Sprintf("%s/api/generate", llm.ollamaURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

// getFallbackContext provides basic context when database query fails
func (llm *LLMService) getFallbackContext() string {
	return `You are Bruno's AI assistant. Answer questions about Bruno Lucena based on this information:

ABOUT BRUNO:
Senior Cloud Native Infrastructure Engineer with extensive experience in designing, implementing, and maintaining scalable, resilient cloud-native infrastructure. Passionate about automation, observability, and modern DevOps practices.

KEY SKILLS:
- Cloud Platforms: AWS, GCP, Azure
- Container Orchestration: Kubernetes, Docker
- Infrastructure as Code: Terraform, Pulumi
- Programming: Go, Python, TypeScript
- Observability: Prometheus, Grafana, Loki
- DevOps: CI/CD, GitOps, SRE practices

CURRENT ROLE:
SRE/DevOps at Notifi (2023-present) - Architecting cloud-native infrastructure, implementing observability solutions, and building serverless applications.

CONTACT:
- Email: bruno@lucena.cloud
- Location: Brazil
- LinkedIn: https://www.linkedin.com/in/bvlucena
- GitHub: https://github.com/brunovlucena

INSTRUCTIONS:
- Answer as Bruno's AI assistant in first person about Bruno
- Be conversational, helpful, and professional
- Use the provided data to give accurate answers
- Keep responses concise but informative

USER QUESTION: `
}

// Health check for LLM service
func (llm *LLMService) HealthCheck() error {
	switch strings.ToLower(llm.provider) {
	case "ollama":
		return llm.checkOllamaHealth()
	case "lmstudio":
		return llm.checkLMStudioHealth()
	default:
		return fmt.Errorf("unsupported provider: %s", llm.provider)
	}
}

func (llm *LLMService) checkOllamaHealth() error {
	resp, err := llm.httpClient.Get(fmt.Sprintf("%s/api/tags", llm.ollamaURL))
	if err != nil {
		return fmt.Errorf("Ollama health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (llm *LLMService) checkLMStudioHealth() error {
	resp, err := llm.httpClient.Get(fmt.Sprintf("%s/v1/models", llm.lmstudioURL))
	if err != nil {
		return fmt.Errorf("LMStudio health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LMStudio health check failed with status: %d", resp.StatusCode)
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
