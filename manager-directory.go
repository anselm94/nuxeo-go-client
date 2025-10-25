package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo/internal"
)

type DirectoryManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (dm *DirectoryManager) FetchDirectories(ctx context.Context, options *nuxeoRequestOptions) (*Directories, error) {
	path := internal.PathApiV1 + "/directory"
	res, err := dm.client.NewRequest(ctx, options).SetResult(&Directories{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directories", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*Directories), nil
}

func (dm *DirectoryManager) FetchDirectoryEntries(ctx context.Context, directoryName string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*DirectoryEntries, error) {
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

func (dm *DirectoryManager) CreateDirectoryEntry(ctx context.Context, directoryName string, entry DirectoryEntry, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName)
	res, err := dm.client.NewRequest(ctx, options).SetBody(entry).SetResult(&DirectoryEntry{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to create directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DirectoryEntry), nil
}

func (dm *DirectoryManager) FetchDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetResult(&DirectoryEntry{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DirectoryEntry), nil
}

func (dm *DirectoryManager) UpdateDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, entry DirectoryEntry, options *nuxeoRequestOptions) (*DirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetBody(entry).SetResult(&DirectoryEntry{}).SetError(&NuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to update directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*DirectoryEntry), nil
}

func (dm *DirectoryManager) DeleteDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to delete directory entry", slog.String("error", err.Error()))
		return err
	}
	return nil
}
