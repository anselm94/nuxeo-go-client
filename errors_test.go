package nuxeo

import "testing"

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found"}
	got := err.Error()
	want := "nuxeo api error: 404 not found"
	if got != want {
		t.Errorf("APIError.Error got %q, want %q", got, want)
	}
}

func TestErrAuthFailed(t *testing.T) {
	if ErrAuthFailed.StatusCode != 401 {
		t.Errorf("ErrAuthFailed.StatusCode got %d, want 401", ErrAuthFailed.StatusCode)
	}
	if ErrAuthFailed.Message != "authentication failed" {
		t.Errorf("ErrAuthFailed.Message got %q, want %q", ErrAuthFailed.Message, "authentication failed")
	}
}
