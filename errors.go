package shopsavvy

import "fmt"

// APIError represents a general API error
type APIError struct {
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// AuthenticationError represents an authentication failure
type AuthenticationError struct {
	Message    string
	StatusCode int
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error (%d): %s", e.StatusCode, e.Message)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Message    string
	StatusCode int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not found error (%d): %s", e.StatusCode, e.Message)
}

// ValidationError represents a request validation error
type ValidationError struct {
	Message    string
	StatusCode int
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error (%d): %s", e.StatusCode, e.Message)
}

// RateLimitError represents a rate limit exceeded error
type RateLimitError struct {
	Message    string
	StatusCode int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit error (%d): %s", e.StatusCode, e.Message)
}

// NetworkError represents a network connectivity error
type NetworkError struct {
	Message string
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network error: %s", e.Message)
}

// TimeoutError represents a request timeout error
type TimeoutError struct {
	Message string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("Timeout error: %s", e.Message)
}