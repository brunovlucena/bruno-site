package security

import (
	"strings"
	"testing"
)

// =============================================================================
// 🧪 SECURITY TESTS
// =============================================================================

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal string",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "String with HTML",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "String with null bytes",
			input:    "Hello\x00World",
			expected: "HelloWorld",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with control characters",
			input:    "Hello\x07World",
			expected: "HelloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateAndSanitizeTitle(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Valid title",
			input:       "My Project",
			expectError: false,
		},
		{
			name:        "Empty title",
			input:       "",
			expectError: true,
		},
		{
			name:        "Title with XSS",
			input:       "<script>alert('xss')</script>",
			expectError: true,
		},
		{
			name:        "Title with JavaScript",
			input:       "javascript:alert('xss')",
			expectError: true,
		},
		{
			name:        "Title too long",
			input:       strings.Repeat("a", MaxTitleLength+1),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateAndSanitizeTitle(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateAndSanitizeTitle(%q) error = %v, expectError %v", tt.input, err, tt.expectError)
			}
		})
	}
}

func TestValidateAndSanitizeURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		fieldName   string
		expectError bool
	}{
		{
			name:        "Valid HTTPS URL",
			input:       "https://example.com",
			fieldName:   "url",
			expectError: false,
		},
		{
			name:        "Valid HTTP URL",
			input:       "http://example.com",
			fieldName:   "url",
			expectError: false,
		},
		{
			name:        "Empty URL",
			input:       "",
			fieldName:   "url",
			expectError: false, // URLs can be optional
		},
		{
			name:        "JavaScript URL",
			input:       "javascript:alert('xss')",
			fieldName:   "url",
			expectError: true,
		},
		{
			name:        "Data URL",
			input:       "data:text/html,<script>alert('xss')</script>",
			fieldName:   "url",
			expectError: true,
		},
		{
			name:        "File URL",
			input:       "file:///etc/passwd",
			fieldName:   "url",
			expectError: true,
		},
		{
			name:        "URL without scheme",
			input:       "example.com",
			fieldName:   "url",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateAndSanitizeURL(tt.input, tt.fieldName)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateAndSanitizeURL(%q, %q) error = %v, expectError %v", tt.input, tt.fieldName, err, tt.expectError)
			}
		})
	}
}

func TestValidateAndSanitizeEmail(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Valid email",
			input:       "user@example.com",
			expectError: false,
		},
		{
			name:        "Valid email with subdomain",
			input:       "user@sub.example.com",
			expectError: false,
		},
		{
			name:        "Empty email",
			input:       "",
			expectError: true,
		},
		{
			name:        "Invalid email - no @",
			input:       "userexample.com",
			expectError: true,
		},
		{
			name:        "Invalid email - no domain",
			input:       "user@",
			expectError: true,
		},
		{
			name:        "Invalid email - no local part",
			input:       "@example.com",
			expectError: true,
		},
		{
			name:        "Email too long",
			input:       string(make([]byte, MaxEmailLength+1)) + "@example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateAndSanitizeEmail(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateAndSanitizeEmail(%q) error = %v, expectError %v", tt.input, err, tt.expectError)
			}
		})
	}
}

func TestValidateInteger(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		fieldName   string
		min         int
		max         int
		expectError bool
	}{
		{
			name:        "Valid integer",
			input:       "42",
			fieldName:   "id",
			min:         1,
			max:         100,
			expectError: false,
		},
		{
			name:        "Empty input",
			input:       "",
			fieldName:   "id",
			min:         1,
			max:         100,
			expectError: true,
		},
		{
			name:        "Non-numeric input",
			input:       "abc",
			fieldName:   "id",
			min:         1,
			max:         100,
			expectError: true,
		},
		{
			name:        "Value too low",
			input:       "0",
			fieldName:   "id",
			min:         1,
			max:         100,
			expectError: true,
		},
		{
			name:        "Value too high",
			input:       "101",
			fieldName:   "id",
			min:         1,
			max:         100,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateInteger(tt.input, tt.fieldName, tt.min, tt.max)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateInteger(%q, %q, %d, %d) error = %v, expectError %v", tt.input, tt.fieldName, tt.min, tt.max, err, tt.expectError)
			}
		})
	}
}

func TestContainsSQLInjectionPattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Normal input",
			input:    "Hello World",
			expected: false,
		},
		{
			name:     "SQL injection - UNION SELECT",
			input:    "UNION SELECT * FROM users",
			expected: true,
		},
		{
			name:     "SQL injection - DROP TABLE",
			input:    "DROP TABLE users",
			expected: true,
		},
		{
			name:     "SQL injection - DELETE FROM",
			input:    "DELETE FROM users",
			expected: true,
		},
		{
			name:     "SQL injection - INSERT INTO",
			input:    "INSERT INTO users",
			expected: true,
		},
		{
			name:     "SQL injection - UPDATE SET",
			input:    "UPDATE users SET",
			expected: true,
		},
		{
			name:     "SQL injection - ALTER TABLE",
			input:    "ALTER TABLE users",
			expected: true,
		},
		{
			name:     "SQL injection - CREATE TABLE",
			input:    "CREATE TABLE users",
			expected: true,
		},
		{
			name:     "SQL injection - EXEC",
			input:    "EXEC sp_help",
			expected: true,
		},
		{
			name:     "SQL injection - Comment",
			input:    "SELECT * FROM users -- comment",
			expected: true,
		},
		{
			name:     "SQL injection - Block comment",
			input:    "SELECT * FROM users /* comment */",
			expected: true,
		},
		{
			name:     "Case insensitive",
			input:    "union select",
			expected: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsSQLInjectionPattern(tt.input)
			if result != tt.expected {
				t.Errorf("containsSQLInjectionPattern(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSecureCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "Equal strings",
			a:        "password123",
			b:        "password123",
			expected: true,
		},
		{
			name:     "Different strings",
			a:        "password123",
			b:        "password456",
			expected: false,
		},
		{
			name:     "Empty strings",
			a:        "",
			b:        "",
			expected: true,
		},
		{
			name:     "One empty string",
			a:        "password123",
			b:        "",
			expected: false,
		},
		{
			name:     "Different lengths",
			a:        "password123",
			b:        "password",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SecureCompare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("SecureCompare(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 🏃‍♂️ BENCHMARK TESTS
// =============================================================================

func BenchmarkSanitizeString(b *testing.B) {
	input := "<script>alert('xss')</script>Hello World"
	for i := 0; i < b.N; i++ {
		SanitizeString(input)
	}
}

func BenchmarkValidateAndSanitizeTitle(b *testing.B) {
	input := "My Project Title"
	for i := 0; i < b.N; i++ {
		ValidateAndSanitizeTitle(input)
	}
}

func BenchmarkContainsSQLInjectionPattern(b *testing.B) {
	input := "UNION SELECT * FROM users"
	for i := 0; i < b.N; i++ {
		containsSQLInjectionPattern(input)
	}
}

func BenchmarkSecureCompare(b *testing.B) {
	a := "password123"
	b_val := "password123"
	for i := 0; i < b.N; i++ {
		SecureCompare(a, b_val)
	}
}
