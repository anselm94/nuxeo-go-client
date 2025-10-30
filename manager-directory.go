package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo-go-client/internal"
)

// directoryManager provides methods to interact with Nuxeo Directory endpoints.
// See: https://doc.nuxeo.com/rest-api/1/directory-endpoint/
type directoryManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

// FetchDirectories retrieves all directories from the Nuxeo server.
// Maps to GET /directory.
func (dm *directoryManager) FetchDirectories(ctx context.Context, options *nuxeoRequestOptions) (*Directories, error) {
	path := internal.PathApiV1 + "/directory"
	res, err := dm.client.NewRequest(ctx, options).SetResult(&Directories{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directories", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Directories), nil
}

// FetchDirectoryEntries retrieves all entries for a given directory.
// Maps to GET /directory/{directoryName}.
// Supports pagination and sorting via SortedPaginationOptions.
func (dm *directoryManager) FetchDirectoryEntries(ctx context.Context, directoryName string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*DirectoryEntries, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName)

	if query := paginationOptions.QueryParams(); query != nil {
		path += "?" + query.Encode()
	}

	res, err := dm.client.NewRequest(ctx, options).SetResult(&DirectoryEntries{}).SetError(&NuxeoError{}).Get(path)
	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directory entries", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DirectoryEntries), nil
}

// CreateDirectoryEntry creates a new entry in the specified directory.
// Maps to POST /directory/{directoryName}.
func (dm *directoryManager) CreateDirectoryEntry(ctx context.Context, directoryName string, entry DirectoryEntry, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName)
	res, err := dm.client.NewRequest(ctx, options).SetBody(entry).SetResult(&DirectoryEntry{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to create directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DirectoryEntry), nil
}

// FetchDirectoryEntry retrieves a specific entry from a directory by id.
// Maps to GET /directory/{directoryName}/{entryId}.
func (dm *directoryManager) FetchDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetResult(&DirectoryEntry{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DirectoryEntry), nil
}

// UpdateDirectoryEntry updates an existing entry in the specified directory.
// Maps to PUT /directory/{directoryName}/{entryId}.
func (dm *directoryManager) UpdateDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, entry DirectoryEntry, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetBody(entry).SetResult(&DirectoryEntry{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to update directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DirectoryEntry), nil
}

// DeleteDirectoryEntry deletes an entry from the specified directory by id.
// Maps to DELETE /directory/{directoryName}/{entryId}.
func (dm *directoryManager) DeleteDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to delete directory entry", slog.String("error", err.Error()))
		return err
	}
	return nil
}
