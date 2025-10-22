package nuxeo

import (
	"context"
	"log/slog"
)

type DirectoryManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (dm *DirectoryManager) FetchDirectories(ctx context.Context) (*Directories, error) {
	return nil, nil
}

func (dm *DirectoryManager) FetchDirectoryEntries(ctx context.Context, directoryName string, paginationOptions *SortedPaginationOptions) (*DirectoryEntries, error) {
	return nil, nil
}

func (dm *DirectoryManager) CreateDirectoryEntry(ctx context.Context, directoryName string, entry DirectoryEntry) (*DirectoryEntry, error) {
	return nil, nil
}

func (dm *DirectoryManager) FetchDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string) (*DirectoryEntry, error) {
	return nil, nil
}

func (dm *DirectoryManager) UpdateDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, entry DirectoryEntry) (*DirectoryEntry, error) {
	return nil, nil
}

func (dm *DirectoryManager) DeleteDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string) error {
	return nil
}
