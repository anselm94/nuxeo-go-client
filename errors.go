package nuxeo

import "fmt"

// APIError represents an error returned by the Nuxeo API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("nuxeo api error: %d %s", e.StatusCode, e.Message)
}

// ErrAuthFailed is returned when authentication fails.
var ErrAuthFailed = &APIError{StatusCode: 401, Message: "authentication failed"}
