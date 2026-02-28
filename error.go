package redcap

import "fmt"

// Error codes for REDCap API errors
const (
	ErrCodeInvalidRequest = "INVALID_REQUEST"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeRateLimit      = "RATE_LIMIT"
	ErrCodeServerError    = "SERVER_ERROR"
	ErrCodeUnknown        = "UNKNOWN"
)

// Error represents a REDCap API error.
type Error struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *Error) Error() string {
	return fmt.Sprintf("redcap: %s (%d): %s", e.Code, e.StatusCode, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// IsRetryable returns true if the error is transient and worth retrying.
func (e *Error) IsRetryable() bool {
	return e.Code == ErrCodeRateLimit || e.Code == ErrCodeServerError
}
