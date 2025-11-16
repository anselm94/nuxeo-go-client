package nuxeo

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/anselm94/nuxeo-go-client/internal"
)

type batchUploadManager struct {

	// internal

	client *NuxeoClient
	logger *slog.Logger
}

// batchUpload represents a file in a Nuxeo batch upload session.
// It contains metadata about the file, chunking status, and upload progress.
// See: https://doc.nuxeo.com/nxdoc/batch-upload-endpoint/
type batchUpload struct {
	Name             string `json:"name"`
	BatchId          string `json:"batchId"`
	FileIdx          string `json:"fileIdx"`
	UploadType       string `json:"uploadType"`
	UploadedSize     string `json:"uploadedSize"`
	UploadedChunkIds []int  `json:"uploadedChunkIds"`
	ChunkCount       int    `json:"chunkCount"`
}

// CreateBatch initializes a new batch upload session with the default handler.
func (bum *batchUploadManager) CreateBatch(ctx context.Context, options *nuxeoRequestOptions) (*batchUpload, error) {
	path := internal.PathApiV1 + "/upload/new/default"
	res, err := bum.client.NewRequest(ctx, options).SetResult(&batchUpload{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to create batch", "error", err, "status", res.StatusCode())
		return nil, err
	}
	return res.Result().(*batchUpload), nil
}

// FetchBatchUploads gets information about all files in a batch.
func (bum *batchUploadManager) FetchBatchUploads(ctx context.Context, batchId string, options *nuxeoRequestOptions) (*[]batchUpload, error) {
	path := internal.PathApiV1 + "/upload/" + batchId
	res, err := bum.client.NewRequest(ctx, options).SetResult(&[]batchUpload{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to fetch batch uploads", "error", err, "status", res.StatusCode())
		return nil, err
	}
	return res.Result().(*[]batchUpload), nil
}

// FetchBatchUpload gets information about a specific file in a batch.
func (bum *batchUploadManager) FetchBatchUpload(ctx context.Context, batchId string, fileIdx string, options *nuxeoRequestOptions) (*batchUpload, error) {
	path := internal.PathApiV1 + "/upload/" + batchId + "/" + fileIdx
	res, err := bum.client.NewRequest(ctx, options).SetResult(&batchUpload{}).SetError(&NuxeoError{}).Get(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to fetch batch upload", "error", err, "status", res.StatusCode())
		return nil, err
	}
	return res.Result().(*batchUpload), nil
}

// CancelBatch deletes a batch upload session and all associated files.
// Maps to DELETE /upload/{batchId}. Returns error if deletion fails.
// See: https://doc.nuxeo.com/nxdoc/batch-upload-endpoint/#delete-a-batch-upload-session
func (bum *batchUploadManager) CancelBatch(ctx context.Context, batchId string, options *nuxeoRequestOptions) error {
	path := internal.PathApiV1 + "/upload/" + batchId
	res, err := bum.client.NewRequest(ctx, options).SetError(&NuxeoError{}).Delete(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to cancel batch", "error", err, "status", res.StatusCode())
		return err
	}
	return nil
}

// ExecuteBatchUploads runs a Nuxeo Automation operation using all blobs in a batch as input.
// Maps to POST /upload/{batchId}/execute/{operationId}. The operation is applied to all files in the batch.
// See: https://doc.nuxeo.com/nxdoc/batch-upload-endpoint/#execute-an-operation-on-batch-blobs
func (bum *batchUploadManager) ExecuteBatchUploads(ctx context.Context, batchId string, operation operation, out any, options *nuxeoRequestOptions) (any, error) {
	path := internal.PathApiV1 + "/upload/" + batchId + "/execute/" + operation.operationId
	res, err := bum.client.NewRequest(ctx, options).SetBody(operation.payload()).SetResult(out).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to execute batch uploads", "error", err, "status", res.StatusCode())
		return nil, err
	}
	return res.Result(), nil
}

// ExecuteBatchUpload executes an Automation operation using a specific file in a batch as input.
func (bum *batchUploadManager) ExecuteBatchUpload(ctx context.Context, batchId string, fileIdx string, operation operation, out any, options *nuxeoRequestOptions) (any, error) {
	path := internal.PathApiV1 + "/upload/" + batchId + "/" + fileIdx + "/execute/" + operation.operationId
	res, err := bum.client.NewRequest(ctx, options).SetBody(operation.payload()).SetResult(out).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to execute batch upload", "error", err, "status", res.StatusCode())
		return nil, err
	}
	return res.Result(), nil
}

// Upload uploads a file, setting all required headers.
func (bum *batchUploadManager) Upload(ctx context.Context, batchId string, fileIdx int, blob *blob, options *nuxeoRequestOptions) (*batchUpload, error) {
	path := internal.PathApiV1 + "/upload/" + batchId + "/" + strconv.Itoa(fileIdx)

	request := bum.client.NewRequest(ctx, options).
		SetHeader("X-Upload-Type", "normal").
		SetHeader("X-File-Name", blob.Filename).
		SetHeader("X-File-Type", blob.MimeType).
		SetHeader("X-File-Size", fmt.Sprintf("%d", blob.Size())).
		SetContentLength(true).
		SetHeader(internal.HeaderContentLength, fmt.Sprintf("%d", blob.Size())).
		SetContentType(internal.HeaderValueOctetStream)

	res, err := request.SetBody(blob).SetResult(&batchUpload{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to upload file to batch", "error", err, "status", res.StatusCode())
		return nil, err
	}
	return res.Result().(*batchUpload), nil
}

// Upload uploads a chunk to a batch, setting all required headers.
func (bum *batchUploadManager) UploadAsChunk(ctx context.Context, batchId string, fileIdx int, chunkIdx int, totalChunks int, blob *blob, options *nuxeoRequestOptions) (*batchUpload, error) {
	path := internal.PathApiV1 + "/upload/" + batchId + "/" + strconv.Itoa(fileIdx)

	request := bum.client.NewRequest(ctx, options).
		SetHeader("X-Upload-Type", "chunked").
		SetHeader("X-File-Name", blob.Filename).
		SetHeader("X-File-Type", blob.MimeType).
		SetHeader("X-File-Size", fmt.Sprintf("%d", blob.Size())).
		SetHeader("X-Upload-Chunk-Index", fmt.Sprintf("%d", chunkIdx)).
		SetHeader("X-Upload-Chunk-Count", fmt.Sprintf("%d", totalChunks)).
		SetContentLength(true).
		SetHeader(internal.HeaderContentLength, fmt.Sprintf("%d", blob.Size())).
		SetContentType(internal.HeaderValueOctetStream)

	res, err := request.SetBody(blob).SetResult(&batchUpload{}).SetError(&NuxeoError{}).Post(path)

	if err := handleNuxeoError(err, res); err != nil {
		bum.logger.Error("Failed to upload file to batch", "error", err, "status", res.StatusCode())
		return nil, err
	}
	return res.Result().(*batchUpload), nil
}
