package nuxeo

import (
	"fmt"

	"resty.dev/v3"
)

// nuxeoError represents an error returned by the Nuxeo API.
// It includes status code, message, and stack trace from the server response.
type nuxeoError struct {
	entity
	Status     int    `json:"status"`
	Message    string `json:"message"`
	StackTrace string `json:"stacktrace"`
}

// Error returns a formatted string describing the Nuxeo error.
func (e *nuxeoError) Error() string {
	return fmt.Sprintf("Nuxeo Exception: %d - %s", e.Status, e.Message)
}

// handleNuxeoError inspects the error and HTTP response, returning a nuxeoError if the response indicates an error.
// Returns nil if no error is present.
func handleNuxeoError(err error, res *resty.Response) error {
	if err != nil {
		return err
	}
	if res.IsError() {
		return res.Error().(*nuxeoError)
	}
	return nil
}
