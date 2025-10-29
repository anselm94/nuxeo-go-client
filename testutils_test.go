package nuxeo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"resty.dev/v3"
)

//////////////
//// MOCK ////
//////////////

// mockTransport intercepts HTTP requests and returns controlled responses.
type mockTransport struct {
	respond func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.respond(req)
}

func newMockNuxeoClient(respond func(req *http.Request) (*http.Response, error)) *NuxeoClient {
	options := DefaultNuxeoClientOptions()
	client := NewClient("http://mock", &options)
	mockResty := resty.New()
	mockResty.SetBaseURL("http://mock")
	mockResty.SetTransport(&mockTransport{respond: respond})
	client.restClient = mockResty
	return client
}

// helper to marshal a value to io.ReadCloser
func testMarshalBody(t *testing.T, v any) io.ReadCloser {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshalBody: %v", err)
	}
	return io.NopCloser(bytes.NewReader(b))
}

type mockRestyResponse struct {
	isError bool
	errVal  any
}

// Implements minimal resty.Response interface for testing
func (s *mockRestyResponse) IsError() bool { return s.isError }
func (s *mockRestyResponse) Error() any    { return s.errVal }
