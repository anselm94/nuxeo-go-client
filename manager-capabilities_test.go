package nuxeo

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"strings"
	"testing"
)

type testEntityCapabilities struct {
	Server     struct{ Version string }
	Cluster    struct{ Enabled bool }
	Repository struct{ Name string }
}

func TestCapabilitiesManager_FetchCapabilities(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		respond  func(req *http.Request) (*http.Response, error)
		wantErr  bool
		wantCaps *testEntityCapabilities
	}{
		{
			name: "success",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `{"server":{"distributionVersion":"10.10"},"cluster":{"enabled":true},"repository":{"default": {"queryBlobKeys":false}}}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr: false,
			wantCaps: &testEntityCapabilities{
				Server:     struct{ Version string }{Version: "10.10"},
				Cluster:    struct{ Enabled bool }{Enabled: true},
				Repository: struct{ Name string }{Name: "default"},
			},
		},
		{
			name: "client error",
			respond: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
			wantErr:  true,
			wantCaps: nil,
		},
		{
			name: "invalid json",
			respond: func(req *http.Request) (*http.Response, error) {
				body := `invalid`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			wantErr:  true,
			wantCaps: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newMockNuxeoClient(tc.respond)
			cm := &capabilitiesManager{client: client, logger: slog.Default()}
			caps, err := cm.FetchCapabilities(context.Background())
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.wantCaps != nil && caps != nil {
				if caps.Server.DistributionVersion != tc.wantCaps.Server.Version ||
					caps.Cluster.Enabled != tc.wantCaps.Cluster.Enabled ||
					slices.Collect(maps.Keys(caps.Repository))[0] != tc.wantCaps.Repository.Name {
					t.Errorf("unexpected capabilities: got %+v, want %+v", caps, tc.wantCaps)
				}
			}
		})
	}
}
