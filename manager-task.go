package nuxeo

import "log/slog"

type TaskManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}
