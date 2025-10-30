package nuxeo

import (
	"fmt"

	"resty.dev/v3"
)

// NuxeoError represents an error returned by the Nuxeo API.
// It includes status code, message, and stack trace from the server response.
type NuxeoError struct {
	entity
	Status     int    `json:"status"`
	Message    string `json:"message"`
	StackTrace string `json:"stacktrace"`
}

// Error returns a formatted string describing the Nuxeo error.
func (e *NuxeoError) Error() string {
	return fmt.Sprintf("Nuxeo Exception: %d - %s", e.Status, e.Message)
}

// handleNuxeoError inspects the error and HTTP response, returning a nuxeoError if the response indicates an error.
// Returns nil if no error is present.
func handleNuxeoError(err error, res *resty.Response) error {
	if err != nil {
		return err
	}
	if res == nil {
		return nil
	}
	if res.IsError() {
		if nuxeoErr, ok := res.Error().(*NuxeoError); ok {
			return nuxeoErr
		}
		if err, ok := res.Error().(error); ok {
			return err
		}
		return fmt.Errorf("unknown error type: %T", res.Error())
	}
	return nil
}
