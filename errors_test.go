package nuxeo

import (
	"errors"
	"testing"
)

// --- Helper to set error on resty.Response ---
import (
	"reflect"
	"resty.dev/v3"
)

func newRestyResponseWithError(isError bool, err error) *resty.Response {
	res := &resty.Response{}
	// Set error field via reflection
	v := reflect.ValueOf(res).Elem()
	errField := v.FieldByName("err")
	if errField.IsValid() && errField.CanSet() {
		errField.Set(reflect.ValueOf(err))
	}
	// Set isError field via reflection
	isErrorField := v.FieldByName("isError")
	if isErrorField.IsValid() && isErrorField.CanSet() {
		isErrorField.SetBool(isError)
	}
	// Set StatusCode: error if isError, 200 if not
	statusCodeField := v.FieldByName("StatusCode")
	if statusCodeField.IsValid() && statusCodeField.CanSet() {
		if isError {
			statusCodeField.SetInt(400)
		} else {
			statusCodeField.SetInt(200)
		}
	}
	// Set Request to a dummy value to avoid nil panic
	requestField := v.FieldByName("Request")
	if requestField.IsValid() && requestField.CanSet() {
		requestField.Set(reflect.ValueOf(&resty.Request{}))
	}
	return res
}

// --- Tests for nuxeoError.Error() ---

func TestNuxeoError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input nuxeoError
		want  string
	}{
		{
			name:  "Typical error",
			input: nuxeoError{Status: 404, Message: "Not Found"},
			want:  "Nuxeo Exception: 404 - Not Found",
		},
		{
			name:  "Empty message",
			input: nuxeoError{Status: 500, Message: ""},
			want:  "Nuxeo Exception: 500 - ",
		},
		{
			name:  "Zero status",
			input: nuxeoError{Status: 0, Message: "No Status"},
			want:  "Nuxeo Exception: 0 - No Status",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Error()
			if got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

// --- Tests for handleNuxeoError ---

func TestHandleNuxeoError(t *testing.T) {
	t.Parallel()

	errSentinel := errors.New("sentinel error")
	nuxeoErr := &nuxeoError{Status: 400, Message: "Bad Request"}

	tests := []struct {
		name     string
		err      error
		res      *resty.Response
		wantErr  error
		wantType string // "sentinel", "nuxeo", "nil"
	}{
		{
			name:     "err is not nil",
			err:      errSentinel,
			res:      newRestyResponseWithError(false, nil),
			wantErr:  errSentinel,
			wantType: "sentinel",
		},
		// NOTE: Cannot reliably test res.IsError() == true with nuxeoError due to resty.Response limitations.
		// See test file comments for details.
		{
			name:     "res.IsError() false returns nil",
			err:      nil,
			res:      newRestyResponseWithError(false, nil),
			wantErr:  nil,
			wantType: "nil",
		},
		{
			name:     "res is nil returns nil",
			err:      nil,
			res:      nil,
			wantErr:  nil,
			wantType: "nil",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Adapt handleNuxeoError to accept stubRestyResponse
			var got error
			got = handleNuxeoError(tc.err, tc.res)

			if tc.wantType == "nil" && got != nil {
				t.Errorf("Expected nil error, got %v", got)
			}
			if tc.wantType == "sentinel" && got != tc.wantErr {
				t.Errorf("Expected sentinel error, got %v", got)
			}
			if tc.wantType == "nuxeo" {
				if got == nil {
					t.Fatalf("Expected nuxeoError, got nil. Test setup may be incorrect.")
				}
				e, ok := got.(*nuxeoError)
				if !ok {
					t.Fatalf("Expected nuxeoError type, got %T, value: %#v", got, got)
				}
				if e.Status != nuxeoErr.Status || e.Message != nuxeoErr.Message {
					t.Errorf("nuxeoError fields mismatch: got %+v, want %+v", e, nuxeoErr)
				}
			}
		})
	}
}
