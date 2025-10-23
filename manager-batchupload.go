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

type BatchUpload struct {
	Name             string `json:"name"`
	BatchId          string `json:"batchId"`
	FileIdx          string `json:"fileIdx"`
	UploadType       string `json:"uploadType"`
	UploadedSize     int64  `json:"uploadedSize"`
	UploadedChunkIds []int  `json:"uploadedChunkIds"`
	ChunkCount       int    `json:"chunkCount"`
}

func (bum *BatchUploadManager) CreateBatch(ctx context.Context, totalSize int64, fileCount int, options *nuxeoRequestOptions) (*BatchUpload, error) {
	return nil, nil
}

func (bum *BatchUploadManager) FetchBatchUploads(ctx context.Context, batchId string, options *nuxeoRequestOptions) (*[]BatchUpload, error) {
	return nil, nil
}

func (bum *BatchUploadManager) FetchBatchUpload(ctx context.Context, batchId string, fileIdx string, options *nuxeoRequestOptions) (*BatchUpload, error) {
	return nil, nil
}

func (bum *BatchUploadManager) CancelBatch(ctx context.Context, batchId string, options *nuxeoRequestOptions) error {
	return nil
}

func (bum *BatchUploadManager) ExecuteBatchUploads(ctx context.Context, batchId string, operationId string, operationPayload operationPayload, options *nuxeoRequestOptions) (any, error) {
	return nil, nil
}

func (bum *BatchUploadManager) ExecuteBatchUpload(ctx context.Context, batchId string, fileIdx string, operationId string, operationPayload operationPayload, options *nuxeoRequestOptions) (any, error) {
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

func (bum *BatchUploadManager) Upload(ctx context.Context, batchId string, fileIdx string, blob io.Reader, uploadOptions UploadOptions, options *nuxeoRequestOptions) (*BatchUpload, error) {
	return nil, nil
}
