package nuxeo

import "fmt"

// NuxeoError represents an error returned by the Nuxeo API.
type NuxeoError struct {
	EntityType string `json:"entity-type"`
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Stack      string `json:"stack"`
}

func (e *NuxeoError) Error() string {
	return fmt.Sprintf("Nuxeo API Exception: %d - %s", e.Status, e.Message)
}

// ErrAuthFailed is returned when authentication fails.
var ErrAuthFailed = &NuxeoError{Status: 401, Message: "authentication failed"}
