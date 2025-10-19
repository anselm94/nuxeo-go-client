package internal

import (
	"net/http"
)

// NewHTTPClient returns a configured http.Client for use by the SDK.
func NewHTTPClient() *http.Client {
	return &http.Client{}
}
