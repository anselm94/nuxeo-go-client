package nuxeo

import "log/slog"

type UploadManager struct {

	// internal

	client *NuxeoClient
	logger *slog.Logger
}
