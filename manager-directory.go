package nuxeo

import "log/slog"

type DirectoryManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}
