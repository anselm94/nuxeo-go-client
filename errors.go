package nuxeo

import (
	"fmt"

	"resty.dev/v3"
)

// NuxeoError represents an error returned by the Nuxeo API.
type NuxeoError struct {
	entity
	Status     int    `json:"status"`
	Message    string `json:"message"`
	StackTrace string `json:"stacktrace"`
}

func (e *NuxeoError) Error() string {
	return fmt.Sprintf("Nuxeo Exception: %d - %s", e.Status, e.Message)
}

func handleNuxeoError(err error, res *resty.Response) error {
	if err != nil {
		return err
	}
	if res.IsError() {
		return res.Error().(*NuxeoError)
	}
	return nil
}
