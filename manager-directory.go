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

func (dm *DirectoryManager) FetchDirectories(ctx context.Context, options *nuxeoRequestOptions) (*Directories, error) {
	return nil, nil
}

func (dm *DirectoryManager) FetchDirectoryEntries(ctx context.Context, directoryName string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*DirectoryEntries, error) {
	return nil, nil
}

func (dm *DirectoryManager) CreateDirectoryEntry(ctx context.Context, directoryName string, entry DirectoryEntry, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	return nil, nil
}

func (dm *DirectoryManager) FetchDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	return nil, nil
}

func (dm *DirectoryManager) UpdateDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, entry DirectoryEntry, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	return nil, nil
}

func (dm *DirectoryManager) DeleteDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) error {
	return nil
}
