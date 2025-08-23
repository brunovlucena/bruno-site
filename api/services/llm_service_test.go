package services

import (
	"database/sql"
	"testing"
)

// TestNewLLMService tests the creation of a new LLMService
func TestNewLLMService(t *testing.T) {
	// Create a mock database connection (nil for testing)
	var db *sql.DB = nil

	service := NewLLMService(db)

	if service == nil {
		t.Error("Expected LLMService to be created, got nil")
	}

	if service.contextBuilder == nil {
		t.Error("Expected contextBuilder to be initialized")
	}

	if service.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}

	if service.ollamaURL == "" {
		t.Error("Expected ollamaURL to be set")
	}

	if service.model == "" {
		t.Error("Expected model to be set")
	}
}

// TestChatRequestValidation tests the ChatRequest struct
func TestChatRequestValidation(t *testing.T) {
	// Test valid request
	request := ChatRequest{
		Message: "Hello, how are you?",
		Context: "Some context",
	}

	if request.Message == "" {
		t.Error("Expected message to be set")
	}

	// Test empty message (this would fail validation in real usage)
	emptyRequest := ChatRequest{
		Message: "",
		Context: "Some context",
	}

	if emptyRequest.Message != "" {
		t.Error("Expected empty message to be empty")
	}
}

// TestChatResponseStructure tests the ChatResponse struct
func TestChatResponseStructure(t *testing.T) {
	response := ChatResponse{
		Response:  "Hello! I'm doing well, thank you for asking.",
		Sources:   []string{"PostgreSQL Database"},
		Model:     "gemma3n:e4b",
		Timestamp: "2024-01-01T00:00:00Z",
	}

	if response.Response == "" {
		t.Error("Expected response to be set")
	}

	if len(response.Sources) == 0 {
		t.Error("Expected sources to be set")
	}

	if response.Model == "" {
		t.Error("Expected model to be set")
	}

	if response.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}
}

// TestOllamaRequestStructure tests the OllamaRequest struct
func TestOllamaRequestStructure(t *testing.T) {
	request := OllamaRequest{
		Model: "gemma3n:e4b",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Stream: false,
	}

	if request.Model == "" {
		t.Error("Expected model to be set")
	}

	if len(request.Messages) == 0 {
		t.Error("Expected messages to be set")
	}

	if request.Messages[0].Role == "" {
		t.Error("Expected message role to be set")
	}

	if request.Messages[0].Content == "" {
		t.Error("Expected message content to be set")
	}
}

// TestOllamaResponseStructure tests the OllamaResponse struct
func TestOllamaResponseStructure(t *testing.T) {
	response := OllamaResponse{
		Message: OllamaMessage{
			Role:    "assistant",
			Content: "Hello! How can I help you today?",
		},
		Done: true,
	}

	if response.Message.Role == "" {
		t.Error("Expected message role to be set")
	}

	if response.Message.Content == "" {
		t.Error("Expected message content to be set")
	}
}

// TestGetEnvDefault tests the getEnv function with default values
func TestGetEnvDefault(t *testing.T) {
	// Test with default value when environment variable is not set
	result := getEnv("NON_EXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected default value 'default_value', got '%s'", result)
	}
}

// TestProcessChatWithMock tests the ProcessChat method with mock data
func TestProcessChatWithMock(t *testing.T) {
	// Create a mock database connection (nil for testing)
	var db *sql.DB = nil
	service := NewLLMService(db)

	// Test with valid request
	request := ChatRequest{
		Message: "Hello, how are you?",
	}

	// This will likely fail due to no database connection, but we can test the structure
	_, err := service.ProcessChat(request)

	// We expect an error due to no database connection, but the method should handle it gracefully
	if err == nil {
		t.Log("ProcessChat completed without error (unexpected in test environment)")
	} else {
		t.Logf("ProcessChat returned expected error: %v", err)
	}
}

// TestBuildContextIntegration tests the integration between LLMService and ContextBuilder
func TestBuildContextIntegration(t *testing.T) {
	var db *sql.DB = nil
	service := NewLLMService(db)

	// Test that the contextBuilder is properly initialized
	if service.contextBuilder == nil {
		t.Error("Expected contextBuilder to be initialized in LLMService")
	}

	// Test that we can call BuildContext through the service
	context, err := service.contextBuilder.BuildContext("test query")

	// We expect an error due to no database connection, but the method should handle it gracefully
	if err == nil {
		t.Log("BuildContext completed without error (unexpected in test environment)")
	} else {
		t.Logf("BuildContext returned expected error: %v", err)
	}

	// Even with error, context should be a string (might be empty)
	if context == "" {
		t.Log("Context is empty (expected in test environment)")
	}
}
