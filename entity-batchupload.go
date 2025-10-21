package nuxeo

type BatchUpload struct {
	Name             string `json:"name"`
	BatchId          string `json:"batchId"`
	FileIdx          string `json:"fileIdx"`
	UploadType       string `json:"uploadType"`
	UploadedSize     int64  `json:"uploadedSize"`
	UploadedChunkIds []int  `json:"uploadedChunkIds"`
	ChunkCount       int    `json:"chunkCount"`
}
