package nuxeo

import (
	"context"
	"io"
	"log/slog"
)

type BatchUploadManager struct {

	// internal

	client *NuxeoClient
	logger *slog.Logger
}

func (bum *BatchUploadManager) CreateBatch(ctx context.Context, totalSize int64, fileCount int) (*BatchUpload, error) {
	return nil, nil
}

func (bum *BatchUploadManager) FetchBatchUploads(ctx context.Context, batchId string) (*[]BatchUpload, error) {
	return nil, nil
}

func (bum *BatchUploadManager) FetchBatchUpload(ctx context.Context, batchId string, fileIdx string) (*BatchUpload, error) {
	return nil, nil
}

func (bum *BatchUploadManager) CancelBatch(ctx context.Context, batchId string) error {
	return nil
}

func (bum *BatchUploadManager) ExecuteBatchUploads(ctx context.Context, batchId string, operationId string, operationPayload OperationPayload) (any, error) {
	return nil, nil
}

func (bum *BatchUploadManager) ExecuteBatchUpload(ctx context.Context, batchId string, fileIdx string, operationId string, operationPayload OperationPayload) (any, error) {
	return nil, nil
}

type UploadOptions struct {
	FileName         string
	FileSize         int64
	FileType         string
	UploadType       string
	UploadChunkIndex int64
	TotalChunkCount  int64
}

func (bum *BatchUploadManager) Upload(ctx context.Context, batchId string, fileIdx string, blob io.Reader, options UploadOptions) (*BatchUpload, error) {
	return nil, nil
}
