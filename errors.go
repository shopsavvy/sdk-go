package shopsavvy

import "fmt"

// ErrorDetails carries the structured error detail the API returns with every error
// (error_detail.code, error_detail.action, error_detail.docs_url, meta.request_id) plus the
// Retry-After response header. It is embedded in every API error type, so the fields are
// available as err.Code, err.Action and so on. Fields are empty when an older server did not
// send them; RetryAfterSeconds is nil unless the response carried a Retry-After header, and is
// only meaningful for ERR_RATE_LIMIT.
type ErrorDetails struct {
	// Code is the machine-readable error code, for example "ERR_RATE_LIMIT". Empty if absent.
	Code string
	// Action is the one-line fix the API suggests.
	Action string
	// DocsURL links to the documentation for this error code.
	DocsURL string
	// RequestID identifies the request in ShopSavvy's logs; quote it when contacting support.
	RequestID string
	// RetryAfterSeconds is the parsed Retry-After header, or nil.
	RetryAfterSeconds *int
}

// APIError represents a general API error
type APIError struct {
	Message    string
	StatusCode int
	ErrorDetails
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// AuthenticationError represents an authentication failure
type AuthenticationError struct {
	Message    string
	StatusCode int
	ErrorDetails
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error (%d): %s", e.StatusCode, e.Message)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Message    string
	StatusCode int
	ErrorDetails
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not found error (%d): %s", e.StatusCode, e.Message)
}

// ValidationError represents a request validation error
type ValidationError struct {
	Message    string
	StatusCode int
	ErrorDetails
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error (%d): %s", e.StatusCode, e.Message)
}

// RateLimitError represents a rate limit exceeded error
type RateLimitError struct {
	Message    string
	StatusCode int
	ErrorDetails
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit error (%d): %s", e.StatusCode, e.Message)
}

// CreditsExhaustedError is returned for ERR_CREDITS_EXHAUSTED (HTTP 429). Do not retry:
// credits come back when the billing period resets or the plan is upgraded. It is deliberately
// a different type from RateLimitError so retry-on-rate-limit code cannot loop on it.
type CreditsExhaustedError struct {
	Message    string
	StatusCode int
	ErrorDetails
}

func (e *CreditsExhaustedError) Error() string {
	return fmt.Sprintf("Credits exhausted error (%d): %s", e.StatusCode, e.Message)
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
