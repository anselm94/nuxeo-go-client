package nuxeo

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/anselm94/nuxeo/internal"
)

type directoryManager struct {
	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (dm *directoryManager) FetchDirectories(ctx context.Context, options *nuxeoRequestOptions) (*entityDirectories, error) {
	path := internal.PathApiV1 + "/directory"
	res, err := dm.client.NewRequest(ctx, options).SetResult(&entityDirectories{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directories", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDirectories), nil
}

func (dm *directoryManager) FetchDirectoryEntries(ctx context.Context, directoryName string, paginationOptions *SortedPaginationOptions, options *nuxeoRequestOptions) (*entityDirectoryEntries, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName)

	if query := paginationOptions.QueryParams(); query != nil {
		path += "?" + query.Encode()
	}

	res, err := dm.client.NewRequest(ctx, options).SetResult(&entityDirectoryEntries{}).SetError(&nuxeoError{}).Get(path)
	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directory entries", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDirectoryEntries), nil
}

func (dm *directoryManager) CreateDirectoryEntry(ctx context.Context, directoryName string, entry entityDirectoryEntry, options *nuxeoRequestOptions) (*entityDirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName)
	res, err := dm.client.NewRequest(ctx, options).SetBody(entry).SetResult(&entityDirectoryEntry{}).SetError(&nuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to create directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDirectoryEntry), nil
}

func (dm *directoryManager) FetchDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) (*entityDirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetResult(&entityDirectoryEntry{}).SetError(&nuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to fetch directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDirectoryEntry), nil
}

func (dm *directoryManager) UpdateDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, entry entityDirectoryEntry, options *nuxeoRequestOptions) (*entityDirectoryEntry, error) {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetBody(entry).SetResult(&entityDirectoryEntry{}).SetError(&nuxeoError{}).Put(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to update directory entry", slog.String("error", err.Error()))
		return nil, err
	}
	return res.Result().(*entityDirectoryEntry), nil
}

func (dm *directoryManager) DeleteDirectoryEntry(ctx context.Context, directoryName string, directoryEntryId string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/directory/" + url.PathEscape(directoryName) + "/" + url.PathEscape(directoryEntryId)
	res, err := dm.client.NewRequest(ctx, options).SetError(&nuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		dm.logger.Error("Failed to delete directory entry", slog.String("error", err.Error()))
		return err
	}
	return nil
}
