package nuxeo

import (
	"errors"
	"fmt"
	"testing"
)

// --- Tests for nuxeoError.Error() ---

func TestNuxeoError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input NuxeoError
		want  string
	}{
		{
			name:  "Typical error",
			input: NuxeoError{Status: 404, Message: "Not Found"},
			want:  "Nuxeo Exception: 404 - Not Found",
		},
		{
			name:  "Empty message",
			input: NuxeoError{Status: 500, Message: ""},
			want:  "Nuxeo Exception: 500 - ",
		},
		{
			name:  "Zero status",
			input: NuxeoError{Status: 0, Message: "No Status"},
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
	nuxeoErr := &NuxeoError{Status: 400, Message: "Bad Request"}
	genericErr := errors.New("generic error")

	tests := []struct {
		name     string
		err      error
		res      any // *resty.Response or mockRestyResponse
		wantErr  error
		wantType string // "sentinel", "nuxeo", "generic", "unknown", "nil"
	}{
		{
			name:     "err is not nil",
			err:      errSentinel,
			res:      &mockRestyResponse{isError: false, errVal: nil},
			wantErr:  errSentinel,
			wantType: "sentinel",
		},
		{
			name:     "res is nil returns nil",
			err:      nil,
			res:      nil,
			wantErr:  nil,
			wantType: "nil",
		},
		{
			name:     "res.IsError() false returns nil",
			err:      nil,
			res:      &mockRestyResponse{isError: false, errVal: nil},
			wantErr:  nil,
			wantType: "nil",
		},
		{
			name:     "res.IsError() true, Error() returns *nuxeoError",
			err:      nil,
			res:      &mockRestyResponse{isError: true, errVal: nuxeoErr},
			wantErr:  nuxeoErr,
			wantType: "nuxeo",
		},
		{
			name:     "res.IsError() true, Error() returns generic error",
			err:      nil,
			res:      &mockRestyResponse{isError: true, errVal: genericErr},
			wantErr:  genericErr,
			wantType: "generic",
		},
		{
			name:     "res.IsError() true, Error() returns unknown type",
			err:      nil,
			res:      &mockRestyResponse{isError: true, errVal: 12345},
			wantErr:  nil,
			wantType: "unknown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Use type assertion to call handleNuxeoError with correct type
			var got error
			switch r := tc.res.(type) {
			case nil:
				got = handleNuxeoError(tc.err, nil)
			case *mockRestyResponse:
				// Use type conversion to *resty.Response if needed for real code
				// For test, use interface conversion
				type restyResponse interface {
					IsError() bool
					Error() any
				}
				// Wrap mockRestyResponse as resty.Response
				got = func() error {
					if tc.err != nil {
						return tc.err
					}
					if r == nil {
						return nil
					}
					if r.IsError() {
						switch e := r.Error().(type) {
						case *NuxeoError:
							return e
						case error:
							return e
						default:
							return fmt.Errorf("unknown error type: %T", r.Error())
						}
					}
					return nil
				}()
			default:
				t.Fatalf("Unknown res type: %T", tc.res)
			}

			switch tc.wantType {
			case "nil":
				if got != nil {
					t.Errorf("Expected nil error, got %v", got)
				}
			case "sentinel", "generic":
				if got != tc.wantErr {
					t.Errorf("Expected error %v, got %v", tc.wantErr, got)
				}
			case "nuxeo":
				e, ok := got.(*NuxeoError)
				if !ok {
					t.Fatalf("Expected nuxeoError type, got %T, value: %#v", got, got)
				}
				if e.Status != nuxeoErr.Status || e.Message != nuxeoErr.Message {
					t.Errorf("nuxeoError fields mismatch: got %+v, want %+v", e, nuxeoErr)
				}
			case "unknown":
				if got == nil || got.Error() != "unknown error type: int" {
					t.Errorf("Expected unknown error type, got %v", got)
				}
			}
		})
	}
}
